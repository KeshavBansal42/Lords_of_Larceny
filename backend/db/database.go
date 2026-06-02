package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDB() {
	connectionString := os.Getenv("DB_URI")

	Conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Println("Error connecting to database")
	}

	SeedDatabase(Conn)
}
