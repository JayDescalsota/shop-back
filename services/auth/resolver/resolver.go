package resolver

import (
	"backend/services/auth/service"
)

type Resolver struct {
	auth service.AuthService
}

func New(auth service.AuthService) *Resolver {
	return &Resolver{auth: auth}
}
