package dtos

import "github.com/KeshavBansal42/Lords_of_Larceny/backend/models"

type GameConfigsResponseDTO struct {
	Buildings []models.BuildingConfig `json:"buildings"`
	Troops    []models.TroopConfig    `json:"troops"`
}
