CREATE OR REPLACE VIEW village_production_stats AS
SELECT 
    vb.village_id,
    COALESCE(SUM(CASE WHEN rgc.name = 'Gold Mine' THEN rgc.production_per_min ELSE 0 END), 0) AS total_gold_rate,
    COALESCE(SUM(CASE WHEN rgc.name = 'Gold Mine' THEN rgc.capacity ELSE 0 END), 0) AS total_gold_cap,
    COALESCE(SUM(CASE WHEN rgc.name = 'Elixir Collector' THEN rgc.production_per_min ELSE 0 END), 0) AS total_elixir_rate,
    COALESCE(SUM(CASE WHEN rgc.name = 'Elixir Collector' THEN rgc.capacity ELSE 0 END), 0) AS total_elixir_cap,
    COALESCE(SUM(CASE WHEN rsc.resource_type = 'gold' OR rsc.resource_type = 'both' THEN rsc.storage_capacity ELSE 0 END), 0) AS max_gold_storage,
    COALESCE(SUM(CASE WHEN rsc.resource_type = 'elixir' OR rsc.resource_type = 'both' THEN rsc.storage_capacity ELSE 0 END), 0) AS max_elixir_storage
FROM village_buildings vb
LEFT JOIN resource_gen_configs rgc ON vb.building_name = rgc.name AND vb.level = rgc.level
LEFT JOIN resource_storage_configs rsc ON vb.building_name = rsc.name AND vb.level = rsc.level
GROUP BY vb.village_id;