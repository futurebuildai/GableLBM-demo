-- Add AI traceability fields to quotes table
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS original_file BYTEA;
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS original_filename TEXT;
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS original_content_type TEXT;
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS parse_map JSONB;
