export const loginUser = async (username: string, password: string) => {
  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const errorText = await response.text(); 
      throw new Error(errorText || 'Failed to log in');
    }

    const data = await response.json();
    return data.token;
    
  } catch (error) {
    console.error('Login Service Error:', error);
    throw error;
  }
};

export const registerUser = async (username: string, password: string) => {
  try {
    const response = await fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const errorText = await response.text(); 
      throw new Error(errorText || 'Failed to register user');
    }
    
  } catch (error) {
    console.error('Register Service Error:', error);
    throw error;
  }
};