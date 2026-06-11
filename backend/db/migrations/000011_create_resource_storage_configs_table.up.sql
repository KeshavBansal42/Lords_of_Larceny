CREATE TABLE resource_storage_configs (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    resource_type VARCHAR(20) NOT NULL,
    storage_capacity INT NOT NULL,
    PRIMARY KEY (name, level)
);