-- Create applications table
CREATE TABLE applications (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  event_id UUID NOT NULL REFERENCES events (id) ON DELETE CASCADE,
  applicant_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  message TEXT,
  status application_status NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Create event_comments table
CREATE TABLE event_comments (
  id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
  event_id UUID NOT NULL REFERENCES events (id) ON DELETE CASCADE,
  author_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  parent_comment_id UUID REFERENCES event_comments (id) ON DELETE CASCADE,
  body TEXT NOT NULL,
  is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);


-- Create indexes for applications table
CREATE INDEX idx_applications_event_status ON applications (event_id, status);


CREATE INDEX idx_applications_user_status ON applications (applicant_user_id, status);


CREATE INDEX idx_applications_event_id ON applications (event_id);


CREATE INDEX idx_applications_applicant_user_id ON applications (applicant_user_id);


-- Create unique constraint to prevent duplicate applications (except canceled ones)
CREATE UNIQUE INDEX idx_applications_unique ON applications (event_id, applicant_user_id)
WHERE
  status != 'CANCELED';


-- Create indexes for event_comments table
CREATE INDEX idx_event_comments_event_id ON event_comments (event_id);


CREATE INDEX idx_event_comments_author_user_id ON event_comments (author_user_id);


CREATE INDEX idx_event_comments_parent_comment_id ON event_comments (parent_comment_id);


CREATE INDEX idx_event_comments_created_at ON event_comments (created_at);


-- Add updated_at triggers
CREATE TRIGGER trigger_applications_updated_at before
UPDATE ON applications FOR each ROW
EXECUTE function update_updated_at_column ();


CREATE TRIGGER trigger_event_comments_updated_at before
UPDATE ON event_comments FOR each ROW
EXECUTE function update_updated_at_column ();
