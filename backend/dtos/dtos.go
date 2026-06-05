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
	Message         string `json:"message"`
	RemainingGold   int    `json:"remaining_gold"`
	RemainingElixir int    `json:"remaining_elixir"`
}

type GetVillageBuildingsResponseDTO struct {
	Buildings []BuildingResponseFromDBDTO `json:"buildings"`
}

type BuildingResponseFromDBDTO struct {
	BuildingId int `json:"building_id"`
	X          int `json:"x"`
	Y          int `json:"y"`
}

type CollectResponseDTO struct {
	Message string `json:"message"`
	Gold    int    `json:"gold"`
	Elixir  int    `json:"elixir"`
}

type UpgradeBuildingRequestDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type UpgradeBuildingResponseDTO struct {
	Message string `json:"message"`
	Gold    int    `json:"gold"`
}

type MoveBuildingRequestDTO struct {
	OldX int `json:"oldx"`
	OldY int `json:"oldy"`
	NewX int `json:"newx"`
	NewY int `json:"newy"`
}

type MoveBuildingResponseDTO struct {
	Message string `json:"message"`
}
