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
    let isMounted = true;

    configApi.getConfig()
      .then((config) => {
        if (isMounted) {
          setOidcEnabled(config.oidcEnabled);
        }
      })
      .catch((err) => {
        console.error('Failed to fetch OIDC config:', err);
        if (isMounted) {
          // Default to false if fetch fails
          setOidcEnabled(false);
        }
      })
      .finally(() => {
        if (isMounted) {
          setOidcLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  return { oidcEnabled, oidcLoading };
}
