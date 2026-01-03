-- Add timezone column to schools table
-- Supports multi-timezone for SaaS across different regions in Indonesia

-- Add timezone column with default WITA (Asia/Makassar)
ALTER TABLE schools 
ADD COLUMN IF NOT EXISTS timezone VARCHAR(50) DEFAULT 'Asia/Makassar';

-- Update existing schools to use WITA timezone
UPDATE schools SET timezone = 'Asia/Makassar' WHERE timezone IS NULL;

-- Add comment for documentation
COMMENT ON COLUMN schools.timezone IS 'School timezone: Asia/Jakarta (WIB), Asia/Makassar (WITA), Asia/Jayapura (WIT)';
