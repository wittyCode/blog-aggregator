-- +goose Up
CREATE TABLE IF NOT EXISTS feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,
  CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_feed_id FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
