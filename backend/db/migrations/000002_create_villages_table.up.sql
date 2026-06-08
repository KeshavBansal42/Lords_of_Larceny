CREATE TABLE villages (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    town_hall_level INT NOT NULL DEFAULT 1,
    gold INT NOT NULL DEFAULT 1000,
    elixir INT NOT NULL DEFAULT 1000,
    last_collected_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_attacked_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_villages_user_id ON villages(user_id);