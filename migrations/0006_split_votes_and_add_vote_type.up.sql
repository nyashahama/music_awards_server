
-- Step 1: Add new vote columns to users table
ALTER TABLE users 
ADD COLUMN free_votes INT,
ADD COLUMN paid_votes INT;

-- Step 2: Migrate existing available_votes data
-- Assume all existing votes are "free" votes, max 3 free votes
UPDATE users 
SET free_votes = LEAST(available_votes, 3),
    paid_votes = GREATEST(available_votes - 3, 0);

-- Step 3: Make new columns NOT NULL with defaults
ALTER TABLE users 
ALTER COLUMN free_votes SET NOT NULL,
ALTER COLUMN free_votes SET DEFAULT 3,
ALTER COLUMN paid_votes SET NOT NULL,
ALTER COLUMN paid_votes SET DEFAULT 0;

-- Step 4: Drop old column
ALTER TABLE users DROP COLUMN available_votes;

-- Step 5: Add vote_type to votes table
ALTER TABLE votes 
ADD COLUMN vote_type VARCHAR(10) NOT NULL DEFAULT 'free';

-- Step 6: Add updated_at to votes if not exists
ALTER TABLE votes 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Step 7: Add performance indexes for votes
CREATE INDEX IF NOT EXISTS idx_votes_user_id ON votes(user_id);
CREATE INDEX IF NOT EXISTS idx_votes_nominee_id ON votes(nominee_id);
CREATE INDEX IF NOT EXISTS idx_votes_created_at ON votes(created_at);
CREATE INDEX IF NOT EXISTS idx_votes_vote_type ON votes(vote_type);

-- Step 8: Composite index for user-category vote checks
CREATE INDEX IF NOT EXISTS idx_votes_user_category_type ON votes(user_id, category_id, vote_type);

