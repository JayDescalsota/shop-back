-- AutoLab: payroll service
CREATE SCHEMA IF NOT EXISTS payroll;

CREATE TABLE payroll.employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    employee_id TEXT, department TEXT, position TEXT,
    hire_date DATE NOT NULL, termination_date DATE,
    employment_type TEXT NOT NULL DEFAULT 'full_time',
    salary DECIMAL(12,2), hourly_rate DECIMAL(10,2), pay_frequency TEXT NOT NULL DEFAULT 'monthly',
    tax_id TEXT, bank_info JSONB, emergency_contact JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payroll.employee_payrolls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    period_start DATE NOT NULL, period_end DATE NOT NULL,
    base_pay DECIMAL(12,2) NOT NULL, overtime_pay DECIMAL(12,2) NOT NULL DEFAULT 0,
    bonuses DECIMAL(12,2) NOT NULL DEFAULT 0, deductions DECIMAL(12,2) NOT NULL DEFAULT 0,
    taxes DECIMAL(12,2) NOT NULL DEFAULT 0, net_pay DECIMAL(12,2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'calculated', paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payroll.time_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    date DATE NOT NULL, start_time TIME NOT NULL, end_time TIME NOT NULL,
    break_min INT NOT NULL DEFAULT 0,
    total_hours DECIMAL(5,2) GENERATED ALWAYS AS (EXTRACT(EPOCH FROM (end_time - start_time)) / 3600 - (break_min::DECIMAL / 60)) STORED,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payroll.attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    date DATE NOT NULL, clock_in TIMESTAMPTZ NOT NULL, clock_out TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'present', notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(employee_id, date)
);

CREATE TABLE payroll.leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES payroll.employees(id) ON DELETE CASCADE,
    leave_type TEXT NOT NULL, start_date DATE NOT NULL, end_date DATE NOT NULL,
    total_days INT NOT NULL, reason TEXT, status TEXT NOT NULL DEFAULT 'pending',
    approved_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
