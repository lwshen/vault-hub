import { apiKeyApi } from '@/apis/api';
import type { VaultAPIKey } from '@lwshen/vault-hub-ts-fetch-client';
import { create } from 'zustand';

interface ApiKeyState {
  apiKeys: VaultAPIKey[];
  isLoading: boolean;
  error: string | null;
}

interface ApiKeyActions {
  fetchApiKeys: () => Promise<void>;
  deleteApiKey: (id: number) => Promise<void>;
  clearError: () => void;
  reset: () => void;
}

type ApiKeyStore = ApiKeyState & ApiKeyActions;

const initialState: ApiKeyState = {
  apiKeys: [],
  isLoading: true,
  error: null,
};

export const useApiKeyStore = create<ApiKeyStore>((set, get) => ({
  ...initialState,

  fetchApiKeys: async () => {
    set({ isLoading: true, error: null });

    try {
      // Fetch first page with a large pageSize to simplify UI for now
      const response = await apiKeyApi.getAPIKeys(100, 1);

      let list: VaultAPIKey[] | undefined;

      if (typeof response === 'object' && response !== null && 'apiKeys' in response) {
        // The response is APIKeysResponse
        list = (response as { apiKeys?: VaultAPIKey[]; }).apiKeys;
      } else if (Array.isArray(response)) {
        list = response as VaultAPIKey[];
      }

      set({ apiKeys: list ?? [] });
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : 'Failed to fetch API keys',
      });
    } finally {
      set({ isLoading: false });
    }
  },

  deleteApiKey: async (id: number) => {
    await apiKeyApi.deleteAPIKey(id);
    // Refetch API keys after deletion
    await get().fetchApiKeys();
  },

  clearError: () => set({ error: null }),

  reset: () => set(initialState),
}));
