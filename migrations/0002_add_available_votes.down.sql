-- Remove the added available_votes column
ALTER TABLE users DROP COLUMN available_votes;

-- Recreate the unique constraint to enforce one vote per category per user
ALTER TABLE votes ADD CONSTRAINT votes_user_id_category_id_key UNIQUE (user_id, category_id);
