CREATE TABLE village_buildings (
    id SERIAL PRIMARY KEY,
    village_id uuid NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    building_name VARCHAR(50) NOT NULL REFERENCES building_configs(name),
    level INT NOT NULL DEFAULT 1,
    x INT NOT NULL,
    y INT NOT NULL,
    UNIQUE (village_id, x, y)
);

CREATE INDEX idx_village_buildings_village_id ON village_buildings(village_id);