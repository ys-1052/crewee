-- Drop teams table and related objects
DROP TRIGGER if EXISTS trigger_teams_updated_at ON teams;


DROP INDEX if EXISTS idx_team_regions_region_code;


DROP INDEX if EXISTS idx_team_regions_team_id;


DROP INDEX if EXISTS idx_teams_sport_id;


DROP INDEX if EXISTS idx_teams_owner_user_id;


DROP TABLE IF EXISTS team_regions;


DROP TABLE IF EXISTS teams;
