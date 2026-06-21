package main

import (
	"embed"
	"log"

	"backend/packages/database"
	"backend/packages/service"
	"backend/services/users/repository"
	"backend/services/users/resolver"
	usersservice "backend/services/users/service"

	"github.com/go-chi/chi/v5"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	service.Serve("users", func(r chi.Router, deps *service.Dependencies) {
		if err := database.RunMigrations(deps.DB, migrationsFS, "migrations"); err != nil {
			log.Printf("warning: migration: %v", err)
		}

		repo := repository.New(deps.DB)
		svc := usersservice.New(repo)
		resv := resolver.New(svc)

		r.Post("/graphql", resv.GraphQLHandler)
		r.Post("/query", resv.GraphQLHandler)
	})
}
