-- AutoLab: auth service
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TYPE auth.user_role AS ENUM ('super_admin', 'tenant_admin', 'manager', 'staff', 'viewer');

CREATE TABLE auth.users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    email           TEXT NOT NULL,
    password_hash   TEXT NOT NULL,
    name            TEXT NOT NULL,
    phone           TEXT,
    avatar_url      TEXT,
    role            auth.user_role NOT NULL DEFAULT 'staff',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    email_verified  BOOLEAN NOT NULL DEFAULT false,
    last_login_at   TIMESTAMPTZ,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, email)
);

CREATE TABLE auth.user_invitations (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    email           TEXT NOT NULL,
    token           TEXT NOT NULL UNIQUE,
    role            auth.user_role NOT NULL DEFAULT 'staff',
    invited_by      UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    expires_at      TIMESTAMPTZ NOT NULL,
    accepted_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auth.refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auth.otp_codes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    purpose     TEXT NOT NULL,
    code        TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auth.password_resets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed: default admin user (password: password123)
-- Password hash generated for 'password123'
INSERT INTO auth.users (id, tenant_id, email, password_hash, name, role)
SELECT uuid_generate_v4(), id, 'admin@autolab.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'Admin', 'super_admin'::auth.user_role
FROM tenants.tenants LIMIT 1
ON CONFLICT DO NOTHING;
