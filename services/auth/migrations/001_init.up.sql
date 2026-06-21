CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.apps (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS auth.tenants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT,
    app_id TEXT NOT NULL REFERENCES auth.apps(id)
);

CREATE TABLE IF NOT EXISTS auth.users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth.user_apps (
    user_id TEXT NOT NULL REFERENCES auth.users(id),
    app_id TEXT NOT NULL REFERENCES auth.apps(id),
    role TEXT NOT NULL DEFAULT 'user',
    PRIMARY KEY (user_id, app_id)
);

CREATE TABLE IF NOT EXISTS auth.user_tenants (
    email TEXT NOT NULL,
    tenant_id TEXT NOT NULL REFERENCES auth.tenants(id),
    PRIMARY KEY (email, tenant_id)
);
