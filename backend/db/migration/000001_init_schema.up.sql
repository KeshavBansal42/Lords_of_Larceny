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
    elixir INT NOT NULL DEFAULT 1000,
    last_collected_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE building_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    hit_points INT NOT NULL,
    damage INT NOT NULL DEFAULT 0,
    build_cost INT NOT NULL,
    build_resource_type VARCHAR(20) NOT NULL DEFAULT 'gold',
    production_per_min INT NOT NULL DEFAULT 0,
    capacity INT NOT NULL DEFAULT 0,
    size INT NOT NULL DEFAULT 2,
    min_thlevel INT NOT NULL DEFAULT 1,
    range INT NOT NULL DEFAULT 4
);

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

CREATE TABLE battles (
    id SERIAL PRIMARY KEY,
    attacker_id INT NOT NULL REFERENCES users(id),
    defender_id INT NOT NULL REFERENCES users(id),
    winner_id INT NOT NULL REFERENCES users(id),
    battle_log JSONB
);

CREATE INDEX idx_villages_user_id ON villages(user_id);
CREATE INDEX idx_village_buildings_village_id ON village_buildings(village_id);
CREATE INDEX idx_village_troops_village_id ON village_troops(village_id);