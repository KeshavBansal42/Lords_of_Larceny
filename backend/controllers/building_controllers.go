package controllers

import (
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
)

func GetAllVillageBuildings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	villageID, _, _, _, err := repository.GetVillageByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buildings, err := repository.GetAllVillageBuildings(villageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.GetVillageBuildingsResponseDTO{
		Buildings: buildings,
	}

	respond(w, http.StatusOK, res)
}

func AddBuilding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.BuildRequestDTO
	check := parseRequest(w, r, &req)
	if !check {
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	gold, elixir, err := repository.AddBuilding(userID, req.BuildingID, req.X, req.Y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.BuildResponseDTO{
		Message:         "Your building has been created.",
		RemainingGold:   gold,
		RemainingElixir: elixir,
	}

	respond(w, http.StatusCreated, res)
}

func UpgradeBuilding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.UpgradeBuildingRequestDTO
	check := parseRequest(w, r, &req)
	if !check {
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	_, _, err = repository.CollectResources(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gold, elixir, err := repository.UpgradeBuilding(userID, req.X, req.Y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.UpgradeBuildingResponseDTO{
		Message: "Building successfully upgraded",
		Gold:    gold,
		Elixir:  elixir,
	}

	respond(w, http.StatusOK, res)
}

func MoveBuilding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.MoveBuildingRequestDTO
	check := parseRequest(w, r, &req)
	if !check {
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	err = repository.MoveBuilding(userID, req.OldX, req.OldY, req.NewX, req.NewY)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.MoveBuildingResponseDTO{
		Message: "Building moved successfully",
	}

	respond(w, http.StatusOK, res)
}
