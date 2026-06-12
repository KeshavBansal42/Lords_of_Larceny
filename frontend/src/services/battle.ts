import type { TroopDrop, BattleResult } from '../types';

export const findMatch = async () => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('/api/battle/matchmake', {
    method: 'GET', 
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to find an opponent');
  }

  return await response.json();
};

export const scoutVillage = async (userId: string) => {
  const token = localStorage.getItem('token');
  
  const response = await fetch(`/api/village/${userId}/scout`, {
    method: 'GET', 
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to scout village');
  }

  return await response.json();
};

export const attackVillage = async (targetUserId: string, drops: TroopDrop[]): Promise<BattleResult> => {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/battle/attack', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ target_user_id: targetUserId, drops }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Battle failed');
  }

  return await response.json(); 
};