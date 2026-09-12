CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(100)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    group_id INTEGER REFERENCES groups(id),
    subgroup SMALLINT NOT NULL DEFAULT 0,
        CHECK (subgroup IN (0, 1, 2)),
    notifications_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);