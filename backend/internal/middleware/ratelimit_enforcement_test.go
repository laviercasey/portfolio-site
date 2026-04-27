package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newRL(t *testing.T, rate, capacity float64) *RateLimiter {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return NewRateLimiter(ctx, rate, capacity)
}

func serve(handler http.Handler, remote string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req.RemoteAddr = remote
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestRateLimit_AllowsUpToCapacity(t *testing.T) {
	t.Parallel()

	rl := newRL(t, 0.0001, 5)
	var served int32
	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&served, 1)
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		rec := serve(handler, "203.0.113.100:1111")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d; want 200", i+1, rec.Code)
		}
	}
	if got := atomic.LoadInt32(&served); got != 5 {
		t.Fatalf("served = %d; want 5", got)
	}
}

func TestRateLimit_RejectsOverCapacity(t *testing.T) {
	t.Parallel()

	rl := newRL(t, 0.0001, 3)
	var served int32
	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&served, 1)
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		rec := serve(handler, "203.0.113.200:2222")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d; want 200", i+1, rec.Code)
		}
	}

	rec := serve(handler, "203.0.113.200:2222")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("over-capacity: status = %d; want 429", rec.Code)
	}
	if got := atomic.LoadInt32(&served); got != 3 {
		t.Errorf("served = %d; want 3 (4th should be blocked)", got)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q; want application/json", ct)
	}
	if ra := rec.Header().Get("Retry-After"); ra != "1" {
		t.Errorf("Retry-After = %q; want 1", ra)
	}
	if !strings.Contains(rec.Body.String(), "rate limit exceeded") {
		t.Errorf("body = %q; want rate limit exceeded", rec.Body.String())
	}
}

func TestRateLimit_SeparateBucketsPerIP(t *testing.T) {
	t.Parallel()

	rl := newRL(t, 0.0001, 2)
	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		if rec := serve(handler, "203.0.113.10:1000"); rec.Code != http.StatusOK {
			t.Fatalf("IP1 req %d: status = %d; want 200", i+1, rec.Code)
		}
	}
	if rec := serve(handler, "203.0.113.10:1000"); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("IP1 third: status = %d; want 429", rec.Code)
	}

	for i := 0; i < 2; i++ {
		if rec := serve(handler, "203.0.113.20:1000"); rec.Code != http.StatusOK {
			t.Fatalf("IP2 req %d (separate bucket): status = %d; want 200", i+1, rec.Code)
		}
	}
	if rec := serve(handler, "203.0.113.20:1000"); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("IP2 third: status = %d; want 429", rec.Code)
	}
}

func TestRateLimit_RefillOverTime(t *testing.T) {
	t.Parallel()

	rl := newRL(t, 100, 2)
	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	const ip = "203.0.113.55:1234"
	for i := 0; i < 2; i++ {
		if rec := serve(handler, ip); rec.Code != http.StatusOK {
			t.Fatalf("initial req %d: status = %d; want 200", i+1, rec.Code)
		}
	}
	if rec := serve(handler, ip); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("immediate third: status = %d; want 429", rec.Code)
	}

	time.Sleep(50 * time.Millisecond)

	if rec := serve(handler, ip); rec.Code != http.StatusOK {
		t.Fatalf("after refill: status = %d; want 200 (tokens should have regenerated)", rec.Code)
	}
}

func TestRateLimit_CapacityCap(t *testing.T) {
	t.Parallel()

	rl := newRL(t, 20, 2)
	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	const ip = "203.0.113.77:1234"
	if rec := serve(handler, ip); rec.Code != http.StatusOK {
		t.Fatalf("initial: status = %d; want 200", rec.Code)
	}

	time.Sleep(500 * time.Millisecond)

	for i := 0; i < 2; i++ {
		if rec := serve(handler, ip); rec.Code != http.StatusOK {
			t.Fatalf("after wait req %d: status = %d; want 200", i+1, rec.Code)
		}
	}
	if rec := serve(handler, ip); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("over cap: status = %d; want 429 (bucket must not exceed capacity)", rec.Code)
	}
}

func TestRateLimit_CleanupStopsOnContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	rl := NewRateLimiter(ctx, 1, 1)

	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	if rec := serve(handler, "203.0.113.88:5678"); rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}

	cancel()
	time.Sleep(20 * time.Millisecond)
}
