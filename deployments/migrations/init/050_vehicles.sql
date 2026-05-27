-- AutoLab: vehicles service
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

ALTER TABLE vehicles.vehicles ADD COLUMN IF NOT EXISTS search_vector TSVECTOR
    GENERATED ALWAYS AS (to_tsvector('english', coalesce(license_plate,'') || ' ' || coalesce(vin,'') || ' ' || coalesce(color,''))) STORED;
CREATE INDEX IF NOT EXISTS idx_vehicles_search ON vehicles.vehicles USING GIN(search_vector);
