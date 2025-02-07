-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS postgis;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create friendships table
CREATE TABLE IF NOT EXISTS friendships (
    user1_id UUID REFERENCES users (id) ON DELETE CASCADE,
    user2_id UUID REFERENCES users (id) ON DELETE CASCADE,
    status VARCHAR(20) CHECK (status IN ('pending', 'accepted', 'blocked')),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (user1_id, user2_id)
);

-- ✅ Create partitioned event table (Fixing Primary Key issue)
CREATE TABLE IF NOT EXISTS event (
    event_id BIGSERIAL NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users (id) ON DELETE SET NULL,
    location GEOMETRY (POINT, 4326) NOT NULL,
    start_date TIMESTAMP
    WITH
        TIME ZONE NOT NULL,
        organizer VARCHAR(255) NOT NULL,
        upvote INTEGER NOT NULL DEFAULT 0,
        downvote INTEGER NOT NULL DEFAULT 0,
        created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (event_id, start_date) -- ✅ Include `start_date` in PK
)
PARTITION BY
    RANGE (start_date);

-- ✅ Create Default Partition
CREATE TABLE IF NOT EXISTS event_default PARTITION OF event DEFAULT;

-- ✅ Create a sample daily partition (partitioning by day)
CREATE TABLE event_2025_02_07 PARTITION OF event FOR
VALUES
FROM
    ('2025-02-07 00:00:00') TO ('2025-02-08 00:00:00');

-- ✅ Create indexes on event table
CREATE INDEX IF NOT EXISTS event_name_index ON event (name);

CREATE INDEX IF NOT EXISTS event_location_index ON event USING GIST (location);

CREATE INDEX IF NOT EXISTS event_date_index ON event (start_date);

-- ✅ Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    booking_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    event_id BIGINT NOT NULL, -- Store event_id manually since we can't have FK to partitioned table
    booking_status VARCHAR(20) CHECK (booking_status IN ('confirmed', 'cancelled')),
    visibility VARCHAR(10) CHECK (visibility IN ('public', 'private')),
    booked_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Add an index for faster event lookups in `bookings`
CREATE INDEX IF NOT EXISTS bookings_event_id_idx ON bookings (event_id);
