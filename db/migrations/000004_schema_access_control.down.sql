-- Revoke privileges on the schema and tables
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA bank FROM ccat;
REVOKE USAGE ON SCHEMA bank FROM ccat;

-- Drop the schema and the user
DROP SCHEMA IF EXISTS bank CASCADE;
DROP USER IF EXISTS ccat;

-- Reset the search path to the default (if necessary)
RESET search_path;
