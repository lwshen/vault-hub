import { useEffect, useMemo } from 'react';
import { Link } from 'wouter';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { PATH } from '@/const/path';

export default function MagicLink() {
  const token = useMemo(() => {
    if (typeof window === 'undefined') {
      return '';
    }
    const search = new URLSearchParams(window.location.search);
    return search.get('token') ?? '';
  }, []);

  useEffect(() => {
    if (!token) {
      return;
    }

    const redirectUrl = `/auth/ml?token=${encodeURIComponent(token)}`;
    // Ensure we trigger a full navigation so the backend can issue the JWT fragment redirect.
    window.location.replace(redirectUrl);
  }, [token]);

  return (
    <div className="flex min-h-full flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Logging you in</CardTitle>
            <CardDescription>
              We&apos;re validating your magic link. This should only take a
              second.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-6 text-sm text-muted-foreground">
            {token ? (
              <>
                <p>
                  You&apos;ll be redirected automatically. If nothing happens,
                  you can{' '}
                  <a
                    href={`/auth/ml?token=${encodeURIComponent(token)}`}
                    className="text-primary underline-offset-2 hover:underline"
                  >
                    continue with this link
                  </a>
                  .
                </p>
                <p>
                  Having trouble? Request a new link from the{' '}
                  <Link href={PATH.LOGIN} className="text-primary underline-offset-2 hover:underline">
                    login page
                  </Link>
                  .
                </p>
              </>
            ) : (
              <>
                <p>This magic link is missing a token or has already been used.</p>
                <div className="grid gap-2">
                  <p>Request a new magic link to continue.</p>
                  <Button asChild variant="outline">
                    <Link href={PATH.LOGIN}>Back to login</Link>
                  </Button>
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
