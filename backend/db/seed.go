package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDatabase(conn *pgxpool.Pool) {
	ctx := context.Background()

	// 0. Insert system ghost user for deleted accounts to preserve battle history
	_, err := conn.Exec(ctx, `
		INSERT INTO users (id, username, password_hash) 
		VALUES ('00000000-0000-0000-0000-000000000000', 'Deleted User', 'system_ghost_hash')
		ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Println("Failed to seed system ghost user: ", err)
	}

	var count int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM troop_configs").Scan(&count)
	if err != nil {
		log.Println("Failed to count troops: ", err)
	}

	if count == 0 {
		troopQuery := `
			INSERT INTO troop_configs (name, level, hit_points, min_thlevel, housing_space, damage, range, speed, airborne) VALUES 
			('Barbarian', 1, 45, 1, 1, 9, 1, 2, FALSE),
			('Archer', 1, 22, 1, 1, 8, 4, 2, FALSE),
			('Goblin', 1, 25, 2, 1, 11, 1, 3, FALSE),
			('Giant', 1, 400, 3, 4, 12, 1, 1, FALSE),
			('Wall Breaker', 1, 20, 4, 2, 10, 1, 2, FALSE);
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
		// 1. Seed base building configurations
		buildingQuery := `
			INSERT INTO building_configs (name, build_resource_type, size) VALUES 
			('Town Hall', 'gold', 6),
			('Cannon', 'gold', 4),
			('Archer Tower', 'gold', 4),
			('Mortar', 'gold', 4),
			('Gold Mine', 'elixir', 4),
			('Elixir Collector', 'gold', 4),
			('Army Camp', 'elixir', 5);
		`
		_, err = conn.Exec(ctx, buildingQuery)
		if err != nil {
			log.Println("Failed to seed base building configs: ", err)
		}

		// 2. Seed Building Level Info (Common Stats)
		levelInfoQuery := `
			INSERT INTO building_level_info (name, level, hit_points, build_cost, min_thlevel, build_time_seconds) VALUES 
			('Cannon', 1, 300, 250, 1, 60),
			('Cannon', 2, 340, 1000, 2, 120),
			('Cannon', 3, 400, 4000, 3, 300),
			('Cannon', 4, 450, 16000, 4, 600),
			('Archer Tower', 1, 380, 1000, 1, 60),
			('Archer Tower', 2, 420, 2000, 2, 120),
			('Archer Tower', 3, 460, 5000, 3, 300),
			('Archer Tower', 4, 500, 20000, 4, 600),
			('Mortar', 1, 400, 8000, 3, 300),
			('Mortar', 2, 450, 32000, 4, 600),
			('Gold Mine', 1, 75, 150, 1, 60),
			('Gold Mine', 2, 150, 300, 2, 120),
			('Gold Mine', 3, 300, 700, 3, 300),
			('Gold Mine', 4, 400, 1400, 4, 600),
			('Elixir Collector', 1, 75, 150, 1, 60),
			('Elixir Collector', 2, 150, 300, 2, 120),
			('Elixir Collector', 3, 300, 700, 3, 300),
			('Elixir Collector', 4, 400, 1400, 4, 600),
			('Town Hall', 1, 400, 0, 1, 0),
			('Town Hall', 2, 800, 1000, 1, 60),
			('Town Hall', 3, 1600, 4000, 1, 300),
			('Town Hall', 4, 2000, 25000, 1, 1200),
			('Army Camp', 1, 100, 200, 1, 60),
			('Army Camp', 2, 150, 2000, 2, 120),
			('Army Camp', 3, 200, 10000, 3, 300),
			('Army Camp', 4, 250, 100000, 4, 600);
		`
		_, err = conn.Exec(ctx, levelInfoQuery)
		if err != nil {
			log.Println("Failed to seed building level info: ", err)
		}

		// 3. Seed Defense Configurations (Stripped down)
		defenseQuery := `
			INSERT INTO defense_configs (name, level, damage, range, single_target, splash_radius, target_type) VALUES 
			('Cannon', 1, 7, 4, TRUE, 0, 'ground'),
			('Cannon', 2, 11, 4, TRUE, 0, 'ground'),
			('Cannon', 3, 15, 4, TRUE, 0, 'ground'),
			('Cannon', 4, 19, 4, TRUE, 0, 'ground'),
			('Archer Tower', 1, 11, 4, TRUE, 0, 'both'),
			('Archer Tower', 2, 15, 4, TRUE, 0, 'both'),
			('Archer Tower', 3, 19, 4, TRUE, 0, 'both'),
			('Archer Tower', 4, 25, 4, TRUE, 0, 'both'),
			('Mortar', 1, 4, 6, FALSE, 2.5, 'ground'),
			('Mortar', 2, 5, 6, FALSE, 2.5, 'ground'),
			('Mortar', 3, 6, 6, FALSE, 2.5, 'ground'),
			('Mortar', 4, 7, 6, FALSE, 2.5, 'ground');
		`
		_, err = conn.Exec(ctx, defenseQuery)
		if err != nil {
			log.Println("Failed to seed defense configs: ", err)
		}

		// 4. Seed Resource Generators (Stripped down)
		resGenQuery := `
			INSERT INTO resource_gen_configs (name, level, production_per_min, capacity, resource_type) VALUES 
			('Gold Mine', 1, 3, 1000, 'gold'),
			('Gold Mine', 2, 6, 2000, 'gold'),
			('Gold Mine', 3, 10, 3000, 'gold'),
			('Gold Mine', 4, 13, 5000, 'gold'),
			('Elixir Collector', 1, 3, 1000, 'elixir'),
			('Elixir Collector', 2, 6, 2000, 'elixir'),
			('Elixir Collector', 3, 10, 3000, 'elixir'),
			('Elixir Collector', 4, 13, 5000, 'elixir');
		`
		_, err = conn.Exec(ctx, resGenQuery)
		if err != nil {
			log.Println("Failed to seed resource generators: ", err)
		}

		// 5. Seed Resource Storages (Stripped down)
		storageQuery := `
			INSERT INTO resource_storage_configs (name, level, resource_type, storage_capacity) VALUES 
			('Town Hall', 1, 'both', 1000),
			('Town Hall', 2, 'both', 2500),
			('Town Hall', 3, 'both', 10000),
			('Town Hall', 4, 'both', 50000);
		`
		_, err = conn.Exec(ctx, storageQuery)
		if err != nil {
			log.Println("Failed to seed resource storages: ", err)
		}

		// 6. Seed Army Camps (Stripped down)
		armyCampQuery := `
			INSERT INTO army_camp_configs (name, level, total_housing_space) VALUES 
			('Army Camp', 1, 20),
			('Army Camp', 2, 30),
			('Army Camp', 3, 35),
			('Army Camp', 4, 40);
		`
		_, err = conn.Exec(ctx, armyCampQuery)
		if err != nil {
			log.Println("Failed to seed army camps: ", err)
		}
	}

	// 7. Seed Fake Users for Matchmaking
	var botCount int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username LIKE 'bot_%'").Scan(&botCount)
	if err != nil {
		log.Println("Failed to count bot users: ", err)
	}

	if botCount == 0 {
		fakeUsersQuery := `
		DO $$
		DECLARE
			bot1_id uuid := gen_random_uuid();
			bot2_id uuid := gen_random_uuid();
			bot3_id uuid := gen_random_uuid();
			village1_id uuid := gen_random_uuid();
			village2_id uuid := gen_random_uuid();
			village3_id uuid := gen_random_uuid();
		BEGIN
			INSERT INTO users (id, username, password_hash) VALUES 
			(bot1_id, 'bot_goblin_king', 'dummy_hash'),
			(bot2_id, 'bot_archer_queen', 'dummy_hash'),
			(bot3_id, 'bot_barbarian_boss', 'dummy_hash');

			INSERT INTO villages (id, user_id, town_hall_level, gold, elixir) VALUES
			(village1_id, bot1_id, 1, 1500, 1500),
			(village2_id, bot2_id, 2, 3000, 3000),
			(village3_id, bot3_id, 3, 8000, 8000);

			-- Bot 1: Basic setup
			INSERT INTO village_buildings (village_id, building_name, level, x, y) VALUES
			(village1_id, 'Town Hall', 1, 16, 16),
			(village1_id, 'Cannon', 1, 10, 10),
			(village1_id, 'Gold Mine', 1, 10, 14);

			-- Bot 2: Moderate setup
			INSERT INTO village_buildings (village_id, building_name, level, x, y) VALUES
			(village2_id, 'Town Hall', 2, 16, 16),
			(village2_id, 'Cannon', 2, 12, 12),
			(village2_id, 'Archer Tower', 1, 20, 12),
			(village2_id, 'Gold Mine', 2, 10, 20);

			-- Bot 3: Advanced setup
			INSERT INTO village_buildings (village_id, building_name, level, x, y) VALUES
			(village3_id, 'Town Hall', 3, 16, 16),
			(village3_id, 'Cannon', 3, 12, 12),
			(village3_id, 'Archer Tower', 2, 20, 12),
			(village3_id, 'Mortar', 1, 16, 12),
			(village3_id, 'Elixir Collector', 3, 20, 20);
		END $$;
		`
		_, err = conn.Exec(ctx, fakeUsersQuery)
		if err != nil {
			log.Println("Failed to seed fake users: ", err)
		}
	}
}
