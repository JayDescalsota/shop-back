-- AutoLab: inventory service
CREATE SCHEMA IF NOT EXISTS inventory;

CREATE TABLE inventory.inventory_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL, parent_id UUID REFERENCES inventory.inventory_categories(id) ON DELETE SET NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, name)
);

CREATE TABLE inventory.inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    category_id UUID REFERENCES inventory.inventory_categories(id),
    sku TEXT NOT NULL, name TEXT NOT NULL, description TEXT, brand TEXT,
    compatible_makes JSONB NOT NULL DEFAULT '[]', compatible_models JSONB NOT NULL DEFAULT '[]',
    unit_price DECIMAL(12,2) NOT NULL, cost_price DECIMAL(12,2),
    quantity_on_hand INT NOT NULL DEFAULT 0, quantity_reserved INT NOT NULL DEFAULT 0,
    quantity_available INT GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
    reorder_point INT NOT NULL DEFAULT 10, reorder_quantity INT NOT NULL DEFAULT 50,
    location TEXT, barcode TEXT, image_url TEXT, weight DECIMAL(10,2),
    is_active BOOLEAN NOT NULL DEFAULT true, metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, sku)
);

CREATE TABLE inventory.inventory_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES inventory.inventory_items(id) ON DELETE CASCADE,
    movement_type TEXT NOT NULL CHECK (movement_type IN ('in','out','adjustment','return','transfer')),
    quantity INT NOT NULL, reference_type TEXT, reference_id UUID,
    unit_cost DECIMAL(12,2), notes TEXT,
    performed_by UUID NOT NULL REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory.suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL, contact_person TEXT, email TEXT, phone TEXT, address TEXT,
    payment_terms TEXT, lead_time_days INT,
    rating INT CHECK (rating BETWEEN 1 AND 5), is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory.purchase_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    supplier_id UUID REFERENCES inventory.suppliers(id),
    order_number TEXT NOT NULL UNIQUE, status TEXT NOT NULL DEFAULT 'draft',
    total_amount DECIMAL(12,2), tax_amount DECIMAL(12,2), notes TEXT,
    ordered_by UUID NOT NULL REFERENCES auth.users(id), approved_by UUID REFERENCES auth.users(id),
    expected_date DATE, received_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory.purchase_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    purchase_order_id UUID NOT NULL REFERENCES inventory.purchase_orders(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES inventory.inventory_items(id),
    quantity INT NOT NULL, unit_price DECIMAL(12,2) NOT NULL, total_price DECIMAL(12,2) NOT NULL,
    quantity_received INT NOT NULL DEFAULT 0, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE inventory.inventory_items ADD COLUMN IF NOT EXISTS search_vector TSVECTOR
    GENERATED ALWAYS AS (to_tsvector('english', coalesce(sku,'') || ' ' || coalesce(name,'') || ' ' || coalesce(description,'') || ' ' || coalesce(brand,'') || ' ' || coalesce(barcode,''))) STORED;
CREATE INDEX IF NOT EXISTS idx_inventory_search ON inventory.inventory_items USING GIN(search_vector);
