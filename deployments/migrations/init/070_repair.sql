-- AutoLab: repair service
CREATE SCHEMA IF NOT EXISTS repair;

CREATE TYPE repair.repair_order_status AS ENUM ('draft','estimate','approved','in_progress','completed','invoiced','paid','cancelled');

CREATE TABLE repair.repair_orders (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    shop_id         UUID NOT NULL,
    booking_id      UUID REFERENCES bookings.bookings(id) ON DELETE SET NULL,
    customer_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    vehicle_id      UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    status          repair.repair_order_status NOT NULL DEFAULT 'draft',
    odometer_in     INT, odometer_out INT,
    customer_complaint TEXT, technician_notes TEXT, diagnosis TEXT,
    estimated_cost DECIMAL(12,2), actual_cost DECIMAL(12,2),
    tax_amount DECIMAL(12,2), total_amount DECIMAL(12,2),
    discount_amount DECIMAL(12,2) DEFAULT 0,
    started_at TIMESTAMPTZ, completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE repair.repair_line_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    repair_order_id UUID NOT NULL REFERENCES repair.repair_orders(id) ON DELETE CASCADE,
    item_type TEXT NOT NULL CHECK (item_type IN ('labor','part','service','fee')),
    description TEXT NOT NULL, quantity DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2) NOT NULL, total_price DECIMAL(12,2) NOT NULL,
    part_id UUID, technician_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE repair.job_cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    repair_order_id UUID NOT NULL REFERENCES repair.repair_orders(id) ON DELETE CASCADE,
    assigned_to UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    task TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'pending', notes TEXT,
    started_at TIMESTAMPTZ, completed_at TIMESTAMPTZ,
    estimated_hours DECIMAL(6,2), actual_hours DECIMAL(6,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE repair.diagnostics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    repair_order_id UUID NOT NULL REFERENCES repair.repair_orders(id) ON DELETE CASCADE,
    diagnostic_code TEXT, description TEXT NOT NULL,
    severity TEXT CHECK (severity IN ('low','medium','high','critical')),
    recommended_action TEXT, performed_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
