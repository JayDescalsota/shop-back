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

-- Get make IDs for use in model inserts
DO $$
DECLARE
    toyota_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Toyota');
    honda_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Honda');
    ford_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Ford');
    chevy_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Chevrolet');
    bmw_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'BMW');
    mbenz_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Mercedes-Benz');
    audi_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Audi');
    vw_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Volkswagen');
    nissan_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Nissan');
    hyundai_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Hyundai');
    kia_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Kia');
    subaru_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Subaru');
    mazda_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Mazda');
    lexus_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Lexus');
    jeep_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Jeep');
    tesla_id UUID := (SELECT id FROM lookup.vehicle_makes WHERE name = 'Tesla');
BEGIN

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (toyota_id, 'Camry', 'camry', 1982, 'Sedan'),
    (toyota_id, 'Corolla', 'corolla', 1966, 'Sedan'),
    (toyota_id, 'RAV4', 'rav4', 1994, 'SUV'),
    (toyota_id, 'Tacoma', 'tacoma', 1995, 'Truck'),
    (toyota_id, 'Tundra', 'tundra', 1999, 'Truck'),
    (toyota_id, 'Highlander', 'highlander', 2000, 'SUV'),
    (toyota_id, '4Runner', '4runner', 1984, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (honda_id, 'Civic', 'civic', 1972, 'Sedan'),
    (honda_id, 'Accord', 'accord', 1976, 'Sedan'),
    (honda_id, 'CR-V', 'cr-v', 1995, 'SUV'),
    (honda_id, 'Pilot', 'pilot', 2002, 'SUV'),
    (honda_id, 'Odyssey', 'odyssey', 1994, 'Minivan')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (ford_id, 'F-150', 'f-150', 1948, 'Truck'),
    (ford_id, 'Mustang', 'mustang', 1964, 'Coupe'),
    (ford_id, 'Explorer', 'explorer', 1990, 'SUV'),
    (ford_id, 'Escape', 'escape', 2000, 'SUV'),
    (ford_id, 'Focus', 'focus', 1998, 'Sedan')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (chevy_id, 'Silverado', 'silverado', 1999, 'Truck'),
    (chevy_id, 'Equinox', 'equinox', 2004, 'SUV'),
    (chevy_id, 'Malibu', 'malibu', 1964, 'Sedan'),
    (chevy_id, 'Tahoe', 'tahoe', 1995, 'SUV'),
    (chevy_id, 'Camaro', 'camaro', 1966, 'Coupe')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (bmw_id, '3 Series', '3-series', 1975, 'Sedan'),
    (bmw_id, '5 Series', '5-series', 1972, 'Sedan'),
    (bmw_id, 'X3', 'x3', 2003, 'SUV'),
    (bmw_id, 'X5', 'x5', 1999, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (mbenz_id, 'C-Class', 'c-class', 1993, 'Sedan'),
    (mbenz_id, 'E-Class', 'e-class', 1953, 'Sedan'),
    (mbenz_id, 'GLC', 'glc', 2015, 'SUV'),
    (mbenz_id, 'GLE', 'gle', 2015, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (audi_id, 'A3', 'a3', 1996, 'Sedan'),
    (audi_id, 'A4', 'a4', 1994, 'Sedan'),
    (audi_id, 'Q5', 'q5', 2008, 'SUV'),
    (audi_id, 'Q7', 'q7', 2005, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (vw_id, 'Golf', 'golf', 1974, 'Hatchback'),
    (vw_id, 'Jetta', 'jetta', 1979, 'Sedan'),
    (vw_id, 'Passat', 'passat', 1973, 'Sedan'),
    (vw_id, 'Tiguan', 'tiguan', 2007, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (nissan_id, 'Altima', 'altima', 1992, 'Sedan'),
    (nissan_id, 'Sentra', 'sentra', 1982, 'Sedan'),
    (nissan_id, 'Rogue', 'rogue', 2007, 'SUV'),
    (nissan_id, 'Frontier', 'frontier', 1997, 'Truck')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (hyundai_id, 'Elantra', 'elantra', 1990, 'Sedan'),
    (hyundai_id, 'Sonata', 'sonata', 1985, 'Sedan'),
    (hyundai_id, 'Tucson', 'tucson', 2004, 'SUV'),
    (hyundai_id, 'Santa Fe', 'santa-fe', 2000, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (kia_id, 'Forte', 'forte', 2008, 'Sedan'),
    (kia_id, 'Sportage', 'sportage', 1993, 'SUV'),
    (kia_id, 'Sorento', 'sorento', 2002, 'SUV'),
    (kia_id, 'Telluride', 'telluride', 2019, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (subaru_id, 'Outback', 'outback', 1994, 'SUV'),
    (subaru_id, 'Forester', 'forester', 1997, 'SUV'),
    (subaru_id, 'Crosstrek', 'crosstrek', 2013, 'SUV'),
    (subaru_id, 'Impreza', 'impreza', 1992, 'Sedan')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (mazda_id, 'CX-5', 'cx-5', 2012, 'SUV'),
    (mazda_id, 'Mazda3', 'mazda3', 2003, 'Sedan'),
    (mazda_id, 'Mazda6', 'mazda6', 2002, 'Sedan'),
    (mazda_id, 'CX-9', 'cx-9', 2006, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (lexus_id, 'RX', 'rx', 1998, 'SUV'),
    (lexus_id, 'ES', 'es', 1989, 'Sedan'),
    (lexus_id, 'NX', 'nx', 2014, 'SUV'),
    (lexus_id, 'IS', 'is', 1999, 'Sedan')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (jeep_id, 'Wrangler', 'wrangler', 1986, 'SUV'),
    (jeep_id, 'Grand Cherokee', 'grand-cherokee', 1992, 'SUV'),
    (jeep_id, 'Cherokee', 'cherokee', 1974, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

INSERT INTO lookup.vehicle_models (make_id, name, slug, year_start, vehicle_type) VALUES
    (tesla_id, 'Model 3', 'model-3', 2017, 'Sedan'),
    (tesla_id, 'Model Y', 'model-y', 2020, 'SUV'),
    (tesla_id, 'Model S', 'model-s', 2012, 'Sedan'),
    (tesla_id, 'Model X', 'model-x', 2015, 'SUV')
ON CONFLICT (make_id, name) DO NOTHING;

END $$;

-- Seed: service types
INSERT INTO lookup.service_types (name, code, category, system, description, estimated_hours, is_global) VALUES
    ('Oil Change', 'oil-change', 'Maintenance', 'Engine', 'Engine oil and filter replacement', 1.0, true),
    ('Tire Rotation', 'tire-rotation', 'Maintenance', 'Tires/Wheels', 'Rotate tires to ensure even wear', 0.5, true),
    ('Brake Pad Replacement', 'brake-pads', 'Repair', 'Brakes', 'Replace front or rear brake pads', 2.0, true),
    ('Brake Rotor Replacement', 'brake-rotors', 'Repair', 'Brakes', 'Replace warped or worn brake rotors', 2.5, true),
    ('Engine Diagnostic', 'engine-diag', 'Diagnostic', 'Engine', 'Computer diagnostic scan of engine systems', 1.0, true),
    ('Check Engine Light', 'cel-diagnostic', 'Diagnostic', 'Engine', 'Diagnose and report check engine light cause', 1.5, true),
    ('Transmission Service', 'trans-service', 'Maintenance', 'Transmission', 'Transmission fluid flush and refill', 2.0, true),
    ('Coolant Flush', 'coolant-flush', 'Maintenance', 'Cooling System', 'Drain and refill engine coolant', 1.5, true),
    ('Battery Replacement', 'battery-replace', 'Repair', 'Electrical', 'Replace vehicle battery', 0.5, true),
    ('Air Filter Replacement', 'air-filter', 'Maintenance', 'Engine', 'Replace engine air filter', 0.25, true),
    ('Cabin Air Filter', 'cabin-filter', 'Maintenance', 'HVAC', 'Replace cabin air filter', 0.25, true),
    ('Spark Plug Replacement', 'spark-plugs', 'Maintenance', 'Engine', 'Replace spark plugs', 1.5, true),
    ('Wheel Alignment', 'wheel-alignment', 'Maintenance', 'Suspension/Steering', 'Align wheels to manufacturer specifications', 1.0, true),
    ('AC Recharge', 'ac-recharge', 'Repair', 'HVAC', 'Recharge air conditioning system', 1.5, true),
    ('Timing Belt Replacement', 'timing-belt', 'Maintenance', 'Engine', 'Replace timing belt and tensioner', 4.0, true),
    ('Serpentine Belt Replacement', 'serpentine-belt', 'Maintenance', 'Engine', 'Replace serpentine/accessory belt', 1.0, true),
    ('Shock/Strut Replacement', 'shock-strut', 'Repair', 'Suspension/Steering', 'Replace worn shocks or struts', 2.5, true),
    ('Exhaust System Repair', 'exhaust', 'Repair', 'Exhaust', 'Repair or replace exhaust components', 2.0, true),
    ('Fuel System Cleaning', 'fuel-system', 'Maintenance', 'Fuel System', 'Fuel injector cleaning and carbon removal', 2.0, true),
    ('State Inspection', 'inspection', 'Inspection', 'Maintenance', 'Complete state vehicle safety inspection', 1.0, true),
    ('Multi-Point Inspection', 'mpi', 'Inspection', 'Maintenance', 'Comprehensive vehicle inspection', 1.0, true),
    ('Headlight Restoration', 'headlight', 'Detailing', 'Electrical', 'Restore foggy or oxidized headlights', 1.0, true),
    ('Deep Interior Cleaning', 'interior-detail', 'Detailing', 'Body/Interior', 'Complete interior detail and shampoo', 3.0, true),
    ('Exterior Wash & Wax', 'exterior-detail', 'Detailing', 'Body/Interior', 'Hand wash and wax exterior', 2.0, true)
ON CONFLICT (name) DO UPDATE SET code=EXCLUDED.code, category=EXCLUDED.category, system=EXCLUDED.system, description=EXCLUDED.description, estimated_hours=EXCLUDED.estimated_hours;

-- Seed: lookup types
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
