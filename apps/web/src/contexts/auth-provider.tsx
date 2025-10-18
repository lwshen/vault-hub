import { authApi, userApi } from '@/apis/api';
import { PATH } from '@/const/path';
import type { GetUserResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react';
import { toast } from 'sonner';
import { navigate } from 'wouter/use-browser-location';
import { AuthContext } from './auth-context';

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
      // Check for OIDC token in URL fragment (hash) first
      // URL fragments are never sent to the server, providing better security
      const fragment = window.location.hash.substring(1); // Remove '#' prefix
      const fragmentParams = new URLSearchParams(fragment);
      const fragmentToken = fragmentParams.get('token');
      const source = fragmentParams.get('source');

      const isMagicLinkSource = source === 'magic' || source === 'magiclink';
      if (fragmentToken && (source === 'oidc' || isMagicLinkSource)) {
        try {
          await setToken(fragmentToken);
          // Remove fragment from URL while preserving any query params
          const cleanUrl = `${window.location.pathname}${window.location.search}`;
          window.history.replaceState({}, document.title, cleanUrl);
          if (isMagicLinkSource) {
            toast.success('You are now signed in with your magic link.');
          }
          setIsLoading(false);
          navigate(PATH.DASHBOARD);
          return; // Skip regular token check since we already set the token
        } catch (error) {
          console.error('Failed to set token from URL fragment:', error);
          // Remove invalid token from localStorage
          localStorage.removeItem('token');
          if (isMagicLinkSource) {
            toast.error('Unable to sign in with this magic link. Please request a new one.');
          }
          // Clear the fragment and continue with regular token check
          const cleanUrl = `${window.location.pathname}${window.location.search}`;
          window.history.replaceState({}, document.title, cleanUrl);
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
        navigate(PATH.DASHBOARD);
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
        navigate(PATH.DASHBOARD);
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

  const requestPasswordReset = useCallback(async (email: string) => {
    try {
      await authApi.requestPasswordReset({ email });
      toast.success(
        "If an account exists with this email, you'll receive password reset instructions shortly.",
      );
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : 'Unable to send reset instructions. Please try again.';
      toast.error(message);
      throw error;
    }
  }, []);

  const requestMagicLink = useCallback(async (email: string) => {
    try {
      await authApi.requestMagicLink({ email });
      toast.success("We've sent you a login link. Please check your email.");
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : 'Unable to send magic link. Please try again.';
      toast.error(message);
      throw error;
    }
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
      requestPasswordReset,
      requestMagicLink,
    }),
    [
      isAuthenticated,
      user,
      login,
      loginWithOidc,
      signup,
      logout,
      isLoading,
      requestPasswordReset,
      requestMagicLink,
    ],
  );

  return <AuthContext value={value}>{children}</AuthContext>;
};
