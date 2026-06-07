package main

import (
	"net/http"
	"os"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/router"
)

func main() {
	db.InitDB()

	r := router.InitRoutes()

	addrString := os.Getenv("SERVER_URL")
	http.ListenAndServe(addrString, r)
}
