CREATE TABLE army_camp_configs (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    build_cost INT NOT NULL,
    min_thlevel INT NOT NULL,
    total_housing_space INT NOT NULL,
    PRIMARY KEY (name, level)
);