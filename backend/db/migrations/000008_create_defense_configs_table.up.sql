CREATE TABLE defense_configs (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    build_cost INT NOT NULL,
    min_thlevel INT NOT NULL,
    damage INT NOT NULL,
    range INT NOT NULL,
    single_target BOOLEAN NOT NULL DEFAULT TRUE,
    splash_radius FLOAT NOT NULL DEFAULT 0,
    target_type VARCHAR(20) NOT NULL,
    PRIMARY KEY (name, level)
);