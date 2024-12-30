CREATE SCHEMA IF NOT EXISTS bank;

ALTER SCHEMA bank OWNER TO ccat;

GRANT USAGE ON SCHEMA bank TO ccat;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA bank TO ccat;

ALTER DEFAULT PRIVILEGES IN SCHEMA bank
GRANT ALL ON TABLES TO ccat;

-- Bank is first so it'll be the current schema for every session.
ALTER DATABASE webgo SET search_path TO bank, public;

