import { useState, useEffect } from 'react';
import { configApi } from '@/apis/api';

/**
 * Hook to fetch OIDC configuration from the backend
 * @returns Object containing public config flags and loading state
 */
export function useOidcConfig() {
  const [oidcEnabled, setOidcEnabled] = useState<boolean>(false);
  const [emailEnabled, setEmailEnabled] = useState<boolean>(false);
  const [oidcLoading, setOidcLoading] = useState(true);

  useEffect(() => {
    let isMounted = true;

    configApi.getConfig()
      .then((config) => {
        if (isMounted) {
          setOidcEnabled(config.oidcEnabled);
          setEmailEnabled(config.emailEnabled);
        }
      })
      .catch((err) => {
        console.error('Failed to fetch OIDC config:', err);
        if (isMounted) {
          // Default to false if fetch fails
          setOidcEnabled(false);
          setEmailEnabled(false);
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

  return { oidcEnabled, emailEnabled, oidcLoading };
}
