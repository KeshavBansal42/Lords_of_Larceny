CREATE TABLE troop_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    damage INT NOT NULL,
    min_thlevel INT NOT NULL DEFAULT 1,
    housing_space INT NOT NULL DEFAULT 1,
    range INT NOT NULL DEFAULT 1,
    speed INT NOT NULL DEFAULT 1
);