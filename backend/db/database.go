package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	connectionString := os.Getenv("DATABASE_URL")

	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Println("Error connecting to database:", err)
		os.Exit(1) // Better to crash on boot if DB is down
	}

	SeedDatabase(Pool)
}
