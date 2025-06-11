ALTER TABLE users
ADD COLUMN available_votes INT NOT NULL DEFAULT 5;

-- Remove unique constraint to allow multiple votes per category
ALTER TABLE votes DROP CONSTRAINT IF EXISTS votes_user_id_category_id_key;
