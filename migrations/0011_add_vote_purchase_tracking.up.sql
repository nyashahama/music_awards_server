CREATE TABLE IF NOT EXISTS vote_purchases (
  purchase_id       UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id           UUID         NOT NULL,
  amount            INT          NOT NULL CHECK (amount > 0),
  price_paid        DECIMAL(10,2) NOT NULL CHECK (price_paid >= 0),
  payment_method    VARCHAR(50),
  payment_reference VARCHAR(255),
  status            VARCHAR(20)  NOT NULL DEFAULT 'pending', -- 'pending', 'completed', 'failed', 'refunded'
  created_at        TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at      TIMESTAMPTZ,
  
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_vote_purchases_user_id ON vote_purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_vote_purchases_status ON vote_purchases(status);
CREATE INDEX IF NOT EXISTS idx_vote_purchases_created_at ON vote_purchases(created_at);

-- Add check constraint for status
ALTER TABLE vote_purchases 
ADD CONSTRAINT vote_purchases_status_check 
  CHECK (status IN ('pending', 'completed', 'failed', 'refunded'));
