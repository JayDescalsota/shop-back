package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/packages/config"
	"backend/packages/database"
	"backend/packages/middleware"
	"backend/services/lookup/repository"
	"backend/services/lookup/resolver"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/uptrace/bun"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := database.NewPostgres(cfg.DSN())
	repo := repository.New(db)
	resv := resolver.New(repo)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Tenant-Id", "X-Request-Id"},
		MaxAge:           86400,
	}))
	r.Use(middleware.JWT(cfg.JWTSecret))
	r.Use(middleware.TenantIsolation)

	r.Get("/healthz", healthHandler)
	r.Get("/readyz", readyHandler(db))

	r.Post("/graphql", resv.GraphQLHandler)
	r.Post("/query", resv.GraphQLHandler)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("lookup service listening on :%s", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("lookup error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"service":   "lookup",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func readyHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ready"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}
}
