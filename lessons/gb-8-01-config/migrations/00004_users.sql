-- +goose Up
-- Foydalanuvchilar (GB-6.03): parolning oʻzi emas, argon2id satri saqlanadi.
CREATE TABLE users (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email         text   NOT NULL CHECK (email <> ''),
    password_hash text   NOT NULL
);
-- Email katta-kichik harfga qaramay bitta (GB-6.03: Aziza@… va aziza@… — bitta odam).
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email));

-- +goose Down
DROP TABLE users;
