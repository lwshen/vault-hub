import type { GetUserResponse } from '@lwshen/vault-hub-ts-fetch-client';
import { createContext } from 'react';

export interface AuthContextType {
  isAuthenticated: boolean;
  user: GetUserResponse | null;
  login: (email: string, password: string) => Promise<void>;
  signup: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
