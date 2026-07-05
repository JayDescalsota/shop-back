package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/packages/middleware"
	"backend/services/auth/model"
	"backend/services/auth/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
		app = "00000000-0000-0000-0000-202605270001"
	}

	user := &repository.User{
		ID:           id,
		TenantID:     "", // Will be set by tenant logic or migration
		Email:        input.Email,
		PasswordHash: string(hashed),
		Name:         input.Name,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	userApps, _ := s.userApps.FindByUserID(ctx, user.ID)
	appIDs := make([]string, len(userApps))
	for i, a := range userApps {
		appIDs[i] = a.AppID
	}

	role := user.Role
	var tenantList []repository.Tenant
	if role == "superAdmin" || role == "super_admin" {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
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
		Name:    &user.Name,
		Role:    user.Role,
		Apps:    appIDs,
		Tenants: toTenantInfo(tenantList),
	}, nil
}

func toGraphQLUser(user *repository.User) *model.User {
	u := &model.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	if user.Name != "" {
		u.Name = user.Name
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
