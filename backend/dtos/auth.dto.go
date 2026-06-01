package dtos

type RegisterRequestDTO struct {
	Username string
	Password string
}

type RegisterResponseDTO struct {
	ID int
}
