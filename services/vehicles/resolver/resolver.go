package resolver

import vehiclesservice "backend/services/vehicles/service"

type Resolver struct {
	svc *vehiclesservice.Service
}

func New(svc *vehiclesservice.Service) *Resolver {
	return &Resolver{svc: svc}
}
