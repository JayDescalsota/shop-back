package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"backend/packages/database"
	"backend/services/auth/generated"
	"backend/services/auth/repository"
	"backend/services/auth/resolver"
	"backend/services/auth/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db := database.NewPostgres(dsn)

	if err := ensureSchema(context.Background(), db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	res := resolver.New(authSvc)
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: res}))

	r := chi.NewRouter()
	r.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	r.Handle("/query", server)
	r.Handle("/graphql", server)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Println("Server running at http://localhost:" + cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, r))
}
