import { authApi } from '@/apis/api';
import { useState } from 'react';

const useAuth = () => {
  // Mock authentication state - in a real app, this would come from your auth provider
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const user = isAuthenticated ? { name: 'Demo User', email: 'user@example.com' } : null;

  function setToken(token: string) {
    localStorage.setItem('token', token);
    setIsAuthenticated(true);
  }

  const login = async (email: string, password: string) => {
    const resp = await authApi.login({
      email,
      password,
    });
    if (resp.token) {
      setToken(resp.token);
    }
  };
  const signup = async (email: string, password: string, name: string) => {
    const resp = await authApi.signup({
      email,
      password,
      name,
    });
    if (resp.token) {
      setToken(resp.token);
    }
  };
  const logout = () => {
    localStorage.removeItem('token');
    setIsAuthenticated(false);
  };
  
  return { isAuthenticated, user, login, signup, logout };
};

export default useAuth;
