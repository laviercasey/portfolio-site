package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/portfolio/backend/internal/service"
)

const testJWTSecret = "test-secret-for-middleware-auth-tests"

func newAuthService(t *testing.T) *service.AuthService {
	t.Helper()
	return service.NewAuthService(testJWTSecret, "$2a$10$7EqJtq98hPqEX7fNZaFWoOa8Q0w8T.2aO6oP1rJ0BpzZg3OKqKj7K")
}

func signToken(t *testing.T, secret string, claims jwt.MapClaims, method jwt.SigningMethod) string {
	t.Helper()
	tok := jwt.NewWithClaims(method, claims)
	s, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func nextOKHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if called != nil {
			*called = true
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})
}

func TestAuth_MissingAuthorizationHeader(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	called := false
	handler := Auth(auth)(nextOKHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d; want 401", rec.Code)
	}
	if called {
		t.Fatal("next handler should not be called")
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q; want application/json", ct)
	}
	if !strings.Contains(rec.Body.String(), "missing authorization header") {
		t.Errorf("body = %q; want missing authorization header", rec.Body.String())
	}
}

func TestAuth_MalformedAuthorizationHeader(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		header string
	}{
		{"single token no scheme", "justatoken"},
		{"wrong scheme basic", "Basic dXNlcjpwYXNz"},
		{"empty bearer value", "Bearer"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			auth := newAuthService(t)
			called := false
			handler := Auth(auth)(nextOKHandler(&called))

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.header)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d; want 401", rec.Code)
			}
			if called {
				t.Fatal("next handler should not be called")
			}
			if !strings.Contains(rec.Body.String(), "invalid authorization format") {
				t.Errorf("body = %q; want invalid authorization format", rec.Body.String())
			}
		})
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	called := false
	handler := Auth(auth)(nextOKHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d; want 401", rec.Code)
	}
	if called {
		t.Fatal("next handler should not be called")
	}
	if !strings.Contains(rec.Body.String(), "invalid or expired token") {
		t.Errorf("body = %q; want invalid or expired token", rec.Body.String())
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	token := signToken(t, testJWTSecret, jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(-1 * time.Hour).Unix(),
		"iat":  time.Now().Add(-2 * time.Hour).Unix(),
	}, jwt.SigningMethodHS256)

	called := false
	handler := Auth(auth)(nextOKHandler(&called))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d; want 401", rec.Code)
	}
	if called {
		t.Fatal("next handler should not be called for expired token")
	}
	if !strings.Contains(rec.Body.String(), "invalid or expired token") {
		t.Errorf("body = %q; want invalid or expired token", rec.Body.String())
	}
}

func TestAuth_WrongSignature(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	token := signToken(t, "different-secret", jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	}, jwt.SigningMethodHS256)

	called := false
	handler := Auth(auth)(nextOKHandler(&called))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d; want 401", rec.Code)
	}
	if called {
		t.Fatal("next handler should not be called when signature invalid")
	}
}

func TestAuth_NonAdminRoleForbidden(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	token := signToken(t, testJWTSecret, jwt.MapClaims{
		"sub":  "user",
		"role": "user",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	}, jwt.SigningMethodHS256)

	called := false
	handler := Auth(auth)(nextOKHandler(&called))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d; want 403", rec.Code)
	}
	if called {
		t.Fatal("next handler should not be called for non-admin")
	}
	if !strings.Contains(rec.Body.String(), "forbidden") {
		t.Errorf("body = %q; want forbidden", rec.Body.String())
	}
}

func TestAuth_MissingRoleClaimForbidden(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	token := signToken(t, testJWTSecret, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}, jwt.SigningMethodHS256)

	called := false
	handler := Auth(auth)(nextOKHandler(&called))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d; want 403", rec.Code)
	}
	if called {
		t.Fatal("next handler should not run when role claim missing")
	}
}

func TestAuth_ValidAdminTokenPassesThrough(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	token := signToken(t, testJWTSecret, jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}, jwt.SigningMethodHS256)

	called := false
	handler := Auth(auth)(nextOKHandler(&called))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler should be called for valid admin token")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("parse body: %v", err)
	}
	if body["ok"] != true {
		t.Errorf("body ok = %v; want true", body["ok"])
	}
}

func TestAuth_LowercaseBearerSchemeAccepted(t *testing.T) {
	t.Parallel()

	auth := newAuthService(t)
	token := signToken(t, testJWTSecret, jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	}, jwt.SigningMethodHS256)

	called := false
	handler := Auth(auth)(nextOKHandler(&called))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next should be called with lowercase bearer scheme")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}
}
