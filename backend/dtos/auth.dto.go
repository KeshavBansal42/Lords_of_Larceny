package dtos

type RegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponseDTO struct {
	ID int `json:"id"`
}
