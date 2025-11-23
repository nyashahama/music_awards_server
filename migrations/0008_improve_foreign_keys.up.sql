-- Step 1: Drop the existing composite foreign key in votes
ALTER TABLE votes 
DROP CONSTRAINT IF EXISTS votes_nominee_id_category_id_fkey;

-- Step 2: Add individual foreign keys for votes
ALTER TABLE votes 
DROP CONSTRAINT IF EXISTS votes_nominee_id_fkey,
ADD CONSTRAINT votes_nominee_id_fkey 
  FOREIGN KEY (nominee_id) 
  REFERENCES nominees(nominee_id) 
  ON DELETE CASCADE;

ALTER TABLE votes 
DROP CONSTRAINT IF EXISTS votes_category_id_fkey,
ADD CONSTRAINT votes_category_id_fkey 
  FOREIGN KEY (category_id) 
  REFERENCES categories(category_id) 
  ON DELETE CASCADE;

-- Step 3: Add check constraint to ensure vote_type is valid
ALTER TABLE votes 
ADD CONSTRAINT votes_vote_type_check 
  CHECK (vote_type IN ('free', 'paid'));

-- Step 4: Improve nominee_categories constraints
ALTER TABLE nominee_categories 
DROP CONSTRAINT IF EXISTS nominee_categories_nominee_id_fkey,
ADD CONSTRAINT nominee_categories_nominee_id_fkey 
  FOREIGN KEY (nominee_id) 
  REFERENCES nominees(nominee_id) 
  ON UPDATE CASCADE 
  ON DELETE CASCADE;

ALTER TABLE nominee_categories 
DROP CONSTRAINT IF EXISTS nominee_categories_category_id_fkey,
ADD CONSTRAINT nominee_categories_category_id_fkey 
  FOREIGN KEY (category_id) 
  REFERENCES categories(category_id) 
  ON UPDATE CASCADE 
  ON DELETE CASCADE;
