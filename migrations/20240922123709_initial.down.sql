-- Drop the index on the extracted event date
DROP INDEX IF EXISTS event_date_index;

-- Drop the function used to extract the date from POINTM
DROP FUNCTION IF EXISTS get_event_date(geom GEOMETRY);

-- Drop indexes on the event table
DROP INDEX IF EXISTS event_location_date_index;
DROP INDEX IF EXISTS event_name_index;

-- Drop the event table
DROP TABLE IF EXISTS event;

-- Drop the index on user email
DROP INDEX IF EXISTS user_email_index;

-- Drop the user table
DROP TABLE IF EXISTS "user";

-- Disable extensions
DROP EXTENSION IF EXISTS postgis;
DROP EXTENSION IF EXISTS "uuid-ossp";
