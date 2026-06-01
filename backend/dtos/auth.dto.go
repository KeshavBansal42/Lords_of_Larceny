package dtos

type RegisterRequestDTO struct {
	username      string
	password_hash string
}

type RegisterResponseDTO struct {
	ID int
}
