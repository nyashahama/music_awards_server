-- Drop indexes on rollback
DROP INDEX IF EXISTS idx_votes_category_id;
DROP INDEX IF EXISTS idx_nominee_categories_category_id;
DROP INDEX IF EXISTS idx_categories_name;
