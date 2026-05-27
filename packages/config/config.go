package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	ServiceName string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	NATSURL string

	JWTSecret         string
	JWTExpiration     string
	RefreshExpiration string

	EncryptionKey string

	AWSRegion    string
	S3Bucket     string
	S3Endpoint   string

	StripeSecretKey string
	StripeWebhookSecret string

	SentryDSN string

	OTLPEndpoint string
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	return Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		AppPort:           getEnv("APP_PORT", "8080"),
		ServiceName:       getEnv("SERVICE_NAME", "service"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "1"),
		DBName:            getEnv("DB_NAME", "postgres"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		NATSURL:           getEnv("NATS_URL", "nats://localhost:4222"),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiration:     getEnv("JWT_EXPIRATION", "24h"),
		RefreshExpiration: getEnv("REFRESH_EXPIRATION", "720h"),
		EncryptionKey:     getEnv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef"),
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		S3Bucket:          getEnv("S3_BUCKET", "uploads"),
		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		StripeSecretKey:   getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		SentryDSN:         getEnv("SENTRY_DSN", ""),
		OTLPEndpoint:      getEnv("OTLP_ENDPOINT", ""),
	}, nil
}

func (c Config) DSN() string {
	sslMode := c.DBSSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, sslMode)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadDotEnv() error {
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
