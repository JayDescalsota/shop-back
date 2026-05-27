-- AutoLab: lookup service (reference data)
CREATE SCHEMA IF NOT EXISTS lookup;

CREATE TABLE lookup.vehicle_makes (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE, logo TEXT);
CREATE TABLE lookup.vehicle_models (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), make_id UUID NOT NULL REFERENCES lookup.vehicle_makes(id) ON DELETE CASCADE, name TEXT NOT NULL, year_start INT, year_end INT, UNIQUE(make_id, name));
CREATE TABLE lookup.service_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE, description TEXT, category TEXT, estimated_hours DECIMAL(6,2), created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE TABLE lookup.diagnostic_codes (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), code TEXT NOT NULL UNIQUE, description TEXT NOT NULL, category TEXT, severity TEXT CHECK (severity IN ('low','medium','high','critical')));
CREATE TABLE lookup.part_categories (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE, description TEXT, parent_id UUID REFERENCES lookup.part_categories(id));
CREATE TABLE lookup.fuel_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE);
CREATE TABLE lookup.transmission_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE);
CREATE TABLE lookup.engine_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE);
CREATE TABLE lookup.labor_rate_tiers (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name TEXT NOT NULL UNIQUE, hourly_rate DECIMAL(10,2) NOT NULL, description TEXT);
CREATE TABLE lookup.countries (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), code TEXT NOT NULL UNIQUE, name TEXT NOT NULL);
CREATE TABLE lookup.currencies (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), code TEXT NOT NULL UNIQUE, name TEXT NOT NULL, symbol TEXT NOT NULL);

-- Seed: vehicle makes
INSERT INTO lookup.vehicle_makes (name) VALUES
    ('Toyota'),('Honda'),('Ford'),('Chevrolet'),('BMW'),('Mercedes-Benz'),('Audi'),
    ('Volkswagen'),('Nissan'),('Hyundai'),('Kia'),('Subaru'),('Mazda'),('Lexus'),
    ('Jeep'),('Dodge'),('Chrysler'),('GMC'),('Cadillac'),('Buick'),('Acura'),
    ('Infiniti'),('Lincoln'),('Volvo'),('Porsche'),('Tesla'),('Rivian'),
    ('Ferrari'),('Lamborghini'),('McLaren')
ON CONFLICT DO NOTHING;

-- Seed: fuel, transmission, engine types
INSERT INTO lookup.fuel_types (name) VALUES ('Gasoline'),('Diesel'),('Electric'),('Hybrid'),('Plug-in Hybrid'),('Hydrogen'),('Ethanol'),('Biodiesel') ON CONFLICT DO NOTHING;
INSERT INTO lookup.transmission_types (name) VALUES ('Automatic'),('Manual'),('CVT'),('DCT'),('Semi-Automatic'),('Electric Drive') ON CONFLICT DO NOTHING;
INSERT INTO lookup.engine_types (name) VALUES ('Inline-4'),('V6'),('V8'),('V10'),('V12'),('Flat-4'),('Flat-6'),('Inline-6'),('W12'),('W16'),('Rotary'),('Electric Motor'),('Twin-Turbo'),('Turbocharged'),('Supercharged'),('Hybrid Powertrain') ON CONFLICT DO NOTHING;

-- Seed: currencies
INSERT INTO lookup.currencies (code, name, symbol) VALUES
    ('USD','US Dollar','$'),('EUR','Euro','€'),('GBP','British Pound','£'),
    ('JPY','Japanese Yen','¥'),('CAD','Canadian Dollar','C$'),('AUD','Australian Dollar','A$'),
    ('CHF','Swiss Franc','Fr'),('CNY','Chinese Yuan','¥'),('MXN','Mexican Peso','Mex$'),
    ('BRL','Brazilian Real','R$')
ON CONFLICT DO NOTHING;

-- Seed: countries
INSERT INTO lookup.countries (code, name) VALUES
    ('US','United States'),('CA','Canada'),('GB','United Kingdom'),('DE','Germany'),
    ('FR','France'),('IT','Italy'),('JP','Japan'),('KR','South Korea'),
    ('MX','Mexico'),('BR','Brazil'),('AU','Australia'),('CN','China')
ON CONFLICT DO NOTHING;
