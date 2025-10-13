import { cn } from '@/lib/utils';
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
import { FaOpenid } from 'react-icons/fa';
import { useLocation } from 'wouter';
import useAuth from '@/hooks/use-auth';
import { useState } from 'react';
import { useOidcConfig } from '@/hooks/use-oidc-config';

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const [, navigate] = useLocation();
  const { login, loginWithOidc } = useAuth();
  const [sending, setSending] = useState<'none' | 'reset' | 'magic'>('none');
  const [form, setForm] = useState({
    email: '',
    password: '',
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { oidcEnabled, oidcLoading } = useOidcConfig();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await login(form.email, form.password);
      navigate(PATH.DASHBOARD);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const navigateToSignup = () => {
    navigate('/signup');
  };

  const handleOidcLogin = () => {
    loginWithOidc();
  };

  const requestReset = async () => {
    setError(null);
    setSending('reset');
    try {
      await fetch('/api/auth/password/reset/request', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: form.email }),
      });
      setError('If your email exists, a reset link has been sent.');
    } catch (e) {
      setError('Failed to request reset');
    } finally {
      setSending('none');
    }
  };

  const requestMagic = async () => {
    setError(null);
    setSending('magic');
    try {
      await fetch('/api/auth/magic-link/request', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: form.email }),
      });
      setError('If your email exists, a magic link has been sent.');
    } catch (e) {
      setError('Failed to request magic link');
    } finally {
      setSending('none');
    }
  };

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            Enter your email below to login to your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="grid gap-6">
              <div className="grid gap-6">
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    name="email"
                    type="email"
                    required
                    value={form.email}
                    onChange={handleChange}
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="password">Password</Label>
                  <Input id="password" name="password" type="password" required value={form.password} onChange={handleChange} />
                </div>
                {error && <div className="text-red-500 text-sm">{error}</div>}
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? 'Logging in...' : 'Login'}
                </Button>
              </div>
              {!oidcLoading && oidcEnabled && (
                <>
                  <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
                    <span className="bg-card text-muted-foreground relative z-10 px-2">
                      Or continue with
                    </span>
                  </div>

                  <div className="flex flex-col gap-4">
                    <Button variant="outline" className="w-full" onClick={handleOidcLogin} aria-label="Login with OpenID Connect">
                      <FaOpenid />
                      Login with OIDC
                    </Button>
                  </div>
                </>
              )}
              <div className="flex items-center justify-between text-sm">
                <Button variant="link" type="button" onClick={requestReset} disabled={!form.email || sending !== 'none'}>
                  {sending === 'reset' ? 'Sending…' : 'Forgot password?'}
                </Button>
                <Button variant="link" type="button" onClick={requestMagic} disabled={!form.email || sending !== 'none'}>
                  {sending === 'magic' ? 'Sending…' : 'Send magic link'}
                </Button>
              </div>
              <div className="text-center text-sm">
                Don&apos;t have an account?{' '}
                <Button variant="link" onClick={navigateToSignup}>
                  Sign up
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
