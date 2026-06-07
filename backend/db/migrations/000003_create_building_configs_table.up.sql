CREATE TABLE building_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    damage INT NOT NULL DEFAULT 0,
    build_cost INT NOT NULL,
    build_resource_type VARCHAR(20) NOT NULL DEFAULT 'gold',
    production_per_min INT NOT NULL DEFAULT 0,
    capacity INT NOT NULL DEFAULT 0,
    size INT NOT NULL DEFAULT 2,
    min_thlevel INT NOT NULL DEFAULT 1,
    range INT NOT NULL DEFAULT 4
);