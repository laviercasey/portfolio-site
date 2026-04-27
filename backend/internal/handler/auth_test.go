package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/portfolio/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func newAuthHandler(t *testing.T, plainPassword string) *AuthHandler {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	svc := service.NewAuthService("test-secret-that-is-long-enough", string(hash))
	return NewAuthHandler(svc)
}

func TestAuth_Login_Happy(t *testing.T) {
	t.Parallel()

	h := newAuthHandler(t, "correct-horse")

	body := strings.NewReader(`{"password":"correct-horse"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	var out service.LoginOutput
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Token == "" {
		t.Errorf("empty token")
	}
	if out.ExpiresAt.IsZero() {
		t.Errorf("zero ExpiresAt")
	}

	cookies := rec.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "auth-token" {
			found = true
			if !c.HttpOnly {
				t.Errorf("cookie should be HttpOnly")
			}
			if c.SameSite != http.SameSiteStrictMode {
				t.Errorf("cookie should be SameSiteStrict")
			}
		}
	}
	if !found {
		t.Errorf("auth-token cookie missing")
	}
}

func TestAuth_Login_WrongPassword(t *testing.T) {
	t.Parallel()

	h := newAuthHandler(t, "correct-horse")

	body := strings.NewReader(`{"password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid credentials") {
		t.Errorf("body: %s", rec.Body.String())
	}
}

func TestAuth_Login_MalformedJSON(t *testing.T) {
	t.Parallel()

	h := newAuthHandler(t, "correct-horse")

	body := strings.NewReader(`{not-json`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid request body") {
		t.Errorf("body: %s", rec.Body.String())
	}
}

func TestAuth_Login_MissingPassword(t *testing.T) {
	t.Parallel()

	h := newAuthHandler(t, "correct-horse")

	body := strings.NewReader(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Password") {
		t.Errorf("expected validation error about Password field: %s", rec.Body.String())
	}
}

func TestAuth_Login_EmptyBody(t *testing.T) {
	t.Parallel()

	h := newAuthHandler(t, "correct-horse")

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}
