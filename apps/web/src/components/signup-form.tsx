import React, { useState, useEffect } from 'react';
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
import { FaOpenid } from 'react-icons/fa';
import { useLocation } from 'wouter';
import { PATH } from '@/const/path';
import useAuth from '@/hooks/use-auth';
import { configApi } from '@/apis/api';

export function SignupForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const [, navigate] = useLocation();
  const { signup, loginWithOidc } = useAuth();

  const [form, setForm] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [oidcEnabled, setOidcEnabled] = useState<boolean>(false);
  const [oidcLoading, setOidcLoading] = useState(true);

  // Fetch OIDC configuration from backend
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

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (form.password !== form.confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    setLoading(true);
    try {
      await signup(form.email, form.password, form.name);
      navigate(PATH.DASHBOARD);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Signup failed');
    } finally {
      setLoading(false);
    }
  };

  const handleOidcLogin = () => {
    loginWithOidc();
  };


  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Create an account</CardTitle>
          <CardDescription>
            Enter your information below to create your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="grid gap-6">
              <div className="grid gap-6">
                <div className="grid gap-3">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    name="name"
                    type="text"
                    required
                    value={form.name}
                    onChange={handleChange}
                  />
                </div>
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
                <div className="grid gap-3">
                  <Label htmlFor="confirm-password">Confirm password</Label>
                  <Input id="confirm-password" name="confirmPassword" type="password" required value={form.confirmPassword} onChange={handleChange} />
                </div>
                {error && <div className="text-red-500 text-sm text-center">{error}</div>}
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? 'Creating account...' : 'Create account'}
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
                    <Button variant="outline" className="w-full" onClick={handleOidcLogin} aria-label="Sign up with OpenID Connect">
                      <FaOpenid />
                      Sign up with OIDC
                    </Button>
                  </div>
                </>
              )}
              <div className="text-center text-sm">
                Already have an account?{' '}
                <Button variant="link" onClick={() => navigate(PATH.LOGIN)}>
                  Sign in
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
