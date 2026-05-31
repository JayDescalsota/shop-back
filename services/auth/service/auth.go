package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/packages/middleware"
	"backend/services/auth/model"
	"backend/services/auth/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	appIDStore      = "00000000-0000-0000-0000-202605270001"
	appIDRepair     = "00000000-0000-0000-0000-202605270002"
	appIDMobile     = "00000000-0000-0000-0000-202605270003"
	appIDAdmin      = "00000000-0000-0000-0000-202605270004"
	appIDAutomobile = "00000000-0000-0000-0000-202605270005"

	TenantStoreDowntown = "00000000-0000-0000-0001-202605270001"
	TenantStoreWestside = "00000000-0000-0000-0001-202605270002"
	TenantRepairMain    = "00000000-0000-0000-0001-202605270003"
	TenantRepairOak     = "00000000-0000-0000-0001-202605270004"
	TenantFleetMain     = "00000000-0000-0000-0001-202605270005"
	TenantMobileMain    = "00000000-0000-0000-0001-202605270006"
)

type UserInfoResponse struct {
	ID      string                `json:"id"`
	Email   string                `json:"email"`
	Name    *string               `json:"name,omitempty"`
	Role    string                `json:"role"`
	Apps    []string              `json:"apps"`
	Tenants []middleware.TenantInfo `json:"tenants"`
}

type AuthService interface {
	Register(ctx context.Context, input model.RegisterInput, app string) (*model.AuthResponse, error)
	Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error)
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	SeedUsers(ctx context.Context) error
	AddUserApp(ctx context.Context, userID, appID, role string) error
	RemoveUserApp(ctx context.Context, userID, appID string) error
	AddUserTenant(ctx context.Context, email, tenantID string) error
	RemoveUserTenant(ctx context.Context, email, tenantID string) error
	GetUserInfo(ctx context.Context, email string) (*UserInfoResponse, error)
}

type authService struct {
	users       repository.UserRepository
	userApps    repository.UserAppRepository
	tenants     repository.TenantRepository
	userTenants repository.UserTenantRepository
	jwtManager  *middleware.JWTManager
}

func NewAuthService(
	users repository.UserRepository,
	userApps repository.UserAppRepository,
	tenants repository.TenantRepository,
	userTenants repository.UserTenantRepository,
	jwtManager *middleware.JWTManager,
) AuthService {
	return &authService{
		users: users, userApps: userApps,
		tenants: tenants, userTenants: userTenants,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, input model.RegisterInput, app string) (*model.AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	id := uuid.New().String()
	role := "user"
	if app == "" {
		app = appIDStore
	}

	user := &repository.User{
		ID:        id,
		Email:     input.Email,
		Password:  string(hashed),
		Name:      &input.Name,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			return nil, err
		}
		return nil, err
	}

	if err := s.userApps.Add(ctx, id, app, role); err != nil {
		return nil, fmt.Errorf("add user app: %w", err)
	}

	appIDs := []string{app}
	tenantList, _ := s.tenants.FindByApps(ctx, appIDs)

	token, _, err := s.jwtManager.GenerateToken(id, "", role, input.Email, appIDs, toTenantInfo(tenantList))
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User:  toGraphQLUser(user),
	}, nil
}

func (s *authService) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	userApps, _ := s.userApps.FindByUserID(ctx, user.ID)
	appIDs := make([]string, len(userApps))
	for i, a := range userApps {
		appIDs[i] = a.AppID
	}

	role := user.Role
	var tenantList []repository.Tenant
	if role == "superAdmin" {
		tenantList, _ = s.tenants.FindByApps(ctx, appIDs)
	} else {
		tenantList, _ = s.userTenants.FindByEmail(ctx, input.Email)
	}

	token, _, err := s.jwtManager.GenerateToken(user.ID, "", role, input.Email, appIDs, toTenantInfo(tenantList))
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User:  toGraphQLUser(user),
	}, nil
}

func (s *authService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return toGraphQLUser(user), nil
}

func (s *authService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.users.UpdatePassword(ctx, user.ID, string(hashed))
}

func (s *authService) AddUserApp(ctx context.Context, userID, appID, role string) error {
	return s.userApps.Add(ctx, userID, appID, role)
}

func (s *authService) RemoveUserApp(ctx context.Context, userID, appID string) error {
	return s.userApps.Remove(ctx, userID, appID)
}

func (s *authService) AddUserTenant(ctx context.Context, email, tenantID string) error {
	return s.userTenants.Add(ctx, email, tenantID)
}

func (s *authService) RemoveUserTenant(ctx context.Context, email, tenantID string) error {
	return s.userTenants.Remove(ctx, email, tenantID)
}

func (s *authService) GetUserInfo(ctx context.Context, email string) (*UserInfoResponse, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	userApps, _ := s.userApps.FindByUserID(ctx, user.ID)
	appIDs := make([]string, len(userApps))
	for i, a := range userApps {
		appIDs[i] = a.AppID
	}
	tenantList, _ := s.userTenants.FindByEmail(ctx, email)

	return &UserInfoResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Role:    user.Role,
		Apps:    appIDs,
		Tenants: toTenantInfo(tenantList),
	}, nil
}

func (s *authService) SeedUsers(ctx context.Context) error {
	type seedUser struct {
		Email    string
		Password string
		Name     string
		Role     string
	}
	type seedApp struct {
		Email string
		App   string
		Role  string
	}

	seedUsers := []seedUser{
		{Email: "admin@autolab.com", Password: "password123", Name: "Admin", Role: "admin"},
		{Email: "user@autolab.com", Password: "password123", Name: "User", Role: "user"},
		{Email: "storeadmin@autolab.com", Password: "password123", Name: "Store Admin", Role: "admin"},
		{Email: "mobile@autolab.com", Password: "password123", Name: "Mobile User", Role: "user"},
		{Email: "superadmin@autolab.com", Password: "password123", Name: "Super Admin", Role: "superAdmin"},
		{Email: "fleet@autolab.com", Password: "password123", Name: "Fleet Manager", Role: "admin"},
	}

	for _, su := range seedUsers {
		if _, err := s.users.FindByEmail(ctx, su.Email); err == nil {
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(su.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash seed password: %w", err)
		}

		user := &repository.User{
			ID:        uuid.New().String(),
			Email:     su.Email,
			Password:  string(hashed),
			Name:      &su.Name,
			Role:      su.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.users.Create(ctx, user); err != nil {
			return fmt.Errorf("create seed user %s: %w", su.Email, err)
		}
	}

	seedApps := []seedApp{
		{Email: "admin@autolab.com", App: appIDStore, Role: "admin"},
		{Email: "admin@autolab.com", App: appIDRepair, Role: "admin"},
		{Email: "admin@autolab.com", App: appIDMobile, Role: "admin"},
		{Email: "admin@autolab.com", App: appIDAutomobile, Role: "admin"},
		{Email: "user@autolab.com", App: appIDRepair, Role: "user"},
		{Email: "storeadmin@autolab.com", App: appIDStore, Role: "admin"},
		{Email: "mobile@autolab.com", App: appIDMobile, Role: "user"},
		{Email: "superadmin@autolab.com", App: appIDAdmin, Role: "superAdmin"},
		{Email: "superadmin@autolab.com", App: appIDStore, Role: "superAdmin"},
		{Email: "superadmin@autolab.com", App: appIDRepair, Role: "superAdmin"},
		{Email: "superadmin@autolab.com", App: appIDMobile, Role: "superAdmin"},
		{Email: "superadmin@autolab.com", App: appIDAutomobile, Role: "superAdmin"},
		{Email: "fleet@autolab.com", App: appIDAutomobile, Role: "admin"},
	}

	for _, sa := range seedApps {
		user, err := s.users.FindByEmail(ctx, sa.Email)
		if err != nil {
			continue
		}
		existing, _ := s.userApps.FindByUserID(ctx, user.ID)
		alreadySeeded := false
		for _, a := range existing {
			if a.AppID == sa.App {
				alreadySeeded = true
				break
			}
		}
		if alreadySeeded {
			continue
		}
		if err := s.userApps.Add(ctx, user.ID, sa.App, sa.Role); err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "unique") &&
				!strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				return fmt.Errorf("seed user app %s/%s: %w", sa.Email, sa.App, err)
			}
		}
	}

	tenantSeeds := []struct {
		Email    string
		TenantID string
	}{
		{Email: "superadmin@autolab.com", TenantID: TenantStoreDowntown},
		{Email: "superadmin@autolab.com", TenantID: TenantStoreWestside},
		{Email: "superadmin@autolab.com", TenantID: TenantRepairMain},
		{Email: "superadmin@autolab.com", TenantID: TenantRepairOak},
		{Email: "admin@autolab.com", TenantID: TenantStoreDowntown},
		{Email: "admin@autolab.com", TenantID: TenantStoreWestside},
		{Email: "admin@autolab.com", TenantID: TenantRepairMain},
		{Email: "admin@autolab.com", TenantID: TenantRepairOak},
		{Email: "user@autolab.com", TenantID: TenantRepairMain},
		{Email: "user@autolab.com", TenantID: TenantRepairOak},
		{Email: "storeadmin@autolab.com", TenantID: TenantStoreDowntown},
		{Email: "storeadmin@autolab.com", TenantID: TenantStoreWestside},
		{Email: "admin@autolab.com", TenantID: TenantFleetMain},
		{Email: "admin@autolab.com", TenantID: TenantMobileMain},
		{Email: "superadmin@autolab.com", TenantID: TenantFleetMain},
		{Email: "superadmin@autolab.com", TenantID: TenantMobileMain},
		{Email: "fleet@autolab.com", TenantID: TenantFleetMain},
		{Email: "mobile@autolab.com", TenantID: TenantMobileMain},
	}

	for _, ts := range tenantSeeds {
		if err := s.userTenants.Add(ctx, ts.Email, ts.TenantID); err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "unique") &&
				!strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				return fmt.Errorf("seed user tenant %s/%s: %w", ts.Email, ts.TenantID, err)
			}
		}
	}

	return nil
}

func toGraphQLUser(user *repository.User) *model.User {
	u := &model.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	if user.Name != nil {
		u.Name = *user.Name
	}
	return u
}

func toTenantInfo(tenants []repository.Tenant) []middleware.TenantInfo {
	infos := make([]middleware.TenantInfo, len(tenants))
	for i, t := range tenants {
		infos[i] = middleware.TenantInfo{ID: t.ID, Name: t.Name, AppID: t.AppID}
	}
	return infos
}
