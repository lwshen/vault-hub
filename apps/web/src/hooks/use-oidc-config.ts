import { useState, useEffect } from 'react';
import { configApi } from '@/apis/api';

/**
 * Hook to fetch OIDC configuration from the backend
 * @returns Object containing oidcEnabled flag and oidcLoading state
 */
export function useOidcConfig() {
  const [oidcEnabled, setOidcEnabled] = useState<boolean>(false);
  const [oidcLoading, setOidcLoading] = useState(true);

  useEffect(() => {
    const fetchOidcConfig = async () => {
      try {
        const config = await configApi.getConfig();
        setOidcEnabled(config.oidcEnabled);
      } catch (err) {
        console.error('Failed to fetch OIDC config:', err);
        // Default to false if fetch fails
        setOidcEnabled(false);
      } finally {
        setOidcLoading(false);
      }
    };

    fetchOidcConfig();
  }, []);

  return { oidcEnabled, oidcLoading };
}
