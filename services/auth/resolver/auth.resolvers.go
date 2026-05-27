package resolver

import (
	"context"
	"encoding/json"

	"backend/packages/middleware"
	"backend/services/auth/generated"
	"backend/services/auth/model"
)

type contextKey string

const (
	AppIDKey    contextKey = "app_id"
	BranchIDKey contextKey = "branch_id"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	app, _ := ctx.Value(AppIDKey).(string)
	return r.auth.Register(ctx, input, app)
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	return r.auth.Login(ctx, input)
}

func (r *queryResolver) Me(ctx context.Context) (string, error) {
	email := middleware.GetEmail(ctx)
	if email == "" {
		return "anonymous", nil
	}
	info, err := r.auth.GetUserInfo(ctx, email)
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(info)
	return string(b), nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver     { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
