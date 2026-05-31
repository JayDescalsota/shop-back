-- AutoLab: lookup service (reference data) - v2 with full column set
CREATE SCHEMA IF NOT EXISTS lookup;

CREATE TABLE IF NOT EXISTS lookup.vehicle_makes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL UNIQUE,
    slug        TEXT NOT NULL DEFAULT '',
    logo_url    TEXT,
    country     TEXT,
    founded_year INT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    sort_order  INT NOT NULL DEFAULT 0,
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.vehicle_models (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    make_id     UUID NOT NULL REFERENCES lookup.vehicle_makes(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL DEFAULT '',
    year_start  INT,
    year_end    INT,
    vehicle_type TEXT,
    trim_levels JSONB NOT NULL DEFAULT '[]',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(make_id, name)
);

CREATE TABLE IF NOT EXISTS lookup.service_types (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id       UUID,
    code            TEXT NOT NULL DEFAULT '',
    name            TEXT NOT NULL UNIQUE,
    description     TEXT,
    category        TEXT,
    estimated_hours DECIMAL(6,2),
    is_global       BOOLEAN NOT NULL DEFAULT true,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.diagnostic_codes (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code                TEXT NOT NULL UNIQUE,
    system              TEXT NOT NULL DEFAULT '',
    description         TEXT NOT NULL,
    category            TEXT,
    severity            TEXT CHECK (severity IN ('low','medium','high','critical')),
    possible_causes     JSONB NOT NULL DEFAULT '[]',
    recommended_actions JSONB NOT NULL DEFAULT '[]',
    related_codes       TEXT[] NOT NULL DEFAULT '{}',
    is_active           BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.part_categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        TEXT NOT NULL DEFAULT '',
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    parent_id   UUID REFERENCES lookup.part_categories(id),
    icon        TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.fuel_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), code TEXT NOT NULL DEFAULT '', name TEXT NOT NULL UNIQUE, description TEXT NOT NULL DEFAULT '', is_active BOOLEAN NOT NULL DEFAULT true, sort_order INT NOT NULL DEFAULT 0);
CREATE TABLE IF NOT EXISTS lookup.transmission_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), code TEXT NOT NULL DEFAULT '', name TEXT NOT NULL UNIQUE, description TEXT NOT NULL DEFAULT '', is_active BOOLEAN NOT NULL DEFAULT true, sort_order INT NOT NULL DEFAULT 0);
CREATE TABLE IF NOT EXISTS lookup.engine_types (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), code TEXT NOT NULL DEFAULT '', name TEXT NOT NULL UNIQUE, fuel_type TEXT, description TEXT NOT NULL DEFAULT '', is_active BOOLEAN NOT NULL DEFAULT true, sort_order INT NOT NULL DEFAULT 0);

CREATE TABLE IF NOT EXISTS lookup.labor_rate_tiers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID,
    name        TEXT NOT NULL UNIQUE,
    hourly_rate DECIMAL(10,2) NOT NULL,
    description TEXT,
    is_global   BOOLEAN NOT NULL DEFAULT false,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lookup.countries (
    code       TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    phone_code TEXT NOT NULL DEFAULT '',
    currency   TEXT NOT NULL DEFAULT 'USD',
    timezone   TEXT NOT NULL DEFAULT 'UTC',
    is_active  BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lookup.currencies (
    code            TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    symbol          TEXT NOT NULL,
    decimal_places  INT NOT NULL DEFAULT 2,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INT NOT NULL DEFAULT 0
);

-- Seed: vehicle makes
INSERT INTO lookup.vehicle_makes (name, slug, country, sort_order) VALUES
    ('Toyota', 'toyota', 'Japan', 1),
    ('Honda', 'honda', 'Japan', 2),
    ('Ford', 'ford', 'United States', 3),
    ('Chevrolet', 'chevrolet', 'United States', 4),
    ('BMW', 'bmw', 'Germany', 5),
    ('Mercedes-Benz', 'mercedes-benz', 'Germany', 6),
    ('Audi', 'audi', 'Germany', 7),
    ('Volkswagen', 'volkswagen', 'Germany', 8),
    ('Nissan', 'nissan', 'Japan', 9),
    ('Hyundai', 'hyundai', 'South Korea', 10),
    ('Kia', 'kia', 'South Korea', 11),
    ('Subaru', 'subaru', 'Japan', 12),
    ('Mazda', 'mazda', 'Japan', 13),
    ('Lexus', 'lexus', 'Japan', 14),
    ('Jeep', 'jeep', 'United States', 15),
    ('Dodge', 'dodge', 'United States', 16),
    ('Chrysler', 'chrysler', 'United States', 17),
    ('GMC', 'gmc', 'United States', 18),
    ('Cadillac', 'cadillac', 'United States', 19),
    ('Buick', 'buick', 'United States', 20),
    ('Acura', 'acura', 'Japan', 21),
    ('Infiniti', 'infiniti', 'Japan', 22),
    ('Lincoln', 'lincoln', 'United States', 23),
    ('Volvo', 'volvo', 'Sweden', 24),
    ('Porsche', 'porsche', 'Germany', 25),
    ('Tesla', 'tesla', 'United States', 26),
    ('Rivian', 'rivian', 'United States', 27),
    ('Ferrari', 'ferrari', 'Italy', 28),
    ('Lamborghini', 'lamborghini', 'Italy', 29),
    ('McLaren', 'mclaren', 'United Kingdom', 30)
ON CONFLICT (name) DO UPDATE SET slug=EXCLUDED.slug, country=EXCLUDED.country, sort_order=EXCLUDED.sort_order;

-- Seed: vehicle models
INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Camry','Corolla','RAV4','Tacoma','Tundra','Highlander','4Runner']), unnest(ARRAY['camry','corolla','rav4','tacoma','tundra','highlander','4runner']), unnest(ARRAY[1982,1966,1994,1995,1999,2000,1984]), unnest(ARRAY['Sedan','Sedan','SUV','Truck','Truck','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'Toyota';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Civic','Accord','CR-V','Pilot','Odyssey']), unnest(ARRAY['civic','accord','cr-v','pilot','odyssey']), unnest(ARRAY[1972,1976,1995,2002,1994]), unnest(ARRAY['Sedan','Sedan','SUV','SUV','Minivan']) FROM lookup.vehicle_makes WHERE name = 'Honda';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['F-150','Mustang','Explorer','Escape','Focus']), unnest(ARRAY['f-150','mustang','explorer','escape','focus']), unnest(ARRAY[1948,1964,1990,2000,1998]), unnest(ARRAY['Truck','Coupe','SUV','SUV','Sedan']) FROM lookup.vehicle_makes WHERE name = 'Ford';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Silverado','Equinox','Malibu','Tahoe','Camaro']), unnest(ARRAY['silverado','equinox','malibu','tahoe','camaro']), unnest(ARRAY[1999,2004,1964,1995,1966]), unnest(ARRAY['Truck','SUV','Sedan','SUV','Coupe']) FROM lookup.vehicle_makes WHERE name = 'Chevrolet';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['3 Series','5 Series','X3','X5']), unnest(ARRAY['3-series','5-series','x3','x5']), unnest(ARRAY[1975,1972,2003,1999]), unnest(ARRAY['Sedan','Sedan','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'BMW';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['C-Class','E-Class','GLC','GLE']), unnest(ARRAY['c-class','e-class','glc','gle']), unnest(ARRAY[1993,1953,2015,2015]), unnest(ARRAY['Sedan','Sedan','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'Mercedes-Benz';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['A3','A4','Q5','Q7']), unnest(ARRAY['a3','a4','q5','q7']), unnest(ARRAY[1996,1994,2008,2005]), unnest(ARRAY['Sedan','Sedan','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'Audi';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Golf','Jetta','Passat','Tiguan']), unnest(ARRAY['golf','jetta','passat','tiguan']), unnest(ARRAY[1974,1979,1973,2007]), unnest(ARRAY['Hatchback','Sedan','Sedan','SUV']) FROM lookup.vehicle_makes WHERE name = 'Volkswagen';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Altima','Sentra','Rogue','Frontier']), unnest(ARRAY['altima','sentra','rogue','frontier']), unnest(ARRAY[1992,1982,2007,1997]), unnest(ARRAY['Sedan','Sedan','SUV','Truck']) FROM lookup.vehicle_makes WHERE name = 'Nissan';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Elantra','Sonata','Tucson','Santa Fe']), unnest(ARRAY['elantra','sonata','tucson','santa-fe']), unnest(ARRAY[1990,1985,2004,2000]), unnest(ARRAY['Sedan','Sedan','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'Hyundai';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Forte','Sportage','Sorento','Telluride']), unnest(ARRAY['forte','sportage','sorento','telluride']), unnest(ARRAY[2008,1993,2002,2019]), unnest(ARRAY['Sedan','SUV','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'Kia';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Outback','Forester','Crosstrek','Impreza']), unnest(ARRAY['outback','forester','crosstrek','impreza']), unnest(ARRAY[1994,1997,2013,1992]), unnest(ARRAY['SUV','SUV','SUV','Sedan']) FROM lookup.vehicle_makes WHERE name = 'Subaru';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['CX-5','Mazda3','Mazda6','CX-9']), unnest(ARRAY['cx-5','mazda3','mazda6','cx-9']), unnest(ARRAY[2012,2003,2002,2006]), unnest(ARRAY['SUV','Sedan','Sedan','SUV']) FROM lookup.vehicle_makes WHERE name = 'Mazda';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['RX','ES','NX','IS']), unnest(ARRAY['rx','es','nx','is']), unnest(ARRAY[1998,1989,2014,1999]), unnest(ARRAY['SUV','Sedan','SUV','Sedan']) FROM lookup.vehicle_makes WHERE name = 'Lexus';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Wrangler','Grand Cherokee','Cherokee']), unnest(ARRAY['wrangler','grand-cherokee','cherokee']), unnest(ARRAY[1986,1992,1974]), unnest(ARRAY['SUV','SUV','SUV']) FROM lookup.vehicle_makes WHERE name = 'Jeep';

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type)
SELECT id, unnest(ARRAY['Model 3','Model Y','Model S','Model X']), unnest(ARRAY['model-3','model-y','model-s','model-x']), unnest(ARRAY[2017,2020,2012,2015]), unnest(ARRAY['Sedan','SUV','Sedan','SUV']) FROM lookup.vehicle_makes WHERE name = 'Tesla';

-- Seed: service types
INSERT INTO lookup.service_types (name, code, category, description, estimated_hours, is_global) VALUES
    ('Oil Change', 'oil-change', 'Maintenance', 'Engine oil and filter replacement', 1.0, true),
    ('Tire Rotation', 'tire-rotation', 'Maintenance', 'Rotate tires to ensure even wear', 0.5, true),
    ('Brake Pad Replacement', 'brake-pads', 'Repair', 'Replace front or rear brake pads', 2.0, true),
    ('Brake Rotor Replacement', 'brake-rotors', 'Repair', 'Replace warped or worn brake rotors', 2.5, true),
    ('Engine Diagnostic', 'engine-diag', 'Diagnostic', 'Computer diagnostic scan of engine systems', 1.0, true),
    ('Check Engine Light', 'cel-diagnostic', 'Diagnostic', 'Diagnose and report check engine light cause', 1.5, true),
    ('Transmission Service', 'trans-service', 'Maintenance', 'Transmission fluid flush and refill', 2.0, true),
    ('Coolant Flush', 'coolant-flush', 'Maintenance', 'Drain and refill engine coolant', 1.5, true),
    ('Battery Replacement', 'battery-replace', 'Repair', 'Replace vehicle battery', 0.5, true),
    ('Air Filter Replacement', 'air-filter', 'Maintenance', 'Replace engine air filter', 0.25, true),
    ('Cabin Air Filter', 'cabin-filter', 'Maintenance', 'Replace cabin air filter', 0.25, true),
    ('Spark Plug Replacement', 'spark-plugs', 'Maintenance', 'Replace spark plugs', 1.5, true),
    ('Wheel Alignment', 'wheel-alignment', 'Maintenance', 'Align wheels to manufacturer specifications', 1.0, true),
    ('AC Recharge', 'ac-recharge', 'Repair', 'Recharge air conditioning system', 1.5, true),
    ('Timing Belt Replacement', 'timing-belt', 'Maintenance', 'Replace timing belt and tensioner', 4.0, true),
    ('Serpentine Belt Replacement', 'serpentine-belt', 'Maintenance', 'Replace serpentine/accessory belt', 1.0, true),
    ('Shock/Strut Replacement', 'shock-strut', 'Repair', 'Replace worn shocks or struts', 2.5, true),
    ('Exhaust System Repair', 'exhaust', 'Repair', 'Repair or replace exhaust components', 2.0, true),
    ('Fuel System Cleaning', 'fuel-system', 'Maintenance', 'Fuel injector cleaning and carbon removal', 2.0, true),
    ('State Inspection', 'inspection', 'Inspection', 'Complete state vehicle safety inspection', 1.0, true),
    ('Multi-Point Inspection', 'mpi', 'Inspection', 'Comprehensive vehicle inspection', 1.0, true),
    ('Headlight Restoration', 'headlight', 'Detailing', 'Restore foggy or oxidized headlights', 1.0, true),
    ('Deep Interior Cleaning', 'interior-detail', 'Detailing', 'Complete interior detail and shampoo', 3.0, true),
    ('Exterior Wash & Wax', 'exterior-detail', 'Detailing', 'Hand wash and wax exterior', 2.0, true)
ON CONFLICT (name) DO UPDATE SET code=EXCLUDED.code, category=EXCLUDED.category, description=EXCLUDED.description, estimated_hours=EXCLUDED.estimated_hours;

-- Seed: fuel, transmission, engine types
INSERT INTO lookup.fuel_types (name, code, sort_order) VALUES ('Gasoline','gasoline',1),('Diesel','diesel',2),('Electric','electric',3),('Hybrid','hybrid',4),('Plug-in Hybrid','plugin-hybrid',5),('Hydrogen','hydrogen',6),('Ethanol','ethanol',7),('Biodiesel','biodiesel',8) ON CONFLICT (name) DO UPDATE SET code=EXCLUDED.code, sort_order=EXCLUDED.sort_order;
INSERT INTO lookup.transmission_types (name, code, sort_order) VALUES ('Automatic','automatic',1),('Manual','manual',2),('CVT','cvt',3),('DCT','dct',4),('Semi-Automatic','semi-auto',5),('Electric Drive','electric',6) ON CONFLICT (name) DO UPDATE SET code=EXCLUDED.code, sort_order=EXCLUDED.sort_order;
INSERT INTO lookup.engine_types (name, code, sort_order) VALUES ('Inline-4','i4',1),('V6','v6',2),('V8','v8',3),('V10','v10',4),('V12','v12',5),('Flat-4','flat4',6),('Flat-6','flat6',7),('Inline-6','i6',8),('W12','w12',9),('W16','w16',10),('Rotary','rotary',11),('Electric Motor','electric',12),('Twin-Turbo','twin-turbo',13),('Turbocharged','turbo',14),('Supercharged','supercharged',15),('Hybrid Powertrain','hybrid',16) ON CONFLICT (name) DO UPDATE SET code=EXCLUDED.code, sort_order=EXCLUDED.sort_order;

-- Seed: currencies
INSERT INTO lookup.currencies (code, name, symbol, decimal_places) VALUES
    ('USD','US Dollar','$',2),('EUR','Euro','€',2),('GBP','British Pound','£',2),
    ('JPY','Japanese Yen','¥',0),('CAD','Canadian Dollar','C$',2),('AUD','Australian Dollar','A$',2),
    ('CHF','Swiss Franc','Fr',2),('CNY','Chinese Yuan','¥',2),('MXN','Mexican Peso','Mex$',2),
    ('BRL','Brazilian Real','R$',2)
ON CONFLICT (code) DO UPDATE SET name=EXCLUDED.name, symbol=EXCLUDED.symbol, decimal_places=EXCLUDED.decimal_places;

-- Seed: countries
INSERT INTO lookup.countries (code, name, phone_code, currency) VALUES
    ('US','United States','+1','USD'),('CA','Canada','+1','CAD'),('GB','United Kingdom','+44','GBP'),
    ('DE','Germany','+49','EUR'),('FR','France','+33','EUR'),('IT','Italy','+39','EUR'),
    ('JP','Japan','+81','JPY'),('KR','South Korea','+82','KRW'),('MX','Mexico','+52','MXN'),
    ('BR','Brazil','+55','BRL'),('AU','Australia','+61','AUD'),('CN','China','+86','CNY')
ON CONFLICT (code) DO UPDATE SET name=EXCLUDED.name, phone_code=EXCLUDED.phone_code, currency=EXCLUDED.currency;
