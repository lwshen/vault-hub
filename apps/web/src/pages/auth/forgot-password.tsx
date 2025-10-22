import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { PATH } from '@/const/path';
import useAuth from '@/hooks/use-auth';
import { MailCheck } from 'lucide-react';
import { useState } from 'react';
import { Link } from 'wouter';

export default function ForgotPassword() {
  const { requestPasswordReset } = useAuth();
  const [email, setEmail] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<'idle' | 'success'>('idle');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await requestPasswordReset(email);
      setStatus('success');
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : 'Unable to send reset instructions. Please try again.';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  const handleStartOver = () => {
    setStatus('idle');
    setError(null);
  };

  return (
    <div className="flex min-h-full flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Reset your password</CardTitle>
          </CardHeader>
          <CardContent>
            {status === 'success' ? (
              <div className="grid gap-6">
                <Alert variant="info">
                  <MailCheck aria-hidden="true" />
                  <AlertTitle>Check your inbox</AlertTitle>
                  <AlertDescription>
                    Reset instructions are on the way to <span className="font-medium text-card-foreground">{email}</span> if it
                    matches an account. The link stays active for 30 minutes—take a glance at spam just in case.
                  </AlertDescription>
                </Alert>
                <Button asChild className="w-full">
                  <Link href={PATH.LOGIN}>Back to login</Link>
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  className="w-full text-sm text-muted-foreground"
                  onClick={handleStartOver}
                >
                  Use a different email
                </Button>
              </div>
            ) : (
              <form className="grid gap-6" onSubmit={handleSubmit}>
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    autoComplete="email"
                    required
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                  />
                </div>
                {error && <p className="text-sm text-red-500">{error}</p>}
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? 'Sending instructions...' : 'Send reset link'}
                </Button>
                <p className="text-center text-sm text-muted-foreground">
                  Remember your password?{' '}
                  <Link href={PATH.LOGIN} className="text-primary underline-offset-2 hover:underline">
                    Go back to login
                  </Link>
                </p>
              </form>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
