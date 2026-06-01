package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func SeedDatabase(conn *pgx.Conn) {
	ctx := context.Background()
	var count int

	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM troop_configs").Scan(&count)
	if err != nil {
		log.Println("Failed to count troops: ", err)
	}

	if count == 0 {
		troopQuery := `
			INSERT INTO troop_configs (name, level, hit_points, damage) VALUES 
			('Barbarian', 1, 45, 9),
			('Archer', 1, 22, 8),
			('Goblin', 1, 25, 11),
			('Giant', 1, 400, 12),
			('Wall Breaker', 1, 20, 10);
		`
		_, err = conn.Exec(ctx, troopQuery)
		if err != nil {
			log.Println("Failed to seed troops: ", err)
		}
	}

	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM building_configs").Scan(&count)
	if err != nil {
		log.Println("Failed to count buildings: ", err)
	}

	if count == 0 {
		buildingQuery := `
			INSERT INTO building_configs (name, level, hit_points, damage, build_cost) VALUES 
			-- Town Hall
			('Town Hall', 1, 400, 0, 0),
			('Town Hall', 2, 800, 0, 1000),
			('Town Hall', 3, 1600, 0, 4000),
			('Town Hall', 4, 2000, 0, 25000),
			
			-- Cannon
			('Cannon', 1, 300, 7, 250),
			('Cannon', 2, 340, 11, 1000),
			('Cannon', 3, 400, 15, 4000),
			('Cannon', 4, 450, 19, 16000),
			
			-- Archer Tower
			('Archer Tower', 1, 380, 11, 1000),
			('Archer Tower', 2, 420, 15, 2000),
			('Archer Tower', 3, 460, 19, 5000),
			('Archer Tower', 4, 500, 25, 20000),
			
			-- Mortar
			('Mortar', 1, 400, 4, 8000),
			('Mortar', 2, 450, 5, 32000),
			('Mortar', 3, 500, 6, 120000),
			('Mortar', 4, 550, 7, 180000);
		`
		_, err = conn.Exec(ctx, buildingQuery)
		if err != nil {
			log.Println("Failed to seed buildings: ", err)
		}
	}
}
