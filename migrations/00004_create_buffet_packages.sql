-- +goose Up
CREATE TABLE buffet_packages (
    id                      UUID PRIMARY KEY DEFAULT uuidv7(),
    name                    TEXT NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    price_per_person        NUMERIC(10, 2) NOT NULL CHECK (price_per_person >= 0),
    time_limit_minutes      INTEGER NOT NULL CHECK (time_limit_minutes > 0),
    free_below_height_cm    INTEGER NOT NULL CHECK (free_below_height_cm >= 0),
    overtime_fee            NUMERIC(10, 2) NOT NULL CHECK (overtime_fee >= 0),
    overtime_block_minutes  INTEGER NOT NULL CHECK (overtime_block_minutes > 0),
    overtime_grace_minutes  INTEGER NOT NULL CHECK (overtime_grace_minutes >= 0),
    waste_fee               NUMERIC(10, 2) NOT NULL CHECK (waste_fee >= 0),
    waste_unit_grams        INTEGER NOT NULL CHECK (waste_unit_grams > 0),
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE package_menus (
    package_id  UUID NOT NULL REFERENCES buffet_packages(id) ON DELETE CASCADE,
    menu_id     UUID NOT NULL REFERENCES menus(id),
    PRIMARY KEY (package_id, menu_id)
);

CREATE INDEX package_menus_menu_id_idx ON package_menus (menu_id);

-- +goose Down
DROP TABLE package_menus;
DROP TABLE buffet_packages;