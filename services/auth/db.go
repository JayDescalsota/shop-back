package main

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

const (
	AppIDStore  = "00000000-0000-0000-0000-202605270001"
	AppIDRepair = "00000000-0000-0000-0000-202605270002"
	AppIDMobile = "00000000-0000-0000-0000-202605270003"
	AppIDAdmin       = "00000000-0000-0000-0000-202605270004"
	AppIDAutomobile  = "00000000-0000-0000-0000-202605270005"

	TenantStoreDowntown = "00000000-0000-0000-0001-202605270001"
	TenantStoreWestside = "00000000-0000-0000-0001-202605270002"
	TenantRepairMain    = "00000000-0000-0000-0001-202605270003"
	TenantRepairOak     = "00000000-0000-0000-0001-202605270004"
	TenantFleetMain     = "00000000-0000-0000-0001-202605270005"
	TenantMobileMain    = "00000000-0000-0000-0001-202605270006"
)

func ensureSchema(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS auth.apps (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			slug TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure apps table: %w", err)
	}

	for _, row := range []struct{ id, name, slug string }{
		{AppIDStore, "Store", "store"},
		{AppIDRepair, "Repair", "repair"},
		{AppIDMobile, "Mobile", "mobile"},
		{AppIDAdmin, "Admin", "admin"},
		{AppIDAutomobile, "Automobile", "automobile"},
	} {
		_, err := db.ExecContext(ctx, fmt.Sprintf(
			`INSERT INTO auth.apps (id, name, slug) VALUES ('%s', '%s', '%s') ON CONFLICT (id) DO NOTHING`,
			row.id, row.name, row.slug,
		))
		if err != nil {
			return fmt.Errorf("seed app %s: %w", row.slug, err)
		}
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS auth.tenants (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			address TEXT,
			app_id TEXT NOT NULL REFERENCES auth.apps(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure tenants table: %w", err)
	}

	for _, row := range []struct{ id, name, address, appID string }{
		{TenantStoreDowntown, "Downtown Auto Parts", "100 Main St, Downtown", AppIDStore},
		{TenantStoreWestside, "Westside Auto Parts", "456 Oak Ave, Westside", AppIDStore},
		{TenantRepairMain, "Main Street Auto Care", "123 Main St, Downtown", AppIDRepair},
		{TenantRepairOak, "Oak Avenue Auto Repair", "789 Pine Rd, Easton", AppIDRepair},
		{TenantFleetMain, "Fleet Operations Center", "500 Fleet St, Metro", AppIDAutomobile},
		{TenantMobileMain, "Mobile Service 1", "200 Mobile Ave, Midtown", AppIDMobile},
	} {
		_, err := db.ExecContext(ctx, fmt.Sprintf(
			`INSERT INTO auth.tenants (id, name, address, app_id) VALUES ('%s', '%s', '%s', '%s') ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, address = EXCLUDED.address, app_id = EXCLUDED.app_id`,
			row.id, row.name, row.address, row.appID,
		))
		if err != nil {
			return fmt.Errorf("seed tenant %s: %w", row.name, err)
		}
	}

	db.ExecContext(ctx, `DROP TABLE IF EXISTS auth.user_tenants CASCADE`)
	db.ExecContext(ctx, `DROP TABLE IF EXISTS auth.user_apps CASCADE`)
	db.ExecContext(ctx, `DROP TABLE IF EXISTS auth.users CASCADE`)
	db.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS auth`)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE auth.users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			name TEXT,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure users table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS auth.user_apps (
			user_id TEXT NOT NULL REFERENCES auth.users(id),
			app_id TEXT NOT NULL REFERENCES auth.apps(id),
			role TEXT NOT NULL DEFAULT 'user',
			PRIMARY KEY (user_id, app_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure user_apps table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS auth.user_tenants (
			email TEXT NOT NULL,
			tenant_id TEXT NOT NULL REFERENCES auth.tenants(id),
			PRIMARY KEY (email, tenant_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure user_tenants table: %w", err)
	}

	return nil
}
