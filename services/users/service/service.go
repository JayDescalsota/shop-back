package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/services/users/model"
	"backend/services/users/repository"

	"github.com/google/uuid"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTenant(ctx context.Context, input map[string]interface{}) (*model.Tenant, error) {
	name, _ := input["name"].(string)
	if name == "" {
		return nil, errors.New("name is required")
	}

	tenantType := model.TenantType(model.TenantTypeAutoOwner)
	if v, ok := input["type"].(string); ok && v != "" {
		tenantType = model.TenantType(v)
	}

	now := time.Now()
	t := &model.Tenant{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      tenantType,
		Status:    model.TenantStatusActive,
		Settings:  "{}",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if v, ok := input["domain"].(string); ok && v != "" {
		t.Domain = &v
	}

	if err := s.repo.CreateTenant(ctx, t); err != nil {
		return nil, fmt.Errorf("create tenant: %w", err)
	}

	return t, nil
}

func (s *Service) GetTenant(ctx context.Context, id string) (*model.Tenant, error) {
	return s.repo.GetTenantByID(ctx, id)
}

func (s *Service) UpdateTenant(ctx context.Context, id string, input map[string]interface{}) (*model.Tenant, error) {
	t, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if v, ok := input["name"].(string); ok && v != "" {
		t.Name = v
	}
	if v, ok := input["domain"].(string); ok {
		t.Domain = &v
	}
	if v, ok := input["settings"].(map[string]interface{}); ok {
		data, _ := jsonString(v)
		t.Settings = data
	}
	t.UpdatedAt = time.Now()

	if err := s.repo.UpdateTenant(ctx, t); err != nil {
		return nil, fmt.Errorf("update tenant: %w", err)
	}

	return t, nil
}

func (s *Service) GetTenantSettings(ctx context.Context, tenantID string) (*model.TenantSettings, error) {
	ts, err := s.repo.GetTenantSettings(ctx, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("tenant settings not found")
		}
		return nil, err
	}
	return ts, nil
}

func (s *Service) UpdateTenantSettings(ctx context.Context, tenantID string, input map[string]interface{}) (*model.TenantSettings, error) {
	now := time.Now()

	ts := &model.TenantSettings{
		TenantID:  tenantID,
		UpdatedAt: now,
	}

	if existing, err := s.repo.GetTenantSettings(ctx, tenantID); err == nil {
		ts.ID = existing.ID
		ts.CreatedAt = existing.CreatedAt
		ts.BusinessHours = existing.BusinessHours
		ts.PaymentConfig = existing.PaymentConfig
		ts.NotificationConfig = existing.NotificationConfig
		ts.Branding = existing.Branding
		ts.Features = existing.Features
	} else {
		ts.CreatedAt = now
	}

	if v, ok := input["businessHours"]; ok {
		ts.BusinessHours = jsonStringOr(v)
	}
	if v, ok := input["paymentConfig"]; ok {
		ts.PaymentConfig = jsonStringOr(v)
	}
	if v, ok := input["notificationConfig"]; ok {
		ts.NotificationConfig = jsonStringOr(v)
	}
	if v, ok := input["branding"]; ok {
		ts.Branding = jsonStringOr(v)
	}
	if v, ok := input["features"]; ok {
		ts.Features = jsonStringOr(v)
	}

	if err := s.repo.UpsertTenantSettings(ctx, ts); err != nil {
		return nil, fmt.Errorf("update tenant settings: %w", err)
	}

	return ts, nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	p, err := s.repo.GetUserProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, input map[string]interface{}) (*model.UserProfile, error) {
	now := time.Now()

	p := &model.UserProfile{
		UserID:    userID,
		Timezone:  "UTC",
		Locale:    "en",
		UpdatedAt: now,
	}

	if existing, err := s.repo.GetUserProfile(ctx, userID); err == nil {
		p.ID = existing.ID
		p.CreatedAt = existing.CreatedAt
		p.Timezone = existing.Timezone
		p.Locale = existing.Locale
		p.Title = existing.Title
		p.Department = existing.Department
		p.NotificationPrefs = existing.NotificationPrefs
	} else {
		p.ID = uuid.New().String()
		p.CreatedAt = now
	}

	if v, ok := input["title"].(string); ok {
		p.Title = &v
	}
	if v, ok := input["department"].(string); ok {
		p.Department = &v
	}
	if v, ok := input["timezone"].(string); ok && v != "" {
		p.Timezone = v
	}
	if v, ok := input["locale"].(string); ok && v != "" {
		p.Locale = v
	}
	if v, ok := input["notificationPrefs"]; ok {
		p.NotificationPrefs = jsonStringOr(v)
	}

	if err := s.repo.UpsertUserProfile(ctx, p); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return p, nil
}

func (s *Service) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *Service) ListRoles(ctx context.Context) ([]model.Role, error) {
	return s.repo.ListRoles(ctx)
}

// ---- SSO ----

func (s *Service) GetUserContext(ctx context.Context, userID, email, name, activeTenantID string) (*model.UserContext, error) {
	branches, err := s.buildBranchInfos(ctx, userID)
	if err != nil {
		return nil, err
	}

	uc := &model.UserContext{
		UserID:   userID,
		Email:    email,
		Name:     name,
		Branches: branches,
	}

	if activeTenantID != "" {
		for _, b := range branches {
			if b.TenantID == activeTenantID {
				uc.ActiveBranch = &b
				break
			}
		}
	}

	if uc.ActiveBranch == nil && len(branches) > 0 {
		uc.ActiveBranch = &branches[0]
	}

	return uc, nil
}

func (s *Service) GetUserBranches(ctx context.Context, userID string) ([]model.BranchInfo, error) {
	return s.buildBranchInfos(ctx, userID)
}

func (s *Service) SwitchBranch(ctx context.Context, userID, email, name, tenantID string) (*model.UserContext, error) {
	branches, err := s.buildBranchInfos(ctx, userID)
	if err != nil {
		return nil, err
	}

	var active *model.BranchInfo
	for _, b := range branches {
		if b.TenantID == tenantID {
			active = &b
			break
		}
	}
	if active == nil {
		return nil, errors.New("access denied to this branch")
	}

	return &model.UserContext{
		UserID:       userID,
		Email:        email,
		Name:         name,
		ActiveBranch: active,
		Branches:     branches,
	}, nil
}

func (s *Service) buildBranchInfos(ctx context.Context, userID string) ([]model.BranchInfo, error) {
	rows, err := s.repo.GetUserBranches(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user branches: %w", err)
	}

	infos := make([]model.BranchInfo, 0, len(rows))
	for _, row := range rows {
		perms, err := s.repo.GetRolePermissions(ctx, row.RoleID)
		if err != nil {
			perms = []string{}
		}

		infos = append(infos, model.BranchInfo{
			TenantID:    row.TenantID,
			TenantName:  row.TenantName,
			TenantType:  model.TenantType(row.TenantType),
			Role:        row.RoleCode,
			Permissions: perms,
		})
	}

	return infos, nil
}

func jsonStringOr(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case map[string]interface{}:
		data, _ := jsonString(val)
		return data
	default:
		return "{}"
	}
}

func jsonString(v map[string]interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}", err
	}
	return string(data), nil
}
