package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type UserApp struct {
	bun.BaseModel `bun:"auth.user_apps"`
	UserID        string `bun:"user_id,pk"`
	AppID         string `bun:"app_id,pk"`
	Role          string `bun:",notnull,default:'user'"`
}

type UserAppRepository interface {
	FindByUserID(ctx context.Context, userID string) ([]UserApp, error)
	Add(ctx context.Context, userID, appID, role string) error
	Remove(ctx context.Context, userID, appID string) error
}

type userAppRepository struct {
	db *bun.DB
}

func NewUserAppRepository(db *bun.DB) UserAppRepository {
	return &userAppRepository{db: db}
}

func (r *userAppRepository) FindByUserID(ctx context.Context, userID string) ([]UserApp, error) {
	var apps []UserApp
	err := r.db.NewSelect().Model(&apps).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("find user apps: %w", err)
	}
	return apps, nil
}

func (r *userAppRepository) Add(ctx context.Context, userID, appID, role string) error {
	_, err := r.db.NewInsert().Model(&UserApp{UserID: userID, AppID: appID, Role: role}).Exec(ctx)
	if err != nil {
		return fmt.Errorf("add user app: %w", err)
	}
	return nil
}

func (r *userAppRepository) Remove(ctx context.Context, userID, appID string) error {
	_, err := r.db.NewDelete().Model(&UserApp{}).Where("user_id = ? AND app_id = ?", userID, appID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove user app: %w", err)
	}
	return nil
}
