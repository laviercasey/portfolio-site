package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL       string   `mapstructure:"DATABASE_URL"`
	JWTSecret         string   `mapstructure:"JWT_SECRET"`
	AdminPasswordHash string   `mapstructure:"ADMIN_PASSWORD_HASH"`
	UploadDir         string   `mapstructure:"UPLOAD_DIR"`
	Port              string   `mapstructure:"PORT"`
	CORSOrigins       []string `mapstructure:"-"`
	corsRaw           string

	AppEnv               string `mapstructure:"APP_ENV"`
	UmamiAPIURL          string `mapstructure:"UMAMI_API_URL"`
	UmamiAPIKey          string `mapstructure:"UMAMI_API_KEY"`
	UmamiWebsiteID       string `mapstructure:"UMAMI_WEBSITE_ID"`
	UmamiCacheTTLSeconds int    `mapstructure:"UMAMI_CACHE_TTL_SECONDS"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	v.SetDefault("UPLOAD_DIR", "./uploads")
	v.SetDefault("PORT", "8080")
	v.SetDefault("CORS_ORIGINS", "http://localhost:3000")

	_ = v.ReadInConfig()

	cfg := &Config{
		DatabaseURL:          v.GetString("DATABASE_URL"),
		JWTSecret:            v.GetString("JWT_SECRET"),
		AdminPasswordHash:    v.GetString("ADMIN_PASSWORD_HASH"),
		UploadDir:            v.GetString("UPLOAD_DIR"),
		Port:                 v.GetString("PORT"),
		corsRaw:              v.GetString("CORS_ORIGINS"),
		AppEnv:               v.GetString("APP_ENV"),
		UmamiAPIURL:          v.GetString("UMAMI_API_URL"),
		UmamiAPIKey:          v.GetString("UMAMI_API_KEY"),
		UmamiWebsiteID:       v.GetString("UMAMI_WEBSITE_ID"),
		UmamiCacheTTLSeconds: v.GetInt("UMAMI_CACHE_TTL_SECONDS"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	origins, err := parseCORSOrigins(cfg.corsRaw)
	if err != nil {
		return nil, err
	}
	cfg.CORSOrigins = origins
	return cfg, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}
	if c.AdminPasswordHash == "" {
		return errors.New("ADMIN_PASSWORD_HASH is required")
	}
	return nil
}

func parseCORSOrigins(raw string) ([]string, error) {
	if raw == "" {
		return []string{"http://localhost:3000"}, nil
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" {
			continue
		}
		if s == "*" {
			return nil, errors.New("CORS_ORIGINS must not contain '*': credentials are enabled, wildcard origin is unsafe")
		}
		origins = append(origins, s)
	}
	return origins, nil
}
