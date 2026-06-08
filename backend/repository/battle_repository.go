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

func Matchmake(userID int) (int, error) {
	ctx := context.Background()

	var villageID int
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
		return 0, errors.New("Error finding user info.")
	}

	if troopCount == 0 {
		return 0, errors.New("You need an army to attack")
	}

	rows, err := db.Conn.Query(ctx, "SELECT user_id FROM villages WHERE id != $1 AND town_hall_level BETWEEN $2 AND $3 AND last_attacked_at < NOW() - INTERVAL '6 hours' LIMIT 50", villageID, (townHallLevel - 1), (townHallLevel + 1))

	if err != nil {
		return 0, errors.New("Error fetching users")
	}
	defer rows.Close()

	villages, err := pgx.CollectRows(rows, pgx.RowToStructByName[dtos.MatchmakeResponseFromDBDTO])
	if err != nil {
		return 0, errors.New("Failed to parse villages")
	}

	if len(villages) == 0 {
		return 0, errors.New("No worthy opponents exist.")
	}

	randomIndex := rand.Intn(len(villages))
	selectedOpponent := villages[randomIndex]

	return selectedOpponent.UserID, nil
}

func Populate(userID int, liveBuildings *map[string]*models.LiveBuilding, liveTroops *map[string]*models.LiveTroop, buildings []dtos.BuildingResponseFromDBDTO, drops []dtos.TroopDropDTO) error {
	buildingConfigsArray, troopConfigsArray, err := GetGameConfigs()
	if err != nil {
		return errors.New("Error fetching config files")
	}

	buildingConfigs := make(map[int]models.BuildingConfig)
	troopConfigs := make(map[int]models.TroopConfig)

	for _, bconfig := range buildingConfigsArray {
		key := bconfig.ID
		buildingConfigs[key] = bconfig
	}
	for _, tconfig := range troopConfigsArray {
		key := tconfig.ID
		troopConfigs[key] = tconfig
	}

	for i, building := range buildings {
		key := fmt.Sprintf("building_%v", (i + 1))
		(*liveBuildings)[key] = &models.LiveBuilding{
			ID:         key,
			BuildingID: building.BuildingId,
			X:          building.X,
			Y:          building.Y,
			MaxHP:      buildingConfigs[building.BuildingId].HitPoints,
			CurrentHP:  buildingConfigs[building.BuildingId].HitPoints,
			Damage:     buildingConfigs[building.BuildingId].Damage,
			TargetID:   "",
			Range:      buildingConfigs[building.BuildingId].Range,
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
			TargetID:  "",
		}
	}

	return nil
}

func Battle(userID int, targetUserID int, drops []dtos.TroopDropDTO) (int, int, int, []dtos.BattleEventDTO, error) {
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
			if distanceFromTarget < float64(building.Range) {
				liveTroops[targetTroop.ID].CurrentHP -= building.Damage

				event := dtos.BattleEventDTO{
					Tick:     tick,
					Action:   "attack",
					EntityID: building.ID,
					TargetID: building.TargetID,
				}
				battleLog = append(battleLog, event)

				if liveTroops[targetTroop.ID].CurrentHP <= 0 {
					event := dtos.BattleEventDTO{
						Tick:     tick,
						Action:   "delete",
						EntityID: building.TargetID,
					}
					battleLog = append(battleLog, event)
					delete(liveTroops, building.TargetID)
				}
			}
		}
	}
	leftBuildings := len(liveBuildings)
	damagePercent := ((totalBuildings - leftBuildings) * 100) / totalBuildings
	lootedGold := (damagePercent * gold) / 100
	lootedElixir := (damagePercent * elixir) / 100

	var winnerID int
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
	_, err = tx.Exec(ctx, "INSERT INTO battles (attacker_id, defender_id, winner_id, battle_log) VALUES ($1, $2, $3, $4)", userID, targetUserID, winnerID, logJSON)
	if err != nil {
		return 0, 0, 0, nil, errors.New("Failed to update battle log in database")
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, 0, nil, err
	}

	return damagePercent, lootedGold, lootedElixir, battleLog, nil
}
