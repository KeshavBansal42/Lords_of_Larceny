package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Village struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	TownHallLevel int    `json:"town_hall_level"`
	Gold          int    `json:"gold"`
	Elixir        int    `json:"elixir"`
}

type BuildingConfig struct {
	Name              string  `json:"name"`
	Level             int     `json:"level"`
	HitPoints         int     `json:"hit_points"`
	Damage            int     `json:"damage"`
	BuildCost         int     `json:"build_cost"`
	BuildResourceType string  `json:"build_resource_type"`
	ProductionPerMin  int     `json:"production_per_min"`
	Capacity          int     `json:"capacity"`
	Size              int     `json:"size"`
	MinThLevel        int     `json:"min_thlevel"`
	Range             int     `json:"range"`
	SingleTarget      bool    `json:"single_target"`
	SplashRadius      float64 `json:"splash_radius"`
	TargetType        string  `json:"target_type"`
	BuildTimeSeconds  int     `json:"build_time_seconds"`
}

type TroopConfig struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Level        int    `json:"level"`
	HitPoints    int    `json:"hit_points"`
	Damage       int    `json:"damage"`
	MinThLevel   int    `json:"min_thlevel"`
	HousingSpace int    `json:"housing_space"`
	Range        int    `json:"range"`
	Speed        int    `json:"speed"`
	Airborne     bool   `json:"airborne"`
}

type VillageBuilding struct {
	ID           int    `json:"id"`
	VillageID    string `json:"village_id"`
	BuildingName string `json:"building_name"`
	Level        int    `json:"level"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

type VillageTroop struct {
	VillageID string `json:"village_id"`
	TroopID   int    `json:"troop_id"`
	Quantity  int    `json:"quantity"`
}

type LiveBuilding struct {
	ID           string
	BuildingName string
	X            int
	Y            int
	MaxHP        int
	CurrentHP    int
	Damage       int
	SingleTarget bool
	TargetID     string
	Range        int
	SplashRadius float64
	TargetType   string
	Status       string
}

type LiveTroop struct {
	ID        string
	TroopID   int
	X         float64
	Y         float64
	MaxHP     int
	CurrentHP int
	Damage    int
	Range     int
	Speed     int
	TargetID  string
	Airborne  bool
}
