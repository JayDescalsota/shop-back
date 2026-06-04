CREATE SCHEMA IF NOT EXISTS repair;

CREATE TABLE IF NOT EXISTS repair.staff_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    appointment_id UUID NOT NULL REFERENCES repair.appointments(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL,
    staff_name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'mechanic',
    status TEXT NOT NULL DEFAULT 'assigned',
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    total_minutes INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_staff_assignments_appointment ON repair.staff_assignments(appointment_id);
CREATE INDEX IF NOT EXISTS idx_staff_assignments_staff ON repair.staff_assignments(staff_id);
CREATE INDEX IF NOT EXISTS idx_staff_assignments_status ON repair.staff_assignments(status);
