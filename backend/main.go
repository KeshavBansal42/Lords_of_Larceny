package main

import (
	"context"
	"log"
	"os"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	connectionString := os.Getenv("DB_URI")

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Println("Error connecting to database")
	}

	db.SeedDatabase(conn)
}
