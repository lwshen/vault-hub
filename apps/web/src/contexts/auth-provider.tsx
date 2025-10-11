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
  }, []);

  const loginWithOidc = useCallback(async () => {
    // Redirect to OIDC login endpoint
    window.location.href = '/api/auth/login/oidc';
  }, []);

  useEffect(() => {
    const initializeAuth = async () => {
      // Check for OIDC token in URL first
      const urlParams = new URLSearchParams(window.location.search);
      const oidcToken = urlParams.get('token');
      const source = urlParams.get('source');

      if (oidcToken && source === 'oidc') {
        // Clean up URL and set token from OIDC
        try {
          await setToken(oidcToken);
          // Remove token from URL
          const newUrl = window.location.pathname;
          window.history.replaceState({}, document.title, newUrl);
          // Navigate to home after successful OIDC login
          setIsLoading(false);
          navigate(PATH.HOME);
          return; // Skip regular token check since we already set the OIDC token
        } catch (error) {
          console.error('Failed to set OIDC token:', error);
          // Continue with regular token check if OIDC fails
        }
      }

      // Regular token check
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
  }, [setToken]);

  const login = useCallback(
    async (email: string, password: string) => {
      const resp = await authApi.login({
        email,
        password,
      });
      if (resp.token) {
        await setToken(resp.token);
        navigate(PATH.HOME);
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
        navigate(PATH.HOME);
      }
    },
    [setToken],
  );

  const logout = useCallback(async () => {
    const token = localStorage.getItem('token');
    // If there is a token, call the backend logout API to record the audit log
    if (token) {
      try {
        await authApi.logout();
      } catch (error) {
        // Even if the API call fails, continue with the logout operation
        console.warn('Failed to call logout API:', error);
      }
    }
    localStorage.removeItem('token');
    setUser(null);
    navigate(PATH.HOME);
  }, []);

  const value = useMemo(
    () => ({
      isAuthenticated,
      user,
      login,
      loginWithOidc,
      signup,
      logout,
      isLoading,
    }),
    [isAuthenticated, user, login, loginWithOidc, signup, logout, isLoading],
  );

  return <AuthContext value={value}>{children}</AuthContext>;
};
