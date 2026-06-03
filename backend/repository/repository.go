package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/jackc/pgx/v5"
)

func CreateUserAndVillage(username, password_hash string) (int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var userID int
	err = tx.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", username, password_hash).Scan(&userID)

	if err != nil {
		return 0, errors.New("username already exists")
	}

	var villageID int
	err = tx.QueryRow(ctx, "INSERT INTO villages (user_id, town_hall_level, gold, elixir) VALUES ($1, 1, 1000, 1000) RETURNING id", userID).Scan(&villageID)

	if err != nil {

		return 0, errors.New("failed to initialise village")
	}

	_, err = tx.Exec(ctx, "INSERT INTO village_buildings (village_id, building_id, x, y) VALUES ($1, 1, 16, 16)", villageID)

	if err != nil {

		return 0, errors.New("failed to add town hall")
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return userID, nil
}

func GetUserByUsername(username string) (int, string, error) {
	var userID int
	var passwordHash string
	err := db.Conn.QueryRow(context.Background(), "SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &passwordHash)

	if err != nil {
		return 0, "", errors.New("Error getting the userID and PasswordHash")
	}

	return userID, passwordHash, nil
}

func GetVillageByUserID(userID int) (int, int, int, int, error) {
	var townHallLevel int
	var gold int
	var elixir int
	var villageID int
	err := db.Conn.QueryRow(context.Background(), "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1", userID).Scan(&villageID, &townHallLevel, &gold, &elixir)

	if err != nil {
		return 0, 0, 0, 0, errors.New("Error fetching the village")
	}

	return villageID, townHallLevel, gold, elixir, nil
}

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

func AddBuilding(userID, buildingID, x, y int) (int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var name string
	var buildCost int
	var size int
	var level int
	var min_thlevel int
	err = tx.QueryRow(ctx, "SELECT name, build_cost, size, level, min_thlevel FROM building_configs WHERE id = $1", buildingID).Scan(&name, &buildCost, &size, &level, &min_thlevel)

	if err != nil {
		return 0, errors.New("Error fetching building details.")
	}

	if level != 1 {
		return 0, fmt.Errorf("Cannot add a level %v machine please upgrade one", level)
	}

	if name == "Town Hall" {
		return 0, errors.New("Cannot build another town hall.")
	}

	var gold int
	var elixir int
	var thlevel int
	var villageID int
	err = tx.QueryRow(ctx, "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &thlevel, &gold, &elixir)

	if err != nil {
		return 0, errors.New("Error fetching village details.")
	}

	if thlevel < min_thlevel {
		return gold, errors.New("Minimum Town Hall Level requirement not met.")
	}

	if gold < buildCost {
		return gold, errors.New("Insufficient balance,")
	}

	if x < 0 || y < 0 || x+size > 36 || y+size > 36 {
		return gold, errors.New("Building is out of bounds.")
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
		return gold, errors.New("Error checking grid.")
	}

	if collisionCount > 0 {
		return gold, errors.New("Cannot place building on an existing building.")
	}

	_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold-$1 WHERE id = $2", buildCost, villageID)

	if err != nil {
		return gold, errors.New("Error updating resources.")
	}

	_, err = tx.Exec(ctx, "INSERT INTO village_buildings (village_id, building_id, x, y) VALUES ($1, $2, $3, $4)", villageID, buildingID, x, y)

	if err != nil {
		return gold, errors.New("Error adding building.")
	}

	if err = tx.Commit(ctx); err != nil {
		return gold, err
	}

	return gold - buildCost, nil
}
