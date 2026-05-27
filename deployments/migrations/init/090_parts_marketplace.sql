-- AutoLab: parts-marketplace service
CREATE SCHEMA IF NOT EXISTS parts_marketplace;

CREATE TYPE parts_marketplace.parts_quote_status AS ENUM ('pending','quoted','accepted','rejected','expired');

CREATE TABLE parts_marketplace.parts_catalog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    oe_number TEXT, sku TEXT NOT NULL, name TEXT NOT NULL, description TEXT, brand TEXT, category TEXT,
    compatible_vehicles JSONB NOT NULL DEFAULT '[]',
    unit_price DECIMAL(12,2) NOT NULL, available_qty INT NOT NULL DEFAULT 0,
    image_urls JSONB NOT NULL DEFAULT '[]', is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, sku)
);

CREATE TABLE parts_marketplace.parts_quotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    requestor_id UUID NOT NULL REFERENCES auth.users(id), shop_id UUID NOT NULL,
    status parts_marketplace.parts_quote_status NOT NULL DEFAULT 'pending',
    notes TEXT, valid_until DATE, total_amount DECIMAL(12,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parts_marketplace.parts_quote_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    quote_id UUID NOT NULL REFERENCES parts_marketplace.parts_quotes(id) ON DELETE CASCADE,
    part_id UUID REFERENCES parts_marketplace.parts_catalog(id),
    description TEXT NOT NULL, quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2), total_price DECIMAL(12,2), supplier_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parts_marketplace.parts_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    buyer_id UUID NOT NULL REFERENCES auth.users(id), seller_id UUID NOT NULL,
    quote_id UUID REFERENCES parts_marketplace.parts_quotes(id),
    order_number TEXT NOT NULL UNIQUE, status TEXT NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(12,2) NOT NULL, shipping_cost DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0, shipping_address JSONB, notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parts_marketplace.parts_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES parts_marketplace.parts_orders(id) ON DELETE CASCADE,
    part_id UUID NOT NULL REFERENCES parts_marketplace.parts_catalog(id),
    quantity INT NOT NULL, unit_price DECIMAL(12,2) NOT NULL, total_price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE parts_marketplace.parts_catalog ADD COLUMN IF NOT EXISTS search_vector TSVECTOR
    GENERATED ALWAYS AS (to_tsvector('english', coalesce(sku,'') || ' ' || coalesce(name,'') || ' ' || coalesce(description,'') || ' ' || coalesce(brand,'') || ' ' || coalesce(oe_number,''))) STORED;
CREATE INDEX IF NOT EXISTS idx_parts_catalog_search ON parts_marketplace.parts_catalog USING GIN(search_vector);
