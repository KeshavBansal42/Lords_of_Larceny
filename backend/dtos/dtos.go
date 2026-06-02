package dtos

type RegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponseDTO struct {
	ID int `json:"id"`
}

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}

type VillageResponseDTO struct {
	TownHallLevel int `json:"town_hall_level"`
	Gold          int `json:"gold"`
	Elixir        int `json:"elixir"`
}
