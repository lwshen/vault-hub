import { userApi, authApi } from '@/apis/api';
import type { GetUserResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { AuthContext } from './auth-context';
import { useState, useEffect, useMemo, type ReactNode, useCallback } from 'react';
import { PATH } from '@/const/path';
import { navigate } from 'wouter/use-browser-location';

export const AuthProvider = ({ children }: { children: ReactNode; }) => {
  const [user, setUser] = useState<GetUserResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const isAuthenticated = useMemo(() => !!user, [user]);

  const setToken = useCallback(async (token: string) => {
    localStorage.setItem('token', token);
    const user = await userApi.getCurrentUser();
    setUser(user);
    navigate(PATH.HOME);
  }, []);

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

  const login = useCallback(
    async (email: string, password: string) => {
      const resp = await authApi.login({
        email,
        password,
      });
      if (resp.token) {
        await setToken(resp.token);
      }
    },
    [setToken],
  );

  const signup = useCallback(
    async (email: string, password: string, name: string) => {
      const resp = await authApi.signup({
        email,
        password,
        name,
      });
      if (resp.token) {
        await setToken(resp.token);
      }
    },
    [setToken],
  );

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({
      isAuthenticated,
      user,
      login,
      signup,
      logout,
      isLoading,
    }),
    [isAuthenticated, user, login, signup, logout, isLoading],
  );

  return <AuthContext value={value}>{children}</AuthContext>;
};
