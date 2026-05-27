-- AutoLab: payments service
CREATE SCHEMA IF NOT EXISTS payments;

CREATE TYPE payments.payment_status AS ENUM ('pending','processing','succeeded','failed','refunded','partially_refunded');
CREATE TYPE payments.invoice_status AS ENUM ('draft','sent','overdue','paid','cancelled','refunded');

CREATE TABLE payments.payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL, currency TEXT NOT NULL DEFAULT 'USD',
    status payments.payment_status NOT NULL DEFAULT 'pending',
    payment_method TEXT, provider TEXT, provider_payment_id TEXT,
    reference_type TEXT NOT NULL, reference_id UUID NOT NULL, description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    paid_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments.invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    invoice_number TEXT NOT NULL UNIQUE, status payments.invoice_status NOT NULL DEFAULT 'draft',
    bill_to_id UUID NOT NULL REFERENCES auth.users(id), bill_to_type TEXT NOT NULL,
    line_items JSONB NOT NULL DEFAULT '[]',
    subtotal DECIMAL(12,2) NOT NULL, tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0, total_amount DECIMAL(12,2) NOT NULL,
    amount_paid DECIMAL(12,2) NOT NULL DEFAULT 0,
    amount_due DECIMAL(12,2) GENERATED ALWAYS AS (total_amount - amount_paid) STORED,
    due_date DATE, notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments.payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL, provider_method_id TEXT NOT NULL, type TEXT NOT NULL,
    last_four TEXT, expiry_month INT, expiry_year INT, is_default BOOLEAN NOT NULL DEFAULT false,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments.subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    plan TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'active',
    current_period_start DATE NOT NULL, current_period_end DATE NOT NULL,
    provider_subscription_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
