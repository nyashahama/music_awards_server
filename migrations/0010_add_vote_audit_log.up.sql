CREATE TABLE IF NOT EXISTS vote_audit_log (
  audit_id      UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  vote_id       UUID        NOT NULL,
  user_id       UUID        NOT NULL,
  action        VARCHAR(20) NOT NULL, -- 'created', 'updated', 'deleted'
  old_nominee_id UUID,
  new_nominee_id UUID,
  category_id   UUID        NOT NULL,
  vote_type     VARCHAR(10) NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_vote_audit_user_id ON vote_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_vote_audit_vote_id ON vote_audit_log(vote_id);
CREATE INDEX IF NOT EXISTS idx_vote_audit_created_at ON vote_audit_log(created_at);

-- Add check constraint for action types
ALTER TABLE vote_audit_log 
ADD CONSTRAINT vote_audit_action_check 
  CHECK (action IN ('created', 'updated', 'deleted'));
