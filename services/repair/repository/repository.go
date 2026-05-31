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

type Customer struct {
	bun.BaseModel `bun:"table:repair.customers"`

	ID        string    `bun:"id,pk,type:uuid"`
	TenantID  string    `bun:"tenant_id,notnull,type:uuid"`
	Name      string    `bun:"name,notnull"`
	Email     string    `bun:"email"`
	Phone     string    `bun:"phone"`
	Address   string    `bun:"address"`
	City      string    `bun:"city"`
	State     string    `bun:"state"`
	Zip       string    `bun:"zip"`
	Notes     string    `bun:"notes"`
	Status    string    `bun:"status,notnull,default:'active'"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

type Vehicle struct {
	bun.BaseModel `bun:"table:repair.vehicles"`

	ID           string    `bun:"id,pk,type:uuid"`
	TenantID     string    `bun:"tenant_id,notnull,type:uuid"`
	CustomerID   *string   `bun:"customer_id,type:uuid,nullzero"`
	Make         string    `bun:"make,notnull"`
	Model        string    `bun:"model,notnull"`
	Year         *int      `bun:"year"`
	VIN          string    `bun:"vin"`
	LicensePlate string    `bun:"license_plate"`
	Color        string    `bun:"color"`
	Notes        string    `bun:"notes"`
	CreatedAt    time.Time `bun:"created_at,notnull"`
	UpdatedAt    time.Time `bun:"updated_at,notnull"`
}

type Appointment struct {
	bun.BaseModel `bun:"table:repair.appointments"`

	ID               string    `bun:"id,pk,type:uuid"`
	TenantID         string    `bun:"tenant_id,notnull,type:uuid"`
	ShopID           string    `bun:"shop_id,notnull"`
	CustomerName     string    `bun:"customer_name,notnull"`
	CustomerPhone    string    `bun:"customer_phone"`
	CustomerEmail    string    `bun:"customer_email"`
	VehicleMake      string    `bun:"vehicle_make,notnull"`
	VehicleModel     string    `bun:"vehicle_model,notnull"`
	VehicleYear      *int      `bun:"vehicle_year"`
	VehiclePlate     string    `bun:"vehicle_plate"`
	ServiceType      string    `bun:"service_type,notnull"`
	Description      string    `bun:"description"`
	Status           string    `bun:"status,notnull"`
	ScheduledDate    time.Time `bun:"scheduled_date,notnull,type:date"`
	StartTime        string    `bun:"start_time,notnull,type:time"`
	EndTime          string    `bun:"end_time,notnull,type:time"`
	AssignedMechanic string    `bun:"assigned_mechanic"`
	Notes            string    `bun:"notes"`
	CreatedAt        time.Time `bun:"created_at,notnull"`
	UpdatedAt        time.Time `bun:"updated_at,notnull"`
}

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE SCHEMA IF NOT EXISTS repair;
		CREATE TABLE IF NOT EXISTS repair.customers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			address TEXT,
			city TEXT,
			state TEXT,
			zip TEXT,
			notes TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		CREATE TABLE IF NOT EXISTS repair.vehicles (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			customer_id UUID,
			make TEXT NOT NULL,
			model TEXT NOT NULL,
			year INTEGER,
			vin TEXT,
			license_plate TEXT,
			color TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		CREATE TABLE IF NOT EXISTS repair.appointments (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			shop_id TEXT NOT NULL DEFAULT '',
			customer_name TEXT NOT NULL,
			customer_phone TEXT,
			customer_email TEXT,
			vehicle_make TEXT NOT NULL,
			vehicle_model TEXT NOT NULL,
			vehicle_year INTEGER,
			vehicle_plate TEXT,
			service_type TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			scheduled_date DATE NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			assigned_mechanic TEXT,
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}

func (r *Repository) Create(ctx context.Context, a *Appointment) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.Status == "" {
		a.Status = "pending"
	}
	if a.ShopID == "" {
		a.ShopID = "default"
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	_, err := r.db.NewInsert().Model(a).Exec(ctx)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Appointment, error) {
	a := &Appointment{}
	err := r.db.NewSelect().Model(a).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get appointment: %w", err)
	}
	return a, nil
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID string) ([]Appointment, error) {
	var apps []Appointment
	err := r.db.NewSelect().Model(&apps).
		Where("tenant_id = ?", tenantID).
		Order("scheduled_date DESC", "start_time DESC").
		Scan(ctx)
	return apps, err
}

func (r *Repository) Update(ctx context.Context, a *Appointment) error {
	a.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(a).Where("id = ?", a.ID).Exec(ctx)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&Appointment{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *Repository) CreateCustomer(ctx context.Context, c *Customer) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Status == "" {
		c.Status = "active"
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	return err
}

func (r *Repository) GetCustomerByID(ctx context.Context, id string) (*Customer, error) {
	c := &Customer{}
	err := r.db.NewSelect().Model(c).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return c, nil
}

func (r *Repository) ListCustomersByTenant(ctx context.Context, tenantID string) ([]Customer, error) {
	var customers []Customer
	err := r.db.NewSelect().Model(&customers).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Scan(ctx)
	return customers, err
}

func (r *Repository) UpdateCustomer(ctx context.Context, c *Customer) error {
	c.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx)
	return err
}

func (r *Repository) DeleteCustomer(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&Customer{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *Repository) CreateVehicle(ctx context.Context, v *Vehicle) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	_, err := r.db.NewInsert().Model(v).Exec(ctx)
	return err
}

func (r *Repository) GetVehicleByID(ctx context.Context, id string) (*Vehicle, error) {
	v := &Vehicle{}
	err := r.db.NewSelect().Model(v).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get vehicle: %w", err)
	}
	return v, nil
}

func (r *Repository) ListVehiclesByTenant(ctx context.Context, tenantID string) ([]Vehicle, error) {
	var vehicles []Vehicle
	err := r.db.NewSelect().Model(&vehicles).
		Where("tenant_id = ?", tenantID).
		Order("make ASC", "model ASC").
		Scan(ctx)
	return vehicles, err
}

func (r *Repository) UpdateVehicle(ctx context.Context, v *Vehicle) error {
	v.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(v).Where("id = ?", v.ID).Exec(ctx)
	return err
}

func (r *Repository) DeleteVehicle(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&Vehicle{}).Where("id = ?", id).Exec(ctx)
	return err
}
