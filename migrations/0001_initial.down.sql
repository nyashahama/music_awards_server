-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS nominee_categories;
DROP TABLE IF EXISTS nominees;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

-- Remove UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";