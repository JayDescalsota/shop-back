package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/services/users/model"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

// ---- Tenant ----

func (r *Repository) CreateTenant(ctx context.Context, t *model.Tenant) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.Status == "" {
		t.Status = model.TenantStatusActive
	}
	_, err := r.db.NewInsert().Model(t).Exec(ctx)
	return err
}

func (r *Repository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	t := &model.Tenant{}
	err := r.db.NewSelect().Model(t).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get tenant: %w", err)
	}
	return t, nil
}

func (r *Repository) UpdateTenant(ctx context.Context, t *model.Tenant) error {
	_, err := r.db.NewUpdate().Model(t).Where("id = ?", t.ID).Exec(ctx)
	return err
}

// ---- TenantSettings ----

func (r *Repository) GetTenantSettings(ctx context.Context, tenantID string) (*model.TenantSettings, error) {
	ts := &model.TenantSettings{}
	err := r.db.NewSelect().Model(ts).Where("tenant_id = ?", tenantID).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get tenant settings: %w", err)
	}
	return ts, nil
}

func (r *Repository) UpsertTenantSettings(ctx context.Context, ts *model.TenantSettings) error {
	if ts.ID == "" {
		ts.ID = uuid.New().String()
	}
	_, err := r.db.NewInsert().Model(ts).
		On("CONFLICT (tenant_id) DO UPDATE").
		Set("business_hours = EXCLUDED.business_hours").
		Set("payment_config = EXCLUDED.payment_config").
		Set("notification_config = EXCLUDED.notification_config").
		Set("branding = EXCLUDED.branding").
		Set("features = EXCLUDED.features").
		Exec(ctx)
	return err
}

// ---- UserProfile ----

func (r *Repository) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	p := &model.UserProfile{}
	err := r.db.NewSelect().Model(p).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user profile: %w", err)
	}
	return p, nil
}

func (r *Repository) UpsertUserProfile(ctx context.Context, p *model.UserProfile) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	_, err := r.db.NewInsert().Model(p).
		On("CONFLICT (user_id) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("department = EXCLUDED.department").
		Set("timezone = EXCLUDED.timezone").
		Set("locale = EXCLUDED.locale").
		Set("notification_prefs = EXCLUDED.notification_prefs").
		Exec(ctx)
	return err
}

// ---- Role ----

func (r *Repository) GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	role := &model.Role{}
	err := r.db.NewSelect().Model(role).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get role: %w", err)
	}
	return role, nil
}

func (r *Repository) ListRoles(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.NewSelect().Model(&roles).Order("code ASC").Scan(ctx)
	return roles, err
}

// ---- Permission ----

func (r *Repository) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.NewSelect().Model(&perms).Order("module ASC", "code ASC").Scan(ctx)
	return perms, err
}

func (r *Repository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	var codes []string
	err := r.db.NewSelect().
		Column("p.code").
		TableExpr("role_permissions rp").
		Join("JOIN permissions p ON p.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Order("p.code ASC").
		Scan(ctx, &codes)
	return codes, err
}

// ---- TenantUser (branch assignments) ----

type userBranchRow struct {
	TenantID   string `bun:"tenant_id"`
	TenantName string `bun:"tenant_name"`
	TenantType string `bun:"tenant_type"`
	RoleID     string `bun:"role_id"`
	RoleCode   string `bun:"role_code"`
}

func (r *Repository) GetUserBranches(ctx context.Context, userID string) ([]userBranchRow, error) {
	var results []userBranchRow
	err := r.db.NewSelect().
		Column("t.id AS tenant_id", "t.name AS tenant_name", "t.type AS tenant_type").
		Column("tu.role_id", "r.code AS role_code").
		TableExpr("tenant_users tu").
		Join("JOIN tenants t ON t.id = tu.tenant_id").
		Join("JOIN roles r ON r.id = tu.role_id").
		Where("tu.user_id = ?", userID).
		Where("t.status = ?", model.TenantStatusActive).
		Scan(ctx, &results)
	return results, err
}
