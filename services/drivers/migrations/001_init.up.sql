CREATE SCHEMA IF NOT EXISTS drivers;

CREATE TABLE IF NOT EXISTS drivers.drivers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    role TEXT NOT NULL DEFAULT 'driver',
    license_number TEXT,
    license_class TEXT,
    license_expiry TEXT,
    date_of_birth TEXT,
    address TEXT,
    emergency_contact TEXT,
    emergency_phone TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    assigned_vehicle_id UUID,
    notes TEXT,
    hire_date TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
