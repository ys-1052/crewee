-- Initial database setup for development environment
-- This file is executed automatically when the PostgreSQL container starts

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Set timezone to UTC
SET timezone = 'UTC';

-- Create basic database configuration
-- Additional setup will be handled by golang-migrate