package main

import (
	"embed"
	"log"

	"backend/packages/database"
	"backend/packages/service"
	"backend/services/repair/generated"
	"backend/services/repair/repository"
	"backend/services/repair/resolver"
	repairservice "backend/services/repair/service"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	service.Serve("repair", func(r chi.Router, deps *service.Dependencies) {
		if err := database.RunMigrations(deps.DB, migrationsFS, "migrations"); err != nil {
			log.Printf("warning: migration: %v", err)
		}

		repo := repository.New(deps.DB)
		svc := repairservice.New(repo)
		resv := resolver.New(svc)
		srv := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resv}))

		r.Post("/graphql", srv.ServeHTTP)
		r.Post("/query", srv.ServeHTTP)
	})
}
