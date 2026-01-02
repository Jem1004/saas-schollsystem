-- Migration: Make students.class_id nullable for bulk import workflow
-- Requirements: 8.1 - THE Student model SHALL allow ClassID to be null
-- Requirements: 8.2 - WHEN ClassID is null, THE Student SHALL have IsActive set to false

-- Step 1: Drop the NOT NULL constraint on class_id column
ALTER TABLE students ALTER COLUMN class_id DROP NOT NULL;

-- Step 2: Update is_active default to false (students without class should be inactive)
ALTER TABLE students ALTER COLUMN is_active SET DEFAULT false;

-- Step 3: Update existing students to ensure data integrity
-- All existing students with class_id should remain active
-- This is a safety measure - existing data should already have class_id set
UPDATE students SET is_active = true WHERE class_id IS NOT NULL;

-- Step 4: Update any students without class_id to be inactive (if any exist)
UPDATE students SET is_active = false WHERE class_id IS NULL;

-- Log migration
DO $
BEGIN
    RAISE NOTICE 'Migration completed: students.class_id is now nullable';
    RAISE NOTICE 'Students with class_id remain active, students without class_id are inactive';
END $;
