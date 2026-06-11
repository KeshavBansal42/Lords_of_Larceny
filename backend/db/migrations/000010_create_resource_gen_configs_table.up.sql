CREATE TABLE resource_gen_configs (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    production_per_min INT NOT NULL,
    capacity INT NOT NULL,
    resource_type VARCHAR(20) NOT NULL,
    PRIMARY KEY (name, level)
);