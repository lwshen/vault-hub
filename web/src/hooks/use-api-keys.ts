import { useState, useEffect } from 'react';
import { apiKeyApi } from '@/apis/api';
import type { APIKey } from '@lwshen/vault-hub-ts-fetch-client';

interface UseApiKeysReturn {
  apiKeys: APIKey[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export const useApiKeys = (): UseApiKeysReturn => {
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchApiKeys = async () => {
    try {
      setIsLoading(true);
      setError(null);
      // Fetch first page with a large pageSize to simplify UI for now
      const response = await apiKeyApi.getAPIKeys(100, 1);

      let list: APIKey[] | undefined;

      if (typeof response === 'object' && response !== null && 'apiKeys' in response) {
        // The response is APIKeysResponse
        list = (response as { apiKeys?: APIKey[]; }).apiKeys;
      } else if (Array.isArray(response)) {
        list = response as APIKey[];
      }

      setApiKeys(list ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch API keys');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchApiKeys();
  }, []);

  return {
    apiKeys,
    isLoading,
    error,
    refetch: fetchApiKeys,
  };
};
