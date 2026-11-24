-- +goose Up
CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    name            TEXT NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Copy every distinct category text that menus already use.
-- btrim keeps " Rice " and "Rice" as one category.
INSERT INTO categories (name)
SELECT DISTINCT btrim(category) FROM menus
ON CONFLICT (name) DO NOTHING;

-- Nullable for now: the API still writes the text column until B2-4.
ALTER TABLE menus ADD COLUMN category_id UUID REFERENCES categories(id);

UPDATE menus m
SET category_id = c.id
FROM categories c
WHERE c.name = btrim(m.category);

-- +goose Down
ALTER TABLE menus DROP COLUMN category_id;
DROP TABLE categories;