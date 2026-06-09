CREATE TABLE troop_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    min_thlevel INT NOT NULL DEFAULT 1,
    housing_space INT NOT NULL DEFAULT 1,
    damage INT NOT NULL,
    range INT NOT NULL DEFAULT 1,
    speed INT NOT NULL DEFAULT 1,
    airborne BOOLEAN NOT NULL DEFAULT FALSE
);