package controllers

import (
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
)

func GetGameConfigs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	buildings, troops, err := repository.GetGameConfigs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.GameConfigsResponseDTO{
		Buildings: buildings,
		Troops:    troops,
	}

	respond(w, http.StatusOK, res)
}
