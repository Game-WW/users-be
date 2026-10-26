CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY,
    username      TEXT,
    email         TEXT,
    phone_number  TEXT,
    password_hash TEXT,
    google_id     TEXT UNIQUE,
    apple_id      TEXT UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
)
