import { vaultApi } from '@/apis/api';
import type { VaultLite } from '@lwshen/vault-hub-ts-fetch-client';
import { create } from 'zustand';

interface VaultState {
  vaults: VaultLite[];
  isLoading: boolean;
  error: string | null;
  isDeleting: boolean;
}

interface VaultActions {
  fetchVaults: () => Promise<void>;
  deleteVault: (uniqueId: string) => Promise<void>;
  clearError: () => void;
  reset: () => void;
}

type VaultStore = VaultState & VaultActions;

const initialState: VaultState = {
  vaults: [],
  isLoading: true,
  error: null,
  isDeleting: false,
};

export const useVaultStore = create<VaultStore>((set, get) => ({
  ...initialState,

  fetchVaults: async () => {
    set({ isLoading: true, error: null });

    try {
      const response = await vaultApi.getVaults();
      set({ vaults: response });
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : 'Failed to fetch vaults',
      });
    } finally {
      set({ isLoading: false });
    }
  },

  deleteVault: async (uniqueId: string) => {
    set({ isDeleting: true });

    try {
      await vaultApi.deleteVault(uniqueId);
      // Refetch vaults after deletion
      await get().fetchVaults();
    } finally {
      set({ isDeleting: false });
    }
  },

  clearError: () => set({ error: null }),

  reset: () => set(initialState),
}));
