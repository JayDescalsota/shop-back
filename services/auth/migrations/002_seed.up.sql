BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS auth.tenants (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    address     TEXT,
    app_id      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO tenants.tenants (id, name, type, status)
VALUES
  ('00000000-0000-0000-0001-202605270001', 'Downtown Auto Parts', 'parts_store', 'active'),
  ('00000000-0000-0000-0001-202605270002', 'Westside Auto Parts', 'parts_store', 'active'),
  ('00000000-0000-0000-0001-202605270003', 'Main Street Auto Care', 'repair_shop', 'active'),
  ('00000000-0000-0000-0001-202605270004', 'Oak Avenue Auto Repair', 'repair_shop', 'active'),
  ('00000000-0000-0000-0001-202605270005', 'Fleet Operations Center', 'auto_owner', 'active'),
  ('00000000-0000-0000-0001-202605270006', 'Mobile Service 1', 'repair_shop', 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO auth.tenants (id, name, address, app_id)
VALUES
  ('00000000-0000-0000-0001-202605270001', 'Downtown Auto Parts', '100 Main St, Downtown', '00000000-0000-0000-0000-202605270001'),
  ('00000000-0000-0000-0001-202605270002', 'Westside Auto Parts', '456 Oak Ave, Westside', '00000000-0000-0000-0000-202605270001'),
  ('00000000-0000-0000-0001-202605270003', 'Main Street Auto Care', '123 Main St, Downtown', '00000000-0000-0000-0000-202605270002'),
  ('00000000-0000-0000-0001-202605270004', 'Oak Avenue Auto Repair', '789 Pine Rd, Easton', '00000000-0000-0000-0000-202605270002'),
  ('00000000-0000-0000-0001-202605270005', 'Fleet Operations Center', '500 Fleet St, Metro', '00000000-0000-0000-0000-202605270005'),
  ('00000000-0000-0000-0001-202605270006', 'Mobile Service 1', '200 Mobile Ave, Midtown', '00000000-0000-0000-0000-202605270003')
ON CONFLICT (id) DO NOTHING;

INSERT INTO auth.users (id, tenant_id, email, password_hash, name, role)
VALUES
  (gen_random_uuid(), '00000000-0000-0000-0001-202605270001', 'admin@autolab.com', '$2a$10$sKYGbavKaA3R77juKJbDS.0hO1dubavHHapVYyl391g7VwaIG7J4i', 'Admin', 'super_admin'),
  (gen_random_uuid(), '00000000-0000-0000-0001-202605270001', 'user@autolab.com', '$2a$10$sKYGbavKaA3R77juKJbDS.0hO1dubavHHapVYyl391g7VwaIG7J4i', 'User', 'staff'),
  (gen_random_uuid(), '00000000-0000-0000-0001-202605270001', 'storeadmin@autolab.com', '$2a$10$sKYGbavKaA3R77juKJbDS.0hO1dubavHHapVYyl391g7VwaIG7J4i', 'Store Admin', 'tenant_admin'),
  (gen_random_uuid(), '00000000-0000-0000-0001-202605270001', 'mobile@autolab.com', '$2a$10$sKYGbavKaA3R77juKJbDS.0hO1dubavHHapVYyl391g7VwaIG7J4i', 'Mobile User', 'staff'),
  (gen_random_uuid(), '00000000-0000-0000-0001-202605270001', 'superadmin@autolab.com', '$2a$10$sKYGbavKaA3R77juKJbDS.0hO1dubavHHapVYyl391g7VwaIG7J4i', 'Super Admin', 'super_admin'),
  (gen_random_uuid(), '00000000-0000-0000-0001-202605270001', 'fleet@autolab.com', '$2a$10$sKYGbavKaA3R77juKJbDS.0hO1dubavHHapVYyl391g7VwaIG7J4i', 'Fleet Manager', 'manager')
ON CONFLICT (tenant_id, email) DO NOTHING;

WITH user_ids AS (
  SELECT id, email FROM auth.users WHERE tenant_id = '00000000-0000-0000-0001-202605270001'
)
INSERT INTO auth.user_apps (user_id, app_id, role)
SELECT u.id, v.app_id, v.role
FROM (VALUES
  ('admin@autolab.com',    '00000000-0000-0000-0000-202605270001', 'admin'),
  ('admin@autolab.com',    '00000000-0000-0000-0000-202605270002', 'admin'),
  ('admin@autolab.com',    '00000000-0000-0000-0000-202605270003', 'admin'),
  ('admin@autolab.com',    '00000000-0000-0000-0000-202605270005', 'admin'),
  ('user@autolab.com',     '00000000-0000-0000-0000-202605270002', 'user'),
  ('storeadmin@autolab.com','00000000-0000-0000-0000-202605270001', 'admin'),
  ('mobile@autolab.com',   '00000000-0000-0000-0000-202605270003', 'user'),
  ('superadmin@autolab.com','00000000-0000-0000-0000-202605270004', 'superAdmin'),
  ('superadmin@autolab.com','00000000-0000-0000-0000-202605270001', 'superAdmin'),
  ('superadmin@autolab.com','00000000-0000-0000-0000-202605270002', 'superAdmin'),
  ('superadmin@autolab.com','00000000-0000-0000-0000-202605270003', 'superAdmin'),
  ('superadmin@autolab.com','00000000-0000-0000-0000-202605270005', 'superAdmin'),
  ('fleet@autolab.com',    '00000000-0000-0000-0000-202605270005', 'admin')
) AS v(email, app_id, role)
JOIN user_ids u ON u.email = v.email
ON CONFLICT (user_id, app_id) DO NOTHING;

INSERT INTO auth.user_tenants (email, tenant_id)
VALUES
  ('admin@autolab.com',      '00000000-0000-0000-0001-202605270001'),
  ('admin@autolab.com',      '00000000-0000-0000-0001-202605270002'),
  ('admin@autolab.com',      '00000000-0000-0000-0001-202605270003'),
  ('admin@autolab.com',      '00000000-0000-0000-0001-202605270004'),
  ('admin@autolab.com',      '00000000-0000-0000-0001-202605270005'),
  ('admin@autolab.com',      '00000000-0000-0000-0001-202605270006'),
  ('user@autolab.com',       '00000000-0000-0000-0001-202605270003'),
  ('user@autolab.com',       '00000000-0000-0000-0001-202605270004'),
  ('storeadmin@autolab.com', '00000000-0000-0000-0001-202605270001'),
  ('storeadmin@autolab.com', '00000000-0000-0000-0001-202605270002'),
  ('mobile@autolab.com',     '00000000-0000-0000-0001-202605270006'),
  ('superadmin@autolab.com', '00000000-0000-0000-0001-202605270001'),
  ('superadmin@autolab.com', '00000000-0000-0000-0001-202605270002'),
  ('superadmin@autolab.com', '00000000-0000-0000-0001-202605270003'),
  ('superadmin@autolab.com', '00000000-0000-0000-0001-202605270004'),
  ('superadmin@autolab.com', '00000000-0000-0000-0001-202605270005'),
  ('superadmin@autolab.com', '00000000-0000-0000-0001-202605270006'),
  ('fleet@autolab.com',      '00000000-0000-0000-0001-202605270005')
ON CONFLICT (email, tenant_id) DO NOTHING;

COMMIT;
