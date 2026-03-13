-- Migration: 040_fleet_logistics_upgrade
-- Description: Enhance vehicles, drivers, deliveries for fleet management and portal visibility.

-- Enhance vehicles with fleet management fields
ALTER TABLE vehicles
ADD COLUMN IF NOT EXISTS vin VARCHAR(17),
ADD COLUMN IF NOT EXISTS year INTEGER,
ADD COLUMN IF NOT EXISTS make VARCHAR(100),
ADD COLUMN IF NOT EXISTS model VARCHAR(100),
ADD COLUMN IF NOT EXISTS insurance_expiry DATE,
ADD COLUMN IF NOT EXISTS next_service_date DATE,
ADD COLUMN IF NOT EXISTS odometer_miles INTEGER,
ADD COLUMN IF NOT EXISTS notes TEXT DEFAULT '';

-- Enhance drivers with HR/compliance fields
ALTER TABLE drivers
ADD COLUMN IF NOT EXISTS cdl_class VARCHAR(5),
ADD COLUMN IF NOT EXISTS cdl_expiry DATE,
ADD COLUMN IF NOT EXISTS hire_date DATE,
ADD COLUMN IF NOT EXISTS email VARCHAR(255);

-- Add geolocation to deliveries if missing
ALTER TABLE deliveries
ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION,
ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION,
ADD COLUMN IF NOT EXISTS estimated_arrival TIMESTAMPTZ;

-- Add route duration/distance if missing
ALTER TABLE delivery_routes
ADD COLUMN IF NOT EXISTS total_duration_mins INTEGER,
ADD COLUMN IF NOT EXISTS total_distance_miles DOUBLE PRECISION;

-- Add delivery time-window scheduling
ALTER TABLE deliveries
ADD COLUMN IF NOT EXISTS scheduled_start TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS scheduled_end TIMESTAMPTZ;
