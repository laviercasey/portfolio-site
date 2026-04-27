package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_PreflightAllowedOrigin(t *testing.T) {
	t.Parallel()

	called := false
	handler := CORS([]string{"https://allowed.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/resource", nil)
	req.Header.Set("Origin", "https://allowed.example")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if called {
		t.Fatal("next handler must not be called on preflight OPTIONS")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://allowed.example" {
		t.Errorf("Access-Control-Allow-Origin = %q; want https://allowed.example", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %q; want true", got)
	}
	methods := rec.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("Access-Control-Allow-Methods should be set on preflight")
	}
}

func TestCORS_PreflightDisallowedOrigin(t *testing.T) {
	t.Parallel()

	called := false
	handler := CORS([]string{"https://allowed.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/resource", nil)
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if called {
		t.Fatal("next handler must not be called on preflight OPTIONS")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got == "https://evil.example" {
		t.Errorf("disallowed origin %q must not be echoed back", got)
	}
}

func TestCORS_SimpleRequestAllowedOrigin(t *testing.T) {
	t.Parallel()

	called := false
	handler := CORS([]string{"https://allowed.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Header().Set("X-Custom", "value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("body"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/resource", nil)
	req.Header.Set("Origin", "https://allowed.example")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler should be called for simple GET")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; want 200", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://allowed.example" {
		t.Errorf("Access-Control-Allow-Origin = %q; want https://allowed.example", got)
	}
	if got := rec.Header().Get("X-Custom"); got != "value" {
		t.Errorf("X-Custom = %q; want value (header passthrough failed)", got)
	}
	if body := rec.Body.String(); body != "body" {
		t.Errorf("body = %q; want body", body)
	}
}

func TestCORS_SimpleRequestDisallowedOrigin(t *testing.T) {
	t.Parallel()

	called := false
	handler := CORS([]string{"https://allowed.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("body"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/resource", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("rs/cors still forwards simple requests to next")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got == "https://evil.example" {
		t.Errorf("CORS must not echo disallowed origin; got %q", got)
	}
}

func TestCORS_NoOriginHeaderPassthrough(t *testing.T) {
	t.Parallel()

	called := false
	handler := CORS([]string{"https://allowed.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("direct"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/resource", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("non-CORS request should pass through")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; want 200", rec.Code)
	}
	if rec.Body.String() != "direct" {
		t.Errorf("body = %q; want direct", rec.Body.String())
	}
}

func TestCORS_MultipleAllowedOrigins(t *testing.T) {
	t.Parallel()

	handler := CORS([]string{"https://a.example", "https://b.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, origin := range []string{"https://a.example", "https://b.example"} {
		origin := origin
		t.Run(origin, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/api/x", nil)
			req.Header.Set("Origin", origin)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != origin {
				t.Errorf("Access-Control-Allow-Origin = %q; want %q", got, origin)
			}
		})
	}
}

func TestCORS_ExposedHeadersAdvertised(t *testing.T) {
	t.Parallel()

	handler := CORS([]string{"https://allowed.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total-Count", "42")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/list", nil)
	req.Header.Set("Origin", "https://allowed.example")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	exposed := rec.Header().Get("Access-Control-Expose-Headers")
	if exposed == "" {
		t.Error("Access-Control-Expose-Headers should be set for allowed origin")
	}
}
