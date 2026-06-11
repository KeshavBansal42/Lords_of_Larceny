CREATE TABLE army_camp_configs (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    total_housing_space INT NOT NULL,
    PRIMARY KEY (name, level)
);