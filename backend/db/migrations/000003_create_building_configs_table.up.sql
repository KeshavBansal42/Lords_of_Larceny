CREATE TABLE building_configs (
    name VARCHAR(50) NOT NULL PRIMARY KEY,
    build_resource_type VARCHAR(20) NOT NULL DEFAULT 'gold',
    size INT NOT NULL DEFAULT 2
);