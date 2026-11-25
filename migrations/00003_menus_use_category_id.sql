-- +goose Up
INSERT INTO categories (name)
SELECT DISTINCT btrim(category) FROM menus
WHERE category_id IS NULL
ON CONFLICT (name) DO NOTHING;

UPDATE menus m SET category_id = c.id
FROM categories c
WHERE m.category_id IS NULL AND c.name = btrim(m.category);

ALTER TABLE menus ALTER COLUMN category_id SET NOT NULL;
ALTER TABLE menus DROP COLUMN category;

-- +goose Down
ALTER TABLE menus ADD COLUMN category TEXT;

UPDATE menus m SET category = c.name
FROM categories c
WHERE c.id = m.category_id;

ALTER TABLE menus ALTER COLUMN category SET NOT NULL;
ALTER TABLE menus ALTER COLUMN category_id DROP NOT NULL;