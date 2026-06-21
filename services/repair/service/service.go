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
	if input.StartTime == "" {
		return nil, errors.New("startTime is required")
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
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:        "queued",
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
	if input.Bay != nil {
		repoApt.Bay = input.Bay
	}
	if input.AssignedMechanic != nil {
		repoApt.AssignedMechanic = *input.AssignedMechanic
	}
	if input.Notes != nil {
		repoApt.Notes = *input.Notes
	}
	if input.ShopID != nil {
		repoApt.ShopID = input.ShopID
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

func (s *Service) UpdateAppointment(ctx context.Context, id string, input model.UpdateAppointmentInput) (*model.Appointment, error) {
	repoApt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.CustomerName != nil {
		repoApt.CustomerName = *input.CustomerName
	}
	if input.CustomerPhone != nil {
		repoApt.CustomerPhone = *input.CustomerPhone
	}
	if input.CustomerEmail != nil {
		repoApt.CustomerEmail = *input.CustomerEmail
	}
	if input.VehicleMake != nil {
		repoApt.VehicleMake = *input.VehicleMake
	}
	if input.VehicleModel != nil {
		repoApt.VehicleModel = *input.VehicleModel
	}
	if input.VehicleYear != nil {
		repoApt.VehicleYear = input.VehicleYear
	}
	if input.VehiclePlate != nil {
		repoApt.VehiclePlate = *input.VehiclePlate
	}
	if input.ServiceType != nil {
		repoApt.ServiceType = *input.ServiceType
	}
	if input.Description != nil {
		repoApt.Description = *input.Description
	}
	if input.AssignedMechanic != nil {
		repoApt.AssignedMechanic = *input.AssignedMechanic
	}
	if input.Bay != nil {
		repoApt.Bay = input.Bay
	}
	if input.Notes != nil {
		repoApt.Notes = *input.Notes
	}
	if input.Status != nil {
		newStatus := *input.Status
		if newStatus == "on_going" {
			if repoApt.AssignedMechanic == "" {
				return nil, errors.New("mechanic must be assigned before starting work")
			}
			if repoApt.Bay == nil || *repoApt.Bay == "" {
				return nil, errors.New("bay must be assigned before starting work")
			}
		}
		repoApt.Status = newStatus
	}

	if err := s.repo.Update(ctx, repoApt); err != nil {
		return nil, fmt.Errorf("update appointment: %w", err)
	}

	return toModel(repoApt), nil
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
	if input.Status != nil {
		v.Status = *input.Status
	}
	if input.RepairStatus != nil {
		v.RepairStatus = *input.RepairStatus
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
	if input.Status != nil {
		repoV.Status = *input.Status
	}
	if input.RepairStatus != nil {
		repoV.RepairStatus = *input.RepairStatus
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

func (s *Service) CreateAssignment(ctx context.Context, input model.CreateStaffAssignmentInput) (*model.StaffAssignment, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if input.AppointmentID == "" {
		return nil, errors.New("appointmentId is required")
	}
	if input.StaffID == "" {
		return nil, errors.New("staffId is required")
	}
	if input.StaffName == "" {
		return nil, errors.New("staffName is required")
	}

	repoA := &repository.StaffAssignment{
		TenantID:      input.TenantID,
		AppointmentID: input.AppointmentID,
		StaffID:       input.StaffID,
		StaffName:     input.StaffName,
		Role:          input.Role,
		Status:        "assigned",
		Notes:         deref(input.Notes),
	}

	if err := s.repo.CreateAssignment(ctx, repoA); err != nil {
		return nil, fmt.Errorf("create assignment: %w", err)
	}

	return assignmentToModel(repoA), nil
}

func (s *Service) GetAssignment(ctx context.Context, id string) (*model.StaffAssignment, error) {
	repoA, err := s.repo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return assignmentToModel(repoA), nil
}

func (s *Service) ListAssignmentsByAppointment(ctx context.Context, appointmentID string) ([]*model.StaffAssignment, error) {
	items, err := s.repo.ListAssignmentsByAppointment(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.StaffAssignment, len(items))
	for i, a := range items {
		ca := a
		result[i] = assignmentToModel(&ca)
	}
	return result, nil
}

func (s *Service) UpdateAssignment(ctx context.Context, id string, input model.UpdateStaffAssignmentInput) (*model.StaffAssignment, error) {
	repoA, err := s.repo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Status != nil && *input.Status != "" {
		repoA.Status = *input.Status
		now := time.Now()
		switch *input.Status {
		case "in_progress":
			repoA.StartedAt = &now
		case "completed":
			repoA.CompletedAt = &now
		}
	}
	if input.TotalMinutes != nil {
		repoA.TotalMinutes = *input.TotalMinutes
	}
	if input.Notes != nil {
		repoA.Notes = *input.Notes
	}

	if err := s.repo.UpdateAssignment(ctx, repoA); err != nil {
		return nil, fmt.Errorf("update assignment: %w", err)
	}

	return assignmentToModel(repoA), nil
}

func (s *Service) DeleteAssignment(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteAssignment(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) ReassignAssignment(ctx context.Context, id string, targetAppointmentID string) (*model.StaffAssignment, error) {
	repoA, err := s.repo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	repoA.AppointmentID = targetAppointmentID
	repoA.Status = "assigned"
	now := time.Now()
	repoA.AssignedAt = now
	repoA.StartedAt = nil
	repoA.CompletedAt = nil
	repoA.TotalMinutes = 0

	if err := s.repo.UpdateAssignment(ctx, repoA); err != nil {
		return nil, fmt.Errorf("reassign assignment: %w", err)
	}

	return assignmentToModel(repoA), nil
}

func (s *Service) StartAssignment(ctx context.Context, id string) (*model.StaffAssignment, error) {
	repoA, err := s.repo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	repoA.Status = "in_progress"
	now := time.Now()
	repoA.StartedAt = &now

	if err := s.repo.UpdateAssignment(ctx, repoA); err != nil {
		return nil, fmt.Errorf("start assignment: %w", err)
	}

	return assignmentToModel(repoA), nil
}

func (s *Service) CompleteAssignment(ctx context.Context, id string, totalMinutes int) (*model.StaffAssignment, error) {
	repoA, err := s.repo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	repoA.Status = "completed"
	now := time.Now()
	repoA.CompletedAt = &now
	if totalMinutes <= 0 && repoA.StartedAt != nil {
		totalMinutes = int(now.Sub(*repoA.StartedAt).Minutes())
	}
	repoA.TotalMinutes = totalMinutes

	if err := s.repo.UpdateAssignment(ctx, repoA); err != nil {
		return nil, fmt.Errorf("complete assignment: %w", err)
	}

	return assignmentToModel(repoA), nil
}

func (s *Service) ListActiveAssignmentsByStaff(ctx context.Context, staffID string) ([]*model.StaffAssignment, error) {
	items, err := s.repo.ListActiveAssignmentsByStaff(ctx, staffID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.StaffAssignment, len(items))
	for i, a := range items {
		ca := a
		result[i] = assignmentToModel(&ca)
	}
	return result, nil
}

func (s *Service) FindAssignmentByID(ctx context.Context, id string) (*model.StaffAssignment, error) {
	return s.GetAssignment(ctx, id)
}

func assignmentToModel(a *repository.StaffAssignment) *model.StaffAssignment {
	var startedAt, completedAt *string
	if a.StartedAt != nil {
		s := a.StartedAt.Format(time.RFC3339)
		startedAt = &s
	}
	if a.CompletedAt != nil {
		s := a.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}
	return &model.StaffAssignment{
		ID:            a.ID,
		TenantID:      a.TenantID,
		AppointmentID: a.AppointmentID,
		StaffID:       a.StaffID,
		StaffName:     a.StaffName,
		Role:          a.Role,
		Status:        a.Status,
		AssignedAt:    a.AssignedAt.Format(time.RFC3339),
		StartedAt:     startedAt,
		CompletedAt:   completedAt,
		TotalMinutes:  &a.TotalMinutes,
		Notes:         strPtr(a.Notes),
	}
}

func toModel(r *repository.Appointment) *model.Appointment {
	scheduledDate := r.ScheduledDate.Format("2006-01-02")

	return &model.Appointment{
		ID:               r.ID,
		TenantID:         r.TenantID,
		ShopID:           coalesce(r.ShopID),
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
		Bay:             r.Bay,
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
		Status:       v.Status,
		RepairStatus: v.RepairStatus,
		CreatedAt:    v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    v.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *Service) CreateShopService(ctx context.Context, input model.CreateShopServiceInput) (*model.ShopService, error) {
	repo := &repository.ShopService{
		TenantID:      input.TenantID,
		ServiceTypeID: input.ServiceTypeID,
		Name:          input.Name,
		IsActive:      true,
	}
	if input.Code != nil {
		repo.Code = *input.Code
	}
	if input.System != nil {
		repo.System = *input.System
	}
	if input.Category != nil {
		repo.Category = *input.Category
	}
	if input.EstimatedHours != nil {
		repo.EstimatedHours = input.EstimatedHours
	}
	if err := s.repo.CreateShopService(ctx, repo); err != nil {
		return nil, fmt.Errorf("create shop service: %w", err)
	}
	return shopServiceToModel(repo), nil
}

func (s *Service) ListShopServices(ctx context.Context, tenantID string) (*model.ShopServiceConnection, error) {
	items, err := s.repo.ListShopServices(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list shop services: %w", err)
	}
	svcs := make([]*model.ShopService, len(items))
	for i := range items {
		svcs[i] = shopServiceToModel(&items[i])
	}
	return &model.ShopServiceConnection{Items: svcs, Total: len(svcs)}, nil
}

func (s *Service) DeleteShopService(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteShopService(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) CreateShopPart(ctx context.Context, input model.CreateShopPartInput) (*model.ShopPart, error) {
	repo := &repository.ShopPart{
		TenantID:  input.TenantID,
		Name:      input.Name,
		Quantity:  0,
		UnitPrice: 0,
	}
	if input.Sku != nil {
		repo.SKU = *input.Sku
	}
	if input.Description != nil {
		repo.Description = *input.Description
	}
	if input.Quantity != nil {
		repo.Quantity = *input.Quantity
	}
	if input.UnitPrice != nil {
		repo.UnitPrice = *input.UnitPrice
	}
	if input.MakeID != nil {
		repo.MakeID = input.MakeID
	}
	if input.ModelID != nil {
		repo.ModelID = input.ModelID
	}
	if input.Year != nil {
		repo.Year = input.Year
	}
	if input.LocationID != nil {
		repo.LocationID = input.LocationID
	}
	if err := s.repo.CreateShopPart(ctx, repo); err != nil {
		return nil, fmt.Errorf("create shop part: %w", err)
	}
	return shopPartToModel(repo), nil
}

func (s *Service) UpdateShopPart(ctx context.Context, id string, input model.UpdateShopPartInput) (*model.ShopPart, error) {
	existing, err := s.repo.GetShopPart(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Sku != nil {
		existing.SKU = *input.Sku
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Quantity != nil {
		existing.Quantity = *input.Quantity
	}
	if input.UnitPrice != nil {
		existing.UnitPrice = *input.UnitPrice
	}
	if input.MakeID != nil {
		existing.MakeID = input.MakeID
	}
	if input.ModelID != nil {
		existing.ModelID = input.ModelID
	}
	if input.Year != nil {
		existing.Year = input.Year
	}
	if input.LocationID != nil {
		existing.LocationID = input.LocationID
	}
	if err := s.repo.UpdateShopPart(ctx, existing); err != nil {
		return nil, fmt.Errorf("update shop part: %w", err)
	}
	return shopPartToModel(existing), nil
}

func (s *Service) ListShopParts(ctx context.Context, tenantID string) (*model.ShopPartConnection, error) {
	items, err := s.repo.ListShopParts(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list shop parts: %w", err)
	}
	parts := make([]*model.ShopPart, len(items))
	for i := range items {
		parts[i] = shopPartToModel(&items[i])
	}
	if err := s.populateBatches(ctx, parts); err != nil {
		return nil, err
	}
	return &model.ShopPartConnection{Items: parts, Total: len(parts)}, nil
}

func (s *Service) populateBatches(ctx context.Context, parts []*model.ShopPart) error {
	for _, p := range parts {
		batches, err := s.repo.ListPartBatches(ctx, p.ID)
		if err != nil {
			return fmt.Errorf("get batches for part %s: %w", p.ID, err)
		}
		for j := range batches {
			p.Batches = append(p.Batches, partBatchToModel(&batches[j]))
		}
	}
	return nil
}

func (s *Service) DeleteShopPart(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteShopPart(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) AddPartBatch(ctx context.Context, input model.CreatePartBatchInput) (*model.PartBatch, error) {
	repo := &repository.PartBatch{
		PartID:   input.PartID,
		TenantID: "", // set from part's tenant
		Quantity: input.Quantity,
		UnitCost: input.UnitCost,
	}
	if repo.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	part, err := s.repo.GetShopPart(ctx, input.PartID)
	if err != nil {
		return nil, fmt.Errorf("get part for batch: %w", err)
	}
	repo.TenantID = part.TenantID
	if err := s.repo.CreatePartBatch(ctx, repo); err != nil {
		return nil, fmt.Errorf("create part batch: %w", err)
	}
	return partBatchToModel(repo), nil
}

func (s *Service) UpdatePartBatch(ctx context.Context, id string, input model.UpdatePartBatchInput) (*model.PartBatch, error) {
	existing, err := s.repo.GetPartBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Quantity != nil {
		existing.Quantity = *input.Quantity
	}
	if input.UnitCost != nil {
		existing.UnitCost = *input.UnitCost
	}
	if err := s.repo.UpdatePartBatch(ctx, existing); err != nil {
		return nil, fmt.Errorf("update part batch: %w", err)
	}
	return partBatchToModel(existing), nil
}

func (s *Service) DeletePartBatch(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeletePartBatch(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) FindShopPartByID(ctx context.Context, id string) (*model.ShopPart, error) {
	p, err := s.repo.GetShopPart(ctx, id)
	if err != nil {
		return nil, err
	}
	return shopPartToModel(p), nil
}

func (s *Service) ListPartBatches(ctx context.Context, partID string) ([]*model.PartBatch, error) {
	items, err := s.repo.ListPartBatches(ctx, partID)
	if err != nil {
		return nil, fmt.Errorf("list part batches: %w", err)
	}
	batchModels := make([]*model.PartBatch, len(items))
	for i := range items {
		batchModels[i] = partBatchToModel(&items[i])
	}
	return batchModels, nil
}

func (s *Service) CreateShopTool(ctx context.Context, input model.CreateShopToolInput) (*model.ShopTool, error) {
	repo := &repository.ShopTool{
		TenantID:  input.TenantID,
		Name:      input.Name,
		Quantity:  0,
		Status:    "available",
	}
	if input.Description != nil {
		repo.Description = *input.Description
	}
	if input.Quantity != nil {
		repo.Quantity = *input.Quantity
	}
	if input.Status != nil {
		repo.Status = *input.Status
	}
	if err := s.repo.CreateShopTool(ctx, repo); err != nil {
		return nil, fmt.Errorf("create shop tool: %w", err)
	}
	return shopToolToModel(repo), nil
}

func (s *Service) UpdateShopTool(ctx context.Context, id string, input model.UpdateShopToolInput) (*model.ShopTool, error) {
	existing, err := s.repo.GetShopTool(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Quantity != nil {
		existing.Quantity = *input.Quantity
	}
	if input.Status != nil {
		existing.Status = *input.Status
	}
	if err := s.repo.UpdateShopTool(ctx, existing); err != nil {
		return nil, fmt.Errorf("update shop tool: %w", err)
	}
	return shopToolToModel(existing), nil
}

func (s *Service) ListShopTools(ctx context.Context, tenantID string) (*model.ShopToolConnection, error) {
	items, err := s.repo.ListShopTools(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list shop tools: %w", err)
	}
	tools := make([]*model.ShopTool, len(items))
	for i := range items {
		tools[i] = shopToolToModel(&items[i])
	}
	return &model.ShopToolConnection{Items: tools, Total: len(tools)}, nil
}

func (s *Service) DeleteShopTool(ctx context.Context, id string) (bool, error) {
	if err := s.repo.DeleteShopTool(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) AddAppointmentPart(ctx context.Context, input model.AddAppointmentPartInput) (*model.AppointmentPart, error) {
	if input.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	part, err := s.repo.GetShopPart(ctx, input.PartID)
	if err != nil {
		return nil, fmt.Errorf("get part: %w", err)
	}

	batches, err := s.repo.ListPartBatches(ctx, input.PartID)
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}

	needed := input.Quantity
	totalAvailable := 0
	for _, b := range batches {
		totalAvailable += b.Quantity
	}
	if totalAvailable < needed {
		return nil, fmt.Errorf("insufficient stock: have %d, need %d", totalAvailable, needed)
	}

	for i := range batches {
		if needed <= 0 {
			break
		}
		if batches[i].Quantity >= needed {
			batches[i].Quantity -= needed
			needed = 0
		} else {
			needed -= batches[i].Quantity
			batches[i].Quantity = 0
		}
		if err := s.repo.UpdatePartBatch(ctx, &batches[i]); err != nil {
			return nil, fmt.Errorf("update batch %s: %w", batches[i].ID, err)
		}
	}

	part.Quantity -= input.Quantity
	if err := s.repo.UpdateShopPart(ctx, part); err != nil {
		return nil, fmt.Errorf("update part quantity: %w", err)
	}

	var unitPrice float64
	if len(batches) > 0 && batches[0].UnitCost > 0 {
		unitPrice = batches[0].UnitCost
	}
	if input.UnitPrice != nil {
		unitPrice = *input.UnitPrice
	}

	repo := &repository.AppointmentPart{
		AppointmentID: input.AppointmentID,
		PartID:        input.PartID,
		PartName:      part.Name,
		Quantity:      input.Quantity,
		UnitPrice:     unitPrice,
	}
	if err := s.repo.CreateAppointmentPart(ctx, repo); err != nil {
		return nil, fmt.Errorf("create appointment part: %w", err)
	}
	return appointmentPartToModel(repo), nil
}

func (s *Service) ListAppointmentParts(ctx context.Context, appointmentID string) ([]*model.AppointmentPart, error) {
	items, err := s.repo.ListAppointmentParts(ctx, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("list appointment parts: %w", err)
	}
	result := make([]*model.AppointmentPart, len(items))
	for i := range items {
		result[i] = appointmentPartToModel(&items[i])
	}
	return result, nil
}

func (s *Service) DeleteAppointmentPart(ctx context.Context, id string) (bool, error) {
	ap, err := s.repo.GetAppointmentPart(ctx, id)
	if err != nil {
		return false, err
	}

	part, err := s.repo.GetShopPart(ctx, ap.PartID)
	if err != nil {
		return false, fmt.Errorf("get part: %w", err)
	}

	if err := s.repo.DeleteAppointmentPart(ctx, id); err != nil {
		return false, fmt.Errorf("delete appointment part: %w", err)
	}

	part.Quantity += ap.Quantity
	if err := s.repo.UpdateShopPart(ctx, part); err != nil {
		return false, fmt.Errorf("restore part quantity: %w", err)
	}

	return true, nil
}

func appointmentPartToModel(r *repository.AppointmentPart) *model.AppointmentPart {
	return &model.AppointmentPart{
		ID:            r.ID,
		AppointmentID: r.AppointmentID,
		PartID:        r.PartID,
		PartName:      r.PartName,
		Quantity:      r.Quantity,
		UnitPrice:     r.UnitPrice,
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     r.UpdatedAt.Format(time.RFC3339),
	}
}

func shopServiceToModel(r *repository.ShopService) *model.ShopService {
	return &model.ShopService{
		ID:             r.ID,
		TenantID:       r.TenantID,
		ServiceTypeID:  r.ServiceTypeID,
		Name:           r.Name,
		Code:           strPtr(r.Code),
		System:         strPtr(r.System),
		Category:       strPtr(r.Category),
		EstimatedHours: r.EstimatedHours,
		IsActive:       r.IsActive,
		CreatedAt:      r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      r.UpdatedAt.Format(time.RFC3339),
	}
}

func shopPartToModel(r *repository.ShopPart) *model.ShopPart {
	up := r.UnitPrice
	return &model.ShopPart{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Sku:         strPtr(r.SKU),
		Description: strPtr(r.Description),
		Quantity:    r.Quantity,
		UnitPrice:   &up,
		MakeID:      r.MakeID,
		ModelID:     r.ModelID,
		Year:        r.Year,
		LocationID:  r.LocationID,
		Batches:     []*model.PartBatch{},
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}
}

func partBatchToModel(r *repository.PartBatch) *model.PartBatch {
	return &model.PartBatch{
		ID:        r.ID,
		PartID:    r.PartID,
		Quantity:  r.Quantity,
		UnitCost:  r.UnitCost,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	}
}

func shopToolToModel(r *repository.ShopTool) *model.ShopTool {
	return &model.ShopTool{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Description: strPtr(r.Description),
		Quantity:    r.Quantity,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func coalesce(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
