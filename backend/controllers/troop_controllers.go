package controllers

import (
	"encoding/json"
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
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid req body.", http.StatusBadRequest)
		return
	}

	value := r.Context().Value("userID")
	userIDFloat, ok := value.(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	userID := int(userIDFloat)

	err = repository.TrainTroops(userID, req.TroopsToTrain)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.TrainTroopsResponseDTO{
		Message: "Troops trained successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetAllVillageTroops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	value := r.Context().Value("userID")
	userIDFloat, ok := value.(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	userID := int(userIDFloat)

	villageID, _, _, _, err := repository.GetVillageByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	troops, err := repository.GetAllVillageTroops(villageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.GetVillageTroopsResponseDTO{
		Troops: troops,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
