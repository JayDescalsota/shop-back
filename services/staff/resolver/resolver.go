package resolver

import staffservice "backend/services/staff/service"

type Resolver struct {
	svc *staffservice.Service
}

func New(svc *staffservice.Service) *Resolver {
	return &Resolver{svc: svc}
}
