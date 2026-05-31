package resolver

import (
	"context"
	"fmt"

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

func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: RefreshToken - refreshToken"))
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (bool, error) {
	panic(fmt.Errorf("not implemented: ForgotPassword - forgotPassword"))
}

func (r *mutationResolver) ResetPassword(ctx context.Context, token string, password string) (bool, error) {
	panic(fmt.Errorf("not implemented: ResetPassword - resetPassword"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, fmt.Errorf("not authenticated")
	}
	return r.auth.FindUserByID(ctx, userID)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

func (r *queryResolver) Users(ctx context.Context, tenantID string, page *int, perPage *int) (*model.UserConnection, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
