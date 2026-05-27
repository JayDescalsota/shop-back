-- AutoLab: global indexes & RLS policies
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

-- RLS policies
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
ALTER TABLE analytics.file_uploads ENABLE ROW LEVEL SECURITY;

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
        'payroll.leave_requests','analytics.audit_logs','analytics.file_uploads'
    ])
    LOOP
        EXECUTE format(
            'DROP POLICY IF EXISTS tenant_isolation ON %I.%s;
             CREATE POLICY tenant_isolation ON %I.%s
                 FOR ALL USING (tenant_id = get_current_tenant_id())
                 WITH CHECK (tenant_id = get_current_tenant_id())',
            split_part(tbl, '.', 1)::regnamespace, split_part(tbl, '.', 2),
            split_part(tbl, '.', 1)::regnamespace, split_part(tbl, '.', 2)
        );
    END LOOP;
END $$;
