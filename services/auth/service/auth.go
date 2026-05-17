package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/services/auth/model"
	"backend/services/auth/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles auth business logic.
type AuthService interface {
	Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error)
	Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error)
	FindUserByID(ctx context.Context, id string) (*model.User, error)
}

type authService struct {
	users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) AuthService {
	return &authService{users: users}
}

func (s *authService) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	id := uuid.New().String()
	user := &repository.User{
		ID:        id,
		Email:     input.Email,
		Password:  string(hashed),
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			return nil, err
		}
		return nil, err
	}

	return &model.AuthResponse{
		Token: fmt.Sprintf("token-%s", id),
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

	return &model.AuthResponse{
		Token: fmt.Sprintf("token-%s", user.ID),
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

func toGraphQLUser(user *repository.User) *model.User {
	return &model.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
