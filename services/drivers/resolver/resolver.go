package resolver

import driverservice "backend/services/drivers/service"

type Resolver struct {
	svc *driverservice.Service
}

func New(svc *driverservice.Service) *Resolver {
	return &Resolver{svc: svc}
}
