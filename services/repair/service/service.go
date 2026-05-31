package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/services/repair/model"
	"backend/services/repair/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Migrate(ctx context.Context) error {
	return s.repo.Migrate(ctx)
}

func (s *Service) CreateAppointment(ctx context.Context, input model.CreateAppointmentInput) (*model.Appointment, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if input.CustomerName == "" {
		return nil, errors.New("customerName is required")
	}
	if input.VehicleMake == "" || input.VehicleModel == "" {
		return nil, errors.New("vehicle make and model are required")
	}
	if input.ServiceType == "" {
		return nil, errors.New("serviceType is required")
	}
	if input.ScheduledDate == "" {
		return nil, errors.New("scheduledDate is required")
	}
	if input.StartTime == "" || input.EndTime == "" {
		return nil, errors.New("startTime and endTime are required")
	}

	scheduledDate, err := time.Parse("2006-01-02", input.ScheduledDate)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduledDate, expected YYYY-MM-DD: %w", err)
	}

	repoApt := &repository.Appointment{
		TenantID:      input.TenantID,
		CustomerName:  input.CustomerName,
		VehicleMake:   input.VehicleMake,
		VehicleModel:  input.VehicleModel,
		ServiceType:   input.ServiceType,
		ScheduledDate: scheduledDate,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		Status:        "pending",
	}

	if input.CustomerPhone != nil {
		repoApt.CustomerPhone = *input.CustomerPhone
	}
	if input.CustomerEmail != nil {
		repoApt.CustomerEmail = *input.CustomerEmail
	}
	if input.VehicleYear != nil {
		repoApt.VehicleYear = input.VehicleYear
	}
	if input.VehiclePlate != nil {
		repoApt.VehiclePlate = *input.VehiclePlate
	}
	if input.Description != nil {
		repoApt.Description = *input.Description
	}
	if input.AssignedMechanic != nil {
		repoApt.AssignedMechanic = *input.AssignedMechanic
	}
	if input.Notes != nil {
		repoApt.Notes = *input.Notes
	}
	if err := s.repo.Create(ctx, repoApt); err != nil {
		return nil, fmt.Errorf("create appointment: %w", err)
	}

	return toModel(repoApt), nil
}

func (s *Service) GetAppointment(ctx context.Context, id string) (*model.Appointment, error) {
	repoApt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toModel(repoApt), nil
}

func (s *Service) ListAppointments(ctx context.Context, tenantID string) (*model.AppointmentConnection, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	appts, err := s.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Appointment, len(appts))
	for i, a := range appts {
		items[i] = toModel(&a)
	}

	return &model.AppointmentConnection{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, id, status string) (*model.Appointment, error) {
	repoApt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	repoApt.Status = status
	if err := s.repo.Update(ctx, repoApt); err != nil {
		return nil, fmt.Errorf("update appointment: %w", err)
	}
	return toModel(repoApt), nil
}

func (s *Service) DeleteAppointment(ctx context.Context, id string) (bool, error) {
	if err := s.repo.Delete(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) FindAppointmentByID(ctx context.Context, id string) (*model.Appointment, error) {
	return s.GetAppointment(ctx, id)
}

func (s *Service) CreateCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.Customer, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	c := &repository.Customer{
		TenantID: input.TenantID,
		Name:     input.Name,
	}
	if input.Email != nil {
		c.Email = *input.Email
	}
	if input.Phone != nil {
		c.Phone = *input.Phone
	}
	if input.Address != nil {
		c.Address = *input.Address
	}
	if input.City != nil {
		c.City = *input.City
	}
	if input.State != nil {
		c.State = *input.State
	}
	if input.Zip != nil {
		c.Zip = *input.Zip
	}
	if input.Notes != nil {
		c.Notes = *input.Notes
	}

	if err := s.repo.CreateCustomer(ctx, c); err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	return customerToModel(c), nil
}

func (s *Service) GetCustomer(ctx context.Context, id string) (*model.Customer, error) {
	repoC, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customerToModel(repoC), nil
}

func (s *Service) ListCustomers(ctx context.Context, tenantID string) (*model.CustomerConnection, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	customers, err := s.repo.ListCustomersByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Customer, len(customers))
	for i, c := range customers {
		items[i] = customerToModel(&c)
	}

	return &model.CustomerConnection{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) UpdateCustomer(ctx context.Context, id string, input model.UpdateCustomerInput) (*model.Customer, error) {
	repoC, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		repoC.Name = *input.Name
	}
	if input.Email != nil {
		repoC.Email = *input.Email
	}
	if input.Phone != nil {
		repoC.Phone = *input.Phone
	}
	if input.Address != nil {
		repoC.Address = *input.Address
	}
	if input.City != nil {
		repoC.City = *input.City
	}
	if input.State != nil {
		repoC.State = *input.State
	}
	if input.Zip != nil {
		repoC.Zip = *input.Zip
	}
	if input.Notes != nil {
		repoC.Notes = *input.Notes
	}
	if input.Status != nil {
		repoC.Status = *input.Status
	}

	if err := s.repo.UpdateCustomer(ctx, repoC); err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}

	return customerToModel(repoC), nil
}

func (s *Service) DeleteCustomer(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteCustomer(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) FindCustomerByID(ctx context.Context, id string) (*model.Customer, error) {
	return s.GetCustomer(ctx, id)
}

func (s *Service) CreateVehicle(ctx context.Context, input model.CreateVehicleInput) (*model.Vehicle, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if input.Make == "" || input.Model == "" {
		return nil, errors.New("make and model are required")
	}

	v := &repository.Vehicle{
		TenantID: input.TenantID,
		Make:     input.Make,
		Model:    input.Model,
	}
	if input.CustomerID != nil {
		v.CustomerID = input.CustomerID
	}
	if input.Year != nil {
		v.Year = input.Year
	}
	if input.Vin != nil {
		v.VIN = *input.Vin
	}
	if input.LicensePlate != nil {
		v.LicensePlate = *input.LicensePlate
	}
	if input.Color != nil {
		v.Color = *input.Color
	}
	if input.Notes != nil {
		v.Notes = *input.Notes
	}

	if err := s.repo.CreateVehicle(ctx, v); err != nil {
		return nil, fmt.Errorf("create vehicle: %w", err)
	}

	return vehicleToModel(v), nil
}

func (s *Service) GetVehicle(ctx context.Context, id string) (*model.Vehicle, error) {
	repoV, err := s.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return vehicleToModel(repoV), nil
}

func (s *Service) ListVehicles(ctx context.Context, tenantID string) (*model.VehicleConnection, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	vehicles, err := s.repo.ListVehiclesByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Vehicle, len(vehicles))
	for i, v := range vehicles {
		items[i] = vehicleToModel(&v)
	}

	return &model.VehicleConnection{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) UpdateVehicle(ctx context.Context, id string, input model.UpdateVehicleInput) (*model.Vehicle, error) {
	repoV, err := s.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Make != nil {
		repoV.Make = *input.Make
	}
	if input.Model != nil {
		repoV.Model = *input.Model
	}
	if input.Year != nil {
		repoV.Year = input.Year
	}
	if input.Vin != nil {
		repoV.VIN = *input.Vin
	}
	if input.LicensePlate != nil {
		repoV.LicensePlate = *input.LicensePlate
	}
	if input.Color != nil {
		repoV.Color = *input.Color
	}
	if input.Notes != nil {
		repoV.Notes = *input.Notes
	}

	if err := s.repo.UpdateVehicle(ctx, repoV); err != nil {
		return nil, fmt.Errorf("update vehicle: %w", err)
	}

	return vehicleToModel(repoV), nil
}

func (s *Service) DeleteVehicle(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteVehicle(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) FindVehicleByID(ctx context.Context, id string) (*model.Vehicle, error) {
	return s.GetVehicle(ctx, id)
}

func toModel(r *repository.Appointment) *model.Appointment {
	scheduledDate := r.ScheduledDate.Format("2006-01-02")

	return &model.Appointment{
		ID:               r.ID,
		TenantID:         r.TenantID,
		ShopID:           r.ShopID,
		CustomerName:     r.CustomerName,
		CustomerPhone:    strPtr(r.CustomerPhone),
		CustomerEmail:    strPtr(r.CustomerEmail),
		VehicleMake:      r.VehicleMake,
		VehicleModel:     r.VehicleModel,
		VehicleYear:      r.VehicleYear,
		VehiclePlate:     strPtr(r.VehiclePlate),
		ServiceType:      r.ServiceType,
		Description:      strPtr(r.Description),
		Status:           r.Status,
		ScheduledDate:    scheduledDate,
		StartTime:        r.StartTime,
		EndTime:          r.EndTime,
		AssignedMechanic: strPtr(r.AssignedMechanic),
		Notes:            strPtr(r.Notes),
		CreatedAt:        r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        r.UpdatedAt.Format(time.RFC3339),
	}
}

func customerToModel(c *repository.Customer) *model.Customer {
	return &model.Customer{
		ID:        c.ID,
		TenantID:  c.TenantID,
		Name:      c.Name,
		Email:     strPtr(c.Email),
		Phone:     strPtr(c.Phone),
		Address:   strPtr(c.Address),
		City:      strPtr(c.City),
		State:     strPtr(c.State),
		Zip:       strPtr(c.Zip),
		Notes:     strPtr(c.Notes),
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

func vehicleToModel(v *repository.Vehicle) *model.Vehicle {
	return &model.Vehicle{
		ID:           v.ID,
		TenantID:     v.TenantID,
		CustomerID:   v.CustomerID,
		Make:         v.Make,
		Model:        v.Model,
		Year:         v.Year,
		Vin:          strPtr(v.VIN),
		LicensePlate: strPtr(v.LicensePlate),
		Color:        strPtr(v.Color),
		Notes:        strPtr(v.Notes),
		CreatedAt:    v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    v.UpdatedAt.Format(time.RFC3339),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
