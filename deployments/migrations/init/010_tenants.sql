-- AutoLab: tenants service
CREATE SCHEMA IF NOT EXISTS tenants;

CREATE TYPE tenants.tenant_type AS ENUM ('auto_owner', 'repair_shop', 'parts_store', 'platform');
CREATE TYPE tenants.tenant_status AS ENUM ('active', 'suspended', 'trial', 'cancelled');

CREATE TABLE tenants.tenants (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    type        tenants.tenant_type NOT NULL,
    status      tenants.tenant_status NOT NULL DEFAULT 'trial',
    domain      TEXT UNIQUE,
    config      JSONB NOT NULL DEFAULT '{}',
    settings    JSONB NOT NULL DEFAULT '{}',
    billing_plan TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
