import { useMemo, useState } from 'react';
import { Link, useLocation } from 'wouter';
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
import { authApi } from '@/apis/api';
import { toast } from 'sonner';
import { CheckCircle2, XCircle } from 'lucide-react';

const requirements = [
  {
    id: 'length',
    label: 'At least 8 characters',
    test: (value: string) => value.length >= 8,
  },
  {
    id: 'uppercase',
    label: 'One uppercase letter',
    test: (value: string) => /[A-Z]/.test(value),
  },
  {
    id: 'lowercase',
    label: 'One lowercase letter',
    test: (value: string) => /[a-z]/.test(value),
  },
  {
    id: 'number',
    label: 'One number',
    test: (value: string) => /\d/.test(value),
  },
  {
    id: 'special',
    label: 'One special character',
    test: (value: string) => /[^A-Za-z0-9]/.test(value),
  },
] as const;

export default function ResetPassword() {
  const [, navigate] = useLocation();
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const token = useMemo(() => {
    if (typeof window === 'undefined') {
      return '';
    }
    const search = new URLSearchParams(window.location.search);
    return search.get('token') ?? '';
  }, []);

  const matches = useMemo(
    () =>
      requirements.map((requirement) => ({
        id: requirement.id,
        label: requirement.label,
        passed: requirement.test(password),
      })),
    [password],
  );

  const passedCount = matches.filter((item) => item.passed).length;
  const strengthPercentage = Math.round((passedCount / requirements.length) * 100);
  const passwordMismatch = confirmPassword.length > 0 && password !== confirmPassword;

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);

    if (!token) {
      setError('This reset link is invalid or has expired. Please request a new one.');
      return;
    }

    if (passwordMismatch) {
      setError('Passwords do not match.');
      return;
    }

    if (passedCount !== requirements.length) {
      setError('Please meet all password requirements before continuing.');
      return;
    }

    setLoading(true);
    try {
      await authApi.confirmPasswordReset({
        token,
        newPassword: password,
      });
      toast.success('Your password has been updated. You can now log in with your new password.');
      navigate(PATH.LOGIN);
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : 'We were unable to reset your password. Please request a new link.';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-full flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Choose a new password</CardTitle>
            <CardDescription>
              Pick a strong password to keep your account secure.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {!token ? (
              <div className="grid gap-4 text-sm text-muted-foreground">
                <p>This reset link is invalid or has already been used.</p>
                <p>
                  Please{' '}
                  <Link href={PATH.FORGOT_PASSWORD} className="text-primary underline-offset-2 hover:underline">
                    request a new password reset email
                  </Link>{' '}
                  to continue.
                </p>
              </div>
            ) : (
              <form className="grid gap-6" onSubmit={handleSubmit}>
                <div className="grid gap-3">
                  <Label htmlFor="password">New password</Label>
                  <Input
                    id="password"
                    type="password"
                    autoComplete="new-password"
                    required
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="confirm-password">Confirm new password</Label>
                  <Input
                    id="confirm-password"
                    type="password"
                    autoComplete="new-password"
                    required
                    value={confirmPassword}
                    onChange={(event) => setConfirmPassword(event.target.value)}
                  />
                  {passwordMismatch && <p className="text-sm text-red-500">Passwords must match.</p>}
                </div>

                <div className="grid gap-2">
                  <Label className="text-xs uppercase tracking-wide text-muted-foreground">Password strength</Label>
                  <div className="h-2 w-full rounded-full bg-muted">
                    <div
                      className="h-full rounded-full bg-primary transition-all"
                      style={{ width: `${strengthPercentage}%` }}
                    />
                  </div>
                  <ul className="grid gap-1 text-sm text-muted-foreground">
                    {matches.map((match) => (
                      <li key={match.id} className="flex items-center gap-2">
                        {match.passed ? (
                          <CheckCircle2 className="h-4 w-4 text-emerald-500" aria-hidden="true" />
                        ) : (
                          <XCircle className="h-4 w-4 text-muted-foreground" aria-hidden="true" />
                        )}
                        <span>{match.label}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                {error && <p className="text-sm text-red-500">{error}</p>}

                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? 'Updating password...' : 'Update password'}
                </Button>
              </form>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
