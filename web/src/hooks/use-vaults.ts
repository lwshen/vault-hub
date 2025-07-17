import { useState, useEffect } from 'react';
import { vaultApi } from '@/apis/api';
import type { Vault } from '@lwshen/vault-hub-ts-fetch-client';

interface UseVaultsReturn {
  vaults: Vault[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export const useVaults = (category?: string): UseVaultsReturn => {
  const [vaults, setVaults] = useState<Vault[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchVaults = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await vaultApi.getVaults(category);
      setVaults(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch vaults');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchVaults();
  }, [category]);

  return {
    vaults,
    isLoading,
    error,
    refetch: fetchVaults,
  };
}; 
