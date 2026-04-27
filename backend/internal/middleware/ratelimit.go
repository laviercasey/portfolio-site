package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type bucket struct {
	tokens    float64
	lastCheck time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*bucket
	rate     float64
	capacity float64
}

func NewRateLimiter(ctx context.Context, rate, capacity float64) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*bucket),
		rate:     rate,
		capacity: capacity,
	}
	go rl.cleanup(ctx)
	return rl
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ClientIP(r)
		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.clients[ip]
	if !ok {
		rl.clients[ip] = &bucket{tokens: rl.capacity - 1, lastCheck: time.Now()}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.capacity {
		b.tokens = rl.capacity
	}
	b.lastCheck = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (rl *RateLimiter) cleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			cutoff := time.Now().Add(-10 * time.Minute)
			for ip, b := range rl.clients {
				if b.lastCheck.Before(cutoff) {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

var trustedProxyCIDRs = func() []*net.IPNet {
	blocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
		"fc00::/7",
	}
	out := make([]*net.IPNet, 0, len(blocks))
	for _, b := range blocks {
		_, ipnet, err := net.ParseCIDR(b)
		if err != nil {
			continue
		}
		out = append(out, ipnet)
	}
	return out
}()

func isTrustedProxy(remoteHost string) bool {
	ip := net.ParseIP(strings.TrimSpace(remoteHost))
	if ip == nil {
		return false
	}
	for _, block := range trustedProxyCIDRs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if real := r.Header.Get("X-Real-IP"); real != "" && isTrustedProxy(host) {
		return strings.TrimSpace(real)
	}
	return host
}
