package repository

import (
	"context"
	"errors"
	"time"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/jackc/pgx/v5"
)

func GetAllVillageBuildings(villageID string) ([]dtos.BuildingResponseFromDBDTO, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT building_name, level, x, y, status, upgrade_complete_at FROM village_buildings WHERE village_id = $1", villageID)

	if err != nil {
		return nil, errors.New("Error fetching buildings")
	}
	defer rows.Close()

	buildings, err := pgx.CollectRows(rows, pgx.RowToStructByName[dtos.BuildingResponseFromDBDTO])
	if err != nil {
		return nil, errors.New("Error parsing buildings")
	}

	return buildings, nil
}

func AddBuilding(userID string, buildingName string, x, y int) (int, int, error) {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	if buildingName == "Town Hall" {
		return 0, 0, errors.New("cannot build another town hall")
	}

	var villageID string
	var thlevel, gold, elixir int
	err = tx.QueryRow(ctx, "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &thlevel, &gold, &elixir)
	if err != nil {
		return 0, 0, errors.New("error fetching village details")
	}

	var size, buildCost, minThLevel, buildTimeSeconds int
	var resourceType string

	configQuery := `
		SELECT b.size, b.build_resource_type, c.build_cost, c.min_thlevel, c.build_time_seconds
		FROM building_configs b
		JOIN (
			SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM defense_configs
			UNION ALL SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM resource_gen_configs
			UNION ALL SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM resource_storage_configs
			UNION ALL SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM army_camp_configs
		) c ON b.name = c.name
		WHERE b.name = $1 AND c.level = 1
	`
	err = tx.QueryRow(ctx, configQuery, buildingName).Scan(&size, &resourceType, &buildCost, &minThLevel, &buildTimeSeconds)
	if err != nil {
		return gold, elixir, errors.New("building configuration not found or cannot be built at level 1")
	}

	if thlevel < minThLevel {
		return gold, elixir, errors.New("minimum Town Hall Level requirement not met")
	}

	if x < 0 || y < 0 || x+size > 36 || y+size > 36 {
		return gold, elixir, errors.New("building is out of bounds")
	}

	var collisionCount int
	overlapQuery := `
        SELECT COUNT(*) 
        FROM village_buildings vb
        JOIN building_configs bc ON vb.building_name = bc.name
        WHERE vb.village_id = $1
          AND $2 < (vb.x + bc.size)
          AND ($2 + $3) > vb.x
          AND $4 < (vb.y + bc.size)
          AND ($4 + $3) > vb.y
    `
	err = tx.QueryRow(ctx, overlapQuery, villageID, x, size, y).Scan(&collisionCount)
	if err != nil {
		return gold, elixir, errors.New("error checking grid")
	}
	if collisionCount > 0 {
		return gold, elixir, errors.New("cannot place building on an existing building")
	}

	if resourceType == "gold" {
		if gold < buildCost {
			return gold, elixir, errors.New("insufficient gold")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold - $1 WHERE id = $2", buildCost, villageID)
		gold -= buildCost
	} else if resourceType == "elixir" {
		if elixir < buildCost {
			return gold, elixir, errors.New("insufficient elixir")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET elixir = elixir - $1 WHERE id = $2", buildCost, villageID)
		elixir -= buildCost
	}

	insertQuery := `
		INSERT INTO village_buildings (village_id, building_name, level, x, y, status, upgrade_complete_at) 
		VALUES ($1, $2, 0, $3, $4, 'upgrading', NOW() + INTERVAL '1 second' * $5)
	`
	_, err = tx.Exec(ctx, insertQuery, villageID, buildingName, x, y, buildTimeSeconds)
	if err != nil {
		return gold, elixir, errors.New("error adding building")
	}

	if err = tx.Commit(ctx); err != nil {
		return gold, elixir, err
	}

	return gold, elixir, nil
}

func UpgradeBuilding(userId string, x, y int) (int, int, error) {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var villageID string
	var thlevel, gold, elixir int

	err = tx.QueryRow(ctx, "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1 FOR UPDATE", userId).Scan(&villageID, &thlevel, &gold, &elixir)

	if err != nil {
		return 0, 0, errors.New("Error fetching village details.")
	}

	var buildingName string
	var currentLevel int
	err = tx.QueryRow(ctx, "SELECT building_name, level FROM village_buildings WHERE village_id = $1 AND x = $2 AND y = $3", villageID, x, y).Scan(&buildingName, &currentLevel)
	if err != nil {
		return gold, elixir, errors.New("error getting building info on grid")
	}

	if buildingName == "Gold Mine" {
		var goldLastCollected time.Time
		err = tx.QueryRow(ctx, "SELECT gold_last_collected_at FROM villages WHERE id = $1", villageID).Scan(&goldLastCollected)
		elapsedMinutes := int(time.Since(goldLastCollected).Minutes())

		var goldGen *int
		err = tx.QueryRow(ctx, "SELECT LEAST(total_gold_cap, total_gold_rate * $2) FROM village_production_stats WHERE village_id = $1", villageID, elapsedMinutes).Scan(&goldGen)
		if goldGen != nil {
			_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold + $1, gold_last_collected_at = NOW() WHERE id = $2", *goldGen, villageID)
		}
	} else if buildingName == "Elixir Collector" {
		var elixirLastCollected time.Time
		err = tx.QueryRow(ctx, "SELECT elixir_last_collected_at FROM villages WHERE id = $1", villageID).Scan(&elixirLastCollected)
		elapsedMinutes := int(time.Since(elixirLastCollected).Minutes())

		var elixirGen *int
		err = tx.QueryRow(ctx, "SELECT LEAST(total_elixir_cap, total_elixir_rate * $2) FROM village_production_stats WHERE village_id = $1", villageID, elapsedMinutes).Scan(&elixirGen)
		if elixirGen != nil {
			_, err = tx.Exec(ctx, "UPDATE villages SET elixir = elixir + $1, elixir_last_collected_at = NOW() WHERE id = $2", *elixirGen, villageID)
		}
	}

	var upgradeCost, minThLevel, buildTimeSeconds int
	var resourceType string

	costQuery := `
		SELECT b.build_resource_type, c.build_cost, c.min_thlevel
		FROM building_configs b
		JOIN (
			SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM defense_configs
			UNION ALL SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM resource_gen_configs
			UNION ALL SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM resource_storage_configs
			UNION ALL SELECT name, level, build_cost, min_thlevel, build_time_seconds FROM army_camp_configs
		) c ON b.name = c.name
		WHERE b.name = $1 AND c.level = $2
	`
	err = tx.QueryRow(ctx, costQuery, buildingName, currentLevel+1).Scan(&resourceType, &upgradeCost, &minThLevel, &buildTimeSeconds)
	if err != nil {
		return gold, elixir, errors.New("building already at max level")
	}

	if thlevel < minThLevel && buildingName != "Town Hall" {
		return gold, elixir, errors.New("minimum town hall level requirement not met")
	}

	if resourceType == "gold" {
		if gold < upgradeCost {
			return gold, elixir, errors.New("insufficient gold")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold - $1 WHERE id = $2", upgradeCost, villageID)
		gold -= upgradeCost
	} else if resourceType == "elixir" {
		if elixir < upgradeCost {
			return gold, elixir, errors.New("insufficient elixir")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET elixir = elixir - $1 WHERE id = $2", upgradeCost, villageID)
		elixir -= upgradeCost
	}

	_, err = tx.Exec(ctx, "UPDATE village_buildings SET status = 'upgrading', upgrade_complete_at = NOW() + INTERVAL '1 second' * $4 WHERE village_id = $1 AND x = $2 AND y = $3", villageID, x, y, buildTimeSeconds)
	if err != nil {
		return gold, elixir, errors.New("couldn't update village building")
	}

	if buildingName == "Town Hall" {
		_, err = tx.Exec(ctx, "UPDATE villages SET town_hall_level = $1 WHERE id = $2", currentLevel+1, villageID)
		if err != nil {
			return gold, elixir, errors.New("error updating town hall level")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return gold, elixir, err
	}

	return gold, elixir, nil
}

func MoveBuilding(userID string, oldX, oldY, newX, newY int) error {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var villageID string
	err = tx.QueryRow(ctx, "SELECT id FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID)

	if err != nil {
		return errors.New("Error fetching villageID")
	}

	var size int
	var villageBuildingID int
	query := `
	SELECT bc.size, vb.id
	FROM building_configs bc
		JOIN village_buildings vb ON vb.building_name = bc.name
		WHERE vb.village_id = $1
		AND vb.x = $2
		AND vb.y = $3
	`
	err = tx.QueryRow(ctx, query, villageID, oldX, oldY).Scan(&size, &villageBuildingID)

	if err != nil {
		return errors.New("No such building exists.")
	}

	if newX < 0 || newY < 0 || newX+size > 36 || newY+size > 36 {
		return errors.New("Out of village bounds.")
	}

	var collisionCount int
	overlapQuery := `
        SELECT COUNT(*) 
        FROM village_buildings vb
        JOIN building_configs bc ON vb.building_name = bc.name
        WHERE vb.village_id = $1
          AND $2 < (vb.x + bc.size)
          AND ($2 + $3) > vb.x
          AND $4 < (vb.y + bc.size)
          AND ($4 + $3) > vb.y
		  AND vb.id != $5
    `
	err = tx.QueryRow(ctx, overlapQuery, villageID, newX, size, newY, villageBuildingID).Scan(&collisionCount)

	if err != nil {
		return errors.New("Error checking grid.")
	}

	if collisionCount > 0 {
		return errors.New("Cannot place building on an existing building.")
	}

	_, err = tx.Exec(ctx, "UPDATE village_buildings SET x = $1, y = $2 WHERE id = $3", newX, newY, villageBuildingID)

	if err != nil {
		return errors.New("Error updating building's co-ordinates")
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func SyncBuildings(userID string) error {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var villageID string
	err = tx.QueryRow(ctx, "SELECT id FROM villages WHERE user_id = $1", userID).Scan(&villageID)
	if err != nil {
		return errors.New("Error fetching village details")
	}

	_, err = tx.Exec(ctx, `
		WITH updated_buildings AS (
			UPDATE village_buildings 
			SET level = level + 1, status = 'active', upgrade_complete_at = NULL 
			WHERE village_id = $1 AND status = 'upgrading' AND upgrade_complete_at <= NOW()
			RETURNING building_name, level
		)
		UPDATE villages v
		SET town_hall_level = ub.level
		FROM updated_buildings ub
		WHERE v.id = $1 AND ub.building_name = 'Town Hall';
	`, villageID)

	if err != nil {
		return errors.New("error syncing buildings")
	}

	return tx.Commit(ctx)
}
