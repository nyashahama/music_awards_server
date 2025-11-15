-- +migrate Up
-- Step 1: Add new columns as nullable first
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS first_name VARCHAR,
ADD COLUMN IF NOT EXISTS last_name VARCHAR,
ADD COLUMN IF NOT EXISTS location VARCHAR(255);

-- Step 2: Populate with default values for any NULL values
UPDATE users SET first_name = 'User' WHERE first_name IS NULL;
UPDATE users SET last_name = user_id::text WHERE last_name IS NULL;

-- Step 3: Make columns NOT NULL
ALTER TABLE users 
ALTER COLUMN first_name SET NOT NULL,
ALTER COLUMN last_name SET NOT NULL;

-- Step 4: Drop the username column (only if it exists)
ALTER TABLE users DROP COLUMN IF EXISTS username;
