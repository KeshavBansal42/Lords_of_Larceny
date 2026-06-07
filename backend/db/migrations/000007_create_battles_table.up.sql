CREATE TABLE battles (
    id SERIAL PRIMARY KEY,
    attacker_id INT NOT NULL REFERENCES users(id),
    defender_id INT NOT NULL REFERENCES users(id),
    winner_id INT NOT NULL REFERENCES users(id),
    battle_log JSONB
);