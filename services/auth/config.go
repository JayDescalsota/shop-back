package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisHost string
	RedisPort string

	JWTSecret          string
	JWTExpiration      time.Duration
	JWTRefreshExpiration time.Duration
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func loadConfig() (Config, error) {
	if err := loadEnvFile(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	jwtExp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse JWT_EXPIRATION: %w", err)
	}

	refreshExp, err := time.ParseDuration(getEnv("REFRESH_EXPIRATION", "720h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse REFRESH_EXPIRATION: %w", err)
	}

	return Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "1"),
		DBName:     getEnv("DB_NAME", "postgres"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		JWTSecret:             getEnv("JWT_SECRET", "mysecretkey"),
		JWTExpiration:         jwtExp,
		JWTRefreshExpiration:  refreshExp,
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadEnvFile() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(envPath); statErr == nil {
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}
