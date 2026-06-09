CREATE TABLE resource_gen_configs (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    build_cost INT NOT NULL,
    min_thlevel INT NOT NULL,
    production_per_min INT NOT NULL,
    capacity INT NOT NULL,
    resource_type VARCHAR(20) NOT NULL,
    build_time_seconds INT NOT NULL DEFAULT 60,
    PRIMARY KEY (name, level)
);