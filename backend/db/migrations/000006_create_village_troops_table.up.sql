CREATE TABLE village_troops (
    village_id uuid NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    troop_id INT NOT NULL REFERENCES troop_configs(id),
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (village_id, troop_id)
);

CREATE INDEX idx_village_troops_village_id ON village_troops(village_id);