import type { VillageStats } from '../types';
import type { Building } from '../types';

export const getVillageStats = async (): Promise<VillageStats> => {
  const token = localStorage.getItem('token');
  
  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to fetch village data');
  }

  return await response.json(); 
};

export const getVillageBuildings = async (): Promise<Building[]> => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/buildings', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) throw new Error("Failed to fetch buildings");
  
  const data = await response.json();
  return data.buildings;
};

export const buildBuilding = async (buildingName: string, x: number, y: number) => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }
  
  const response = await fetch('/api/village/buildings/build', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ building_name: buildingName, x, y }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to build');
  }

  return await response.json(); 
};

export const getVillageTroops = async () => {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/village/troops', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) throw new Error("Failed to fetch troops");
  
  const data = await response.json();
  return data.troops || [];
};

export const trainTroops = async (troopsToTrain: Record<number, number>) => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/troops/train', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ troopstotrain: troopsToTrain }), 
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to train troops');
  }

  return await response.json(); 
};

export const collectGold = async () => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/collect/gold', {
    method: 'PUT',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!response.ok) throw new Error("Failed to collect gold");
  return await response.json();
};

export const collectElixir = async () => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/collect/elixir', {
    method: 'PUT',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!response.ok) throw new Error("Failed to collect elixir");
  return await response.json();
};

export const upgradeBuilding = async (x: number, y: number) => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/buildings/upgrade', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ x, y }),
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to upgrade building');
  }
  return await response.json();
};

export const moveBuilding = async (oldX: number, oldY: number, newX: number, newY: number) => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/village/buildings/move', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ oldx: oldX, oldy: oldY, newx: newX, newy: newY }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to move building');
  }

  return await response.json(); 
};

export const getBattleHistory = async () => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/battle/history', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) throw new Error("Failed to fetch battle history");
  
  const data = await response.json();
  return data.battles; 
};