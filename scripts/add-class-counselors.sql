-- Migration: Add class_counselors table for mapping BK teachers to classes
-- This allows multiple BK teachers to be assigned to different classes

CREATE TABLE IF NOT EXISTS class_counselors (
    id SERIAL PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    counselor_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(class_id, counselor_id)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_class_counselors_class_id ON class_counselors(class_id);
CREATE INDEX IF NOT EXISTS idx_class_counselors_counselor_id ON class_counselors(counselor_id);
CREATE INDEX IF NOT EXISTS idx_class_counselors_school_id ON class_counselors(school_id);

-- Add comment to table
COMMENT ON TABLE class_counselors IS 'Mapping table for BK teachers (counselors) to classes. Allows multiple BK teachers per class and multiple classes per BK teacher.';
