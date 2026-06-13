export const getGameConfigs = async () => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error("No token found. User is not logged in.");
  }

  const response = await fetch('/api/game/configs', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) throw new Error("Failed to fetch game configs");
  
  const data = await response.json();
  return {
    buildings: data.buildings,
    troops: data.troops
  };
};