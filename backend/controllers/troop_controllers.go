package controllers

import (
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
)

func TrainTroops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.TrainTroopsRequestDTO
	check := parseRequest(w, r, &req)
	if !check {
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	err = repository.TrainTroops(r.Context(), userID, req.TroopsToTrain)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.TrainTroopsResponseDTO{
		Message: "Troops trained successfully",
	}

	respond(w, http.StatusOK, res)
}

func GetAllVillageTroops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	villageID, _, _, _, err := repository.GetVillageByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	troops, err := repository.GetAllVillageTroops(r.Context(), villageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.GetVillageTroopsResponseDTO{
		Troops: troops,
	}

	respond(w, http.StatusOK, res)
}
