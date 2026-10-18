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

type VehicleMake struct {
	bun.BaseModel `bun:"table:vehicles.vehicle_makes"`

	ID        string    `bun:"id,pk,type:uuid"`
	Name      string    `bun:"name,notnull"`
	Logo      *string   `bun:"logo"`
	CreatedAt time.Time `bun:"created_at,notnull"`
}

type VehicleModel struct {
	bun.BaseModel `bun:"table:vehicles.vehicle_models"`

	ID        string    `bun:"id,pk,type:uuid"`
	MakeID    string    `bun:"make_id,notnull,type:uuid"`
	Name      string    `bun:"name,notnull"`
	YearStart *int      `bun:"year_start"`
	YearEnd   *int      `bun:"year_end"`
	CreatedAt time.Time `bun:"created_at,notnull"`
}

type Vehicle struct {
	bun.BaseModel `bun:"table:vehicles.vehicles"`

	ID             string    `bun:"id,pk,type:uuid"`
	TenantID       string    `bun:"tenant_id,notnull,type:uuid"`
	OwnerID        string    `bun:"owner_id,notnull,type:uuid"`
	MakeID         *string   `bun:"make_id,type:uuid"`
	ModelID        *string   `bun:"model_id,type:uuid"`
	Year           *int      `bun:"year"`
	Vin            *string   `bun:"vin"`
	LicensePlate   *string   `bun:"license_plate"`
	Color          *string   `bun:"color"`
	EngineType     *string   `bun:"engine_type"`
	Transmission   *string   `bun:"transmission"`
	CurrentMileage int       `bun:"current_mileage,notnull"`
	FuelType       *string   `bun:"fuel_type"`
	PurchaseDate   *string   `bun:"purchase_date"`
	PurchasePrice  *float64  `bun:"purchase_price"`
	ImageURL       *string   `bun:"image_url"`
	Notes          *string   `bun:"notes"`
	Status         string    `bun:"status,notnull"`
	RepairStatus   string    `bun:"repair_status,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

type MaintenanceRecord struct {
	bun.BaseModel `bun:"table:vehicles.maintenance_records"`

	ID               string    `bun:"id,pk,type:uuid"`
	VehicleID        string    `bun:"vehicle_id,notnull,type:uuid"`
	ServiceType      string    `bun:"service_type,notnull"`
	Description      *string   `bun:"description"`
	MileageAtService *int      `bun:"mileage_at_service"`
	Cost             *float64  `bun:"cost"`
	ServiceDate      string    `bun:"service_date,notnull"`
	ProviderName     *string   `bun:"provider_name"`
	Status           string    `bun:"status,notnull"`
	NextDueMileage   *int      `bun:"next_due_mileage"`
	NextDueDate      *string   `bun:"next_due_date"`
	Documents        string    `bun:"documents,type:jsonb"`
	CreatedAt        time.Time `bun:"created_at,notnull"`
	UpdatedAt        time.Time `bun:"updated_at,notnull"`
}

type Expense struct {
	bun.BaseModel `bun:"table:vehicles.expenses"`

	ID          string    `bun:"id,pk,type:uuid"`
	VehicleID   string    `bun:"vehicle_id,notnull,type:uuid"`
	Category    string    `bun:"category,notnull"`
	Amount      float64   `bun:"amount,notnull"`
	Date        string    `bun:"date,notnull"`
	Description *string   `bun:"description"`
	ReceiptURL  *string   `bun:"receipt_url"`
	Recurring   bool      `bun:"recurring,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at,notnull"`
}

type Document struct {
	bun.BaseModel `bun:"table:vehicles.documents"`

	ID         string    `bun:"id,pk,type:uuid"`
	VehicleID  string    `bun:"vehicle_id,notnull,type:uuid"`
	DocType    string    `bun:"doc_type,notnull"`
	Title      string    `bun:"title,notnull"`
	FileURL    string    `bun:"file_url,notnull"`
	ExpiryDate *string   `bun:"expiry_date"`
	CreatedAt  time.Time `bun:"created_at,notnull"`
	UpdatedAt  time.Time `bun:"updated_at,notnull"`
}

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

// ---- Vehicle ----

func (r *Repository) CreateVehicle(ctx context.Context, v *Vehicle) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	if v.Status == "" {
		v.Status = "active"
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

func (r *Repository) ListVehicles(ctx context.Context, tenantID string, ownerID string, page, perPage int) ([]Vehicle, int, error) {
	var vehicles []Vehicle
	q := r.db.NewSelect().Model(&vehicles)

	if tenantID != "" {
		q = q.Where("tenant_id = ?", tenantID)
	}
	if ownerID != "" {
		q = q.Where("owner_id = ?", ownerID)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count vehicles: %w", err)
	}

	if page > 0 && perPage > 0 {
		q = q.Offset((page - 1) * perPage).Limit(perPage)
	}

	err = q.Order("created_at DESC").Scan(ctx)
	return vehicles, total, err
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

// ---- VehicleMake ----

func (r *Repository) ListVehicleMakes(ctx context.Context) ([]VehicleMake, error) {
	var makes []VehicleMake
	err := r.db.NewSelect().Model(&makes).Order("name ASC").Scan(ctx)
	return makes, err
}

func (r *Repository) GetVehicleMakeByID(ctx context.Context, id string) (*VehicleMake, error) {
	m := &VehicleMake{}
	err := r.db.NewSelect().Model(m).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get vehicle make: %w", err)
	}
	return m, nil
}

// ---- VehicleModel ----

func (r *Repository) ListVehicleModels(ctx context.Context, makeID string) ([]VehicleModel, error) {
	var models []VehicleModel
	q := r.db.NewSelect().Model(&models)
	if makeID != "" {
		q = q.Where("make_id = ?", makeID)
	}
	err := q.Order("name ASC").Scan(ctx)
	return models, err
}

func (r *Repository) GetVehicleModelByID(ctx context.Context, id string) (*VehicleModel, error) {
	m := &VehicleModel{}
	err := r.db.NewSelect().Model(m).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get vehicle model: %w", err)
	}
	return m, nil
}

// ---- Maintenance ----

func (r *Repository) CreateMaintenance(ctx context.Context, m *MaintenanceRecord) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Status == "" {
		m.Status = "completed"
	}
	if m.Documents == "" {
		m.Documents = "[]"
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	_, err := r.db.NewInsert().Model(m).Exec(ctx)
	return err
}

func (r *Repository) ListMaintenance(ctx context.Context, vehicleID string) ([]MaintenanceRecord, error) {
	var records []MaintenanceRecord
	err := r.db.NewSelect().Model(&records).
		Where("vehicle_id = ?", vehicleID).
		Order("service_date DESC").
		Scan(ctx)
	return records, err
}

// ---- Expenses ----

func (r *Repository) CreateExpense(ctx context.Context, e *Expense) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	_, err := r.db.NewInsert().Model(e).Exec(ctx)
	return err
}

func (r *Repository) ListExpenses(ctx context.Context, vehicleID string) ([]Expense, error) {
	var expenses []Expense
	err := r.db.NewSelect().Model(&expenses).
		Where("vehicle_id = ?", vehicleID).
		Order("date DESC").
		Scan(ctx)
	return expenses, err
}

func (r *Repository) GetExpenseSummary(ctx context.Context, vehicleID string) (map[string]float64, float64, error) {
	type result struct {
		Category string
		Total    float64
	}
	var results []result
	err := r.db.NewSelect().
		TableExpr("vehicles.expenses").
		Column("category").
		ColumnExpr("SUM(amount) AS total").
		Where("vehicle_id = ?", vehicleID).
		Group("category").
		Scan(ctx, &results)
	if err != nil {
		return nil, 0, fmt.Errorf("get expense summary: %w", err)
	}

	summary := make(map[string]float64)
	var grandTotal float64
	for _, r := range results {
		summary[r.Category] = r.Total
		grandTotal += r.Total
	}

	return summary, grandTotal, nil
}

// ---- Documents ----

func (r *Repository) CreateDocument(ctx context.Context, d *Document) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.NewInsert().Model(d).Exec(ctx)
	return err
}

func (r *Repository) ListDocuments(ctx context.Context, vehicleID string) ([]Document, error) {
	var docs []Document
	err := r.db.NewSelect().Model(&docs).
		Where("vehicle_id = ?", vehicleID).
		Order("created_at DESC").
		Scan(ctx)
	return docs, err
}
