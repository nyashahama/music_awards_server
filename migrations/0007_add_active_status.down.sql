-- Drop indexes
DROP INDEX IF EXISTS idx_categories_active;
DROP INDEX IF EXISTS idx_nominees_active;
DROP INDEX IF EXISTS idx_nominees_name;

-- Drop columns
ALTER TABLE categories DROP COLUMN is_active;
ALTER TABLE nominees DROP COLUMN is_active;

-- Drop timestamps from nominee_categories
ALTER TABLE nominee_categories 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;
