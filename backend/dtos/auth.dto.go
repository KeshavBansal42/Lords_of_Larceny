package dtos

type RegisterRequestDTO struct {
	Username     string
	PasswordHash string
}

type RegisterResponseDTO struct {
	ID int
}
