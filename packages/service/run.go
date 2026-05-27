package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/packages/config"
	"backend/packages/database"
	"backend/packages/middleware"
	natss "backend/packages/nats"
	backendredis "backend/packages/redis"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	goredis "github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type Dependencies struct {
	Config config.Config
	DB     *bun.DB
	Redis  *goredis.Client
	NATS   *nats.Conn
	JS     jetstream.JetStream
}

type SetupFunc func(r chi.Router, deps *Dependencies)

func Serve(serviceName string, setup SetupFunc) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := database.NewPostgres(cfg.DSN())

	var rdb *goredis.Client
	if cfg.RedisHost != "" {
		rdb = backendredis.NewClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	}

	var nc *nats.Conn
	var js jetstream.JetStream
	if cfg.NATSURL != "" {
		nc, js, err = natss.NewConnection(cfg.NATSURL)
		if err != nil {
			log.Printf("nats not available (non-fatal): %v", err)
		}
	}

	deps := &Dependencies{
		Config: cfg,
		DB:     db,
		Redis:  rdb,
		NATS:   nc,
		JS:     js,
	}

	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(middleware.CORSMiddleware([]string{"*"}))

	if cfg.JWTSecret != "" {
		r.Use(middleware.JWT(cfg.JWTSecret))
	}
	r.Use(middleware.TenantIsolation)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"` + serviceName + `"}`))
	})

	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if err := db.PingContext(context.Background()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	setup(r, deps)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("%s listening on :%s", serviceName, cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s error: %v", serviceName, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
