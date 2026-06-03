DROP TABLE IF EXISTS village_troops;
DROP TABLE IF EXISTS village_buildings;
DROP TABLE IF EXISTS troop_configs;
DROP TABLE IF EXISTS building_configs;
DROP TABLE IF EXISTS villages;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_villages_user_id;
DROP INDEX IF EXISTS idx_village_buildings_village_id;
DROP INDEX IF EXISTS idx_village_troops_village_id;