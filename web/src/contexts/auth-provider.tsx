import { userApi, authApi } from '@/apis/api';
import type { GetUserResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { AuthContext } from './auth-context';
import { useState, useEffect, useMemo, type ReactNode } from 'react';
import { PATH } from '@/const/path';
import { navigate } from 'wouter/use-browser-location';

export const AuthProvider = ({ children }: { children: ReactNode; }) => {
  const [user, setUser] = useState<GetUserResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const isAuthenticated = useMemo(() => !!user, [user]);

  useEffect(() => {
    const initializeAuth = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        try {
          const user = await userApi.getCurrentUser();
          setUser(user);
        } catch {
          localStorage.removeItem('token');
        }
      }
      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const setToken = async (token: string) => {
    localStorage.setItem('token', token);
    const user = await userApi.getCurrentUser();
    setUser(user);
    navigate(PATH.HOME);
  };

  const login = async (email: string, password: string) => {
    const resp = await authApi.login({
      email,
      password,
    });
    if (resp.token) {
      await setToken(resp.token);
    }
  };

  const signup = async (email: string, password: string, name: string) => {
    const resp = await authApi.signup({
      email,
      password,
      name,
    });
    if (resp.token) {
      await setToken(resp.token);
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  return (
    <AuthContext value={{ isAuthenticated, user, login, signup, logout, isLoading }}>
      {children}
    </AuthContext>
  );
};
