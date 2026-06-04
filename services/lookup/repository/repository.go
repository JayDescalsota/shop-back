package repository

import (
	"context"

	"backend/services/lookup/model"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func createUUID() string {
	return uuid.New().String()
}

func (r *Repository) GetVehicleMakes(ctx context.Context, search string, isActive *bool) ([]model.VehicleMake, error) {
	q := r.db.NewSelect().Model(&model.VehicleMake{}).Order("name ASC")
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if isActive != nil {
		q = q.Where("is_active = ?", *isActive)
	}
	var makes []model.VehicleMake
	err := q.Scan(ctx, &makes)
	return makes, err
}

func (r *Repository) GetVehicleMake(ctx context.Context, id string) (*model.VehicleMake, error) {
	make := &model.VehicleMake{}
	err := r.db.NewSelect().Model(make).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return make, nil
}

func (r *Repository) GetVehicleModels(ctx context.Context, makeID, search string, year int, vehicleType string) ([]model.VehicleModel, error) {
	q := r.db.NewSelect().Model(&model.VehicleModel{}).Where("make_id = ?", makeID)
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if year > 0 {
		q = q.Where("(year_start IS NULL OR year_start <= ?) AND (year_end IS NULL OR year_end >= ?)", year, year)
	}
	if vehicleType != "" {
		q = q.Where("vehicle_type = ?", vehicleType)
	}
	q = q.Order("name ASC")
	var models []model.VehicleModel
	err := q.Scan(ctx, &models)
	return models, err
}

func (r *Repository) GetVehicleModel(ctx context.Context, id string) (*model.VehicleModel, error) {
	m := &model.VehicleModel{}
	err := r.db.NewSelect().Model(m).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Repository) SearchVehicles(ctx context.Context, query string, limit int) ([]model.VehicleModel, error) {
	var models []model.VehicleModel
	err := r.db.NewSelect().
		Model(&models).
		Relation("Make").
		Where("name ILIKE ? OR slug ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Scan(ctx)
	return models, err
}

func (r *Repository) CreateVehicleMake(ctx context.Context, make *model.VehicleMake) error {
	_, err := r.db.NewInsert().Model(make).Exec(ctx)
	return err
}

func (r *Repository) GetServiceTypes(ctx context.Context, category string, isGlobal bool, tenantID string) ([]model.ServiceType, error) {
	q := r.db.NewSelect().Model(&model.ServiceType{}).Where("is_active = ?", true)
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if tenantID != "" {
		q = q.Where("tenant_id = ? OR is_global = ?", tenantID, true)
	} else if isGlobal {
		q = q.Where("is_global = ?", true)
	}
	q = q.Order("name ASC")
	var types []model.ServiceType
	err := q.Scan(ctx, &types)
	return types, err
}

func (r *Repository) GetServiceType(ctx context.Context, id string) (*model.ServiceType, error) {
	st := &model.ServiceType{}
	err := r.db.NewSelect().Model(st).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (r *Repository) GetDiagnosticCodes(ctx context.Context, system, search string) ([]model.DiagnosticCode, error) {
	q := r.db.NewSelect().Model(&model.DiagnosticCode{}).Where("is_active = ?", true)
	if system != "" {
		q = q.Where("system = ?", system)
	}
	if search != "" {
		q = q.Where("description ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q = q.Order("code ASC")
	var codes []model.DiagnosticCode
	err := q.Scan(ctx, &codes)
	return codes, err
}

func (r *Repository) GetDiagnosticCode(ctx context.Context, id string) (*model.DiagnosticCode, error) {
	dc := &model.DiagnosticCode{}
	err := r.db.NewSelect().Model(dc).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

func (r *Repository) SearchDiagnosticCodes(ctx context.Context, query string) ([]model.DiagnosticCode, error) {
	var codes []model.DiagnosticCode
	err := r.db.NewSelect().
		Model(&codes).
		Where("is_active = ?", true).
		Where("description ILIKE ? OR code ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(20).
		Scan(ctx)
	return codes, err
}

func (r *Repository) GetPartCategories(ctx context.Context, parentID string) ([]model.PartCategory, error) {
	q := r.db.NewSelect().Model(&model.PartCategory{}).Where("is_active = ?", true)
	if parentID != "" {
		q = q.Where("parent_id = ?", parentID)
	} else {
		q = q.Where("parent_id IS NULL")
	}
	q = q.Order("sort_order ASC")
	var cats []model.PartCategory
	err := q.Scan(ctx, &cats)
	return cats, err
}

func (r *Repository) GetPartCategory(ctx context.Context, id string) (*model.PartCategory, error) {
	cat := &model.PartCategory{}
	err := r.db.NewSelect().Model(cat).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *Repository) CheckPartCompatibility(ctx context.Context, partSKU, makeID, modelID string, year int) ([]model.PartCompatibility, error) {
	q := r.db.NewSelect().Model(&model.PartCompatibility{}).Where("part_sku = ?", partSKU)
	if makeID != "" {
		q = q.Where("(make_id = ? OR make_id IS NULL)", makeID)
	}
	if modelID != "" {
		q = q.Where("(model_id = ? OR model_id IS NULL)", modelID)
	}
	if year > 0 {
		q = q.Where("(year_start IS NULL OR year_start <= ?) AND (year_end IS NULL OR year_end >= ?)", year, year)
	}
	var compats []model.PartCompatibility
	err := q.Scan(ctx, &compats)
	return compats, err
}

func (r *Repository) GetFuelTypes(ctx context.Context) ([]model.FuelType, error) {
	var types []model.FuelType
	err := r.db.NewSelect().Model(&types).Order("sort_order ASC").Scan(ctx)
	return types, err
}

func (r *Repository) GetTransmissionTypes(ctx context.Context) ([]model.TransmissionType, error) {
	var types []model.TransmissionType
	err := r.db.NewSelect().Model(&types).Order("sort_order ASC").Scan(ctx)
	return types, err
}

func (r *Repository) GetEngineTypes(ctx context.Context) ([]model.EngineType, error) {
	var types []model.EngineType
	err := r.db.NewSelect().Model(&types).Order("sort_order ASC").Scan(ctx)
	return types, err
}

func (r *Repository) GetLaborRateTiers(ctx context.Context, tenantID string, isGlobal *bool) ([]model.LaborRateTier, error) {
	q := r.db.NewSelect().Model(&model.LaborRateTier{}).Where("is_active = ?", true)
	if tenantID != "" {
		q = q.Where("(tenant_id = ? OR is_global = ?)", tenantID, true)
	} else if isGlobal != nil && *isGlobal {
		q = q.Where("is_global = ?", true)
	}
	q = q.Order("hourly_rate ASC")
	var tiers []model.LaborRateTier
	err := q.Scan(ctx, &tiers)
	return tiers, err
}

func (r *Repository) GetLaborRateTier(ctx context.Context, id string) (*model.LaborRateTier, error) {
	tier := &model.LaborRateTier{}
	err := r.db.NewSelect().Model(tier).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (r *Repository) GetCountries(ctx context.Context) ([]model.Country, error) {
	var countries []model.Country
	err := r.db.NewSelect().Model(&countries).Where("is_active = ?", true).Order("sort_order ASC").Scan(ctx)
	return countries, err
}

func (r *Repository) GetCurrencies(ctx context.Context) ([]model.Currency, error) {
	var currencies []model.Currency
	err := r.db.NewSelect().Model(&currencies).Where("is_active = ?", true).Order("sort_order ASC").Scan(ctx)
	return currencies, err
}
