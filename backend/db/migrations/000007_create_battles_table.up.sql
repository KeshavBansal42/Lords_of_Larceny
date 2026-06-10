CREATE TABLE battles (
    id SERIAL PRIMARY KEY,
    attacker_id uuid NOT NULL REFERENCES users(id),
    defender_id uuid NOT NULL REFERENCES users(id),
    winner_id uuid NOT NULL REFERENCES users(id),
    damage_percent INT NOT NULL,
    occurred_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    battle_log JSONB
);