package main

import (
	"backend/packages/service"
	"backend/services/users/repository"
	"backend/services/users/resolver"
	usersservice "backend/services/users/service"

	"github.com/go-chi/chi/v5"
)

func main() {
	service.Serve("users", func(r chi.Router, deps *service.Dependencies) {
		repo := repository.New(deps.DB)
		svc := usersservice.New(repo)
		resv := resolver.New(svc)

		r.Post("/graphql", resv.GraphQLHandler)
		r.Post("/query", resv.GraphQLHandler)
	})
}
