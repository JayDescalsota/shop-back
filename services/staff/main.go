package main

import (
	"context"
	"log"

	"backend/packages/service"
	"backend/services/staff/generated"
	"backend/services/staff/repository"
	"backend/services/staff/resolver"
	staffservice "backend/services/staff/service"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	service.Serve("staff", func(r chi.Router, deps *service.Dependencies) {
		repo := repository.New(deps.DB)
		svc := staffservice.New(repo)

		if err := svc.Migrate(context.Background()); err != nil {
			log.Printf("warning: migration: %v", err)
		}

		resv := resolver.New(svc)
		srv := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resv}))

		r.Post("/graphql", srv.ServeHTTP)
		r.Post("/query", srv.ServeHTTP)
	})
}
