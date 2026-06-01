package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/router"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB()

	r := router.InitRoutes()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	addrString := os.Getenv("SERVER_URI")

	http.ListenAndServe(addrString, r)
}
