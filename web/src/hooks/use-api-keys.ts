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
      // If the API returns an object with apiKeys array
      if (Array.isArray((response as any).apiKeys)) {
        setApiKeys((response as any).apiKeys);
      } else if (Array.isArray(response)) {
        setApiKeys(response as APIKey[]);
      } else {
        setApiKeys([]);
      }
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