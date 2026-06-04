package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.RegisterRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid req body.", http.StatusBadRequest)
		return
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
		return
	}

	userID, err := repository.CreateUserAndVillage(req.Username, passwordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := dtos.RegisterResponseDTO{
		Message: "User successfully registered.",
		ID:      userID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.LoginRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid req body.", http.StatusBadRequest)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	userID, passwordHash, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Wrong credentials.", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	res := dtos.LoginResponseDTO{
		Message: "User successfully logged in",
		Token:   tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetVillage(w http.ResponseWriter, r *http.Request) {
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

	_, thlevel, gold, elixir, err := repository.GetVillageByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.VillageResponseDTO{
		TownHallLevel: thlevel,
		Gold:          gold,
		Elixir:        elixir,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetAllVillageBuildings(w http.ResponseWriter, r *http.Request) {
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

	buildings, err := repository.GetAllVillageBuildings(villageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.GetVillageBuildingsResponseDTO{
		Buildings: buildings,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func AddBuilding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.BuildRequestDTO
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

	gold, err := repository.AddBuilding(userID, req.BuildingID, req.X, req.Y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.BuildResponseDTO{
		Message:       "Your building has been created.",
		RemainingGold: gold,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func CollectResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	gold, elixir, err := repository.CollectResources(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.CollectResponseDTO{
		Message: "Resources succesfully collected",
		Gold:    gold,
		Elixir:  elixir,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func UpgradeBuilding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.UpgradeBuildingRequestDTO
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

	_, _, err = repository.CollectResources(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gold, err := repository.UpgradeBuilding(userID, req.X, req.Y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.UpgradeBuildingResponseDTO{
		Message: "Building successfully upgraded",
		Gold:    gold,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
