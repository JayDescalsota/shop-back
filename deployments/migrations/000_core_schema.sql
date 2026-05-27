-- ============================================================================
-- AutoLab Core Schema
-- Each service owns a PostgreSQL schema. Tables are created within their
-- respective service schema and fully qualified (e.g. auth.users).
-- ============================================================================

-- ============================================================================
-- FOUNDATION (public schema)
-- ============================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS public.global_id_sequences (
    entity_type TEXT PRIMARY KEY,
    last_id     BIGINT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: tenants
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS tenants;

CREATE TYPE tenants.tenant_type AS ENUM (
    'auto_owner', 'repair_shop', 'parts_store', 'platform'
);

CREATE TYPE tenants.tenant_status AS ENUM (
    'active', 'suspended', 'trial', 'cancelled'
);

CREATE TABLE tenants.tenants (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    type        tenants.tenant_type NOT NULL,
    status      tenants.tenant_status NOT NULL DEFAULT 'trial',
    domain      TEXT UNIQUE,
    config      JSONB NOT NULL DEFAULT '{}',
    settings    JSONB NOT NULL DEFAULT '{}',
    billing_plan TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: auth
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TYPE auth.user_role AS ENUM (
    'super_admin', 'tenant_admin', 'manager', 'staff', 'viewer'
);

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

-- ============================================================================
-- SERVICE: users
-- ============================================================================
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

-- ============================================================================
-- SERVICE: lookup
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS lookup;

CREATE TABLE lookup.vehicle_makes (
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name    TEXT NOT NULL UNIQUE,
    logo    TEXT
);

CREATE TABLE lookup.vehicle_models (
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    make_id UUID NOT NULL REFERENCES lookup.vehicle_makes(id) ON DELETE CASCADE,
    name    TEXT NOT NULL,
    year_start INT,
    year_end   INT,
    UNIQUE(make_id, name)
);

CREATE TABLE lookup.service_types (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    category    TEXT,
    estimated_hours DECIMAL(6,2),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE lookup.diagnostic_codes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    category    TEXT,
    severity    TEXT CHECK (severity IN ('low', 'medium', 'high', 'critical'))
);

CREATE TABLE lookup.part_categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    parent_id   UUID REFERENCES lookup.part_categories(id)
);

CREATE TABLE lookup.fuel_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE);
CREATE TABLE lookup.transmission_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE);
CREATE TABLE lookup.engine_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE);

CREATE TABLE lookup.labor_rate_tiers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    hourly_rate DECIMAL(10,2) NOT NULL,
    description TEXT
);

CREATE TABLE lookup.countries (
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code    TEXT NOT NULL UNIQUE,
    name    TEXT NOT NULL
);

CREATE TABLE lookup.currencies (
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code    TEXT NOT NULL UNIQUE,
    name    TEXT NOT NULL,
    symbol  TEXT NOT NULL
);

-- ============================================================================
-- SERVICE: vehicles
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS vehicles;

CREATE TABLE vehicles.vehicles (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    owner_id        UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    make_id         UUID REFERENCES lookup.vehicle_makes(id),
    model_id        UUID REFERENCES lookup.vehicle_models(id),
    year            INT,
    vin             TEXT,
    license_plate   TEXT,
    color           TEXT,
    engine_type     TEXT,
    transmission    TEXT,
    current_mileage INT NOT NULL DEFAULT 0,
    fuel_type       TEXT,
    purchase_date   DATE,
    purchase_price  DECIMAL(12,2),
    image_url       TEXT,
    notes           TEXT,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vehicles.vehicle_mileage_log (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    vehicle_id  UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    mileage     INT NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    source      TEXT NOT NULL DEFAULT 'manual',
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vehicles.vehicle_maintenance (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    vehicle_id      UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    service_type    TEXT NOT NULL,
    description     TEXT,
    mileage_at_service INT,
    cost            DECIMAL(12,2),
    service_date    DATE NOT NULL,
    provider_name   TEXT,
    provider_shop_id UUID,
    status          TEXT NOT NULL DEFAULT 'completed',
    notes           TEXT,
    documents       JSONB NOT NULL DEFAULT '[]',
    next_due_mileage INT,
    next_due_date   DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vehicles.vehicle_expenses (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    vehicle_id  UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    category    TEXT NOT NULL,
    amount      DECIMAL(12,2) NOT NULL,
    date        DATE NOT NULL,
    description TEXT,
    receipt_url TEXT,
    recurring   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vehicles.vehicle_documents (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    vehicle_id  UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    doc_type    TEXT NOT NULL,
    title       TEXT NOT NULL,
    file_url    TEXT NOT NULL,
    expiry_date DATE,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: bookings
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS bookings;

CREATE TYPE bookings.booking_status AS ENUM (
    'pending', 'confirmed', 'in_progress', 'completed', 'cancelled', 'no_show'
);

CREATE TABLE bookings.bookings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    shop_id         UUID NOT NULL,
    customer_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    vehicle_id      UUID REFERENCES vehicles.vehicles(id) ON DELETE SET NULL,
    service_type    TEXT NOT NULL,
    description     TEXT,
    status          bookings.booking_status NOT NULL DEFAULT 'pending',
    scheduled_date  DATE NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    assigned_staff  UUID,
    notes           TEXT,
    cancellation_reason TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE bookings.booking_slots (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    shop_id     UUID NOT NULL,
    date        DATE NOT NULL,
    start_time  TIME NOT NULL,
    end_time    TIME NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    staff_id    UUID,
    booking_id  UUID REFERENCES bookings.bookings(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, shop_id, date, start_time, staff_id)
);

CREATE TABLE bookings.shop_availability (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    shop_id     UUID NOT NULL,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    open_time   TIME NOT NULL,
    close_time  TIME NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, shop_id, day_of_week)
);

-- ============================================================================
-- SERVICE: repair
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS repair;

CREATE TYPE repair.repair_order_status AS ENUM (
    'draft', 'estimate', 'approved', 'in_progress', 'completed', 'invoiced', 'paid', 'cancelled'
);

CREATE TABLE repair.repair_orders (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    shop_id         UUID NOT NULL,
    booking_id      UUID REFERENCES bookings.bookings(id) ON DELETE SET NULL,
    customer_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    vehicle_id      UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    status          repair.repair_order_status NOT NULL DEFAULT 'draft',
    odometer_in     INT,
    odometer_out    INT,
    customer_complaint TEXT,
    technician_notes TEXT,
    diagnosis       TEXT,
    estimated_cost  DECIMAL(12,2),
    actual_cost     DECIMAL(12,2),
    tax_amount      DECIMAL(12,2),
    total_amount    DECIMAL(12,2),
    discount_amount DECIMAL(12,2) DEFAULT 0,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE repair.repair_line_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    repair_order_id UUID NOT NULL REFERENCES repair.repair_orders(id) ON DELETE CASCADE,
    item_type       TEXT NOT NULL CHECK (item_type IN ('labor', 'part', 'service', 'fee')),
    description     TEXT NOT NULL,
    quantity        DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit_price      DECIMAL(12,2) NOT NULL,
    total_price     DECIMAL(12,2) NOT NULL,
    part_id         UUID,
    technician_id   UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE repair.job_cards (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    repair_order_id UUID NOT NULL REFERENCES repair.repair_orders(id) ON DELETE CASCADE,
    assigned_to     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    task            TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending',
    notes           TEXT,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    estimated_hours DECIMAL(6,2),
    actual_hours    DECIMAL(6,2),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE repair.diagnostics (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    repair_order_id UUID NOT NULL REFERENCES repair.repair_orders(id) ON DELETE CASCADE,
    diagnostic_code TEXT,
    description     TEXT NOT NULL,
    severity        TEXT CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    recommended_action TEXT,
    performed_by    UUID REFERENCES auth.users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: inventory
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS inventory;

CREATE TABLE inventory.inventory_categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    parent_id   UUID REFERENCES inventory.inventory_categories(id) ON DELETE SET NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, name)
);

CREATE TABLE inventory.inventory_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    category_id     UUID REFERENCES inventory.inventory_categories(id),
    sku             TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    brand           TEXT,
    compatible_makes JSONB NOT NULL DEFAULT '[]',
    compatible_models JSONB NOT NULL DEFAULT '[]',
    unit_price      DECIMAL(12,2) NOT NULL,
    cost_price      DECIMAL(12,2),
    quantity_on_hand INT NOT NULL DEFAULT 0,
    quantity_reserved INT NOT NULL DEFAULT 0,
    quantity_available INT GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
    reorder_point   INT NOT NULL DEFAULT 10,
    reorder_quantity INT NOT NULL DEFAULT 50,
    location        TEXT,
    barcode         TEXT,
    image_url       TEXT,
    weight          DECIMAL(10,2),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, sku)
);

CREATE TABLE inventory.inventory_movements (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    item_id         UUID NOT NULL REFERENCES inventory.inventory_items(id) ON DELETE CASCADE,
    movement_type   TEXT NOT NULL CHECK (movement_type IN ('in', 'out', 'adjustment', 'return', 'transfer')),
    quantity        INT NOT NULL,
    reference_type  TEXT,
    reference_id    UUID,
    unit_cost       DECIMAL(12,2),
    notes           TEXT,
    performed_by    UUID NOT NULL REFERENCES auth.users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory.suppliers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    contact_person TEXT,
    email       TEXT,
    phone       TEXT,
    address     TEXT,
    payment_terms TEXT,
    lead_time_days INT,
    rating      INT CHECK (rating BETWEEN 1 AND 5),
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory.purchase_orders (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    supplier_id     UUID REFERENCES inventory.suppliers(id),
    order_number    TEXT NOT NULL UNIQUE,
    status          TEXT NOT NULL DEFAULT 'draft',
    total_amount    DECIMAL(12,2),
    tax_amount      DECIMAL(12,2),
    notes           TEXT,
    ordered_by      UUID NOT NULL REFERENCES auth.users(id),
    approved_by     UUID REFERENCES auth.users(id),
    expected_date   DATE,
    received_date   DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory.purchase_order_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    purchase_order_id UUID NOT NULL REFERENCES inventory.purchase_orders(id) ON DELETE CASCADE,
    item_id         UUID NOT NULL REFERENCES inventory.inventory_items(id),
    quantity        INT NOT NULL,
    unit_price      DECIMAL(12,2) NOT NULL,
    total_price     DECIMAL(12,2) NOT NULL,
    quantity_received INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: parts-marketplace
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS parts_marketplace;

CREATE TYPE parts_marketplace.parts_quote_status AS ENUM (
    'pending', 'quoted', 'accepted', 'rejected', 'expired'
);

CREATE TABLE parts_marketplace.parts_catalog (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    oe_number       TEXT,
    sku             TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    brand           TEXT,
    category        TEXT,
    compatible_vehicles JSONB NOT NULL DEFAULT '[]',
    unit_price      DECIMAL(12,2) NOT NULL,
    available_qty   INT NOT NULL DEFAULT 0,
    image_urls      JSONB NOT NULL DEFAULT '[]',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, sku)
);

CREATE TABLE parts_marketplace.parts_quotes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    requestor_id    UUID NOT NULL REFERENCES auth.users(id),
    shop_id         UUID NOT NULL,
    status          parts_marketplace.parts_quote_status NOT NULL DEFAULT 'pending',
    notes           TEXT,
    valid_until     DATE,
    total_amount    DECIMAL(12,2),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parts_marketplace.parts_quote_items (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    quote_id    UUID NOT NULL REFERENCES parts_marketplace.parts_quotes(id) ON DELETE CASCADE,
    part_id     UUID REFERENCES parts_marketplace.parts_catalog(id),
    description TEXT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    unit_price  DECIMAL(12,2),
    total_price DECIMAL(12,2),
    supplier_id UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parts_marketplace.parts_orders (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    buyer_id        UUID NOT NULL REFERENCES auth.users(id),
    seller_id       UUID NOT NULL,
    quote_id        UUID REFERENCES parts_marketplace.parts_quotes(id),
    order_number    TEXT NOT NULL UNIQUE,
    status          TEXT NOT NULL DEFAULT 'pending',
    total_amount    DECIMAL(12,2) NOT NULL,
    shipping_cost   DECIMAL(12,2) DEFAULT 0,
    tax_amount      DECIMAL(12,2) DEFAULT 0,
    shipping_address JSONB,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parts_marketplace.parts_order_items (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    order_id    UUID NOT NULL REFERENCES parts_marketplace.parts_orders(id) ON DELETE CASCADE,
    part_id     UUID NOT NULL REFERENCES parts_marketplace.parts_catalog(id),
    quantity    INT NOT NULL,
    unit_price  DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: payments
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS payments;

CREATE TYPE payments.payment_status AS ENUM (
    'pending', 'processing', 'succeeded', 'failed', 'refunded', 'partially_refunded'
);

CREATE TYPE payments.invoice_status AS ENUM (
    'draft', 'sent', 'overdue', 'paid', 'cancelled', 'refunded'
);

CREATE TABLE payments.payments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    amount          DECIMAL(12,2) NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'USD',
    status          payments.payment_status NOT NULL DEFAULT 'pending',
    payment_method  TEXT,
    provider        TEXT,
    provider_payment_id TEXT,
    reference_type  TEXT NOT NULL,
    reference_id    UUID NOT NULL,
    description     TEXT,
    metadata        JSONB NOT NULL DEFAULT '{}',
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments.invoices (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    invoice_number  TEXT NOT NULL UNIQUE,
    status          payments.invoice_status NOT NULL DEFAULT 'draft',
    bill_to_id      UUID NOT NULL REFERENCES auth.users(id),
    bill_to_type    TEXT NOT NULL,
    line_items      JSONB NOT NULL DEFAULT '[]',
    subtotal        DECIMAL(12,2) NOT NULL,
    tax_amount      DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_amount    DECIMAL(12,2) NOT NULL,
    amount_paid     DECIMAL(12,2) NOT NULL DEFAULT 0,
    amount_due      DECIMAL(12,2) GENERATED ALWAYS AS (total_amount - amount_paid) STORED,
    due_date        DATE,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments.payment_methods (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    provider        TEXT NOT NULL,
    provider_method_id TEXT NOT NULL,
    type            TEXT NOT NULL,
    last_four       TEXT,
    expiry_month    INT,
    expiry_year     INT,
    is_default      BOOLEAN NOT NULL DEFAULT false,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments.subscriptions (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id           UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    plan                TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'active',
    current_period_start DATE NOT NULL,
    current_period_end   DATE NOT NULL,
    provider_subscription_id TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: payroll
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS payroll;

CREATE TABLE payroll.employees (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    employee_id     TEXT,
    department      TEXT,
    position        TEXT,
    hire_date       DATE NOT NULL,
    termination_date DATE,
    employment_type TEXT NOT NULL DEFAULT 'full_time',
    salary          DECIMAL(12,2),
    hourly_rate     DECIMAL(10,2),
    pay_frequency   TEXT NOT NULL DEFAULT 'monthly',
    tax_id          TEXT,
    bank_info       JSONB,
    emergency_contact JSONB,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payroll.employee_payrolls (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    period_start    DATE NOT NULL,
    period_end      DATE NOT NULL,
    base_pay        DECIMAL(12,2) NOT NULL,
    overtime_pay    DECIMAL(12,2) NOT NULL DEFAULT 0,
    bonuses         DECIMAL(12,2) NOT NULL DEFAULT 0,
    deductions      DECIMAL(12,2) NOT NULL DEFAULT 0,
    taxes           DECIMAL(12,2) NOT NULL DEFAULT 0,
    net_pay         DECIMAL(12,2) NOT NULL,
    status          TEXT NOT NULL DEFAULT 'calculated',
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payroll.time_shifts (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    start_time  TIME NOT NULL,
    end_time    TIME NOT NULL,
    break_min   INT NOT NULL DEFAULT 0,
    total_hours DECIMAL(5,2) GENERATED ALWAYS AS (
        EXTRACT(EPOCH FROM (end_time - start_time)) / 3600 - (break_min::DECIMAL / 60)
    ) STORED,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payroll.attendance (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    clock_in    TIMESTAMPTZ NOT NULL,
    clock_out   TIMESTAMPTZ,
    status      TEXT NOT NULL DEFAULT 'present',
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(employee_id, date)
);

CREATE TABLE payroll.leave_requests (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    leave_type      TEXT NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    total_days      INT NOT NULL,
    reason          TEXT,
    status          TEXT NOT NULL DEFAULT 'pending',
    approved_by     UUID REFERENCES auth.users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: notifications
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS notifications;

CREATE TABLE notifications.notifications (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    title       TEXT NOT NULL,
    body        TEXT,
    data        JSONB NOT NULL DEFAULT '{}',
    is_read     BOOLEAN NOT NULL DEFAULT false,
    read_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE notifications.notification_templates (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    channel     TEXT NOT NULL,
    subject     TEXT,
    template    TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- SERVICE: analytics
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE analytics.audit_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id     UUID REFERENCES auth.users(id),
    action      TEXT NOT NULL,
    resource    TEXT NOT NULL,
    resource_id UUID,
    details     JSONB NOT NULL DEFAULT '{}',
    ip_address  TEXT,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE analytics.data_access_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id),
    user_id     UUID NOT NULL REFERENCES auth.users(id),
    resource    TEXT NOT NULL,
    resource_id UUID NOT NULL,
    permission  TEXT NOT NULL,
    granted     BOOLEAN NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE analytics.file_uploads (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES auth.users(id),
    filename    TEXT NOT NULL,
    original_name TEXT NOT NULL,
    mime_type   TEXT NOT NULL,
    size_bytes  BIGINT NOT NULL,
    bucket      TEXT NOT NULL,
    key         TEXT NOT NULL,
    url         TEXT,
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- FULL-TEXT SEARCH SUPPORT
-- ============================================================================
ALTER TABLE vehicles.vehicles ADD COLUMN IF NOT EXISTS search_vector TSVECTOR
    GENERATED ALWAYS AS (
        to_tsvector('english',
            coalesce(license_plate,'') || ' ' ||
            coalesce(vin,'') || ' ' ||
            coalesce(color,'')
        )
    ) STORED;

ALTER TABLE inventory.inventory_items ADD COLUMN IF NOT EXISTS search_vector TSVECTOR
    GENERATED ALWAYS AS (
        to_tsvector('english',
            coalesce(sku,'') || ' ' ||
            coalesce(name,'') || ' ' ||
            coalesce(description,'') || ' ' ||
            coalesce(brand,'') || ' ' ||
            coalesce(barcode,'')
        )
    ) STORED;

ALTER TABLE parts_marketplace.parts_catalog ADD COLUMN IF NOT EXISTS search_vector TSVECTOR
    GENERATED ALWAYS AS (
        to_tsvector('english',
            coalesce(sku,'') || ' ' ||
            coalesce(name,'') || ' ' ||
            coalesce(description,'') || ' ' ||
            coalesce(brand,'') || ' ' ||
            coalesce(oe_number,'')
        )
    ) STORED;

CREATE INDEX IF NOT EXISTS idx_vehicles_search ON vehicles.vehicles USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_inventory_search ON inventory.inventory_items USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_parts_catalog_search ON parts_marketplace.parts_catalog USING GIN(search_vector);

-- ============================================================================
-- INDEXES
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_users_tenant ON auth.users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_tenant ON vehicles.vehicles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_owner ON vehicles.vehicles(owner_id);
CREATE INDEX IF NOT EXISTS idx_bookings_tenant ON bookings.bookings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_bookings_shop ON bookings.bookings(shop_id);
CREATE INDEX IF NOT EXISTS idx_bookings_customer ON bookings.bookings(customer_id);
CREATE INDEX IF NOT EXISTS idx_bookings_date ON bookings.bookings(tenant_id, scheduled_date);
CREATE INDEX IF NOT EXISTS idx_repair_orders_tenant ON repair.repair_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_repair_orders_customer ON repair.repair_orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_repair_orders_vehicle ON repair.repair_orders(vehicle_id);
CREATE INDEX IF NOT EXISTS idx_inventory_items_tenant ON inventory.inventory_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_item ON inventory.inventory_movements(item_id);
CREATE INDEX IF NOT EXISTS idx_parts_catalog_tenant ON parts_marketplace.parts_catalog(tenant_id);
CREATE INDEX IF NOT EXISTS idx_parts_quotes_tenant ON parts_marketplace.parts_quotes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_parts_orders_buyer ON parts_marketplace.parts_orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_parts_orders_seller ON parts_marketplace.parts_orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_payments_tenant ON payments.payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_reference ON payments.payments(reference_type, reference_id);
CREATE INDEX IF NOT EXISTS idx_invoices_tenant ON payments.invoices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_employees_tenant ON payroll.employees(tenant_id);
CREATE INDEX IF NOT EXISTS idx_attendance_employee ON payroll.attendance(employee_id, date);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant ON analytics.audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON analytics.audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications.notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_user_profiles_user ON users.user_profiles(user_id);

-- ============================================================================
-- ROW-LEVEL SECURITY
-- ============================================================================
ALTER TABLE auth.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicles.vehicles ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicles.vehicle_mileage_log ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicles.vehicle_maintenance ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicles.vehicle_expenses ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicles.vehicle_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE bookings.bookings ENABLE ROW LEVEL SECURITY;
ALTER TABLE bookings.booking_slots ENABLE ROW LEVEL SECURITY;
ALTER TABLE bookings.shop_availability ENABLE ROW LEVEL SECURITY;
ALTER TABLE repair.repair_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE repair.repair_line_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE repair.job_cards ENABLE ROW LEVEL SECURITY;
ALTER TABLE repair.diagnostics ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.inventory_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.inventory_movements ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.suppliers ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.purchase_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.purchase_order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE parts_marketplace.parts_catalog ENABLE ROW LEVEL SECURITY;
ALTER TABLE parts_marketplace.parts_quotes ENABLE ROW LEVEL SECURITY;
ALTER TABLE parts_marketplace.parts_quote_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE parts_marketplace.parts_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE parts_marketplace.parts_order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments.payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments.invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments.payment_methods ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments.subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE payroll.employees ENABLE ROW LEVEL SECURITY;
ALTER TABLE payroll.employee_payrolls ENABLE ROW LEVEL SECURITY;
ALTER TABLE payroll.time_shifts ENABLE ROW LEVEL SECURITY;
ALTER TABLE payroll.attendance ENABLE ROW LEVEL SECURITY;
ALTER TABLE payroll.leave_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE analytics.audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE analytics.notifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE analytics.file_uploads ENABLE ROW LEVEL SECURITY;

CREATE OR REPLACE FUNCTION public.get_current_tenant_id() RETURNS UUID
    LANGUAGE SQL STABLE
AS $$ SELECT NULLIF(current_setting('app.current_tenant_id', true), '')::UUID; $$;

CREATE OR REPLACE FUNCTION public.get_current_user_id() RETURNS UUID
    LANGUAGE SQL STABLE
AS $$ SELECT NULLIF(current_setting('app.current_user_id', true), '')::UUID; $$;

CREATE OR REPLACE FUNCTION public.get_current_user_role() RETURNS TEXT
    LANGUAGE SQL STABLE
AS $$ SELECT NULLIF(current_setting('app.current_user_role', true), ''); $$;

CREATE OR REPLACE FUNCTION public.tenant_isolation_policy() RETURNS TEXT
    LANGUAGE SQL STABLE
AS $$ SELECT 'tenant_id = get_current_tenant_id()'; $$;

DO $$ DECLARE tbl TEXT;
BEGIN
    FOR tbl IN SELECT unnest(ARRAY[
        'auth.users','vehicles.vehicles','vehicles.vehicle_mileage_log',
        'vehicles.vehicle_maintenance','vehicles.vehicle_expenses','vehicles.vehicle_documents',
        'bookings.bookings','bookings.booking_slots','bookings.shop_availability',
        'repair.repair_orders','repair.repair_line_items','repair.job_cards','repair.diagnostics',
        'inventory.inventory_items','inventory.inventory_movements','inventory.suppliers',
        'inventory.purchase_orders','inventory.purchase_order_items',
        'parts_marketplace.parts_catalog','parts_marketplace.parts_quotes',
        'parts_marketplace.parts_quote_items','parts_marketplace.parts_orders',
        'parts_marketplace.parts_order_items','payments.payments','payments.invoices',
        'payments.payment_methods','payments.subscriptions','payroll.employees',
        'payroll.employee_payrolls','payroll.time_shifts','payroll.attendance',
        'payroll.leave_requests','analytics.audit_logs','analytics.notifications',
        'analytics.file_uploads'
    ])
    LOOP
        EXECUTE format(
            'DROP POLICY IF EXISTS tenant_isolation ON %I.%s;
             CREATE POLICY tenant_isolation ON %I.%s
                 FOR ALL USING (tenant_id = get_current_tenant_id())
                 WITH CHECK (tenant_id = get_current_tenant_id())',
            split_part(tbl, '.', 1)::regnamespace,
            split_part(tbl, '.', 2),
            split_part(tbl, '.', 1)::regnamespace,
            split_part(tbl, '.', 2)
        );
    END LOOP;
END $$;

CREATE OR REPLACE FUNCTION public.is_super_admin() RETURNS BOOLEAN
    LANGUAGE SQL STABLE
AS $$ SELECT get_current_user_role() = 'super_admin'; $$;
