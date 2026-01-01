-- Migration: Add attendance_schedules and display_tokens tables
-- Requirements: 3.1, 3.2, 3.3, 5.1, 6.1, 3.5

-- ============================================
-- Task 1.1: Create attendance_schedules table
-- ============================================
-- Requirements: 3.1, 3.2, 3.3 - Multi-schedule support for different activities

CREATE TABLE IF NOT EXISTS attendance_schedules (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    late_threshold INT NOT NULL DEFAULT 15,
    very_late_threshold INT,
    days_of_week VARCHAR(20) DEFAULT '1,2,3,4,5',
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add index on school_id for efficient queries
CREATE INDEX IF NOT EXISTS idx_attendance_schedules_school ON attendance_schedules(school_id);

-- Add index for finding active schedules
CREATE INDEX IF NOT EXISTS idx_attendance_schedules_active ON attendance_schedules(school_id, is_active);

COMMENT ON TABLE attendance_schedules IS 'Configurable attendance time slots for different activities (morning entry, dismissal, prayer times, etc.)';
COMMENT ON COLUMN attendance_schedules.name IS 'Schedule name, e.g., "Masuk Pagi", "Pulang", "Sholat Dzuhur"';
COMMENT ON COLUMN attendance_schedules.start_time IS 'Schedule start time in HH:MM format';
COMMENT ON COLUMN attendance_schedules.end_time IS 'Schedule end time in HH:MM format';
COMMENT ON COLUMN attendance_schedules.late_threshold IS 'Minutes after start_time to be considered late';
COMMENT ON COLUMN attendance_schedules.very_late_threshold IS 'Minutes after start_time to be considered very late (optional)';
COMMENT ON COLUMN attendance_schedules.days_of_week IS 'Comma-separated day numbers (1=Monday to 7=Sunday or 0=Sunday to 6=Saturday)';
COMMENT ON COLUMN attendance_schedules.is_default IS 'Default schedule for the school when no specific schedule matches';

-- ============================================
-- Task 1.2: Create display_tokens table
-- ============================================
-- Requirements: 5.1, 6.1 - Token-based access for public display screens

CREATE TABLE IF NOT EXISTS display_tokens (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    token VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    last_accessed_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add index on school_id for efficient queries
CREATE INDEX IF NOT EXISTS idx_display_tokens_school ON display_tokens(school_id);

-- Add unique index on token for fast lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_display_tokens_token ON display_tokens(token);

COMMENT ON TABLE display_tokens IS 'Tokens for public display access without authentication';
COMMENT ON COLUMN display_tokens.token IS 'Cryptographically secure 64-character hex token';
COMMENT ON COLUMN display_tokens.name IS 'Display location name, e.g., "Display Pintu Utama"';
COMMENT ON COLUMN display_tokens.last_accessed_at IS 'Timestamp of last successful access';
COMMENT ON COLUMN display_tokens.expires_at IS 'Optional expiration date for the token';

-- ============================================
-- Task 1.3: Add schedule_id to attendances table
-- ============================================
-- Requirements: 3.5 - Associate attendance records with schedules

-- Add nullable schedule_id column
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS schedule_id INT REFERENCES attendance_schedules(id);

-- Add index on schedule_id for efficient queries
CREATE INDEX IF NOT EXISTS idx_attendances_schedule ON attendances(schedule_id);

COMMENT ON COLUMN attendances.schedule_id IS 'Reference to the attendance schedule this record belongs to';

-- Log migration completion
DO $$
BEGIN
    RAISE NOTICE 'Migration completed: attendance_schedules and display_tokens tables created, schedule_id added to attendances';
END $$;
