-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

-- Create user table
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user email
CREATE INDEX IF NOT EXISTS user_email_index ON "user"(email);

-- Create event table
CREATE TABLE IF NOT EXISTS event (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES "user"(id),
    location_date GEOMETRY(POINTM, 4326) NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    upvote INTEGER NOT NULL DEFAULT 0,
    downvote INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on event table
CREATE INDEX IF NOT EXISTS event_name_index ON event(name);
CREATE INDEX IF NOT EXISTS event_location_date_index
ON event
USING GIST(location_date gist_geometry_ops_nd);

-- Create an IMMUTABLE version of the function to extract date from POINTM
CREATE OR REPLACE FUNCTION get_event_date(geom GEOMETRY)
RETURNS TIMESTAMP WITH TIME ZONE IMMUTABLE AS $$
BEGIN
    RETURN to_timestamp(ST_M(geom)) AT TIME ZONE 'UTC';
END;
$$ LANGUAGE plpgsql;

-- Create an index on the extracted date
CREATE INDEX IF NOT EXISTS event_date_index
ON event (get_event_date(location_date));
