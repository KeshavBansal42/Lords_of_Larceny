package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Village struct {
	ID            int `json:"id"`
	UserID        int `json:"user_id"`
	TownHallLevel int `json:"town_hall_level"`
	Gold          int `json:"gold"`
	Elixir        int `json:"elixir"`
}

type BuildingConfig struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Level             int    `json:"level"`
	HitPoints         int    `json:"hit_points"`
	Damage            int    `json:"damage"`
	BuildCost         int    `json:"build_cost"`
	BuildResourceType string `json:"build_resource_type"`
	ProductionPerMin  int    `json:"production_per_min"`
	Capacity          int    `json:"capacity"`
	Size              int    `json:"size"`
	MinThLevel        int    `json:"min_thlevel"`
}

type TroopConfig struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Level        int    `json:"level"`
	HitPoints    int    `json:"hit_points"`
	Damage       int    `json:"damage"`
	MinThLevel   int    `json:"min_thlevel"`
	HousingSpace int    `json:"housing_space"`
}

type VillageBuilding struct {
	ID         int `json:"id"`
	VillageID  int `json:"village_id"`
	BuildingID int `json:"building_id"`
	X          int `json:"x"`
	Y          int `json:"y"`
}

type VillageTroop struct {
	VillageID int `json:"village_id"`
	TroopID   int `json:"troop_id"`
	Quantity  int `json:"quantity"`
}
