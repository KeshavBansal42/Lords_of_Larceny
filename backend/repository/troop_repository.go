package repository

import (
	"context"
	"errors"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/jackc/pgx/v5"
)

func TrainTroops(userID string, troopsToTrain map[int]int) error {
	ctx := context.Background()

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var villageID string
	var thLevel int
	err = tx.QueryRow(ctx, "SELECT id, town_hall_level FROM villages WHERE user_id = $1 FOR UPDATE", userID).Scan(&villageID, &thLevel)
	if err != nil {
		return errors.New("Error fetching village details")
	}

	var maxCapacity int
	err = tx.QueryRow(ctx, `
        SELECT COALESCE(SUM(ac.total_housing_space), 0) 
        FROM village_buildings vb 
        JOIN army_camp_configs ac ON vb.building_name = ac.name AND vb.level = ac.level
        WHERE vb.village_id = $1 AND vb.building_name = 'Army Camp'
    `, villageID).Scan(&maxCapacity)
	if err != nil {
		return errors.New("Error calculating max army capacity")
	}

	var currentSpace int
	err = tx.QueryRow(ctx, `
        SELECT COALESCE(SUM(vt.quantity * tc.housing_space), 0) 
        FROM village_troops vt 
        JOIN troop_configs tc ON vt.troop_id = tc.id 
        WHERE vt.village_id = $1
    `, villageID).Scan(&currentSpace)
	if err != nil {
		return errors.New("Error calculating current army space")
	}

	requestedSpace := 0
	for troopID, quantity := range troopsToTrain {
		var housingSpace int
		var minThLevel int
		err = tx.QueryRow(ctx, "SELECT housing_space, min_thlevel FROM troop_configs WHERE id = $1", troopID).Scan(&housingSpace, &minThLevel)
		if err != nil {
			return errors.New("Invalid troop ID")
		}
		if thLevel < minThLevel {
			return errors.New("Minimum town hall level requirement not met for requested troop")
		}
		requestedSpace += (housingSpace * quantity)
	}

	if currentSpace+requestedSpace > maxCapacity {
		return errors.New("Insufficient army camp capacity")
	}

	for troopID, quantity := range troopsToTrain {
		_, err = tx.Exec(ctx, `
            INSERT INTO village_troops (village_id, troop_id, quantity) 
            VALUES ($1, $2, $3) 
            ON CONFLICT (village_id, troop_id) 
            DO UPDATE SET quantity = village_troops.quantity + EXCLUDED.quantity
        `, villageID, troopID, quantity)
		if err != nil {
			return errors.New("Error training troops")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func GetAllVillageTroops(villageID string) ([]dtos.TroopResponseFromDBDTO, error) {
	rows, err := db.Conn.Query(context.Background(), "SELECT troop_id, quantity FROM village_troops WHERE village_id = $1", villageID)

	if err != nil {
		return nil, errors.New("Error fetching troops")
	}
	defer rows.Close()

	troops, err := pgx.CollectRows(rows, pgx.RowToStructByName[dtos.TroopResponseFromDBDTO])
	if err != nil {
		return nil, errors.New("Error parsing troops")
	}

	return troops, nil
}
