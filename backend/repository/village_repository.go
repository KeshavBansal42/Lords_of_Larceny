package repository

import (
	"context"
	"errors"
	"time"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/jackc/pgx/v5"
)

func CreateUserAndVillage(username, password_hash string) (string, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var userID string
	err = tx.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", username, password_hash).Scan(&userID)

	if err != nil {
		return "", errors.New("username already exists")
	}

	var villageID string
	err = tx.QueryRow(ctx, "INSERT INTO villages (user_id, town_hall_level, gold, elixir) VALUES ($1, 1, 1000, 1000) RETURNING id", userID).Scan(&villageID)

	if err != nil {

		return "", errors.New("failed to initialise village")
	}

	_, err = tx.Exec(ctx, "INSERT INTO village_buildings (village_id, building_name, level, x, y) VALUES ($1, 'Town Hall', 1, 16, 16)", villageID)

	if err != nil {

		return "", errors.New("failed to add town hall")
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return userID, nil
}

func GetUserByUsername(username string) (string, string, error) {
	var userID string
	var passwordHash string
	err := db.Conn.QueryRow(context.Background(), "SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &passwordHash)

	if err != nil {
		return "", "", errors.New("Error getting the userID and PasswordHash")
	}

	return userID, passwordHash, nil
}

func GetVillageByUserID(userID string) (string, int, int, int, error) {
	var townHallLevel int
	var gold int
	var elixir int
	var villageID string
	err := db.Conn.QueryRow(context.Background(), "SELECT id, town_hall_level, gold, elixir FROM villages WHERE user_id = $1", userID).Scan(&villageID, &townHallLevel, &gold, &elixir)

	if err != nil {
		return "", 0, 0, 0, errors.New("Error fetching the village")
	}

	return villageID, townHallLevel, gold, elixir, nil
}

func CollectGold(userID string) (int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var villageID string
	var currentGold int
	var goldLastCollected time.Time

	err = tx.QueryRow(ctx, "SELECT id, gold, gold_last_collected_at FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &currentGold, &goldLastCollected)
	if err != nil {
		return 0, errors.New("Error fetching village details.")
	}

	var goldRate, goldCap, maxStorage int
	query := `
		SELECT total_gold_rate, total_gold_cap, max_gold_storage
		FROM village_production_stats
		WHERE village_id = $1;
	`
	err = tx.QueryRow(ctx, query, villageID).Scan(&goldRate, &goldCap, &maxStorage)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			goldRate, goldCap, maxStorage = 0, 0, 0
		} else {
			return 0, errors.New("error calculating resources")
		}
	}

	if goldRate == 0 {
		if err = tx.Commit(ctx); err != nil {
			return 0, err
		}
		return currentGold, nil
	}

	now := time.Now()
	elapsedMinutes := now.Sub(goldLastCollected).Minutes()

	totalGenerated := int(float64(goldRate) * elapsedMinutes)
	if totalGenerated > goldCap {
		totalGenerated = goldCap
	}

	availableSpace := maxStorage - currentGold
	if availableSpace < 0 {
		availableSpace = 0
	}

	goldToAdd := totalGenerated
	if goldToAdd > availableSpace {
		goldToAdd = availableSpace
	}

	leftoverGold := totalGenerated - goldToAdd
	var newLastCollected time.Time

	if leftoverGold > 0 {
		leftoverMinutes := float64(leftoverGold) / float64(goldRate)
		newLastCollected = now.Add(-time.Duration(leftoverMinutes * float64(time.Minute)))
	} else {
		newLastCollected = now
	}

	var newGold int
	err = tx.QueryRow(ctx, "UPDATE villages SET gold = gold + $1, gold_last_collected_at = $2 WHERE id = $3 RETURNING gold", goldToAdd, newLastCollected, villageID).Scan(&newGold)
	if err != nil {
		return 0, errors.New("error updating resources")
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return newGold, nil
}

func CollectElixir(userID string) (int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var villageID string
	var currentElixir int
	var elixirLastCollected time.Time

	err = tx.QueryRow(ctx, "SELECT id, elixir, elixir_last_collected_at FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &currentElixir, &elixirLastCollected)
	if err != nil {
		return 0, errors.New("Error fetching village details.")
	}

	var elixirRate, elixirCap, maxStorage int
	query := `
		SELECT total_elixir_rate, total_elixir_cap, max_elixir_storage
		FROM village_production_stats
		WHERE village_id = $1;
	`
	err = tx.QueryRow(ctx, query, villageID).Scan(&elixirRate, &elixirCap, &maxStorage)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			elixirRate, elixirCap, maxStorage = 0, 0, 0
		} else {
			return 0, errors.New("error calculating resources")
		}
	}

	if elixirRate == 0 {
		if err = tx.Commit(ctx); err != nil {
			return 0, err
		}
		return currentElixir, nil
	}

	now := time.Now()
	elapsedMinutes := now.Sub(elixirLastCollected).Minutes()

	totalGenerated := int(float64(elixirRate) * elapsedMinutes)
	if totalGenerated > elixirCap {
		totalGenerated = elixirCap
	}

	availableSpace := maxStorage - currentElixir
	if availableSpace < 0 {
		availableSpace = 0
	}

	elixirToAdd := totalGenerated
	if elixirToAdd > availableSpace {
		elixirToAdd = availableSpace
	}

	leftoverElixir := totalGenerated - elixirToAdd
	var newLastCollected time.Time

	if leftoverElixir > 0 {
		leftoverMinutes := float64(leftoverElixir) / float64(elixirRate)
		newLastCollected = now.Add(-time.Duration(leftoverMinutes * float64(time.Minute)))
	} else {
		newLastCollected = now
	}

	var newElixir int
	err = tx.QueryRow(ctx, "UPDATE villages SET elixir = elixir + $1, elixir_last_collected_at = $2 WHERE id = $3 RETURNING elixir", elixirToAdd, newLastCollected, villageID).Scan(&newElixir)
	if err != nil {
		return 0, errors.New("error updating resources")
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return newElixir, nil
}

func ScoutVillage(targetUserID string) (string, int, int, int, []dtos.BuildingResponseFromDBDTO, error) {
	ctx := context.Background()

	var username string
	var villageID string
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

	buildingRows, err := db.Conn.Query(ctx, "SELECT building_name, level, x, y, status FROM village_buildings WHERE village_id = $1", villageID)
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

func DeleteAccount(userID string) error {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return errors.New("Error deleting user account")
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
