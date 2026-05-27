-- AutoLab: notifications service
CREATE SCHEMA IF NOT EXISTS notifications;

CREATE TABLE notifications.notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    type TEXT NOT NULL, title TEXT NOT NULL, body TEXT,
    data JSONB NOT NULL DEFAULT '{}', is_read BOOLEAN NOT NULL DEFAULT false, read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE notifications.notification_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    type TEXT NOT NULL, channel TEXT NOT NULL, subject TEXT, template TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
