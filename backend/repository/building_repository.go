package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/jackc/pgx/v5"
)

func GetAllVillageBuildings(villageID int) ([]dtos.BuildingResponseFromDBDTO, error) {
	rows, err := db.Conn.Query(context.Background(), "SELECT building_id, x, y FROM village_buildings WHERE village_id = $1", villageID)

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

func AddBuilding(userID, buildingID, x, y int) (int, int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var name string
	var buildCost int
	var size int
	var level int
	var min_thlevel int
	var resourceType string

	err = tx.QueryRow(ctx, "SELECT name, build_cost, size, level, min_thlevel, build_resource_type FROM building_configs WHERE id = $1", buildingID).Scan(&name, &buildCost, &size, &level, &min_thlevel, &resourceType)

	if err != nil {
		return 0, 0, errors.New("Error fetching building details.")
	}

	if level != 1 {
		return 0, 0, fmt.Errorf("Cannot add a level %v machine please upgrade one", level)
	}

	if name == "Town Hall" {
		return 0, 0, errors.New("Cannot build another town hall.")
	}

	var gold int
	var elixir int
	var thlevel int
	var villageID int
	err = tx.QueryRow(ctx, "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &thlevel, &gold, &elixir)

	if err != nil {
		return 0, 0, errors.New("Error fetching village details.")
	}

	if thlevel < min_thlevel {
		return gold, elixir, errors.New("Minimum Town Hall Level requirement not met.")
	}

	if x < 0 || y < 0 || x+size > 36 || y+size > 36 {
		return gold, elixir, errors.New("Building is out of bounds.")
	}

	var collisionCount int
	overlapQuery := `
        SELECT COUNT(*) 
        FROM village_buildings vb
        JOIN building_configs bc ON vb.building_id = bc.id
        WHERE vb.village_id = $1
          AND $2 < (vb.x + bc.size)
          AND ($2 + $3) > vb.x
          AND $4 < (vb.y + bc.size)
          AND ($4 + $3) > vb.y
    `
	err = tx.QueryRow(ctx, overlapQuery, villageID, x, size, y).Scan(&collisionCount)

	if err != nil {
		return gold, elixir, errors.New("Error checking grid.")
	}

	if collisionCount > 0 {
		return gold, elixir, errors.New("Cannot place building on an existing building.")
	}

	if resourceType == "gold" {
		if gold < buildCost {
			return gold, elixir, errors.New("Insufficient gold.")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold-$1 WHERE id = $2", buildCost, villageID)
		gold -= buildCost
	} else if resourceType == "elixir" {
		if elixir < buildCost {
			return gold, elixir, errors.New("Insufficient elixir.")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET elixir = elixir-$1 WHERE id = $2", buildCost, villageID)
		elixir -= buildCost
	} else {
		return gold, elixir, errors.New("Invalid resource type configured for this building.")
	}

	if err != nil {
		return gold, elixir, errors.New("Error updating resources.")
	}

	_, err = tx.Exec(ctx, "INSERT INTO village_buildings (village_id, building_id, x, y) VALUES ($1, $2, $3, $4)", villageID, buildingID, x, y)

	if err != nil {
		return gold, elixir, errors.New("Error adding building.")
	}

	if err = tx.Commit(ctx); err != nil {
		return gold, elixir, err
	}

	return gold, elixir, nil
}

func UpgradeBuilding(userId, x, y int) (int, int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var villageID int
	var thlevel int
	var gold int
	var elixir int

	err = tx.QueryRow(ctx, "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1 FOR UPDATE", userId).Scan(&villageID, &thlevel, &gold, &elixir)

	if err != nil {
		return 0, 0, errors.New("Error fetching village details.")
	}

	var buildingName string
	var buildingLevel int
	var min_thlevel int
	query := `
        SELECT name, level, min_thlevel
        FROM building_configs bc
        JOIN village_buildings vb ON vb.building_id = bc.id
        WHERE vb.village_id = $1
        AND vb.x = $2
        AND vb.y = $3
    `
	err = tx.QueryRow(ctx, query, villageID, x, y).Scan(&buildingName, &buildingLevel, &min_thlevel)

	if err != nil {
		return gold, elixir, errors.New("Error getting building info")
	}

	if thlevel < min_thlevel {
		return gold, elixir, errors.New("Minimum town hall level requirement not met")
	}

	var upgradeCost int
	var newID int
	var resourceType string

	err = tx.QueryRow(ctx, "SELECT build_cost, id, build_resource_type FROM building_configs WHERE name = $1 AND level = $2", buildingName, (buildingLevel+1)).Scan(&upgradeCost, &newID, &resourceType)

	if err != nil {
		return gold, elixir, errors.New("Building already at max level.")
	}

	if resourceType == "gold" {
		if gold < upgradeCost {
			return gold, elixir, errors.New("Insufficient gold.")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold - $1 WHERE id = $2", upgradeCost, villageID)
		gold -= upgradeCost
	} else if resourceType == "elixir" {
		if elixir < upgradeCost {
			return gold, elixir, errors.New("Insufficient elixir.")
		}
		_, err = tx.Exec(ctx, "UPDATE villages SET elixir = elixir - $1 WHERE id = $2", upgradeCost, villageID)
		elixir -= upgradeCost
	} else {
		return gold, elixir, errors.New("Invalid resource type configured for this building.")
	}

	if err != nil {
		return gold, elixir, errors.New("Couldn't update balance")
	}

	_, err = tx.Exec(ctx, "UPDATE village_buildings SET building_id = $1 WHERE village_id = $2 AND x = $3 AND y = $4", newID, villageID, x, y)

	if err != nil {
		return gold, elixir, errors.New("Couldn't update village building")
	}

	if buildingName == "Town Hall" {
		_, err = tx.Exec(ctx, "UPDATE villages SET town_hall_level = $1 WHERE id = $2", (buildingLevel + 1), villageID)

		if err != nil {
			return gold, elixir, errors.New("Error updating town hall level")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return gold, elixir, err
	}

	return gold, elixir, nil
}

func MoveBuilding(userID, oldX, oldY, newX, newY int) error {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var villageID int
	err = tx.QueryRow(ctx, "SELECT id FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID)

	if err != nil {
		return errors.New("Error fetching villageID")
	}

	var size int
	var villageBuildingID int
	query := `
	SELECT bc.size, vb.id
	FROM building_configs bc
		JOIN village_buildings vb ON vb.building_id = bc.id
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
        JOIN building_configs bc ON vb.building_id = bc.id
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
