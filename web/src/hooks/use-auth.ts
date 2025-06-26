import { authApi, userApi } from '@/apis/api';
import type { GetUserResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { useMemo, useState } from 'react';

const useAuth = () => {
  const [user, setUser] = useState<GetUserResponse | null>(null);
  const isAuthenticated = useMemo(() => !!user, [user]);

  const setToken = async (token: string) => {
    localStorage.setItem('token', token);
    console.log('isAuthenticated', isAuthenticated);
    const user = await userApi.getCurrentUser();
    setUser(user);
    console.log('isAuthenticated', isAuthenticated);
  };

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
    setUser(null);
  };

  return { isAuthenticated, user, login, signup, logout };
};

export default useAuth;
