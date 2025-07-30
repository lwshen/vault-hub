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

      let list: unknown = null;

      if (response && typeof response === 'object' && 'apiKeys' in response) {
        // @ts-ignore dynamic
        list = (response as any).apiKeys;
      } else {
        list = response;
      }

      if (Array.isArray(list)) {
        setApiKeys(list as APIKey[]);
      } else {
        // response may include null, undefined or other structure
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