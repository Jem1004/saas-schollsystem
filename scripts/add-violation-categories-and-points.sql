-- Migration: Add violation categories table and point field to violations
-- This migration adds:
-- 1. violation_categories table for customizable categories per school
-- 2. point field to violations table for tracking penalty points
-- 3. category_id foreign key to link violations to categories

-- Create violation_categories table
CREATE TABLE IF NOT EXISTS violation_categories (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    default_point INTEGER NOT NULL DEFAULT -5,
    default_level VARCHAR(20) NOT NULL DEFAULT 'ringan',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_default_point CHECK (default_point <= 0),
    CONSTRAINT chk_default_level CHECK (default_level IN ('ringan', 'sedang', 'berat'))
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_violation_categories_school_id ON violation_categories(school_id);
CREATE INDEX IF NOT EXISTS idx_violation_categories_is_active ON violation_categories(school_id, is_active);

-- Add point and category_id columns to violations table
ALTER TABLE violations 
ADD COLUMN IF NOT EXISTS point INTEGER NOT NULL DEFAULT -5,
ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES violation_categories(id) ON DELETE SET NULL;

-- Create index for category_id
CREATE INDEX IF NOT EXISTS idx_violations_category_id ON violations(category_id);

-- Update existing violations with default points based on level
UPDATE violations SET point = CASE 
    WHEN level = 'ringan' THEN -5
    WHEN level = 'sedang' THEN -15
    WHEN level = 'berat' THEN -30
    ELSE -5
END WHERE point = -5 OR point IS NULL;

-- Insert default categories for existing schools
INSERT INTO violation_categories (school_id, name, default_point, default_level, description, is_active)
SELECT s.id, cat.name, cat.default_point, cat.default_level, cat.description, true
FROM schools s
CROSS JOIN (VALUES 
    ('Keterlambatan', -5, 'ringan', 'Terlambat masuk sekolah'),
    ('Bolos', -15, 'sedang', 'Tidak masuk tanpa keterangan'),
    ('Seragam', -5, 'ringan', 'Tidak memakai seragam sesuai aturan'),
    ('Perilaku', -10, 'sedang', 'Perilaku tidak sopan'),
    ('Kekerasan', -30, 'berat', 'Melakukan kekerasan fisik'),
    ('Bullying', -25, 'berat', 'Melakukan perundungan'),
    ('Merokok', -20, 'berat', 'Merokok di lingkungan sekolah'),
    ('Narkoba', -50, 'berat', 'Terlibat narkoba'),
    ('Pencurian', -30, 'berat', 'Melakukan pencurian'),
    ('Vandalisme', -20, 'sedang', 'Merusak fasilitas sekolah'),
    ('Lainnya', -5, 'ringan', 'Pelanggaran lainnya')
) AS cat(name, default_point, default_level, description)
WHERE NOT EXISTS (
    SELECT 1 FROM violation_categories vc 
    WHERE vc.school_id = s.id AND vc.name = cat.name
);

-- Link existing violations to their categories where possible
UPDATE violations v
SET category_id = vc.id
FROM violation_categories vc
JOIN students st ON st.id = v.student_id
WHERE vc.school_id = st.school_id 
AND LOWER(vc.name) = LOWER(v.category)
AND v.category_id IS NULL;

-- Add comment for documentation
COMMENT ON TABLE violation_categories IS 'Customizable violation categories per school with default points';
COMMENT ON COLUMN violations.point IS 'Penalty points for this violation (negative value)';
COMMENT ON COLUMN violations.category_id IS 'Reference to violation_categories for standardized categories';
