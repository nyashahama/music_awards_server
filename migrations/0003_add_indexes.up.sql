-- indexes to optimize JOIN and WHERE clause performance
CREATE INDEX idx_votes_category_id ON votes (category_id);
CREATE INDEX idx_nominee_categories_category_id ON nominee_categories (category_id);
CREATE INDEX idx_categories_name ON categories (name);
