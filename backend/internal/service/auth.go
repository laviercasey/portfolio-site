package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret    []byte
	passwordHash string
}

func NewAuthService(secret, passwordHash string) *AuthService {
	return &AuthService{
		jwtSecret:    []byte(secret),
		passwordHash: passwordHash,
	}
}

type LoginInput struct {
	Password string `json:"password" validate:"required"`
}

type LoginOutput struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (s *AuthService) Login(_ context.Context, password string) (*LoginOutput, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(s.passwordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	exp := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  exp.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to sign token")
	}

	return &LoginOutput{Token: signed, ExpiresAt: exp}, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
