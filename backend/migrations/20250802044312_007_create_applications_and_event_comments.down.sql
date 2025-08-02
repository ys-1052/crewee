-- Drop applications and event_comments tables and related objects
DROP TRIGGER if EXISTS trigger_event_comments_updated_at ON event_comments;


DROP TRIGGER if EXISTS trigger_applications_updated_at ON applications;


DROP INDEX if EXISTS idx_event_comments_created_at;


DROP INDEX if EXISTS idx_event_comments_parent_comment_id;


DROP INDEX if EXISTS idx_event_comments_author_user_id;


DROP INDEX if EXISTS idx_event_comments_event_id;


DROP INDEX if EXISTS idx_applications_unique;


DROP INDEX if EXISTS idx_applications_applicant_user_id;


DROP INDEX if EXISTS idx_applications_event_id;


DROP INDEX if EXISTS idx_applications_user_status;


DROP INDEX if EXISTS idx_applications_event_status;


DROP TABLE IF EXISTS event_comments;


DROP TABLE IF EXISTS applications;
