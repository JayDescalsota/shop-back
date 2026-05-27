package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ContextKeyTenantID contextKey = "tenant_id"
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUserRole contextKey = "user_role"
	ContextKeyEmail    contextKey = "email"
	ContextKeyToken    contextKey = "token"
)

func GetTenantID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyTenantID).(string); ok {
		return v
	}
	return ""
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return v
	}
	return ""
}

func GetUserRole(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyUserRole).(string); ok {
		return v
	}
	return ""
}

func GetEmail(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyEmail).(string); ok {
		return v
	}
	return ""
}

func SetTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, ContextKeyTenantID, tenantID)
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

func SetUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ContextKeyUserRole, role)
}

func SetEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, ContextKeyEmail, email)
}

func AuthFromRequest(r *http.Request) (userID, tenantID, role, email string, err error) {
	userID = GetUserID(r.Context())
	tenantID = GetTenantID(r.Context())
	role = GetUserRole(r.Context())
	email = GetEmail(r.Context())
	return
}
