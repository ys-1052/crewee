-- Drop events table and related objects
DROP TRIGGER if EXISTS trigger_events_updated_at ON events;


DROP INDEX if EXISTS idx_events_visibility_status;


DROP INDEX if EXISTS idx_events_start_at;


DROP INDEX if EXISTS idx_events_created_by;


DROP INDEX if EXISTS idx_events_team_status;


DROP INDEX if EXISTS idx_events_search;


DROP TABLE IF EXISTS events;
