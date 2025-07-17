import { useState } from 'react';
import { vaultApi } from '@/apis/api';
import type { VaultLite } from '@lwshen/vault-hub-ts-fetch-client';

interface UseVaultsReturn {
  vaults: VaultLite[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export const useVaults = (): UseVaultsReturn => {
  const [vaults, setVaults] = useState<VaultLite[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchVaults = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await vaultApi.getVaults();
      setVaults(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch vaults');
    } finally {
      setIsLoading(false);
    }
  };

  return {
    vaults,
    isLoading,
    error,
    refetch: fetchVaults,
  };
}; 
