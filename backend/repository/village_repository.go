package repository

import (
	"context"
	"errors"
	"time"

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

func CollectResources(userID int) (int, int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var villageID int
	var lastCollectedAt time.Time

	err = tx.QueryRow(ctx, "SELECT id, last_collected_at FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &lastCollectedAt)

	if err != nil {
		return 0, 0, errors.New("Error fetching village details.")
	}

	now := time.Now()
	elapsedMinutes := int(now.Sub(lastCollectedAt).Minutes())

	var goldGen *int
	var elixirGen *int
	query := `
        SELECT 
            LEAST(total_gold_cap, total_gold_rate * $2) AS total_gold_generated,
            LEAST(total_elixir_cap, total_elixir_rate * $2) AS total_elixir_generated
        FROM village_production_stats
        WHERE village_id = $1;
    `
	err = tx.QueryRow(ctx, query, villageID, elapsedMinutes).Scan(&goldGen, &elixirGen)
	if err != nil {
		return 0, 0, errors.New("error calculating resources")
	}

	goldToAdd := 0
	if goldGen != nil {
		goldToAdd = *goldGen
	}

	elixirToAdd := 0
	if elixirGen != nil {
		elixirToAdd = *elixirGen
	}

	var newGold int
	var newElixir int
	err = tx.QueryRow(ctx, "UPDATE villages SET gold = gold + $1, elixir = elixir + $2, last_collected_at = $3 WHERE id = $4 RETURNING gold, elixir", goldToAdd, elixirToAdd, now, villageID).Scan(&newGold, &newElixir)
	if err != nil {
		return 0, 0, errors.New("error updating resources")
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, err
	}

	return newGold, newElixir, nil
}

func ScoutVillage(targetUserID int) (string, int, int, int, []dtos.BuildingResponseFromDBDTO, error) {
	ctx := context.Background()

	var username string
	var villageID int
	var thLevel int
	var gold int
	var elixir int

	userQuery := `
		SELECT u.username, v.id, v.town_hall_level, v.gold, v.elixir 
		FROM villages v
		JOIN users u ON v.user_id = u.id
		WHERE v.user_id = $1
	`
	err := db.Conn.QueryRow(ctx, userQuery, targetUserID).Scan(&username, &villageID, &thLevel, &gold, &elixir)
	if err != nil {
		return "", 0, 0, 0, nil, errors.New("Village not found")
	}

	buildingRows, err := db.Conn.Query(ctx, "SELECT building_id, x, y FROM village_buildings WHERE village_id = $1", villageID)
	if err != nil {
		return "", 0, 0, 0, nil, errors.New("Error fetching enemy buildings")
	}
	defer buildingRows.Close()

	buildings, err := pgx.CollectRows(buildingRows, pgx.RowToStructByName[dtos.BuildingResponseFromDBDTO])
	if err != nil {
		return "", 0, 0, 0, nil, errors.New("Error parsing enemy buildings")
	}

	return username, thLevel, gold, elixir, buildings, nil
}
