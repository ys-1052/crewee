-- Create events table
CREATE TABLE events (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  team_id UUID NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
  created_by UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  sport_id UUID NOT NULL REFERENCES sports (id) ON DELETE RESTRICT,
  region_code VARCHAR(6) NOT NULL REFERENCES regions (region_code) ON DELETE RESTRICT,
  title VARCHAR(100) NOT NULL,
  start_at TIMESTAMPTZ NOT NULL,
  duration_min INTEGER NOT NULL CHECK (duration_min > 0),
  level skill_level NOT NULL,
  capacity INTEGER NOT NULL CHECK (capacity > 0),
  fee TEXT,
  note TEXT,
  status event_status NOT NULL DEFAULT 'OPEN',
  visibility event_visibility NOT NULL DEFAULT 'PUBLIC',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Create search indexes for performance optimization
CREATE INDEX idx_events_search ON events (status, start_at, sport_id, region_code);


CREATE INDEX idx_events_team_status ON events (team_id, status);


CREATE INDEX idx_events_created_by ON events (created_by);


CREATE INDEX idx_events_start_at ON events (start_at);


CREATE INDEX idx_events_visibility_status ON events (visibility, status);


-- Add updated_at trigger for events table
CREATE TRIGGER trigger_events_updated_at before
UPDATE ON events FOR each ROW
EXECUTE function update_updated_at_column ();
