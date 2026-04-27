package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func makeAuthService(t *testing.T, password string) *AuthService {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return NewAuthService("super-secret-key-123", string(hash))
}

func TestAuthService_Login_Success(t *testing.T) {
	t.Parallel()

	svc := makeAuthService(t, "correct-horse-battery")
	out, err := svc.Login(context.Background(), "correct-horse-battery")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected non-nil output")
	}
	if out.Token == "" {
		t.Error("expected non-empty token")
	}
	if strings.Count(out.Token, ".") != 2 {
		t.Errorf("JWT should have three parts, got %q", out.Token)
	}
	if out.ExpiresAt.Before(time.Now()) {
		t.Errorf("expiry in the past: %v", out.ExpiresAt)
	}
	if out.ExpiresAt.After(time.Now().Add(25 * time.Hour)) {
		t.Errorf("expiry too far in the future: %v", out.ExpiresAt)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pw   string
	}{
		{"empty password", ""},
		{"wrong password", "nope"},
		{"prefix match", "correct-horse"},
		{"case mismatch", "CORRECT-HORSE-BATTERY"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := makeAuthService(t, "correct-horse-battery")
			out, err := svc.Login(context.Background(), tc.pw)
			if err == nil {
				t.Fatalf("expected error, got out=%+v", out)
			}
			if out != nil {
				t.Errorf("expected nil output, got %+v", out)
			}
			if !strings.Contains(err.Error(), "invalid credentials") {
				t.Errorf("error should be 'invalid credentials', got %v", err)
			}
		})
	}
}

func TestAuthService_Login_InvalidHash(t *testing.T) {
	t.Parallel()

	svc := NewAuthService("secret", "not-a-bcrypt-hash")
	_, err := svc.Login(context.Background(), "whatever")
	if err == nil {
		t.Fatal("expected error for invalid hash")
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	t.Parallel()

	svc := makeAuthService(t, "pw")
	out, err := svc.Login(context.Background(), "pw")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	claims, err := svc.ValidateToken(out.Token)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if claims["sub"] != "admin" {
		t.Errorf("sub: got %v, want admin", claims["sub"])
	}
	if claims["role"] != "admin" {
		t.Errorf("role: got %v, want admin", claims["role"])
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	t.Parallel()

	svc := makeAuthService(t, "pw")

	tests := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"garbage", "not.a.token"},
		{"only one dot", "a.b"},
		{"random text", "hello-world"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.ValidateToken(tc.token)
			if err == nil {
				t.Fatalf("expected error for %q", tc.token)
			}
		})
	}
}

func TestAuthService_ValidateToken_WrongSecret(t *testing.T) {
	t.Parallel()

	a := makeAuthService(t, "pw")
	out, err := a.Login(context.Background(), "pw")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	b := NewAuthService("different-secret", a.passwordHash)
	if _, err := b.ValidateToken(out.Token); err == nil {
		t.Fatal("expected error when validating with different secret")
	}
}

func TestAuthService_ValidateToken_RejectsNoneAlg(t *testing.T) {
	t.Parallel()

	claims := jwt.MapClaims{"sub": "admin", "exp": time.Now().Add(time.Hour).Unix()}
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}

	svc := makeAuthService(t, "pw")
	if _, err := svc.ValidateToken(signed); err == nil {
		t.Fatal("expected validate to reject unsigned token")
	}
}

func TestAuthService_ValidateToken_Expired(t *testing.T) {
	t.Parallel()

	svc := makeAuthService(t, "pw")
	claims := jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(-time.Hour).Unix(),
		"iat":  time.Now().Add(-2 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(svc.jwtSecret)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := svc.ValidateToken(signed); err == nil {
		t.Fatal("expected error on expired token")
	}
}

func TestNewAuthService_StoresFields(t *testing.T) {
	t.Parallel()

	s := NewAuthService("s", "h")
	if string(s.jwtSecret) != "s" {
		t.Errorf("jwtSecret: got %q, want s", string(s.jwtSecret))
	}
	if s.passwordHash != "h" {
		t.Errorf("passwordHash: got %q, want h", s.passwordHash)
	}
}
