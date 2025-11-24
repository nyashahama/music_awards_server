-- Remove check constraint
ALTER TABLE votes DROP CONSTRAINT IF EXISTS votes_vote_type_check;
