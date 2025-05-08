-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tables in dependency order
CREATE TABLE users (
  user_id        UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  -- first_name     VARCHAR(255) NOT NULL,
  -- last_name      VARCHAR(255) NOT NULL,
  username       VARCHAR(255) NOT NULL UNIQUE,
  password_hash  VARCHAR(255) NOT NULL,
  email          VARCHAR(255) NOT NULL UNIQUE,
  role           VARCHAR(50)  NOT NULL,
  created_at     TIMESTAMP   NOT NULL DEFAULT now(),
  updated_at     TIMESTAMP   NOT NULL DEFAULT now()
);

--SUDO -I -U POSTGRES
--PSQL -U POSTGRES -D MUSIX

CREATE TABLE categories (
  category_id   UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          VARCHAR(255) NOT NULL UNIQUE,
  description   TEXT,
  created_at    TIMESTAMP   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE TABLE nominees (
  nominee_id    UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          VARCHAR(255) NOT NULL,
  description   TEXT,
  sample_works  JSONB,
  image_url     VARCHAR(500),
  created_at    TIMESTAMP   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE TABLE nominee_categories (
  nominee_id    UUID NOT NULL,
  category_id   UUID NOT NULL,
  PRIMARY KEY (nominee_id, category_id),
  FOREIGN KEY (nominee_id) REFERENCES nominees(nominee_id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);

CREATE TABLE votes (
  vote_id       UUID       PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id       UUID       NOT NULL,
  nominee_id    UUID       NOT NULL,
  category_id   UUID       NOT NULL,
  created_at    TIMESTAMP  NOT NULL DEFAULT now(),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY (nominee_id, category_id) REFERENCES nominee_categories(nominee_id, category_id) ON DELETE CASCADE,
  UNIQUE (user_id, category_id)
);