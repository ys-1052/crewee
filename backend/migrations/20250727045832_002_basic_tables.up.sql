-- Basic tables creation
-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL,
  email_verified_at TIMESTAMP WITH TIME ZONE,
  home_region_code VARCHAR(10) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);


-- Sports master table
CREATE TABLE sports (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  code VARCHAR(20) NOT NULL UNIQUE,
  name VARCHAR(50) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);


-- Regions master table (JIS code compliant)
CREATE TABLE regions (
  region_code VARCHAR(10) PRIMARY KEY, -- JIS code
  name VARCHAR(100) NOT NULL,
  level region_level NOT NULL,
  parent_code VARCHAR(10),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
  FOREIGN key (parent_code) REFERENCES regions (region_code)
);


-- Create indexes
CREATE INDEX idx_users_email ON users (email);


CREATE INDEX idx_users_home_region ON users (home_region_code);


CREATE INDEX idx_sports_code ON sports (code);


CREATE INDEX idx_sports_active ON sports (is_active);


CREATE INDEX idx_regions_level ON regions (level);


CREATE INDEX idx_regions_parent ON regions (parent_code);


-- Foreign key constraints
ALTER TABLE users
ADD CONSTRAINT fk_users_home_region FOREIGN key (home_region_code) REFERENCES regions (region_code);


-- Function for auto-updating updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column () returns trigger AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


-- Setup triggers for updated_at
CREATE TRIGGER update_users_updated_at before
UPDATE ON users FOR each ROW
EXECUTE function update_updated_at_column ();


CREATE TRIGGER update_sports_updated_at before
UPDATE ON sports FOR each ROW
EXECUTE function update_updated_at_column ();


CREATE TRIGGER update_regions_updated_at before
UPDATE ON regions FOR each ROW
EXECUTE function update_updated_at_column ();
