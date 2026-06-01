package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}

	var req dtos.RegisterRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid req body.", http.StatusBadRequest)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
	}

	userID, err := repository.CreateUserAndVillage(req.Username, passwordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	res := dtos.RegisterResponseDTO{
		ID: userID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
