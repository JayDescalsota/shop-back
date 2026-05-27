-- AutoLab: users service
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.user_profiles (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    title           TEXT,
    department      TEXT,
    employee_id     TEXT,
    timezone        TEXT NOT NULL DEFAULT 'UTC',
    locale          TEXT NOT NULL DEFAULT 'en',
    notification_prefs JSONB NOT NULL DEFAULT '{"email":true,"sms":false,"push":true}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users.permissions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT,
    module      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users.role_permissions (
    role            auth.user_role NOT NULL,
    permission_id   UUID NOT NULL REFERENCES users.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role, permission_id)
);

CREATE TABLE users.user_permissions (
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    permission_id   UUID NOT NULL REFERENCES users.permissions(id) ON DELETE CASCADE,
    granted_by      UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE users.tenant_settings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL UNIQUE REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    business_hours  JSONB NOT NULL DEFAULT '{}',
    payment_config  JSONB NOT NULL DEFAULT '{}',
    notification_config JSONB NOT NULL DEFAULT '{}',
    branding        JSONB NOT NULL DEFAULT '{}',
    features        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
