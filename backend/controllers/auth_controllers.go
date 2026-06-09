package controllers

import (
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
	check := parseRequest(w, r, &req)
	if !check {
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

	respond(w, http.StatusCreated, res)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.LoginRequestDTO
	check := parseRequest(w, r, &req)
	if !check {
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

	respond(w, http.StatusOK, res)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	err = repository.DeleteAccount(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.DeleteAccountResponseDTO{
		Message: "Account deleted successfully",
	}

	respond(w, http.StatusOK, res)
}
