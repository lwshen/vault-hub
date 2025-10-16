import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { PATH } from '@/const/path';
import useAuth from '@/hooks/use-auth';
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

  const handleResend = () => {
    setStatus('idle');
  };

  return (
    <div className="flex min-h-full flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Reset your password</CardTitle>
            <CardDescription>
              Enter the email associated with your account and we&apos;ll email
              you reset instructions.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {status === 'success' ? (
              <div className="grid gap-5 text-sm text-muted-foreground">
                <p>
                  If an account exists with <span className="font-medium">{email}</span>, you&apos;ll receive an email with
                  instructions to reset your password. The link expires in 30
                  minutes.
                </p>
                <p>
                  Can&apos;t find the email? Check your spam folder or{' '}
                  <button
                    type="button"
                    onClick={handleResend}
                    className="text-primary underline-offset-2 hover:underline"
                  >
                    send another email
                  </button>
                  .
                </p>
                <p>
                  Ready to try signing in again? Return to{' '}
                  <Link href={PATH.LOGIN} className="text-primary underline-offset-2 hover:underline">
                    the login page
                  </Link>
                  .
                </p>
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
