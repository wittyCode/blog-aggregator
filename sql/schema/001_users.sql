-- +goose Up
CREATE TABLE IF NOT EXISTS users(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  name TEXT NOT NULL,
  UNIQUE(name)
);

-- +goose Down
DROP TABLE users;
