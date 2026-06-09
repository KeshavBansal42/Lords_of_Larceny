package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/models"
	"github.com/jackc/pgx/v5"
)

func Matchmake(userID string) (string, error) {
	ctx := context.Background()

	var villageID string
	var townHallLevel int
	var troopCount int
	query := `
        SELECT v.id, v.town_hall_level, COALESCE(SUM(vt.quantity), 0)
        FROM villages v
        LEFT JOIN village_troops vt ON v.id = vt.village_id
        WHERE v.user_id = $1
        GROUP BY v.id, v.town_hall_level
    `
	err := db.Conn.QueryRow(ctx, query, userID).Scan(&villageID, &townHallLevel, &troopCount)

	if err != nil {
		return "", errors.New("Error finding user info.")
	}

	if troopCount == 0 {
		return "", errors.New("You need an army to attack")
	}

	matchQuery := `
        SELECT user_id 
        FROM villages 
        WHERE id != $1 
          AND town_hall_level BETWEEN $2 AND $3 
          AND (last_attacked_at IS NULL OR last_attacked_at < NOW() - INTERVAL '6 hours')
        LIMIT 50
    `
	rows, err := db.Conn.Query(ctx, matchQuery, villageID, (townHallLevel - 1), (townHallLevel + 1))

	if err != nil {
		return "", errors.New("Error fetching users")
	}
	defer rows.Close()

	villages, err := pgx.CollectRows(rows, pgx.RowToStructByName[dtos.MatchmakeResponseFromDBDTO])
	if err != nil {
		return "", errors.New("Failed to parse villages")
	}

	if len(villages) == 0 {
		return "", errors.New("No worthy opponents exist.")
	}

	randomIndex := rand.Intn(len(villages))
	selectedOpponent := villages[randomIndex]

	return selectedOpponent.UserID, nil
}

func Populate(userID string, liveBuildings *map[string]*models.LiveBuilding, liveTroops *map[string]*models.LiveTroop, buildings []dtos.BuildingResponseFromDBDTO, drops []dtos.TroopDropDTO) error {
	buildingConfigsArray, troopConfigsArray, err := GetGameConfigs()
	if err != nil {
		return errors.New("Error fetching config files")
	}

	buildingConfigs := make(map[string]map[int]models.BuildingConfig)
	for _, bconfig := range buildingConfigsArray {
		if buildingConfigs[bconfig.Name] == nil {
			buildingConfigs[bconfig.Name] = make(map[int]models.BuildingConfig)
		}
		buildingConfigs[bconfig.Name][bconfig.Level] = bconfig
	}

	troopConfigs := make(map[int]models.TroopConfig)
	for _, tconfig := range troopConfigsArray {
		key := tconfig.ID
		troopConfigs[key] = tconfig
	}

	for i, building := range buildings {
		key := fmt.Sprintf("building_%v", (i + 1))
		config := buildingConfigs[building.BuildingName][building.Level]

		(*liveBuildings)[key] = &models.LiveBuilding{
			ID:           key,
			BuildingName: building.BuildingName,
			X:            building.X,
			Y:            building.Y,
			MaxHP:        config.HitPoints,
			CurrentHP:    config.HitPoints,
			Damage:       config.Damage,
			TargetID:     "",
			Range:        config.Range,
			SingleTarget: config.SingleTarget,
			SplashRadius: config.SplashRadius,
			TargetType:   config.TargetType,
		}
	}
	for i, troop := range drops {
		key := fmt.Sprintf("troop_%v", (i + 1))
		(*liveTroops)[key] = &models.LiveTroop{
			ID:        key,
			TroopID:   troop.TroopID,
			X:         float64(troop.X),
			Y:         float64(troop.Y),
			MaxHP:     troopConfigs[troop.TroopID].HitPoints,
			CurrentHP: troopConfigs[troop.TroopID].HitPoints,
			Damage:    troopConfigs[troop.TroopID].Damage,
			Range:     troopConfigs[troop.TroopID].Range,
			Speed:     troopConfigs[troop.TroopID].Speed,
			Airborne:  troopConfigs[troop.TroopID].Airborne,
			TargetID:  "",
		}
	}

	return nil
}

func Battle(userID string, targetUserID string, drops []dtos.TroopDropDTO) (int, int, int, []dtos.BattleEventDTO, error) {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return 0, 0, 0, nil, err
	}
	defer tx.Rollback(ctx)

	troopQuantity := make(map[int]int)
	for _, troop := range drops {
		key := troop.TroopID
		troopQuantity[key]++
	}

	villageID, _, _, _, err := GetVillageByUserID(userID)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Error fetching village id")
	}
	villageTroops, err := GetAllVillageTroops(villageID)

	dbTroopMap := make(map[int]int)
	for _, dbTroop := range villageTroops {
		dbTroopMap[dbTroop.TroopID] = dbTroop.Quantity
	}

	for requestedTroopID, requestedAmount := range troopQuantity {
		if dbTroopMap[requestedTroopID] < requestedAmount {
			return 0, 0, 0, nil, errors.New("Cannot attack with untrained troops.")
		}
	}

	_, _, gold, elixir, buildings, err := ScoutVillage(targetUserID)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Error fetching enemy village")
	}

	liveBuildings := make(map[string]*models.LiveBuilding)
	liveTroops := make(map[string]*models.LiveTroop)

	err = Populate(targetUserID, &liveBuildings, &liveTroops, buildings, drops)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	maxTicks := 1800
	var battleLog []dtos.BattleEventDTO

	totalBuildings := len(liveBuildings)

	for tick := 0; tick < maxTicks; tick++ {
		if len(liveTroops) == 0 || len(liveBuildings) == 0 {
			break
		}

		for _, troop := range liveTroops {
			minDistance := math.MaxFloat64
			minDisID := ""
			for _, building := range liveBuildings {
				distance := math.Sqrt((troop.X-float64(building.X))*(troop.X-float64(building.X)) + (troop.Y-float64(building.Y))*(troop.Y-float64(building.Y)))
				if distance < minDistance {
					minDistance = distance
					minDisID = building.ID
				}
			}
			troop.TargetID = minDisID
		}

		for _, troop := range liveTroops {
			targetBuilding := liveBuildings[troop.TargetID]
			if targetBuilding == nil {
				continue
			}
			distanceFromTarget := math.Sqrt((troop.X-float64(targetBuilding.X))*(troop.X-float64(targetBuilding.X)) + (troop.Y-float64(targetBuilding.Y))*(troop.Y-float64(targetBuilding.Y)))
			if distanceFromTarget <= float64(troop.Range) {
				liveBuildings[troop.TargetID].CurrentHP -= troop.Damage
				event := dtos.BattleEventDTO{
					Tick:     tick,
					Action:   "attack",
					EntityID: troop.ID,
					TargetID: troop.TargetID,
				}
				battleLog = append(battleLog, event)
				if liveBuildings[troop.TargetID].CurrentHP <= 0 {
					event := dtos.BattleEventDTO{
						Tick:     tick,
						Action:   "delete",
						EntityID: troop.TargetID,
					}
					battleLog = append(battleLog, event)
					delete(liveBuildings, troop.TargetID)
				}
			} else {
				diffX := targetBuilding.X - int(troop.X)
				diffY := targetBuilding.Y - int(troop.Y)

				if diffX != 0 && diffY != 0 {
					troop.X += (float64(diffX) / math.Abs(float64(diffX))) * math.Min(math.Abs(float64(troop.Speed)), math.Abs(float64(diffX)))
					troop.Y += (float64(diffY) / math.Abs(float64(diffY))) * math.Min(math.Abs(float64(troop.Speed)), math.Abs(float64(diffY)))

					event := dtos.BattleEventDTO{
						Tick:     tick,
						Action:   "move",
						EntityID: troop.ID,
						TargetID: troop.TargetID,
						X:        int(troop.X),
						Y:        int(troop.Y),
					}
					battleLog = append(battleLog, event)
				} else if diffX != 0 {
					troop.X += (float64(diffX) / math.Abs(float64(diffX))) * math.Min(math.Abs(float64(troop.Speed)), math.Abs(float64(diffX)))

					event := dtos.BattleEventDTO{
						Tick:     tick,
						Action:   "move",
						EntityID: troop.ID,
						TargetID: troop.TargetID,
						X:        int(troop.X),
						Y:        int(troop.Y),
					}
					battleLog = append(battleLog, event)
				} else if diffY != 0 {
					troop.Y += (float64(diffY) / math.Abs(float64(diffY))) * math.Min(math.Abs(float64(troop.Speed)), math.Abs(float64(diffY)))

					event := dtos.BattleEventDTO{
						Tick:     tick,
						Action:   "move",
						EntityID: troop.ID,
						TargetID: troop.TargetID,
						X:        int(troop.X),
						Y:        int(troop.Y),
					}
					battleLog = append(battleLog, event)
				}
			}
		}

		for _, building := range liveBuildings {
			minDistance := math.MaxFloat64
			minDisID := ""
			if building.Damage == 0 {
				continue
			}
			for _, troop := range liveTroops {
				if building.TargetType == "ground" && troop.Airborne {
					continue
				}
				if building.TargetType == "air" && !troop.Airborne {
					continue
				}

				distance := math.Sqrt((troop.X-float64(building.X))*(troop.X-float64(building.X)) + (troop.Y-float64(building.Y))*(troop.Y-float64(building.Y)))
				if distance < minDistance {
					minDistance = distance
					minDisID = troop.ID
				}
			}
			building.TargetID = minDisID
		}

		for _, building := range liveBuildings {
			targetTroop := liveTroops[building.TargetID]
			if targetTroop == nil {
				continue
			}
			distanceFromTarget := math.Sqrt((targetTroop.X-float64(building.X))*(targetTroop.X-float64(building.X)) + (targetTroop.Y-float64(building.Y))*(targetTroop.Y-float64(building.Y)))
			if distanceFromTarget <= float64(building.Range) {
				event := dtos.BattleEventDTO{
					Tick:     tick,
					Action:   "attack",
					EntityID: building.ID,
					TargetID: building.TargetID,
				}
				battleLog = append(battleLog, event)

				if building.SingleTarget {
					liveTroops[targetTroop.ID].CurrentHP -= building.Damage
					if liveTroops[targetTroop.ID].CurrentHP <= 0 {
						deleteEvent := dtos.BattleEventDTO{
							Tick:     tick,
							Action:   "delete",
							EntityID: targetTroop.ID,
						}
						battleLog = append(battleLog, deleteEvent)
						delete(liveTroops, targetTroop.ID)
					}
				} else {
					splashRadius := building.SplashRadius
					impactX := targetTroop.X
					impactY := targetTroop.Y

					for troopID, troop := range liveTroops {
						distFromImpact := math.Sqrt((troop.X-impactX)*(troop.X-impactX) + (troop.Y-impactY)*(troop.Y-impactY))

						if distFromImpact <= splashRadius {
							liveTroops[troopID].CurrentHP -= building.Damage

							if liveTroops[troopID].CurrentHP <= 0 {
								deleteEvent := dtos.BattleEventDTO{
									Tick:     tick,
									Action:   "delete",
									EntityID: troopID,
								}
								battleLog = append(battleLog, deleteEvent)
								delete(liveTroops, troopID)
							}
						}
					}
				}
			}
		}
	}

	leftBuildings := len(liveBuildings)
	damagePercent := ((totalBuildings - leftBuildings) * 100) / totalBuildings
	lootedGold := (damagePercent * gold) / 100
	lootedElixir := (damagePercent * elixir) / 100

	var winnerID string
	if damagePercent > 50 {
		winnerID = userID
	} else {
		winnerID = targetUserID
	}

	_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold - $1, elixir = elixir - $2, last_attacked_at = NOW() WHERE user_id = $3", lootedGold, lootedElixir, targetUserID)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Error deducting resources")
	}

	_, err = tx.Exec(ctx, "UPDATE villages SET gold = gold + $1, elixir = elixir + $2 WHERE user_id = $3", lootedGold, lootedElixir, userID)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Error adding resources")
	}

	for key, troop := range troopQuantity {
		_, err = tx.Exec(ctx, "UPDATE village_troops SET quantity = quantity - $1 WHERE village_id = $2 AND troop_id = $3", troop, villageID, key)
		if err != nil {
			return 0, 0, 0, nil, errors.New("Error updating troops")
		}
	}

	logJSON, err := json.Marshal(battleLog)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Failed to format battle log")
	}

	_, err = tx.Exec(ctx, "INSERT INTO battles (attacker_id, defender_id, winner_id, damage_percent, battle_log) VALUES ($1, $2, $3, $4, $5)", userID, targetUserID, winnerID, damagePercent, logJSON)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Failed to update battle log in database")
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, 0, nil, err
	}

	return damagePercent, lootedGold, lootedElixir, battleLog, nil
}

func GetBattleHistory(userID string) ([]dtos.BattleRecordDTO, error) {
	ctx := context.Background()

	query := `
		SELECT id, attacker_id, defender_id, winner_id, damage_percent, occurred_at
		FROM battles
		WHERE attacker_id = $1 OR defender_id = $1
		ORDER BY occurred_at DESC
		LIMIT 20
	`
	rows, err := db.Conn.Query(ctx, query, userID)

	if err != nil {
		return nil, errors.New("Error fetching battle history")
	}
	defer rows.Close()

	battles, err := pgx.CollectRows(rows, pgx.RowToStructByName[dtos.BattleRecordDTO])
	if err != nil {
		return nil, errors.New("Error parsing battle history")
	}

	if battles == nil {
		battles = []dtos.BattleRecordDTO{}
	}

	return battles, nil
}
