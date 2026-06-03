package dtos

type RegisterRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponseDTO struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
}

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type VillageResponseDTO struct {
	TownHallLevel int `json:"town_hall_level"`
	Gold          int `json:"gold"`
	Elixir        int `json:"elixir"`
}

type BuildRequestDTO struct {
	BuildingID int `json:"building_id"`
	X          int `json:"x"`
	Y          int `json:"y"`
}

type BuildResponseDTO struct {
	Message       string `json:"message"`
	RemainingGold int    `json:"remaining_gold"`
}
