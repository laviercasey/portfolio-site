package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"DATABASE_URL",
		"JWT_SECRET",
		"ADMIN_PASSWORD_HASH",
		"UPLOAD_DIR",
		"PORT",
		"CORS_ORIGINS",
		"APP_ENV",
		"UMAMI_API_URL",
		"UMAMI_API_KEY",
		"UMAMI_WEBSITE_ID",
		"UMAMI_CACHE_TTL_SECONDS",
	}
	for _, k := range keys {
		t.Setenv(k, "")
		_ = os.Unsetenv(k)
	}
}

func chdirToTemp(t *testing.T) string {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(orig)
	})
	return dir
}

func setValidSecrets(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/portfolio?sslmode=disable")
	t.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-32chars+")
	t.Setenv("ADMIN_PASSWORD_HASH", "$2a$10$abcdefghijklmnopqrstuv")
}

func TestLoad_DefaultsApplied(t *testing.T) {
	clearConfigEnv(t)
	chdirToTemp(t)
	setValidSecrets(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.UploadDir != "./uploads" {
		t.Errorf("UploadDir default = %q, want %q", cfg.UploadDir, "./uploads")
	}
	if cfg.Port != "8080" {
		t.Errorf("Port default = %q, want %q", cfg.Port, "8080")
	}
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://localhost:3000" {
		t.Errorf("CORSOrigins default = %v, want [http://localhost:3000]", cfg.CORSOrigins)
	}
	if cfg.UmamiCacheTTLSeconds != 0 {
		t.Errorf("UmamiCacheTTLSeconds default = %d, want 0", cfg.UmamiCacheTTLSeconds)
	}
	if cfg.AppEnv != "" {
		t.Errorf("AppEnv default = %q, want empty", cfg.AppEnv)
	}
}

func TestLoad_CustomValuesParsed(t *testing.T) {
	clearConfigEnv(t)
	chdirToTemp(t)
	setValidSecrets(t)
	t.Setenv("UPLOAD_DIR", "/var/uploads")
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", "production")
	t.Setenv("UMAMI_API_URL", "https://umami.example.com")
	t.Setenv("UMAMI_API_KEY", "key-xyz")
	t.Setenv("UMAMI_WEBSITE_ID", "website-1")
	t.Setenv("UMAMI_CACHE_TTL_SECONDS", "300")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.UploadDir != "/var/uploads" {
		t.Errorf("UploadDir = %q, want /var/uploads", cfg.UploadDir)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want 9090", cfg.Port)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("AppEnv = %q, want production", cfg.AppEnv)
	}
	if cfg.UmamiAPIURL != "https://umami.example.com" {
		t.Errorf("UmamiAPIURL = %q", cfg.UmamiAPIURL)
	}
	if cfg.UmamiAPIKey != "key-xyz" {
		t.Errorf("UmamiAPIKey = %q", cfg.UmamiAPIKey)
	}
	if cfg.UmamiWebsiteID != "website-1" {
		t.Errorf("UmamiWebsiteID = %q", cfg.UmamiWebsiteID)
	}
	if cfg.UmamiCacheTTLSeconds != 300 {
		t.Errorf("UmamiCacheTTLSeconds = %d, want 300", cfg.UmamiCacheTTLSeconds)
	}
}

func TestLoad_RequiredMissing(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr string
	}{
		{
			name: "DATABASE_URL missing",
			setup: func(t *testing.T) {
				t.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-32chars+")
				t.Setenv("ADMIN_PASSWORD_HASH", "$2a$10$abcdefghijklmnopqrstuv")
			},
			wantErr: "DATABASE_URL is required",
		},
		{
			name: "JWT_SECRET missing",
			setup: func(t *testing.T) {
				t.Setenv("DATABASE_URL", "postgres://u:p@h/db")
				t.Setenv("ADMIN_PASSWORD_HASH", "$2a$10$abcdefghijklmnopqrstuv")
			},
			wantErr: "JWT_SECRET is required",
		},
		{
			name: "JWT_SECRET too short",
			setup: func(t *testing.T) {
				t.Setenv("DATABASE_URL", "postgres://u:p@h/db")
				t.Setenv("JWT_SECRET", "short-secret")
				t.Setenv("ADMIN_PASSWORD_HASH", "$2a$10$abcdefghijklmnopqrstuv")
			},
			wantErr: "JWT_SECRET must be at least 32 characters",
		},
		{
			name: "ADMIN_PASSWORD_HASH missing",
			setup: func(t *testing.T) {
				t.Setenv("DATABASE_URL", "postgres://u:p@h/db")
				t.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-32chars+")
			},
			wantErr: "ADMIN_PASSWORD_HASH is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearConfigEnv(t)
			chdirToTemp(t)
			tc.setup(t)

			cfg, err := Load()
			if err == nil {
				t.Fatalf("Load() expected error, got cfg=%+v", cfg)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("Load() error = %v, want contains %q", err, tc.wantErr)
			}
		})
	}
}

func TestLoad_CORSOriginsParsed(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    []string
		wantErr bool
	}{
		{
			name: "single origin",
			raw:  "https://example.com",
			want: []string{"https://example.com"},
		},
		{
			name: "multiple origins trimmed",
			raw:  "https://a.com, https://b.com ,https://c.com",
			want: []string{"https://a.com", "https://b.com", "https://c.com"},
		},
		{
			name: "empty parts skipped",
			raw:  "https://a.com,,https://b.com",
			want: []string{"https://a.com", "https://b.com"},
		},
		{
			name:    "wildcard rejected",
			raw:     "https://a.com,*",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearConfigEnv(t)
			chdirToTemp(t)
			setValidSecrets(t)
			t.Setenv("CORS_ORIGINS", tc.raw)

			cfg, err := Load()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Load() expected error for raw=%q, got cfg=%+v", tc.raw, cfg)
				}
				if !strings.Contains(err.Error(), "CORS_ORIGINS") {
					t.Errorf("error = %v, want CORS_ORIGINS mention", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if len(cfg.CORSOrigins) != len(tc.want) {
				t.Fatalf("CORSOrigins len = %d, want %d (%v)", len(cfg.CORSOrigins), len(tc.want), cfg.CORSOrigins)
			}
			for i := range tc.want {
				if cfg.CORSOrigins[i] != tc.want[i] {
					t.Errorf("CORSOrigins[%d] = %q, want %q", i, cfg.CORSOrigins[i], tc.want[i])
				}
			}
		})
	}
}

func TestLoad_IntCoercion(t *testing.T) {
	clearConfigEnv(t)
	chdirToTemp(t)
	setValidSecrets(t)
	t.Setenv("UMAMI_CACHE_TTL_SECONDS", "600")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.UmamiCacheTTLSeconds != 600 {
		t.Errorf("UmamiCacheTTLSeconds = %d, want 600", cfg.UmamiCacheTTLSeconds)
	}
}

func TestLoad_IntCoercionInvalid(t *testing.T) {
	clearConfigEnv(t)
	chdirToTemp(t)
	setValidSecrets(t)
	t.Setenv("UMAMI_CACHE_TTL_SECONDS", "not-a-number")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.UmamiCacheTTLSeconds != 0 {
		t.Errorf("UmamiCacheTTLSeconds = %d, want 0 for invalid input", cfg.UmamiCacheTTLSeconds)
	}
}

func TestLoad_EnvFileRead(t *testing.T) {
	clearConfigEnv(t)
	dir := chdirToTemp(t)

	envContent := "DATABASE_URL=postgres://u:p@h/db\n" +
		"JWT_SECRET=this-is-a-very-long-jwt-secret-key-32chars+\n" +
		"ADMIN_PASSWORD_HASH=$2a$10$abcdefghijklmnopqrstuv\n" +
		"PORT=7070\n" +
		"APP_ENV=staging\n"

	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(envContent), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != "7070" {
		t.Errorf("Port = %q, want 7070", cfg.Port)
	}
	if cfg.AppEnv != "staging" {
		t.Errorf("AppEnv = %q, want staging", cfg.AppEnv)
	}
	if cfg.DatabaseURL == "" {
		t.Errorf("DatabaseURL should be loaded from .env")
	}
}

func TestLoad_StructShape(t *testing.T) {
	clearConfigEnv(t)
	chdirToTemp(t)
	setValidSecrets(t)
	t.Setenv("PORT", "8080")
	t.Setenv("UPLOAD_DIR", "./uploads")
	t.Setenv("CORS_ORIGINS", "http://localhost:3000")
	t.Setenv("APP_ENV", "test")
	t.Setenv("UMAMI_API_URL", "https://u.example.com")
	t.Setenv("UMAMI_API_KEY", "k")
	t.Setenv("UMAMI_WEBSITE_ID", "w")
	t.Setenv("UMAMI_CACHE_TTL_SECONDS", "42")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	checks := map[string]string{
		"DatabaseURL":       cfg.DatabaseURL,
		"JWTSecret":         cfg.JWTSecret,
		"AdminPasswordHash": cfg.AdminPasswordHash,
		"UploadDir":         cfg.UploadDir,
		"Port":              cfg.Port,
		"AppEnv":            cfg.AppEnv,
		"UmamiAPIURL":       cfg.UmamiAPIURL,
		"UmamiAPIKey":       cfg.UmamiAPIKey,
		"UmamiWebsiteID":    cfg.UmamiWebsiteID,
	}
	for field, val := range checks {
		if val == "" {
			t.Errorf("field %s is empty, want populated", field)
		}
	}
	if cfg.UmamiCacheTTLSeconds != 42 {
		t.Errorf("UmamiCacheTTLSeconds = %d, want 42", cfg.UmamiCacheTTLSeconds)
	}
	if len(cfg.CORSOrigins) == 0 {
		t.Errorf("CORSOrigins is empty, want populated")
	}
}
