package dtos

type MatchmakeResponseFromDBDTO struct {
	UserID int `db:"user_id"`
}

type MatchmakeResponseDTO struct {
	Message   string `json:"message"`
	VillageID int    `json:"villageid"`
}

type TroopDropDTO struct {
	TroopID int `json:"troop_id"`
	X       int `json:"x"`
	Y       int `json:"y"`
}

type AttackRequestDTO struct {
	TargetUserID int            `json:"target_user_id"`
	Drops        []TroopDropDTO `json:"drops"`
}

type BattleEventDTO struct {
	Tick     int    `json:"tick"`
	Action   string `json:"action"`
	EntityID string `json:"entity_id"`
	TargetID string `json:"target_id,omitempty"`
	X        int    `json:"x,omitempty"`
	Y        int    `json:"y,omitempty"`
}

type BattleResultDTO struct {
	PercentageDestroyed int              `json:"percentage_destroyed"`
	GoldStolen          int              `json:"gold_stolen"`
	ElixirStolen        int              `json:"elixir_stolen"`
	Log                 []BattleEventDTO `json:"log"`
}
