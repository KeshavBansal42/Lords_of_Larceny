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

	buildingRows, err := db.Conn.Query(ctx, `
		SELECT id, name, level, hit_points, damage, build_cost, build_resource_type, production_per_min, capacity, size, min_thlevel, range 
		FROM building_configs
		ORDER BY MIN(id) OVER (PARTITION BY name) ASC, level ASC
	`)
	if err != nil {
		return nil, nil, errors.New("failed to fetch building configs")
	}
	defer buildingRows.Close()

	buildings, err := pgx.CollectRows(buildingRows, pgx.RowToStructByName[models.BuildingConfig])
	if err != nil {
		return nil, nil, errors.New("failed to parse building configs")
	}

	troopRows, err := db.Conn.Query(ctx, `
		SELECT id, name, level, hit_points, damage, min_thlevel, housing_space, range, speed 
		FROM troop_configs
		ORDER BY MIN(id) OVER (PARTITION BY name) ASC, level ASC
	`)
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
