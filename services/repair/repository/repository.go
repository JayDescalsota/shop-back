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
	Status       string    `bun:"status,notnull,default:'running'"`
	RepairStatus string    `bun:"repair_status,notnull,default:'none'"`
	CreatedAt    time.Time `bun:"created_at,notnull"`
	UpdatedAt    time.Time `bun:"updated_at,notnull"`
}

type Appointment struct {
	bun.BaseModel `bun:"table:repair.appointments"`

	ID               string    `bun:"id,pk,type:uuid"`
	TenantID         string    `bun:"tenant_id,notnull,type:uuid"`
	ShopID           *string   `bun:"shop_id,nullzero"`
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
	EndTime          *string   `bun:"end_time,nullzero,type:time"`
	AssignedMechanic string    `bun:"assigned_mechanic"`
	Bay              *string   `bun:"bay,nullzero"`
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
			status TEXT NOT NULL DEFAULT 'running',
			repair_status TEXT NOT NULL DEFAULT 'none',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		ALTER TABLE repair.vehicles ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'running';
		ALTER TABLE repair.vehicles ADD COLUMN IF NOT EXISTS repair_status TEXT NOT NULL DEFAULT 'none';
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
		);
		CREATE TABLE IF NOT EXISTS repair.staff_assignments (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			appointment_id UUID NOT NULL REFERENCES repair.appointments(id) ON DELETE CASCADE,
			staff_id UUID NOT NULL,
			staff_name TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'mechanic',
			status TEXT NOT NULL DEFAULT 'assigned',
			assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			started_at TIMESTAMPTZ,
			completed_at TIMESTAMPTZ,
			total_minutes INTEGER DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		CREATE TABLE IF NOT EXISTS repair.shop_services (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			service_type_id UUID NOT NULL,
			name TEXT NOT NULL,
			code TEXT,
			system TEXT,
			category TEXT,
			estimated_hours DECIMAL(6,2),
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			UNIQUE(tenant_id, service_type_id)
		);
		CREATE TABLE IF NOT EXISTS repair.shop_parts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			name TEXT NOT NULL,
			sku TEXT,
			description TEXT,
			quantity INTEGER NOT NULL DEFAULT 0,
			unit_price DECIMAL(10,2) DEFAULT 0,
			make_id UUID,
			model_id UUID,
			year INT,
			location_id UUID,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		ALTER TABLE repair.shop_parts ADD COLUMN IF NOT EXISTS make_id UUID;
		ALTER TABLE repair.shop_parts ADD COLUMN IF NOT EXISTS model_id UUID;
		ALTER TABLE repair.shop_parts ADD COLUMN IF NOT EXISTS year INT;
		ALTER TABLE repair.shop_parts ADD COLUMN IF NOT EXISTS location_id UUID;
		CREATE TABLE IF NOT EXISTS repair.part_batches (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			part_id UUID NOT NULL REFERENCES repair.shop_parts(id) ON DELETE CASCADE,
			tenant_id UUID NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 0,
			unit_cost DECIMAL(10,2) NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		CREATE TABLE IF NOT EXISTS repair.shop_tools (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			tenant_id UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			quantity INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'available',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		CREATE TABLE IF NOT EXISTS repair.appointment_parts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			appointment_id UUID NOT NULL REFERENCES repair.appointments(id) ON DELETE CASCADE,
			part_id UUID NOT NULL REFERENCES repair.shop_parts(id),
			part_name TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			unit_price DECIMAL(10,2) NOT NULL DEFAULT 0,
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
		a.Status = "queued"
	}
	if a.ShopID == nil || *a.ShopID == "" {
		a.ShopID = strPtr("default")
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
	if v.Status == "" {
		v.Status = "running"
	}
	if v.RepairStatus == "" {
		v.RepairStatus = "none"
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

type StaffAssignment struct {
	bun.BaseModel `bun:"table:repair.staff_assignments"`

	ID            string    `bun:"id,pk,type:uuid"`
	TenantID      string    `bun:"tenant_id,notnull,type:uuid"`
	AppointmentID string    `bun:"appointment_id,notnull,type:uuid"`
	StaffID       string    `bun:"staff_id,notnull,type:uuid"`
	StaffName     string    `bun:"staff_name,notnull"`
	Role          string    `bun:"role,notnull"`
	Status        string    `bun:"status,notnull"`
	AssignedAt    time.Time `bun:"assigned_at,notnull"`
	StartedAt     *time.Time `bun:"started_at,nullzero"`
	CompletedAt   *time.Time `bun:"completed_at,nullzero"`
	TotalMinutes  int       `bun:"total_minutes"`
	Notes         string    `bun:"notes"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"`
}

func (r *Repository) CreateAssignment(ctx context.Context, a *StaffAssignment) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.Status == "" {
		a.Status = "assigned"
	}
	now := time.Now()
	a.AssignedAt = now
	a.CreatedAt = now
	a.UpdatedAt = now
	_, err := r.db.NewInsert().Model(a).Exec(ctx)
	return err
}

func (r *Repository) GetAssignmentByID(ctx context.Context, id string) (*StaffAssignment, error) {
	a := &StaffAssignment{}
	err := r.db.NewSelect().Model(a).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get assignment: %w", err)
	}
	return a, nil
}

func (r *Repository) ListAssignmentsByAppointment(ctx context.Context, appointmentID string) ([]StaffAssignment, error) {
	var as []StaffAssignment
	err := r.db.NewSelect().Model(&as).
		Where("appointment_id = ?", appointmentID).
		Order("assigned_at ASC").
		Scan(ctx)
	return as, err
}

func (r *Repository) ListActiveAssignmentsByStaff(ctx context.Context, staffID string) ([]StaffAssignment, error) {
	var as []StaffAssignment
	err := r.db.NewSelect().Model(&as).
		Where("staff_id = ? AND status IN ('assigned', 'in_progress')", staffID).
		Scan(ctx)
	return as, err
}

func (r *Repository) UpdateAssignment(ctx context.Context, a *StaffAssignment) error {
	a.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(a).Where("id = ?", a.ID).Exec(ctx)
	return err
}

func (r *Repository) DeleteAssignment(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&StaffAssignment{}).Where("id = ?", id).Exec(ctx)
	return err
}

type ShopService struct {
	bun.BaseModel `bun:"table:repair.shop_services"`

	ID             string    `bun:"id,pk,type:uuid"`
	TenantID       string    `bun:"tenant_id,notnull,type:uuid"`
	ServiceTypeID  string    `bun:"service_type_id,notnull,type:uuid"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code"`
	System         string    `bun:"system"`
	Category       string    `bun:"category"`
	EstimatedHours *float64  `bun:"estimated_hours"`
	IsActive       bool      `bun:"is_active,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

type ShopPart struct {
	bun.BaseModel `bun:"table:repair.shop_parts"`

	ID          string    `bun:"id,pk,type:uuid"`
	TenantID    string    `bun:"tenant_id,notnull,type:uuid"`
	Name        string    `bun:"name,notnull"`
	SKU         string    `bun:"sku"`
	Description string    `bun:"description"`
	Quantity    int       `bun:"quantity,notnull"`
	UnitPrice   float64   `bun:"unit_price"`
	MakeID      *string   `bun:"make_id,type:uuid,nullzero"`
	ModelID     *string   `bun:"model_id,type:uuid,nullzero"`
	Year        *int      `bun:"year,nullzero"`
	LocationID  *string   `bun:"location_id,type:uuid,nullzero"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at,notnull"`
}

type PartBatch struct {
	bun.BaseModel `bun:"table:repair.part_batches"`

	ID        string    `bun:"id,pk,type:uuid"`
	PartID    string    `bun:"part_id,notnull,type:uuid"`
	TenantID  string    `bun:"tenant_id,notnull,type:uuid"`
	Quantity  int       `bun:"quantity,notnull"`
	UnitCost  float64   `bun:"unit_cost,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

type ShopTool struct {
	bun.BaseModel `bun:"table:repair.shop_tools"`

	ID          string    `bun:"id,pk,type:uuid"`
	TenantID    string    `bun:"tenant_id,notnull,type:uuid"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	Quantity    int       `bun:"quantity,notnull"`
	Status      string    `bun:"status,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at,notnull"`
}

func (r *Repository) CreateShopService(ctx context.Context, s *ShopService) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	return err
}

func (r *Repository) ListShopServices(ctx context.Context, tenantID string) ([]ShopService, error) {
	var items []ShopService
	err := r.db.NewSelect().Model(&items).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Scan(ctx)
	return items, err
}

func (r *Repository) GetShopService(ctx context.Context, id string) (*ShopService, error) {
	s := &ShopService{}
	err := r.db.NewSelect().Model(s).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get shop service: %w", err)
	}
	return s, nil
}

func (r *Repository) DeleteShopService(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&ShopService{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *Repository) CreateShopPart(ctx context.Context, p *ShopPart) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	return err
}

func (r *Repository) ListShopParts(ctx context.Context, tenantID string) ([]ShopPart, error) {
	var items []ShopPart
	err := r.db.NewSelect().Model(&items).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Scan(ctx)
	return items, err
}

func (r *Repository) GetShopPart(ctx context.Context, id string) (*ShopPart, error) {
	p := &ShopPart{}
	err := r.db.NewSelect().Model(p).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get shop part: %w", err)
	}
	return p, nil
}

func (r *Repository) UpdateShopPart(ctx context.Context, p *ShopPart) error {
	p.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(p).WherePK().Exec(ctx)
	return err
}

func (r *Repository) DeleteShopPart(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&ShopPart{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *Repository) CreatePartBatch(ctx context.Context, b *PartBatch) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	_, err := r.db.NewInsert().Model(b).Exec(ctx)
	return err
}

func (r *Repository) ListPartBatches(ctx context.Context, partID string) ([]PartBatch, error) {
	var items []PartBatch
	err := r.db.NewSelect().Model(&items).
		Where("part_id = ?", partID).
		Order("created_at ASC").
		Scan(ctx)
	return items, err
}

func (r *Repository) GetPartBatch(ctx context.Context, id string) (*PartBatch, error) {
	b := &PartBatch{}
	err := r.db.NewSelect().Model(b).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get part batch: %w", err)
	}
	return b, nil
}

func (r *Repository) UpdatePartBatch(ctx context.Context, b *PartBatch) error {
	b.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(b).WherePK().Exec(ctx)
	return err
}

func (r *Repository) DeletePartBatch(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&PartBatch{}).Where("id = ?", id).Exec(ctx)
	return err
}

type AppointmentPart struct {
	bun.BaseModel `bun:"table:repair.appointment_parts"`

	ID            string    `bun:"id,pk,type:uuid"`
	AppointmentID string    `bun:"appointment_id,notnull,type:uuid"`
	PartID        string    `bun:"part_id,notnull,type:uuid"`
	PartName      string    `bun:"part_name,notnull"`
	Quantity      int       `bun:"quantity,notnull"`
	UnitPrice     float64   `bun:"unit_price,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"`
}

func (r *Repository) CreateAppointmentPart(ctx context.Context, ap *AppointmentPart) error {
	if ap.ID == "" {
		ap.ID = uuid.New().String()
	}
	now := time.Now()
	ap.CreatedAt = now
	ap.UpdatedAt = now
	_, err := r.db.NewInsert().Model(ap).Exec(ctx)
	return err
}

func (r *Repository) ListAppointmentParts(ctx context.Context, appointmentID string) ([]AppointmentPart, error) {
	var items []AppointmentPart
	err := r.db.NewSelect().Model(&items).
		Where("appointment_id = ?", appointmentID).
		Order("created_at ASC").
		Scan(ctx)
	return items, err
}

func (r *Repository) GetAppointmentPart(ctx context.Context, id string) (*AppointmentPart, error) {
	ap := &AppointmentPart{}
	err := r.db.NewSelect().Model(ap).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get appointment part: %w", err)
	}
	return ap, nil
}

func (r *Repository) DeleteAppointmentPart(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&AppointmentPart{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *Repository) CreateShopTool(ctx context.Context, t *ShopTool) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.Status == "" {
		t.Status = "available"
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	_, err := r.db.NewInsert().Model(t).Exec(ctx)
	return err
}

func (r *Repository) ListShopTools(ctx context.Context, tenantID string) ([]ShopTool, error) {
	var items []ShopTool
	err := r.db.NewSelect().Model(&items).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Scan(ctx)
	return items, err
}

func (r *Repository) GetShopTool(ctx context.Context, id string) (*ShopTool, error) {
	t := &ShopTool{}
	err := r.db.NewSelect().Model(t).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get shop tool: %w", err)
	}
	return t, nil
}

func (r *Repository) UpdateShopTool(ctx context.Context, t *ShopTool) error {
	t.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(t).WherePK().Exec(ctx)
	return err
}

func (r *Repository) DeleteShopTool(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&ShopTool{}).Where("id = ?", id).Exec(ctx)
	return err
}

func strPtr(s string) *string {
	return &s
}
