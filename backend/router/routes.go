package router

import (
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/controllers"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/middleware"
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/village", middleware.RequireAuth(controllers.GetVillage)).Methods("GET")
	router.HandleFunc("/village/buildings", middleware.RequireAuth(controllers.GetAllVillageBuildings)).Methods("GET")
	router.HandleFunc("/village/buildings/build", middleware.RequireAuth(controllers.AddBuilding)).Methods("POST")
	router.HandleFunc("/village/buildings/upgrade", middleware.RequireAuth(controllers.UpgradeBuilding)).Methods("PUT")
	router.HandleFunc("/village/buildings/move", middleware.RequireAuth(controllers.MoveBuilding)).Methods("PUT")
	router.HandleFunc("/village/troops", middleware.RequireAuth(controllers.GetAllVillageTroops)).Methods("GET")
	router.HandleFunc("/village/troops/train", middleware.RequireAuth(controllers.TrainTroops)).Methods("PUT")
	router.HandleFunc("/village/resources/collect", middleware.RequireAuth(controllers.CollectResources)).Methods("PUT")
	router.HandleFunc("/game/configs", controllers.GetGameConfigs).Methods("GET")

	return router
}
