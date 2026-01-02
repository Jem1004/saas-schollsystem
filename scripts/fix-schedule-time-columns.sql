-- Fix attendance_schedules time columns
-- The columns were incorrectly created as timestamp with time zone by GORM auto-migrate
-- They should be TIME type for storing only time of day (HH:MM:SS)

-- First, backup existing data and alter columns
DO $$
BEGIN
    -- Check if columns exist and are wrong type
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'attendance_schedules' 
        AND column_name = 'start_time' 
        AND data_type = 'timestamp with time zone'
    ) THEN
        -- Alter start_time column
        ALTER TABLE attendance_schedules 
        ALTER COLUMN start_time TYPE TIME 
        USING start_time::time;
        
        RAISE NOTICE 'Fixed start_time column type to TIME';
    END IF;
    
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'attendance_schedules' 
        AND column_name = 'end_time' 
        AND data_type = 'timestamp with time zone'
    ) THEN
        -- Alter end_time column
        ALTER TABLE attendance_schedules 
        ALTER COLUMN end_time TYPE TIME 
        USING end_time::time;
        
        RAISE NOTICE 'Fixed end_time column type to TIME';
    END IF;
END $$;

-- Verify the fix
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'attendance_schedules' 
AND column_name IN ('start_time', 'end_time');
