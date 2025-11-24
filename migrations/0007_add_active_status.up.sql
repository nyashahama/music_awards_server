-- Step 1: Add is_active to categories
ALTER TABLE categories 
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Step 2: Add is_active to nominees
ALTER TABLE nominees 
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Step 3: Add indexes for active status queries
CREATE INDEX IF NOT EXISTS idx_categories_active ON categories(is_active);
CREATE INDEX IF NOT EXISTS idx_nominees_active ON nominees(is_active);

-- Step 4: Add index on nominees name for search
CREATE INDEX IF NOT EXISTS idx_nominees_name ON nominees(name);

-- Step 5: Add timestamps to nominee_categories if not exists
ALTER TABLE nominee_categories 
ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
