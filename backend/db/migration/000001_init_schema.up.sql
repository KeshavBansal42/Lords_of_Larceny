CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE villages (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    town_hall_level INT NOT NULL DEFAULT 1,
    gold INT NOT NULL DEFAULT 1000,
    elixir INT NOT NULL DEFAULT 1000
);

CREATE TABLE building_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    damage INT NOT NULL DEFAULT 0,
    build_cost INT NOT NULL
);

CREATE TABLE troop_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    damage INT NOT NULL,
);

CREATE TABLE village_buildings (
    id SERIAL PRIMARY KEY,
    village_id INT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    building_id INT NOT NULL REFERENCES building_configs(id),
    x INT NOT NULL,
    y INT NOT NULL,
    UNIQUE (village_id, x, y)
);

CREATE TABLE village_troops (
    village_id INT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    troop_id INT NOT NULL REFERENCES troop_configs(id),
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (village_id, troop_id)
);