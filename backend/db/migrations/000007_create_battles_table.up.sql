CREATE TABLE battles (
    id SERIAL PRIMARY KEY,
    attacker_id uuid NOT NULL REFERENCES users(id),
    defender_id uuid NOT NULL REFERENCES users(id),
    winner_id uuid NOT NULL REFERENCES users(id),
    damage_percent INT NOT NULL,
    battle_log JSONB
);