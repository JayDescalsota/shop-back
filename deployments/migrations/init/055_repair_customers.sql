-- AutoLab: repair customers
CREATE SCHEMA IF NOT EXISTS repair;

CREATE TABLE IF NOT EXISTS repair.customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    zip TEXT,
    notes TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed: customer records for app users
INSERT INTO repair.customers (tenant_id, name, email)
SELECT '00000000-0000-0000-0001-202605270006', 'Mobile User', 'mobile@autolab.com'
WHERE NOT EXISTS (SELECT 1 FROM repair.customers WHERE email = 'mobile@autolab.com');

INSERT INTO repair.customers (tenant_id, name, email)
SELECT '00000000-0000-0000-0001-202605270005', 'Fleet Manager', 'fleet@autolab.com'
WHERE NOT EXISTS (SELECT 1 FROM repair.customers WHERE email = 'fleet@autolab.com');
