package repository

import (
	"context"
	"errors"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/models"
	"github.com/jackc/pgx/v5"
)

func GetGameConfigs() ([]models.BuildingConfig, []models.TroopConfig, error) {
	ctx := context.Background()

	buildingQuery := `
		SELECT 
			b.name, 
			COALESCE(d.level, rg.level, rs.level, ac.level) AS level, 
			COALESCE(d.hit_points, rg.hit_points, rs.hit_points, ac.hit_points) AS hit_points, 
			COALESCE(d.damage, 0) AS damage, 
			COALESCE(d.build_cost, rg.build_cost, rs.build_cost, ac.build_cost) AS build_cost, 
			b.build_resource_type, 
			COALESCE(rg.production_per_min, 0) AS production_per_min, 
			COALESCE(rg.capacity, rs.storage_capacity, ac.total_housing_space, 0) AS capacity, 
			b.size, 
			COALESCE(d.min_thlevel, rg.min_thlevel, rs.min_thlevel, ac.min_thlevel) AS min_thlevel, 
			COALESCE(d.range, 0) AS range,
			COALESCE(d.single_target, TRUE) AS single_target,
			COALESCE(d.splash_radius, 0) AS splash_radius,
			COALESCE(d.target_type, 'none') AS target_type
		FROM building_configs b
		LEFT JOIN defense_configs d ON b.name = d.name
		LEFT JOIN resource_gen_configs rg ON b.name = rg.name
		LEFT JOIN resource_storage_configs rs ON b.name = rs.name
		LEFT JOIN army_camp_configs ac ON b.name = ac.name
		WHERE COALESCE(d.level, rg.level, rs.level, ac.level) IS NOT NULL
		ORDER BY b.name ASC, level ASC
	`
	buildingRows, err := db.Pool.Query(ctx, buildingQuery)
	if err != nil {
		return nil, nil, errors.New("failed to fetch building configs")
	}
	defer buildingRows.Close()

	buildings, err := pgx.CollectRows(buildingRows, pgx.RowToStructByName[models.BuildingConfig])
	if err != nil {
		return nil, nil, errors.New("failed to parse building configs")
	}

	troopQuery := `
		SELECT id, name, level, hit_points, damage, min_thlevel, housing_space, range, speed, airborne 
		FROM troop_configs
		ORDER BY MIN(id) OVER (PARTITION BY name) ASC, level ASC
	`
	troopRows, err := db.Pool.Query(ctx, troopQuery)
	if err != nil {
		return nil, nil, errors.New("failed to fetch troop configs")
	}
	defer troopRows.Close()

	troops, err := pgx.CollectRows(troopRows, pgx.RowToStructByName[models.TroopConfig])
	if err != nil {
		return nil, nil, errors.New("failed to parse troop configs")
	}

	return buildings, troops, nil
}
