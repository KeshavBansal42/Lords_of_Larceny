package dtos

type TrainTroopsRequestDTO struct {
	TroopsToTrain map[int]int `json:"troopstotrain"`
}

type TrainTroopsResponseDTO struct {
	Message string `json:"message"`
}

type TroopResponseFromDBDTO struct {
	TroopID  int `json:"troopid"`
	Quantity int `json:"quantity"`
}

type GetVillageTroopsResponseDTO struct {
	Troops []TroopResponseFromDBDTO `json:"troops"`
}
