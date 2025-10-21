import { CheckCircle2, Loader2, XCircle } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Link } from 'wouter';

import { authApi } from '@/apis/api';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { PATH } from '@/const/path';
import useAuth from '@/hooks/use-auth';
import { ResponseError } from '@lwshen/vault-hub-ts-fetch-client';

type Status = 'idle' | 'processing' | 'success' | 'error';

type MagicLinkResponseBody = {
  redirectUrl?: string;
  token?: string;
  code?: string;
};

const INVALID_TOKEN_MESSAGE = 'This magic link token is malformed. Please request a new one.';
const MISSING_TOKEN_MESSAGE = 'This magic link is missing a token or has already been used.';
const TOKEN_PATTERN = /^[A-Za-z0-9._-]+$/;
const MIN_TOKEN_LENGTH = 16;
const MAX_TOKEN_LENGTH = 2048;

export default function MagicLink() {
  const { authenticateWithToken } = useAuth();
  const token = useMemo(() => {
    if (typeof window === 'undefined') {
      return '';
    }
    const search = new URLSearchParams(window.location.search);
    return search.get('token') ?? '';
  }, []);

  const isTokenValid = useMemo(() => {
    if (!token) {
      return false;
    }
    if (token.length < MIN_TOKEN_LENGTH || token.length > MAX_TOKEN_LENGTH) {
      return false;
    }
    return TOKEN_PATTERN.test(token);
  }, [token]);

  const [status, setStatus] = useState<Status>(
    token && isTokenValid ? 'processing' : 'error',
  );
  const [errorMessage, setErrorMessage] = useState<string | null>(() => {
    if (!token) {
      return MISSING_TOKEN_MESSAGE;
    }
    if (!isTokenValid) {
      return INVALID_TOKEN_MESSAGE;
    }
    return null;
  });
  const [redirectUrl, setRedirectUrl] = useState<string | null>(null);
  const [attempt, setAttempt] = useState(0);
  const [responseStatus, setResponseStatus] = useState<number | null>(null);

  const resolveRedirectPath = (raw?: string | null) => {
    if (!raw) {
      return PATH.DASHBOARD;
    }

    if (typeof window === 'undefined') {
      return PATH.DASHBOARD;
    }

    try {
      const targetUrl = raw.startsWith('http')
        ? new URL(raw)
        : new URL(raw, window.location.origin);
      return `${targetUrl.pathname}${targetUrl.search}`;
    } catch {
      return raw.startsWith('/') ? raw : `/${raw}`;
    }
  };

  useEffect(() => {
    if (!token || !isTokenValid) {
      return;
    }

    let isCancelled = false;

    const verifyMagicLink = async () => {
      setStatus('processing');
      setErrorMessage(null);
      setResponseStatus(null);

      try {
        const apiResponse = await authApi.consumeMagicLinkRaw(
          { token },
          {
            headers: {
              Accept: 'application/json',
            },
          },
        );

        if (isCancelled) {
          return;
        }

        const { raw } = apiResponse;
        setResponseStatus(raw.status);

        let data: MagicLinkResponseBody | null = null;
        try {
          data = (await raw.clone().json()) as MagicLinkResponseBody;
        } catch {
          // No JSON body available; fall through to error handling.
        }

        if (!data || typeof data.token !== 'string' || data.token.length === 0) {
          throw new Error('Unable to verify this magic link. Please request a new one.');
        }

        const nextPath = resolveRedirectPath(data.redirectUrl);
        setRedirectUrl(nextPath);
        setStatus('success');
        await authenticateWithToken(data.token, {
          source: 'magic',
          redirectTo: nextPath,
        });
        return;
      } catch (error) {
        if (isCancelled) {
          return;
        }

        if (error instanceof ResponseError) {
          const { response } = error;
          setResponseStatus(response.status);

          if (response.status === 302) {
            const location = response.headers.get('location');
            const defaultDestination = `${window.location.origin}${PATH.LOGIN}`;
            let destination = defaultDestination;
            if (location) {
              try {
                const resolved = new URL(location, window.location.origin);
                if (resolved.origin === window.location.origin) {
                  destination = resolved.toString();
                }
              } catch {
                // Ignore invalid URLs and fall back to default destination.
              }
            }
            setRedirectUrl(destination);
            setStatus('success');
            return;
          }

          let message = 'This magic link is invalid or has expired. Please request a new one.';
          try {
            const errorJson = await response.clone().json();
            if (typeof errorJson?.error === 'string' && errorJson.error.length > 0) {
              message = errorJson.error;
            } else if (typeof errorJson?.message === 'string' && errorJson.message.length > 0) {
              message = errorJson.message;
            } else if (typeof errorJson?.code === 'string' && errorJson.code.length > 0) {
              message = errorJson.code;
            } else {
              const text = await response.clone().text();
              if (text) {
                message = text;
              }
            }
          } catch {
            try {
              const text = await response.clone().text();
              if (text) {
                message = text;
              }
            } catch {
              // Ignore body parsing errors and fall back to default message.
            }
          }
          setErrorMessage(message);
        } else if (error instanceof Error) {
          setErrorMessage(error.message);
        } else {
          setErrorMessage('We were unable to verify this magic link. Check your connection and try again.');
        }
        setStatus('error');
      }
    };

    void verifyMagicLink();

    return () => {
      isCancelled = true;
    };
  }, [authenticateWithToken, isTokenValid, token, attempt]);

  useEffect(() => {
    if (status === 'success' && redirectUrl) {
      const timeout = window.setTimeout(() => {
        window.location.href = redirectUrl;
      }, 1500);
      return () => {
        window.clearTimeout(timeout);
      };
    }
    return undefined;
  }, [status, redirectUrl]);

  const description = (() => {
    if (!token) {
      return 'We could not find a valid magic link token in this URL.';
    }
    switch (status) {
      case 'processing':
        return 'Hang tight while we verify your magic link.';
      case 'success':
        return 'Magic link verified. You will be redirected shortly.';
      case 'error':
        return 'We could not sign you in with this magic link.';
      default:
        return 'Follow the steps below to finish signing in.';
    }
  })();

  const canRetry =
    Boolean(token) && isTokenValid && (responseStatus == null || responseStatus >= 500);

  const handleContinue = () => {
    if (redirectUrl) {
      window.location.href = redirectUrl;
    } else {
      window.location.href = PATH.LOGIN;
    }
  };

  const handleRetry = () => {
    if (!isTokenValid) {
      return;
    }
    setAttempt((prev) => prev + 1);
  };

  return (
    <div className="flex min-h-full flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Magic link sign in</CardTitle>
            <CardDescription>{description}</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-6 text-sm text-muted-foreground">
            {!token ? (
              <>
                <div className="flex items-center gap-2 text-red-500">
                  <XCircle className="h-5 w-5" aria-hidden="true" />
                  <span>This link is missing a token or has already been used.</span>
                </div>
                <div className="grid gap-2">
                  <p>Request a new magic link to continue.</p>
                  <Button asChild variant="outline">
                    <Link href={PATH.LOGIN}>Back to login</Link>
                  </Button>
                </div>
              </>
            ) : null}

            {token && status === 'processing' && (
              <div className="flex flex-col items-center gap-4">
                <Loader2 className="h-6 w-6 animate-spin text-primary" aria-hidden="true" />
                <p>Verifying your magic link…</p>
              </div>
            )}

            {token && status === 'success' && (
              <div className="grid gap-4">
                <div className="flex items-center gap-2 text-emerald-600">
                  <CheckCircle2 className="h-5 w-5" aria-hidden="true" />
                  <span>Magic link verified successfully.</span>
                </div>
                <p>
                  We&apos;re finishing up your sign in. If you&apos;re not redirected automatically, continue
                  below.
                </p>
                <Button onClick={handleContinue}>Continue to Vault Hub</Button>
              </div>
            )}

            {token && status === 'error' && (
              <div className="grid gap-4">
                <div className="flex items-start gap-2 text-red-500">
                  <XCircle className="h-5 w-5" aria-hidden="true" />
                  <span>{errorMessage}</span>
                </div>
                <div className="grid gap-2">
                  {canRetry && <Button onClick={handleRetry}>Try verification again</Button>}
                  <Button asChild variant="outline">
                    <Link href={PATH.LOGIN}>Request a new magic link</Link>
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
