CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.user_profiles (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL,
    user_id         UUID NOT NULL UNIQUE,
    title           TEXT,
    department      TEXT,
    employee_id     TEXT,
    timezone        TEXT NOT NULL DEFAULT 'UTC',
    locale          TEXT NOT NULL DEFAULT 'en',
    notification_prefs JSONB NOT NULL DEFAULT '{"email":true,"sms":false,"push":true}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users.permissions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT,
    module      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users.role_permissions (
    role            TEXT NOT NULL,
    permission_id   UUID NOT NULL REFERENCES users.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role, permission_id)
);

CREATE TABLE IF NOT EXISTS users.user_permissions (
    user_id         UUID NOT NULL,
    permission_id   UUID NOT NULL REFERENCES users.permissions(id) ON DELETE CASCADE,
    granted_by      UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE IF NOT EXISTS users.tenant_settings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL UNIQUE,
    business_hours  JSONB NOT NULL DEFAULT '{}',
    payment_config  JSONB NOT NULL DEFAULT '{}',
    notification_config JSONB NOT NULL DEFAULT '{}',
    branding        JSONB NOT NULL DEFAULT '{}',
    features        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users.roles (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users.tenants (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    type        TEXT NOT NULL DEFAULT 'branch',
    status      TEXT NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users.tenant_users (
    user_id   UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES users.tenants(id) ON DELETE CASCADE,
    role_id   UUID NOT NULL REFERENCES users.roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tenant_id)
);
