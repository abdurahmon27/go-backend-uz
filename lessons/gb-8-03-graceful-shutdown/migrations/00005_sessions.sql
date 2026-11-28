-- +goose Up
CREATE TABLE sessions (
    token_hash bytea       PRIMARY KEY,
    user_id    bigint      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL
);
-- "Hamma qurilmadan chiqish" (DELETE … WHERE user_id = $1) jadvalni oxirigacha oʻqimasin.
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

-- +goose Down
DROP TABLE sessions;
