CREATE OR REPLACE VIEW village_production_stats AS
SELECT 
    vb.village_id,
    COALESCE(SUM(CASE WHEN bc.name = 'Gold Mine' THEN bc.production_per_min ELSE 0 END), 0) AS total_gold_rate,
    COALESCE(SUM(CASE WHEN bc.name = 'Gold Mine' THEN bc.capacity ELSE 0 END), 0) AS total_gold_cap,
    COALESCE(SUM(CASE WHEN bc.name = 'Elixir Collector' THEN bc.production_per_min ELSE 0 END), 0) AS total_elixir_rate,
    COALESCE(SUM(CASE WHEN bc.name = 'Elixir Collector' THEN bc.capacity ELSE 0 END), 0) AS total_elixir_cap
FROM village_buildings vb
JOIN building_configs bc ON vb.building_id = bc.id
GROUP BY vb.village_id;