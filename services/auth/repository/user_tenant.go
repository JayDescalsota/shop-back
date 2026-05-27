package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type UserTenant struct {
	bun.BaseModel `bun:"auth.user_tenants"`
	Email         string `bun:",pk"`
	TenantID      string `bun:"tenant_id,pk"`
}

type UserTenantRepository interface {
	FindByEmail(ctx context.Context, email string) ([]Tenant, error)
	Add(ctx context.Context, email, tenantID string) error
	Remove(ctx context.Context, email, tenantID string) error
}

type userTenantRepository struct {
	db *bun.DB
}

func NewUserTenantRepository(db *bun.DB) UserTenantRepository {
	return &userTenantRepository{db: db}
}

func (r *userTenantRepository) FindByEmail(ctx context.Context, email string) ([]Tenant, error) {
	var tenants []Tenant
	err := r.db.NewRaw(
		`SELECT t.* FROM auth.tenants t
		 JOIN auth.user_tenants ut ON ut.tenant_id = t.id
		 WHERE ut.email = ?`, email,
	).Scan(ctx, &tenants)
	if err != nil {
		return nil, fmt.Errorf("find user tenants: %w", err)
	}
	return tenants, nil
}

func (r *userTenantRepository) Add(ctx context.Context, email, tenantID string) error {
	_, err := r.db.NewInsert().Model(&UserTenant{Email: email, TenantID: tenantID}).Exec(ctx)
	if err != nil {
		return fmt.Errorf("add user tenant: %w", err)
	}
	return nil
}

func (r *userTenantRepository) Remove(ctx context.Context, email, tenantID string) error {
	_, err := r.db.NewDelete().
		Model(&UserTenant{}).
		Where("email = ? AND tenant_id = ?", email, tenantID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove user tenant: %w", err)
	}
	return nil
}
