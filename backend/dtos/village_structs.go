package dtos

type VillageResponseDTO struct {
	TownHallLevel int `json:"town_hall_level"`
	Gold          int `json:"gold"`
	Elixir        int `json:"elixir"`
}

type CollectGoldResponseDTO struct {
	Message string `json:"message"`
	Gold    int    `json:"gold"`
}

type CollectElixirResponseDTO struct {
	Message string `json:"message"`
	Elixir  int    `json:"elixir"`
}

type ScoutVillageResponseDTO struct {
	Username      string                      `json:"username"`
	TownHallLevel int                         `json:"town_hall_level"`
	Gold          int                         `json:"gold"`
	Elixir        int                         `json:"elixir"`
	Buildings     []BuildingResponseFromDBDTO `json:"buildings"`
}
