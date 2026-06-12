export interface VillageStats {
  town_hall_level: number;
  gold: number;
  elixir: number;
}

export interface BuildingConfig {
  name: string;
  level: number;
  build_cost: number;
  build_resource_type: string;
  size: number;
  min_thlevel: number;
  capacity: number;
}

export interface TroopConfig {
  id: number;
  name: string;
  level: number;
  housing_space: number;
  damage: number;
  hit_points: number;
  min_thlevel: number;
}

export interface Building {
  building_name: string;
  level: number;
  x: number;
  y: number;
  status: string;
}

export interface GameState {
  townHallLevel: number;
  gold: number;
  elixir: number;
  buildingConfigs: BuildingConfig[];
  troopConfigs: TroopConfig[];
  army: Record<number, number>;
  setVillageStats: (thLevel: number, gold: number, elixir: number) => void;
  setBuildingConfigs: (configs: BuildingConfig[]) => void;
  setTroopConfigs: (configs: TroopConfig[]) => void;
  setArmy: (troops: { troopid: number; quantity: number }[]) => void;
  addTroopsToArmy: (trainedTroops: Record<number, number>) => void;
  
  spendGold: (amount: number) => void;
  spendElixir: (amount: number) => void;
}