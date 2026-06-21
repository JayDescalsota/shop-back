CREATE SCHEMA IF NOT EXISTS lookup;

CREATE TABLE IF NOT EXISTS lookup.vehicle_makes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    slug        TEXT NOT NULL DEFAULT '',
    logo_url    TEXT,
    country     TEXT,
    founded_year INT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    sort_order  INT NOT NULL DEFAULT 0,
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.vehicle_models (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    make_id     UUID NOT NULL REFERENCES lookup.vehicle_makes(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL DEFAULT '',
    year_start  INT,
    year_end    INT,
    vehicle_type TEXT,
    trim_levels JSONB NOT NULL DEFAULT '[]',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(make_id, name)
);

CREATE TABLE IF NOT EXISTS lookup.service_types (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID,
    code            TEXT NOT NULL DEFAULT '',
    name            TEXT NOT NULL UNIQUE,
    description     TEXT,
    category        TEXT,
    system          TEXT NOT NULL DEFAULT '',
    estimated_hours DECIMAL(6,2),
    is_global       BOOLEAN NOT NULL DEFAULT true,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.diagnostic_codes (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code                TEXT NOT NULL UNIQUE,
    system              TEXT NOT NULL DEFAULT '',
    description         TEXT NOT NULL,
    category            TEXT,
    severity            TEXT CHECK (severity IN ('low','medium','high','critical')),
    possible_causes     JSONB NOT NULL DEFAULT '[]',
    recommended_actions JSONB NOT NULL DEFAULT '[]',
    related_codes       TEXT[] NOT NULL DEFAULT '{}',
    is_active           BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.part_categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        TEXT NOT NULL DEFAULT '',
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    parent_id   UUID REFERENCES lookup.part_categories(id),
    icon        TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.fuel_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lookup.transmission_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lookup.engine_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL UNIQUE,
    fuel_type TEXT,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lookup.labor_rate_tiers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID,
    name        TEXT NOT NULL UNIQUE,
    hourly_rate DECIMAL(10,2) NOT NULL,
    description TEXT,
    is_global   BOOLEAN NOT NULL DEFAULT false,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.countries (
    code       TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    phone_code TEXT NOT NULL DEFAULT '',
    currency   TEXT NOT NULL DEFAULT 'USD',
    timezone   TEXT NOT NULL DEFAULT 'UTC',
    is_active  BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lookup.currencies (
    code            TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    symbol          TEXT NOT NULL,
    decimal_places  INT NOT NULL DEFAULT 2,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lookup.parts (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    category    TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.storage_locations (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
