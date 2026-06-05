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
			INSERT INTO troop_configs (name, level, hit_points, damage, min_thlevel, housing_space) VALUES 
			('Barbarian', 1, 45, 9, 1, 1),
			('Archer', 1, 22, 8, 1, 1),
			('Goblin', 1, 25, 11, 2, 1),
			('Giant', 1, 400, 12, 3, 4),
			('Wall Breaker', 1, 20, 10, 4, 2);
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
			INSERT INTO building_configs (name, level, hit_points, damage, build_cost, build_resource_type, production_per_min, capacity, size, min_thlevel) VALUES 
			-- Town Hall
			('Town Hall', 1, 400, 0, 0, 'gold', 0, 1000, 4, 1),
			('Town Hall', 2, 800, 0, 1000, 'gold', 0, 2500, 4, 1),
			('Town Hall', 3, 1600, 0, 4000, 'gold', 0, 10000, 4, 1),
			('Town Hall', 4, 2000, 0, 25000, 'gold', 0, 50000, 4, 1),

			-- Cannon
			('Cannon', 1, 300, 7, 250, 'gold', 0, 0, 2, 1),
			('Cannon', 2, 340, 11, 1000, 'gold', 0, 0, 2, 2),
			('Cannon', 3, 400, 15, 4000, 'gold', 0, 0, 2, 3),
			('Cannon', 4, 450, 19, 16000, 'gold', 0, 0, 2, 4),

			-- Archer Tower
			('Archer Tower', 1, 380, 11, 1000, 'gold', 0, 0, 2, 1),
			('Archer Tower', 2, 420, 15, 2000, 'gold', 0, 0, 2, 2),
			('Archer Tower', 3, 460, 19, 5000, 'gold', 0, 0, 2, 3),
			('Archer Tower', 4, 500, 25, 20000, 'gold', 0, 0, 2, 4),

			-- Mortar
			('Mortar', 1, 400, 4, 8000, 'gold', 0, 0, 2, 3),
			('Mortar', 2, 450, 5, 32000, 'gold', 0, 0, 2, 3),
			('Mortar', 3, 500, 6, 120000, 'gold', 0, 0, 2, 4),
			('Mortar', 4, 550, 7, 180000, 'gold', 0, 0, 2, 4),

			-- Gold Mine
			('Gold Mine', 1, 75, 0, 150, 'elixir', 3, 1000, 2, 1),
			('Gold Mine', 2, 150, 0, 300, 'elixir', 6, 2000, 2, 2),
			('Gold Mine', 3, 300, 0, 700, 'elixir', 10, 3000, 2, 3),
			('Gold Mine', 4, 400, 0, 1400, 'elixir', 13, 5000, 2, 4),

			-- Elixir Collector
			('Elixir Collector', 1, 75, 0, 150, 'gold', 3, 1000, 2, 1),
			('Elixir Collector', 2, 150, 0, 300, 'gold', 6, 2000, 2, 2),
			('Elixir Collector', 3, 300, 0, 700, 'gold', 10, 3000, 2, 3),
			('Elixir Collector', 4, 400, 0, 1400, 'gold', 13, 5000, 2, 4),

			--Army Camp
			('Army Camp', 1, 100, 0, 200, 'elixir', 0, 20, 3, 1),
			('Army Camp', 2, 150, 0, 2000, 'elixir', 0, 30, 3, 2),
			('Army Camp', 3, 200, 0, 10000, 'elixir', 0, 35, 3, 3),
			('Army Camp', 4, 250, 0, 100000, 'elixir', 0, 40, 3, 4);
		`
		_, err = conn.Exec(ctx, buildingQuery)
		if err != nil {
			log.Println("Failed to seed buildings: ", err)
		}
	}
}
