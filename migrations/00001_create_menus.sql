-- +goose Up
CREATE TABLE menus (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    name            TEXT NOT NULL CHECK (btrim(name) <> ''),
    price           NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    category        TEXT NOT NULL CHECK (btrim(category) <> ''),
    is_available    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE menus;
