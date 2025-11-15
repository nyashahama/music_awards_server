-- +migrate Down
-- Step 1: Add back username column
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS username VARCHAR;

-- Step 2: Populate username from first and last name
UPDATE users SET username = first_name || ' ' || last_name;

-- Step 3: Make username NOT NULL
ALTER TABLE users 
ALTER COLUMN username SET NOT NULL;

-- Step 4: Drop the new columns
ALTER TABLE users 
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS location;
