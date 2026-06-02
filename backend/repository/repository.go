package repository

import (
	"context"
	"errors"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
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

	_, err = tx.Exec(ctx, "INSERT INTO villages (user_id, town_hall_level, gold, elixir) VALUES ($1, 1, 1000, 1000)", userID)

	if err != nil {

		return 0, errors.New("failed to initialise village")
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
