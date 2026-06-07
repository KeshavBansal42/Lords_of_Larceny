package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	connectionString := os.Getenv("DATABASE_URL")

	Conn, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Println("Error connecting to database")
	}

	SeedDatabase(Conn)
}
