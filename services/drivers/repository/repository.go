package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrNotFound = errors.New("not found")

type Driver struct {
	bun.BaseModel `bun:"table:drivers.drivers"`

	ID                string    `bun:"id,pk,type:uuid"`
	TenantID          string    `bun:"tenant_id,notnull,type:uuid"`
	Name              string    `bun:"name,notnull"`
	Email             string    `bun:"email"`
	Phone             string    `bun:"phone"`
	Role              string    `bun:"role,notnull,default:'driver'"`
	LicenseNumber     string    `bun:"license_number"`
	LicenseClass      string    `bun:"license_class"`
	LicenseExpiry     string    `bun:"license_expiry"`
	DateOfBirth       string    `bun:"date_of_birth"`
	Address           string    `bun:"address"`
	EmergencyContact  string    `bun:"emergency_contact"`
	EmergencyPhone    string    `bun:"emergency_phone"`
	Status            string    `bun:"status,notnull,default:'active'"`
	AssignedVehicleID *string   `bun:"assigned_vehicle_id,type:uuid,nullzero"`
	Notes             string    `bun:"notes"`
	HireDate          string    `bun:"hire_date"`
	CreatedAt         time.Time `bun:"created_at,notnull"`
	UpdatedAt         time.Time `bun:"updated_at,notnull"`
}

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE SCHEMA IF NOT EXISTS drivers;
		CREATE TABLE IF NOT EXISTS drivers.drivers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			role TEXT NOT NULL DEFAULT 'driver',
			license_number TEXT,
			license_class TEXT,
			license_expiry TEXT,
			date_of_birth TEXT,
			address TEXT,
			emergency_contact TEXT,
			emergency_phone TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			assigned_vehicle_id UUID,
			notes TEXT,
			hire_date TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		ALTER TABLE drivers.drivers ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'driver'
	`)
	return err
}

func (r *Repository) Create(ctx context.Context, d *Driver) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Role == "" {
		d.Role = "driver"
	}
	if d.Status == "" {
		d.Status = "active"
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.NewInsert().Model(d).Exec(ctx)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Driver, error) {
	d := &Driver{}
	err := r.db.NewSelect().Model(d).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get driver: %w", err)
	}
	return d, nil
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID string) ([]Driver, error) {
	var drivers []Driver
	err := r.db.NewSelect().Model(&drivers).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Scan(ctx)
	return drivers, err
}

func (r *Repository) Update(ctx context.Context, d *Driver) error {
	d.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(d).Where("id = ?", d.ID).Exec(ctx)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&Driver{}).Where("id = ?", id).Exec(ctx)
	return err
}
