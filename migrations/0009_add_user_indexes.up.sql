-- Add index on email for faster login lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Add index on role for filtering by role
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Add composite index for active user searches
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
