package service

import (
	"context"
	"errors"
	"fmt"

	"backend/services/vehicles/model"
	"backend/services/vehicles/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// ---- Vehicle ----

func (s *Service) CreateVehicle(ctx context.Context, input model.CreateVehicleInput) (*model.Vehicle, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if input.OwnerID == "" {
		return nil, errors.New("ownerId is required")
	}

	v := &repository.Vehicle{
		TenantID:       input.TenantID,
		OwnerID:        input.OwnerID,
		MakeID:         input.MakeID,
		ModelID:        input.ModelID,
		Year:           input.Year,
		Vin:            input.Vin,
		LicensePlate:   input.LicensePlate,
		Color:          input.Color,
		EngineType:     input.EngineType,
		Transmission:   input.Transmission,
		CurrentMileage: 0,
		FuelType:       input.FuelType,
		PurchaseDate:   input.PurchaseDate,
		PurchasePrice:  input.PurchasePrice,
		Notes:          input.Notes,
		Status:         "active",
		RepairStatus:   "none",
	}

	if input.CurrentMileage != nil {
		v.CurrentMileage = *input.CurrentMileage
	}

	if err := s.repo.CreateVehicle(ctx, v); err != nil {
		return nil, fmt.Errorf("create vehicle: %w", err)
	}

	return toModel(v), nil
}

func (s *Service) GetVehicle(ctx context.Context, id string) (*model.Vehicle, error) {
	v, err := s.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toModel(v), nil
}

func (s *Service) ListVehicles(ctx context.Context, tenantID string, ownerID string, page, perPage int) (*model.VehicleConnection, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	vehicles, total, err := s.repo.ListVehicles(ctx, tenantID, ownerID, page, perPage)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Vehicle, len(vehicles))
	for i, v := range vehicles {
		items[i] = toModel(&v)
	}

	return &model.VehicleConnection{
		Items:   items,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

func (s *Service) UpdateVehicle(ctx context.Context, id string, input model.UpdateVehicleInput) (*model.Vehicle, error) {
	v, err := s.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.MakeID != nil {
		v.MakeID = input.MakeID
	}
	if input.ModelID != nil {
		v.ModelID = input.ModelID
	}
	if input.Year != nil {
		v.Year = input.Year
	}
	if input.Vin != nil {
		v.Vin = input.Vin
	}
	if input.LicensePlate != nil {
		v.LicensePlate = input.LicensePlate
	}
	if input.Color != nil {
		v.Color = input.Color
	}
	if input.CurrentMileage != nil {
		v.CurrentMileage = *input.CurrentMileage
	}
	if input.Notes != nil {
		v.Notes = input.Notes
	}

	if err := s.repo.UpdateVehicle(ctx, v); err != nil {
		return nil, fmt.Errorf("update vehicle: %w", err)
	}

	return toModel(v), nil
}

func (s *Service) DeleteVehicle(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteVehicle(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) LogMileage(ctx context.Context, vehicleID string, mileage int, notes string) (*model.Vehicle, error) {
	v, err := s.repo.GetVehicleByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}

	if mileage < v.CurrentMileage {
		return nil, errors.New("new mileage cannot be less than current mileage")
	}

	v.CurrentMileage = mileage
	if notes != "" {
		v.Notes = &notes
	}

	if err := s.repo.UpdateVehicle(ctx, v); err != nil {
		return nil, fmt.Errorf("log mileage: %w", err)
	}

	return toModel(v), nil
}

// ---- VehicleMake ----

func (s *Service) ListVehicleMakes(ctx context.Context) ([]*model.VehicleMake, error) {
	makes, err := s.repo.ListVehicleMakes(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*model.VehicleMake, len(makes))
	for i, m := range makes {
		result[i] = toModelMake(&m)
	}
	return result, nil
}

func (s *Service) GetVehicleMake(ctx context.Context, id string) (*model.VehicleMake, error) {
	m, err := s.repo.GetVehicleMakeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toModelMake(m), nil
}

// ---- VehicleModel ----

func (s *Service) ListVehicleModels(ctx context.Context, makeID string) ([]*model.VehicleModel, error) {
	models, err := s.repo.ListVehicleModels(ctx, makeID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.VehicleModel, len(models))
	for i, m := range models {
		result[i] = toModelModel(&m)
	}
	return result, nil
}

// ---- Maintenance ----

func (s *Service) AddMaintenance(ctx context.Context, input model.AddMaintenanceInput) (*model.VehicleMaintenance, error) {
	if input.VehicleID == "" {
		return nil, errors.New("vehicleId is required")
	}
	if input.ServiceType == "" {
		return nil, errors.New("serviceType is required")
	}

	m := &repository.MaintenanceRecord{
		VehicleID:        input.VehicleID,
		ServiceType:      input.ServiceType,
		Description:      input.Description,
		MileageAtService: input.MileageAtService,
		Cost:             input.Cost,
		ServiceDate:      input.ServiceDate,
		ProviderName:     input.ProviderName,
		Status:           "completed",
		NextDueMileage:   input.NextDueMileage,
		NextDueDate:      input.NextDueDate,
		Documents:        "[]",
	}

	if err := s.repo.CreateMaintenance(ctx, m); err != nil {
		return nil, fmt.Errorf("add maintenance: %w", err)
	}

	return toModelMaintenance(m), nil
}

func (s *Service) ListMaintenance(ctx context.Context, vehicleID string) ([]*model.VehicleMaintenance, error) {
	records, err := s.repo.ListMaintenance(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.VehicleMaintenance, len(records))
	for i, r := range records {
		result[i] = toModelMaintenance(&r)
	}
	return result, nil
}

// ---- Expenses ----

func (s *Service) AddExpense(ctx context.Context, input model.AddExpenseInput) (*model.VehicleExpense, error) {
	if input.VehicleID == "" {
		return nil, errors.New("vehicleId is required")
	}
	if input.Category == "" {
		return nil, errors.New("category is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	e := &repository.Expense{
		VehicleID:   input.VehicleID,
		Category:    input.Category,
		Amount:      input.Amount,
		Date:        input.Date,
		Description: input.Description,
		ReceiptURL:  input.ReceiptURL,
		Recurring:   input.Recurring,
	}

	if err := s.repo.CreateExpense(ctx, e); err != nil {
		return nil, fmt.Errorf("add expense: %w", err)
	}

	return toModelExpense(e), nil
}

func (s *Service) ListExpenses(ctx context.Context, vehicleID string) ([]*model.VehicleExpense, error) {
	expenses, err := s.repo.ListExpenses(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.VehicleExpense, len(expenses))
	for i, e := range expenses {
		result[i] = toModelExpense(&e)
	}
	return result, nil
}

func (s *Service) GetExpenseSummary(ctx context.Context, vehicleID string) (*model.ExpenseSummary, error) {
	summary, grandTotal, err := s.repo.GetExpenseSummary(ctx, vehicleID)
	if err != nil {
		return nil, err
	}

	monthlyAverage := grandTotal / 12

	return &model.ExpenseSummary{
		TotalRepairs:     summary["repair"],
		TotalMaintenance: summary["maintenance"],
		TotalFuel:        summary["fuel"],
		TotalInsurance:   summary["insurance"],
		TotalOther:       summary["other"],
		GrandTotal:       grandTotal,
		MonthlyAverage:   monthlyAverage,
	}, nil
}

// ---- Documents ----

func (s *Service) AddDocument(ctx context.Context, input model.AddDocumentInput) (*model.VehicleDocument, error) {
	if input.VehicleID == "" {
		return nil, errors.New("vehicleId is required")
	}
	if input.DocType == "" {
		return nil, errors.New("docType is required")
	}
	if input.Title == "" {
		return nil, errors.New("title is required")
	}
	if input.FileURL == "" {
		return nil, errors.New("fileUrl is required")
	}

	d := &repository.Document{
		VehicleID:  input.VehicleID,
		DocType:    input.DocType,
		Title:      input.Title,
		FileURL:    input.FileURL,
		ExpiryDate: input.ExpiryDate,
	}

	if err := s.repo.CreateDocument(ctx, d); err != nil {
		return nil, fmt.Errorf("add document: %w", err)
	}

	return toModelDocument(d), nil
}

func (s *Service) ListDocuments(ctx context.Context, vehicleID string) ([]*model.VehicleDocument, error) {
	docs, err := s.repo.ListDocuments(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.VehicleDocument, len(docs))
	for i, d := range docs {
		result[i] = toModelDocument(&d)
	}
	return result, nil
}

// ---- FindEntity ----

func (s *Service) FindVehicleByID(ctx context.Context, id string) (*model.Vehicle, error) {
	return s.GetVehicle(ctx, id)
}

func (s *Service) FindVehicleMakeByID(ctx context.Context, id string) (*model.VehicleMake, error) {
	return s.GetVehicleMake(ctx, id)
}

func (s *Service) FindVehicleModelByID(ctx context.Context, id string) (*model.VehicleModel, error) {
	m, err := s.repo.GetVehicleModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toModelModel(m), nil
}

// ---- Mappers ----

func toModel(v *repository.Vehicle) *model.Vehicle {
	return &model.Vehicle{
		ID:             v.ID,
		TenantID:       v.TenantID,
		OwnerID:        v.OwnerID,
		MakeID:         v.MakeID,
		ModelID:        v.ModelID,
		Year:           v.Year,
		Vin:            v.Vin,
		LicensePlate:   v.LicensePlate,
		Color:          v.Color,
		EngineType:     v.EngineType,
		Transmission:   v.Transmission,
		CurrentMileage: v.CurrentMileage,
		FuelType:       v.FuelType,
		PurchaseDate:   v.PurchaseDate,
		PurchasePrice:  v.PurchasePrice,
		ImageURL:       v.ImageURL,
		Notes:          v.Notes,
		Status:         v.Status,
		RepairStatus:   v.RepairStatus,
		CreatedAt:      v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toModelMake(m *repository.VehicleMake) *model.VehicleMake {
	return &model.VehicleMake{
		ID:   m.ID,
		Name: m.Name,
		Logo: m.Logo,
	}
}

func toModelModel(m *repository.VehicleModel) *model.VehicleModel {
	return &model.VehicleModel{
		ID:        m.ID,
		MakeID:    m.MakeID,
		Name:      m.Name,
		YearStart: m.YearStart,
		YearEnd:   m.YearEnd,
	}
}

func toModelMaintenance(m *repository.MaintenanceRecord) *model.VehicleMaintenance {
	return &model.VehicleMaintenance{
		ID:               m.ID,
		VehicleID:        m.VehicleID,
		ServiceType:      m.ServiceType,
		Description:      m.Description,
		MileageAtService: m.MileageAtService,
		Cost:             m.Cost,
		ServiceDate:      m.ServiceDate,
		ProviderName:     m.ProviderName,
		Status:           m.Status,
		NextDueMileage:   m.NextDueMileage,
		NextDueDate:      m.NextDueDate,
		Documents:        m.Documents,
	}
}

func toModelExpense(e *repository.Expense) *model.VehicleExpense {
	return &model.VehicleExpense{
		ID:          e.ID,
		VehicleID:   e.VehicleID,
		Category:    e.Category,
		Amount:      e.Amount,
		Date:        e.Date,
		Description: e.Description,
		ReceiptURL:  e.ReceiptURL,
		Recurring:   e.Recurring,
	}
}

func toModelDocument(d *repository.Document) *model.VehicleDocument {
	return &model.VehicleDocument{
		ID:         d.ID,
		VehicleID:  d.VehicleID,
		DocType:    d.DocType,
		Title:      d.Title,
		FileURL:    d.FileURL,
		ExpiryDate: d.ExpiryDate,
	}
}
