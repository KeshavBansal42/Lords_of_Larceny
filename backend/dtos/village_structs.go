package dtos

type VillageResponseDTO struct {
	TownHallLevel int `json:"town_hall_level"`
	Gold          int `json:"gold"`
	Elixir        int `json:"elixir"`
}

type CollectResponseDTO struct {
	Message string `json:"message"`
	Gold    int    `json:"gold"`
	Elixir  int    `json:"elixir"`
}

type ScoutVillageResponseDTO struct {
	Username      string                      `json:"username"`
	TownHallLevel int                         `json:"town_hall_level"`
	Gold          int                         `json:"gold"`
	Elixir        int                         `json:"elixir"`
	Buildings     []BuildingResponseFromDBDTO `json:"buildings"`
}
