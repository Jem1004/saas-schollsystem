-- Initial database setup for School Management SaaS
-- This script runs automatically when PostgreSQL container starts for the first time

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schema for better organization (optional)
-- CREATE SCHEMA IF NOT EXISTS school;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE school_management TO school_admin;

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'School Management database initialized successfully!';
END $$;
