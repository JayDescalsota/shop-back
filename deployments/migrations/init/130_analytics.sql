-- AutoLab: analytics service
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE analytics.audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id),
    action TEXT NOT NULL, resource TEXT NOT NULL, resource_id UUID,
    details JSONB NOT NULL DEFAULT '{}', ip_address TEXT, user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE analytics.data_access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id),
    user_id UUID NOT NULL REFERENCES auth.users(id),
    resource TEXT NOT NULL, resource_id UUID NOT NULL, permission TEXT NOT NULL, granted BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE analytics.file_uploads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id),
    filename TEXT NOT NULL, original_name TEXT NOT NULL, mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL, bucket TEXT NOT NULL, key TEXT NOT NULL, url TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
