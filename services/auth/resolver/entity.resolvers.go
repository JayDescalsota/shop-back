package resolver

import (
	"context"

	"backend/services/auth/generated"
	"backend/services/auth/model"
)

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	return r.auth.FindUserByID(ctx, id)
}

func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
