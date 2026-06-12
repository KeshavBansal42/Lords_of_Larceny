import { create } from 'zustand';
import type { GameState } from '../types';

export const useGameStore = create<GameState>((set) => ({
  townHallLevel: 0,
  gold: 0,
  elixir: 0,
  buildingConfigs: [],
  troopConfigs: [],
  army: {},

  setVillageStats: (thLevel, gold, elixir) => set({ townHallLevel: thLevel, gold, elixir }),
  setBuildingConfigs: (configs) => set({ buildingConfigs: configs }), 
  setTroopConfigs: (configs) => set({ troopConfigs: configs }), 
  
  setArmy: (troops) => set(() => {
    const newArmy: Record<number, number> = {};
    troops.forEach(t => { newArmy[t.troopid] = t.quantity; });
    return { army: newArmy };
  }),

  addTroopsToArmy: (trainedTroops) => set((state) => {
    const newArmy = { ...state.army };
    for (const [id, qty] of Object.entries(trainedTroops)) {
      const numId = Number(id);
      newArmy[numId] = (newArmy[numId] || 0) + qty;
    }
    return { army: newArmy };
  }),
  
  spendGold: (amount) => set((state) => ({ gold: state.gold - amount })),
  spendElixir: (amount) => set((state) => ({ elixir: state.elixir - amount })),
}));