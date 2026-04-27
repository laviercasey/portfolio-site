package service

import (
	"strings"
	"testing"
)

func TestValidateUmamiURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		raw       string
		env       string
		wantErr   bool
		wantHost  string
		errSubstr string
	}{
		{
			name:     "http umami internal dns",
			raw:      "http://umami:3000",
			env:      "",
			wantErr:  false,
			wantHost: "umami:3000",
		},
		{
			name:     "https external fqdn in prod",
			raw:      "https://umami.lavier.tech",
			env:      "production",
			wantErr:  false,
			wantHost: "umami.lavier.tech",
		},
		{
			name:     "https external fqdn with path",
			raw:      "https://analytics.example.com/umami",
			env:      "",
			wantErr:  false,
			wantHost: "analytics.example.com",
		},
		{
			name:     "localhost allowed in dev",
			raw:      "http://localhost:3000",
			env:      "dev",
			wantErr:  false,
			wantHost: "localhost:3000",
		},
		{
			name:     "127.0.0.1 allowed in dev",
			raw:      "http://127.0.0.1:3000",
			env:      "dev",
			wantErr:  false,
			wantHost: "127.0.0.1:3000",
		},
		{
			name:     "ipv6 loopback allowed in dev",
			raw:      "http://[::1]:3000",
			env:      "dev",
			wantErr:  false,
			wantHost: "[::1]:3000",
		},

		{
			name:      "localhost forbidden in prod",
			raw:       "http://localhost:3000",
			env:       "",
			wantErr:   true,
			errSubstr: "loopback",
		},
		{
			name:      "127.0.0.1 forbidden in prod",
			raw:       "http://127.0.0.1",
			env:       "production",
			wantErr:   true,
			errSubstr: "loopback",
		},
		{
			name:      "ipv6 loopback forbidden in prod",
			raw:       "http://[::1]:3000",
			env:       "",
			wantErr:   true,
			errSubstr: "loopback",
		},
		{
			name:      "0.0.0.0 forbidden in prod",
			raw:       "http://0.0.0.0:8080",
			env:       "",
			wantErr:   true,
			errSubstr: "loopback",
		},

		{
			name:      "AWS metadata ip forbidden even in dev",
			raw:       "http://169.254.169.254/latest/meta-data",
			env:       "dev",
			wantErr:   true,
			errSubstr: "metadata",
		},
		{
			name:      "GCP metadata hostname forbidden",
			raw:       "http://metadata.google.internal",
			env:       "",
			wantErr:   true,
			errSubstr: "metadata",
		},
		{
			name:      "GCP short metadata hostname forbidden",
			raw:       "http://metadata",
			env:       "",
			wantErr:   true,
			errSubstr: "metadata",
		},
		{
			name:      "link-local ipv4 forbidden",
			raw:       "http://169.254.1.2",
			env:       "dev",
			wantErr:   true,
			errSubstr: "link-local",
		},

		{
			name:      "file scheme rejected",
			raw:       "file:///etc/passwd",
			env:       "dev",
			wantErr:   true,
			errSubstr: "scheme",
		},
		{
			name:      "ftp scheme rejected",
			raw:       "ftp://umami",
			env:       "",
			wantErr:   true,
			errSubstr: "scheme",
		},
		{
			name:      "gopher scheme rejected",
			raw:       "gopher://umami:70/",
			env:       "",
			wantErr:   true,
			errSubstr: "scheme",
		},

		{
			name:      "empty url",
			raw:       "",
			env:       "dev",
			wantErr:   true,
			errSubstr: "empty",
		},
		{
			name:      "missing host",
			raw:       "http://",
			env:       "dev",
			wantErr:   true,
			errSubstr: "host",
		},
		{
			name:      "unparseable url with control char",
			raw:       "ht\x7ftp://umami",
			env:       "",
			wantErr:   true,
			errSubstr: "parse",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			u, err := validateUmamiURL(tc.raw, tc.env)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (url=%v)", u)
				}
				if tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.errSubstr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if u == nil {
				t.Fatalf("expected non-nil URL on success")
			}
			if u.Host != tc.wantHost {
				t.Fatalf("Host mismatch: got %q, want %q", u.Host, tc.wantHost)
			}
		})
	}
}
