import { authApi } from '@/apis/api';
import { useState } from 'react';

const useAuth = () => {
  // Mock authentication state - in a real app, this would come from your auth provider
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const user = isAuthenticated ? { name: 'Demo User', email: 'user@example.com' } : null;
  
  // const login = () => setIsAuthenticated(true);
  const login = async (email: string, password: string) => {
    await authApi.login({
      email,
      password
    });
  };
  const signup = async (email: string, password: string, name: string) => {
    await authApi.signup({
      email,
      password,
      name
    });
  };
  const logout = () => setIsAuthenticated(false);
  
  return { isAuthenticated, user, login, signup, logout };
};

export default useAuth;
