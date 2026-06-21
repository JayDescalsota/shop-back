package service

import (
	"context"
	"errors"
	"fmt"

	"backend/services/staff/model"
	"backend/services/staff/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateStaff(ctx context.Context, input model.CreateStaffInput) (*model.Staff, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	d := &repository.Staff{
		TenantID: input.TenantID,
		Name:     input.Name,
	}
	if input.Email != nil {
		d.Email = *input.Email
	}
	if input.Phone != nil {
		d.Phone = *input.Phone
	}
	if input.Role != nil {
		d.Role = *input.Role
	}
	if input.LicenseNumber != nil {
		d.LicenseNumber = *input.LicenseNumber
	}
	if input.LicenseClass != nil {
		d.LicenseClass = *input.LicenseClass
	}
	if input.LicenseExpiry != nil {
		d.LicenseExpiry = *input.LicenseExpiry
	}
	if input.DateOfBirth != nil {
		d.DateOfBirth = *input.DateOfBirth
	}
	if input.Address != nil {
		d.Address = *input.Address
	}
	if input.EmergencyContact != nil {
		d.EmergencyContact = *input.EmergencyContact
	}
	if input.EmergencyPhone != nil {
		d.EmergencyPhone = *input.EmergencyPhone
	}
	if input.Status != nil {
		d.Status = *input.Status
	}
	if input.AssignedVehicleID != nil {
		d.AssignedVehicleID = input.AssignedVehicleID
	}
	if input.Notes != nil {
		d.Notes = *input.Notes
	}
	if input.HireDate != nil {
		d.HireDate = *input.HireDate
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create staff: %w", err)
	}

	return toModel(d), nil
}

func (s *Service) GetStaff(ctx context.Context, id string) (*model.Staff, error) {
	repoD, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toModel(repoD), nil
}

func (s *Service) ListStaff(ctx context.Context, tenantID string) (*model.StaffConnection, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	staff, err := s.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Staff, len(staff))
	for i, d := range staff {
		items[i] = toModel(&d)
	}

	return &model.StaffConnection{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) UpdateStaff(ctx context.Context, id string, input model.UpdateStaffInput) (*model.Staff, error) {
	repoD, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		repoD.Name = *input.Name
	}
	if input.Email != nil {
		repoD.Email = *input.Email
	}
	if input.Phone != nil {
		repoD.Phone = *input.Phone
	}
	if input.Role != nil {
		repoD.Role = *input.Role
	}
	if input.LicenseNumber != nil {
		repoD.LicenseNumber = *input.LicenseNumber
	}
	if input.LicenseClass != nil {
		repoD.LicenseClass = *input.LicenseClass
	}
	if input.LicenseExpiry != nil {
		repoD.LicenseExpiry = *input.LicenseExpiry
	}
	if input.DateOfBirth != nil {
		repoD.DateOfBirth = *input.DateOfBirth
	}
	if input.Address != nil {
		repoD.Address = *input.Address
	}
	if input.EmergencyContact != nil {
		repoD.EmergencyContact = *input.EmergencyContact
	}
	if input.EmergencyPhone != nil {
		repoD.EmergencyPhone = *input.EmergencyPhone
	}
	if input.Status != nil {
		repoD.Status = *input.Status
	}
	if input.AssignedVehicleID != nil {
		repoD.AssignedVehicleID = input.AssignedVehicleID
	}
	if input.Notes != nil {
		repoD.Notes = *input.Notes
	}
	if input.HireDate != nil {
		repoD.HireDate = *input.HireDate
	}

	if err := s.repo.Update(ctx, repoD); err != nil {
		return nil, fmt.Errorf("update staff: %w", err)
	}

	return toModel(repoD), nil
}

func (s *Service) DeleteStaff(ctx context.Context, id string) (bool, error) {
	if err := s.repo.Delete(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) FindStaffByID(ctx context.Context, id string) (*model.Staff, error) {
	return s.GetStaff(ctx, id)
}

func toModel(d *repository.Staff) *model.Staff {
	return &model.Staff{
		ID:                d.ID,
		TenantID:          d.TenantID,
		Name:              d.Name,
		Email:             strPtr(d.Email),
		Phone:             strPtr(d.Phone),
		Role:              d.Role,
		LicenseNumber:     strPtr(d.LicenseNumber),
		LicenseClass:      strPtr(d.LicenseClass),
		LicenseExpiry:     strPtr(d.LicenseExpiry),
		DateOfBirth:       strPtr(d.DateOfBirth),
		Address:           strPtr(d.Address),
		EmergencyContact:  strPtr(d.EmergencyContact),
		EmergencyPhone:    strPtr(d.EmergencyPhone),
		Status:            d.Status,
		AssignedVehicleID: d.AssignedVehicleID,
		Notes:             strPtr(d.Notes),
		HireDate:          strPtr(d.HireDate),
		CreatedAt:         d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
