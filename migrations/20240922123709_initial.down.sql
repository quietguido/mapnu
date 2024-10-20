-- Drop the index on the extracted event date
DROP INDEX IF EXISTS event_date_index;

-- Drop the function used to extract the date from POINTM
DROP FUNCTION IF EXISTS get_event_date(geom GEOMETRY);

-- Drop indexes on the event table
DROP INDEX IF EXISTS event_location_date_index;
DROP INDEX IF EXISTS event_name_index;

-- Drop the event table
DROP TABLE IF EXISTS event;

-- Drop the index on client email
DROP INDEX IF EXISTS client_email_index;

-- Drop the client table (formerly "user")
DROP TABLE IF EXISTS client;

-- Disable extensions
DROP EXTENSION IF EXISTS postgis;
DROP EXTENSION IF EXISTS "uuid-ossp";
