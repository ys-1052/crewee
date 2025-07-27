-- Drop basic tables (reverse order of creation)

-- Drop triggers
DROP TRIGGER IF EXISTS update_regions_updated_at ON regions;
DROP TRIGGER IF EXISTS update_sports_updated_at ON sports;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (considering foreign key constraints)
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sports;
DROP TABLE IF EXISTS regions;