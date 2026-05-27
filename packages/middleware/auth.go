package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type TenantInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	AppID string `json:"app_id"`
}

type Claims struct {
	UserID   string       `json:"user_id"`
	TenantID string       `json:"tenant_id"`
	Role     string       `json:"role"`
	Email    string       `json:"email"`
	App      []string     `json:"app"`
	Tenants  []TenantInfo `json:"tenants"`
	jwt.RegisteredClaims
}

func JWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = SetUserID(ctx, claims.UserID)
			ctx = SetTenantID(ctx, claims.TenantID)
			ctx = SetUserRole(ctx, claims.Role)
			ctx = SetEmail(ctx, claims.Email)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r.Context())
			if role == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			for _, allowed := range roles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}
