package dtos

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

type UpgradeBuildingRequestDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type UpgradeBuildingResponseDTO struct {
	Message string `json:"message"`
	Gold    int    `json:"gold"`
	Elixir  int    `json:"elixir"`
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
