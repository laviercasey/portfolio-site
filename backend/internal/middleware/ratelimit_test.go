package middleware

import (
	"net/http"
	"testing"
)

func TestClientIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		xRealIP    string
		want       string
	}{
		{
			name:       "public peer + X-Real-IP spoof → header ignored, return peer host",
			remoteAddr: "1.2.3.4:5678",
			xRealIP:    "9.9.9.9",
			want:       "1.2.3.4",
		},
		{
			name:       "Docker RFC1918 peer + X-Real-IP → trust header",
			remoteAddr: "172.18.0.5:5678",
			xRealIP:    "9.9.9.9",
			want:       "9.9.9.9",
		},
		{
			name:       "Docker RFC1918 peer, no X-Real-IP → return peer host",
			remoteAddr: "172.18.0.5:5678",
			xRealIP:    "",
			want:       "172.18.0.5",
		},
		{
			name:       "IPv4 loopback peer + X-Real-IP → trust header (local dev)",
			remoteAddr: "127.0.0.1:5678",
			xRealIP:    "203.0.113.1",
			want:       "203.0.113.1",
		},
		{
			name:       "IPv6 loopback peer + X-Real-IP → trust header",
			remoteAddr: "[::1]:5678",
			xRealIP:    "203.0.113.1",
			want:       "203.0.113.1",
		},
		{
			name:       "malformed RemoteAddr → fallback to raw value",
			remoteAddr: "garbage",
			xRealIP:    "",
			want:       "garbage",
		},
		{
			name:       "public peer + empty X-Real-IP → return peer host",
			remoteAddr: "1.2.3.4:5678",
			xRealIP:    "",
			want:       "1.2.3.4",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(http.MethodGet, "http://example/", nil)
			if err != nil {
				t.Fatalf("build request: %v", err)
			}
			req.RemoteAddr = tc.remoteAddr
			if tc.xRealIP != "" {
				req.Header.Set("X-Real-IP", tc.xRealIP)
			}

			got := ClientIP(req)
			if got != tc.want {
				t.Errorf("ClientIP(remote=%q, X-Real-IP=%q) = %q; want %q",
					tc.remoteAddr, tc.xRealIP, got, tc.want)
			}
		})
	}
}
