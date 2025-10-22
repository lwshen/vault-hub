import type { GetUserResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { createContext } from 'react';

export type AuthTokenSource = 'magic' | 'oidc' | string;

export interface AuthenticateOptions {
  source?: AuthTokenSource;
  redirectTo?: string;
  showToast?: boolean;
}

export interface AuthContextType {
  isAuthenticated: boolean;
  user: GetUserResponse | null;
  login: (email: string, password: string) => Promise<void>;
  loginWithOidc: () => Promise<void>;
  signup: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
  requestMagicLink: (email: string) => Promise<void>;
  requestPasswordReset: (email: string) => Promise<void>;
  authenticateWithToken: (token: string, options?: AuthenticateOptions) => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
