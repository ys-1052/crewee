-- ENUM type definitions
-- Event status (OPEN -> CLOSED -> CANCELED transition order)
CREATE TYPE event_status AS ENUM ('OPEN', 'CLOSED', 'CANCELED');

-- Event visibility (public/unlisted)
CREATE TYPE event_visibility AS ENUM ('PUBLIC', 'UNLISTED');

-- Skill level (beginner/intermediate/advanced)
CREATE TYPE skill_level AS ENUM ('BEGINNER', 'INTERMEDIATE', 'ADVANCED');

-- Region level (prefecture/municipality)
CREATE TYPE region_level AS ENUM ('PREFECTURE', 'MUNICIPALITY');

-- Application status (approval workflow)
CREATE TYPE application_status AS ENUM ('PENDING', 'ACCEPTED', 'DECLINED', 'CANCELED', 'WAITLISTED');