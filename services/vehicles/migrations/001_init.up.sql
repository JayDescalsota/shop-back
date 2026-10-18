CREATE SCHEMA IF NOT EXISTS vehicles;

CREATE TABLE IF NOT EXISTS vehicles.vehicle_makes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    logo TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicles.vehicle_models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    make_id UUID NOT NULL REFERENCES vehicles.vehicle_makes(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    year_start INT,
    year_end INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicles.vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    make_id UUID REFERENCES vehicles.vehicle_makes(id),
    model_id UUID REFERENCES vehicles.vehicle_models(id),
    year INT,
    vin TEXT,
    license_plate TEXT,
    color TEXT,
    engine_type TEXT,
    transmission TEXT,
    current_mileage INT NOT NULL DEFAULT 0,
    fuel_type TEXT,
    purchase_date TEXT,
    purchase_price DECIMAL(12,2),
    image_url TEXT,
    notes TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    repair_status TEXT NOT NULL DEFAULT 'none',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicles.maintenance_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    service_type TEXT NOT NULL,
    description TEXT,
    mileage_at_service INT,
    cost DECIMAL(12,2),
    service_date TEXT NOT NULL,
    provider_name TEXT,
    status TEXT NOT NULL DEFAULT 'completed',
    next_due_mileage INT,
    next_due_date TEXT,
    documents JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicles.expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    category TEXT NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    date TEXT NOT NULL,
    description TEXT,
    receipt_url TEXT,
    recurring BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicles.documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL REFERENCES vehicles.vehicles(id) ON DELETE CASCADE,
    doc_type TEXT NOT NULL,
    title TEXT NOT NULL,
    file_url TEXT NOT NULL,
    expiry_date TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_vehicles_tenant ON vehicles.vehicles(tenant_id);
CREATE INDEX idx_vehicles_owner ON vehicles.vehicles(owner_id);
CREATE INDEX idx_vehicles_make ON vehicles.vehicles(make_id);
CREATE INDEX idx_vehicles_model ON vehicles.vehicles(model_id);
CREATE INDEX idx_maintenance_vehicle ON vehicles.maintenance_records(vehicle_id);
CREATE INDEX idx_expenses_vehicle ON vehicles.expenses(vehicle_id);
CREATE INDEX idx_documents_vehicle ON vehicles.documents(vehicle_id);
