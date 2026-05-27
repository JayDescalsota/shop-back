-- AutoLab: bookings service
CREATE SCHEMA IF NOT EXISTS bookings;

CREATE TYPE bookings.booking_status AS ENUM ('pending','confirmed','in_progress','completed','cancelled','no_show');

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
