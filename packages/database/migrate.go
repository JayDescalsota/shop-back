package database

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/uptrace/bun"
)

func RunMigrations(db *bun.DB, migrationsFS fs.FS, path string) error {
	sqldb := db.DB

	driver, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres driver: %w", err)
	}

	src, err := iofs.New(migrationsFS, path)
	if err != nil {
		return fmt.Errorf("iofs source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate instance: %w", err)
	}

	m.Log = &migrateLogger{}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	fmt.Printf("[migrate] "+format, v...)
}

func (l *migrateLogger) Verbose() bool {
	return true
}
