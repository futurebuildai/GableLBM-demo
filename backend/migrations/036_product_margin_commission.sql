-- Add unit cost, target margin, and commission rate to products

ALTER TABLE products 
ADD COLUMN average_unit_cost NUMERIC(15,4) DEFAULT 0,
ADD COLUMN target_margin NUMERIC(5,2) DEFAULT 0,
ADD COLUMN commission_rate NUMERIC(5,2) DEFAULT 0;
