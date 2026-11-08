CREATE SEQUENCE IF NOT EXISTS user_nickname_seq;

CREATE TABLE IF NOT EXISTS users (
    id                  UUID PRIMARY KEY,
    nickname            TEXT NOT NULL DEFAULT ('Guest' || nextval('user_nickname_seq')::TEXT),
    location            TEXT NOT NULL DEFAULT '',
    show_online_status  BOOLEAN NOT NULL DEFAULT TRUE,
    password_hash       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_nickname_lower_unique
    ON users (lower(nickname));
