package middleware

import (
	"context"
	"net/http"

	"backend/packages/errors"
)

func TenantIsolation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-Id")
		if tenantID != "" {
			ctx := SetTenantID(r.Context(), tenantID)
			r = r.WithContext(ctx)
		}
		if appID := r.Header.Get("X-App-Id"); appID != "" {
			r = r.WithContext(SetAppID(r.Context(), appID))
		}
		if branchID := r.Header.Get("X-Branch-Id"); branchID != "" {
			r = r.WithContext(SetBranchID(r.Context(), branchID))
		}
		next.ServeHTTP(w, r)
	})
}

func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := GetTenantID(r.Context())
		if tenantID == "" {
			http.Error(w, "tenant context required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func TenantScoped(ctx context.Context, tenantID string) error {
	if GetTenantID(ctx) != "" && GetTenantID(ctx) != tenantID {
		return errors.TenantMismatch()
	}
	return nil
}
