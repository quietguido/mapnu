-- ❌ Drop indexes
DROP INDEX IF EXISTS bookings_event_id_idx;

DROP INDEX IF EXISTS event_date_index;

DROP INDEX IF EXISTS event_location_index;

DROP INDEX IF EXISTS event_name_index;

-- ❌ Drop partitions
DROP TABLE IF EXISTS event_2025_02_07;

DROP TABLE IF EXISTS event_default;

-- ❌ Drop partitioned table
DROP TABLE IF EXISTS event;

-- ❌ Drop bookings table
DROP TABLE IF EXISTS bookings;

-- ❌ Drop friendships table
DROP TABLE IF EXISTS friendships;

-- ❌ Drop users table
DROP TABLE IF EXISTS users;

-- ❌ Disable required extensions
DROP EXTENSION IF EXISTS postgis;

DROP EXTENSION IF EXISTS "uuid-ossp";
