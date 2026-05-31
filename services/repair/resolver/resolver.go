package resolver

import repairservice "backend/services/repair/service"

type Resolver struct {
	svc *repairservice.Service
}

func New(svc *repairservice.Service) *Resolver {
	return &Resolver{svc: svc}
}
