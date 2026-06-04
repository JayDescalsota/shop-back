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
	ContextKeyAppID    contextKey = "app_id"
	ContextKeyBranchID contextKey = "branch_id"
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

func GetAppID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyAppID).(string); ok {
		return v
	}
	return ""
}

func GetBranchID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyBranchID).(string); ok {
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

func SetAppID(ctx context.Context, appID string) context.Context {
	return context.WithValue(ctx, ContextKeyAppID, appID)
}

func SetBranchID(ctx context.Context, branchID string) context.Context {
	return context.WithValue(ctx, ContextKeyBranchID, branchID)
}

func AuthFromRequest(r *http.Request) (userID, tenantID, role, email string, err error) {
	userID = GetUserID(r.Context())
	tenantID = GetTenantID(r.Context())
	role = GetUserRole(r.Context())
	email = GetEmail(r.Context())
	return
}
