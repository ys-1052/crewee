-- Create teams table
CREATE TABLE teams (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  owner_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  sport_id UUID NOT NULL REFERENCES sports (id) ON DELETE RESTRICT,
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Create team_regions table (many-to-many relationship between teams and regions)
CREATE TABLE team_regions (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  team_id UUID NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
  region_code VARCHAR(6) NOT NULL REFERENCES regions (region_code) ON DELETE CASCADE,
  note TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (team_id, region_code)
);


-- Create indexes for performance
CREATE INDEX idx_teams_owner_user_id ON teams (owner_user_id);


CREATE INDEX idx_teams_sport_id ON teams (sport_id);


CREATE INDEX idx_team_regions_team_id ON team_regions (team_id);


CREATE INDEX idx_team_regions_region_code ON team_regions (region_code);


-- Add updated_at trigger for teams table
CREATE TRIGGER trigger_teams_updated_at before
UPDATE ON teams FOR each ROW
EXECUTE function update_updated_at_column ();
