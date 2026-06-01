package repository

import (
	"context"
	"errors"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
)

func CreateUserAndVillage(req dtos.RegisterRequestDTO) (int, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var userID int
	err = tx.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", req.Username, req.PasswordHash).Scan(&userID)

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
