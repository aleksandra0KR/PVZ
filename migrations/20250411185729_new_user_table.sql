-- +goose Up
-- +goose StatementBegin
CREATE TYPE role AS ENUM ('employee', 'moderator');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role role NOT NULL
);
CREATE INDEX idx_users_email_password ON users (email, password);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
