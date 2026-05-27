-- AutoLab: Shared bootstrap (extensions, public helpers)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS public.global_id_sequences (
    entity_type TEXT PRIMARY KEY,
    last_id     BIGINT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- RLS helpers
CREATE OR REPLACE FUNCTION public.get_current_tenant_id() RETURNS UUID
    LANGUAGE SQL STABLE AS $$ SELECT NULLIF(current_setting('app.current_tenant_id', true), '')::UUID; $$;
CREATE OR REPLACE FUNCTION public.get_current_user_id() RETURNS UUID
    LANGUAGE SQL STABLE AS $$ SELECT NULLIF(current_setting('app.current_user_id', true), '')::UUID; $$;
CREATE OR REPLACE FUNCTION public.get_current_user_role() RETURNS TEXT
    LANGUAGE SQL STABLE AS $$ SELECT NULLIF(current_setting('app.current_user_role', true), ''); $$;
CREATE OR REPLACE FUNCTION public.tenant_isolation_policy() RETURNS TEXT
    LANGUAGE SQL STABLE AS $$ SELECT 'tenant_id = get_current_tenant_id()'; $$;
CREATE OR REPLACE FUNCTION public.is_super_admin() RETURNS BOOLEAN
    LANGUAGE SQL STABLE AS $$ SELECT get_current_user_role() = 'super_admin'; $$;
