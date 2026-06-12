import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { loginUser } from '../services/auth';

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault(); 
    setErrorMsg("");
    
    try {
      const token = await loginUser(username, password);
      
      localStorage.setItem('token', token);
      
      navigate('/village'); 
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '400px', margin: '0 auto' }}>
      <h2>Login to Lords of Larceny</h2>
      
      {errorMsg && <p style={{ color: 'red' }}>{errorMsg}</p>}
      
      <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
        <input 
          type="text" 
          placeholder="Username" 
          value={username} 
          onChange={(e) => setUsername(e.target.value)} 
        />
        <input 
          type="password" 
          placeholder="Password" 
          value={password} 
          onChange={(e) => setPassword(e.target.value)} 
        />
        <button type="submit">Login</button>
      </form>

      <button onClick={() => navigate('/register')} style={{ marginTop: '20px' }}>
        Don't have an account? Register
      </button>
    </div>
  );
}