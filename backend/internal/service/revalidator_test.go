package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/portfolio/backend/internal/testutil"
)

func TestRevalidateClient_Disabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		url    string
		secret string
	}{
		{"both empty", "", ""},
		{"only url", "http://x", ""},
		{"only secret", "", "s"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := NewRevalidateService(tc.url, tc.secret, testutil.SilentLogger())
			if c.Enabled() {
				t.Errorf("Enabled() = true, want false")
			}
			if err := c.Revalidate(context.Background(), []string{"/"}); err != nil {
				t.Errorf("disabled Revalidate should be noop, got %v", err)
			}
		})
	}
}

func TestRevalidateClient_HappyPath(t *testing.T) {
	t.Parallel()

	var gotAuth, gotCT string
	var gotBody map[string][]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotCT = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		_, _ = w.Write([]byte(`{"revalidated":true}`))
	}))
	defer ts.Close()

	c := NewRevalidateService(ts.URL, "secret-xyz", testutil.SilentLogger())
	c.client = &http.Client{Timeout: 2 * time.Second}

	if err := c.Revalidate(context.Background(), []string{"/", "/about"}); err != nil {
		t.Fatalf("Revalidate: %v", err)
	}
	if gotAuth != "Bearer secret-xyz" {
		t.Errorf("auth = %q", gotAuth)
	}
	if gotCT != "application/json" {
		t.Errorf("content-type = %q", gotCT)
	}
	paths := gotBody["paths"]
	if len(paths) != 2 || paths[0] != "/" || paths[1] != "/about" {
		t.Errorf("paths = %v", paths)
	}
}

func TestRevalidateClient_NonOKStatus(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c := NewRevalidateService(ts.URL, "s", testutil.SilentLogger())
	c.client = &http.Client{Timeout: 2 * time.Second}

	err := c.Revalidate(context.Background(), []string{"/"})
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("err = %v, want 401", err)
	}
}

func TestRevalidateClient_NetworkError(t *testing.T) {
	t.Parallel()

	c := NewRevalidateService("http://127.0.0.1:1", "s", testutil.SilentLogger())
	c.client = &http.Client{Timeout: 100 * time.Millisecond}
	if err := c.Revalidate(context.Background(), []string{"/"}); err == nil {
		t.Fatal("expected network error")
	}
}

func TestRevalidateClient_EmptyPaths(t *testing.T) {
	t.Parallel()

	hit := false
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		hit = true
	}))
	defer ts.Close()

	c := NewRevalidateService(ts.URL, "s", testutil.SilentLogger())
	if err := c.Revalidate(context.Background(), nil); err != nil {
		t.Fatalf("nil paths: %v", err)
	}
	if hit {
		t.Error("should not call upstream for empty paths")
	}
}
