package main

import (
	"context"
	"log"
	"net/http"

	"backend/packages/database"
	"backend/packages/middleware"
	"backend/services/auth/generated"
	"backend/services/auth/handler"
	"backend/services/auth/repository"
	"backend/services/auth/resolver"
	authservice "backend/services/auth/service"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := database.NewPostgres(cfg.DSN())

	if err := ensureSchema(context.Background(), db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userAppRepo := repository.NewUserAppRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	userTenantRepo := repository.NewUserTenantRepository(db)
	jwtManager := middleware.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration, cfg.JWTRefreshExpiration)
	authSvc := authservice.NewAuthService(userRepo, userAppRepo, tenantRepo, userTenantRepo, jwtManager)

	if err := authSvc.SeedUsers(context.Background()); err != nil {
		log.Printf("warning: seed users: %v", err)
	}

	res := resolver.New(authSvc)
	server := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: res}))

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appID := r.Header.Get("X-App-Id")
			if appID != "" {
				ctx := context.WithValue(r.Context(), resolver.AppIDKey, appID)
				r = r.WithContext(ctx)
			}
			branchID := r.Header.Get("X-Branch-Id")
			if branchID != "" {
				ctx := context.WithValue(r.Context(), resolver.BranchIDKey, branchID)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.JWT(cfg.JWTSecret))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-App-Id", "X-Branch-Id", "X-User-Id", "X-Tenant-Id", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	rest := handler.NewRESTHandler(authSvc)
	r.Post("/api/auth/login", rest.Login)
	r.Post("/api/auth/register", rest.Register)
	r.Put("/api/auth/password", rest.ChangePassword)
	r.Post("/api/auth/users/{userId}/apps/{appId}", rest.AddUserApp)
	r.Delete("/api/auth/users/{userId}/apps/{appId}", rest.RemoveUserApp)
	r.Post("/api/auth/users/{email}/tenants/{tenantId}", rest.AddUserTenant)
	r.Delete("/api/auth/users/{email}/tenants/{tenantId}", rest.RemoveUserTenant)

	r.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	r.Handle("/query", server)
	r.Handle("/graphql", server)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Println("auth service running at http://localhost:" + cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, r))
}
