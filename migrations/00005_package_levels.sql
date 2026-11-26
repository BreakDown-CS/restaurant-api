-- +goose Up
-- Packages now include menus by level: a menu is in every package whose level
-- is at least the menu's min_package_level. NULL = not in any buffet (sold separately).
DROP TABLE package_menus;

ALTER TABLE buffet_packages
    ADD COLUMN level INTEGER NOT NULL DEFAULT 1 CHECK (level > 0);

ALTER TABLE menus
    ADD COLUMN min_package_level INTEGER CHECK (min_package_level > 0);

-- +goose Down
ALTER TABLE menus DROP COLUMN min_package_level;

ALTER TABLE buffet_packages DROP COLUMN level;

CREATE TABLE package_menus (
    package_id  UUID NOT NULL REFERENCES buffet_packages(id) ON DELETE CASCADE,
    menu_id     UUID NOT NULL REFERENCES menus(id),
    PRIMARY KEY (package_id, menu_id)
);

CREATE INDEX package_menus_menu_id_idx ON package_menus (menu_id);