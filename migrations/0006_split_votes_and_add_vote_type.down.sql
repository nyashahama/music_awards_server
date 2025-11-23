-- Step 1: Add back available_votes column
ALTER TABLE users 
ADD COLUMN available_votes INT;

-- Step 2: Migrate data back (combine free and paid votes)
UPDATE users 
SET available_votes = COALESCE(free_votes, 0) + COALESCE(paid_votes, 0);

-- Step 3: Make available_votes NOT NULL with default
ALTER TABLE users 
ALTER COLUMN available_votes SET NOT NULL,
ALTER COLUMN available_votes SET DEFAULT 5;

-- Step 4: Drop new columns
ALTER TABLE users 
DROP COLUMN free_votes,
DROP COLUMN paid_votes;

-- Step 5: Drop vote_type from votes
ALTER TABLE votes DROP COLUMN vote_type;

-- Step 6: Drop indexes
DROP INDEX IF EXISTS idx_votes_user_id;
DROP INDEX IF EXISTS idx_votes_nominee_id;
DROP INDEX IF EXISTS idx_votes_created_at;
DROP INDEX IF EXISTS idx_votes_vote_type;
DROP INDEX IF EXISTS idx_votes_user_category_type;
