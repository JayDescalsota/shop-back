package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Tenant struct {
	bun.BaseModel `bun:"auth.tenants"`
	ID            string `bun:",pk"`
	Name          string `bun:",notnull"`
	Address       string
	AppID         string `bun:"app_id,notnull"`
}

type TenantRepository interface {
	FindByApps(ctx context.Context, appIDs []string) ([]Tenant, error)
}

type tenantRepository struct {
	db *bun.DB
}

func NewTenantRepository(db *bun.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) FindByApps(ctx context.Context, appIDs []string) ([]Tenant, error) {
	var tenants []Tenant
	err := r.db.NewSelect().Model(&tenants).Where("app_id IN (?)", bun.In(appIDs)).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("find tenants by apps: %w", err)
	}
	return tenants, nil
}
