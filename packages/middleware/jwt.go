package middleware

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret       string
	expiration   time.Duration
	refreshExp   time.Duration
	mu           sync.RWMutex
}

func NewJWTManager(secret string, expiration, refreshExp time.Duration) *JWTManager {
	return &JWTManager{
		secret:     secret,
		expiration: expiration,
		refreshExp: refreshExp,
	}
}

func (m *JWTManager) GenerateToken(userID, tenantID, role, email string, app []string, tenants []TenantInfo) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiration)

	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Email:    email,
		App:      app,
		Tenants:  tenants,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "autolab",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (m *JWTManager) GenerateRefreshToken(userID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.refreshExp)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "autolab",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	return tokenString, expiresAt, err
}
