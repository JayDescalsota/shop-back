package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already registered")
)

type User struct {
	bun.BaseModel `bun:"auth.users"`
	ID            string    `bun:",pk"`
	Email         string    `bun:",notnull,unique"`
	Password      string    `bun:",notnull"`
	Name          *string
	Role          string    `bun:",notnull,default:'user'"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	UpdatePassword(ctx context.Context, id, hashedPassword string) error
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		if isDuplicateKey(err) {
			return ErrEmailTaken
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	_, err := r.db.NewUpdate().
		Model(&User{}).
		Set("password = ?", hashedPassword).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func isDuplicateKey(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
