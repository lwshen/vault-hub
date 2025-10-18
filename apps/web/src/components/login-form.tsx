import { Lock, Mail } from 'lucide-react';
import { useState } from 'react';
import { FaOpenid } from 'react-icons/fa';
import { useLocation } from 'wouter';

import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { PATH } from '@/const/path';
import useAuth from '@/hooks/use-auth';
import { useOidcConfig } from '@/hooks/use-oidc-config';
import { cn } from '@/lib/utils';

const STEP = {
  EMAIL: 'email',
  OPTIONS: 'options',
  PASSWORD: 'password',
} as const;

type Step = (typeof STEP)[keyof typeof STEP];
type MagicLinkStatus = 'success' | 'error';

interface MagicLinkFeedback {
  message: string;
  status: MagicLinkStatus;
}

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const [, navigate] = useLocation();
  const { login, loginWithOidc, requestMagicLink } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [step, setStep] = useState<Step>(STEP.EMAIL);
  const [magicLinkFeedback, setMagicLinkFeedback] = useState<MagicLinkFeedback | null>(null);
  const [magicLinkLoading, setMagicLinkLoading] = useState(false);
  const { oidcEnabled, oidcLoading } = useOidcConfig();
  const magicLinkMessage = magicLinkFeedback?.message;
  const magicLinkStatus = magicLinkFeedback?.status;

  const resetMagicLinkFeedback = () => setMagicLinkFeedback(null);

  const showMagicLinkFeedback = (status: MagicLinkStatus, message: string) => {
    setMagicLinkFeedback({
      status,
      message,
    });
  };

  const renderMagicLinkFeedback = () => {
    if (!magicLinkMessage) {
      return null;
    }

    const status = magicLinkStatus ?? 'error';

    return (
      <p
        className={cn(
          'text-xs',
          status === 'success' ? 'text-emerald-600' : 'text-red-500',
        )}
      >
        {magicLinkMessage}
      </p>
    );
  };

  const handleEmailSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    const normalizedEmail = email.trim();
    if (!normalizedEmail) {
      setError('Please enter your email to continue.');
      return;
    }
    setEmail(normalizedEmail);
    setStep(STEP.OPTIONS);
    resetMagicLinkFeedback();
  };

  const handlePasswordSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await login(email, password);
      navigate(PATH.DASHBOARD);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const navigateToSignup = () => {
    navigate(PATH.SIGNUP);
  };

  const handleOidcLogin = () => {
    loginWithOidc();
  };

  const navigateToForgotPassword = () => {
    navigate(PATH.FORGOT_PASSWORD);
  };

  const handleMagicLink = async () => {
    const normalizedEmail = email.trim();
    if (!normalizedEmail) {
      showMagicLinkFeedback('error', 'Enter your email address to receive a magic link.');
      return;
    }
    setError(null);
    resetMagicLinkFeedback();
    setMagicLinkLoading(true);
    try {
      await requestMagicLink(normalizedEmail);
      setEmail(normalizedEmail);
      showMagicLinkFeedback(
        'success',
        "We've sent a magic link to your email. It expires in 15 minutes.",
      );
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : 'Unable to send a magic link right now. Please try again.';
      showMagicLinkFeedback('error', message);
    } finally {
      setMagicLinkLoading(false);
    }
  };

  const handleChangeEmail = () => {
    setStep(STEP.EMAIL);
    setPassword('');
    setError(null);
    resetMagicLinkFeedback();
  };

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6">
            {step === STEP.EMAIL && (
              <form className="grid gap-6" onSubmit={handleEmailSubmit}>
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
                {error && <div className="text-red-500 text-sm">{error}</div>}
                <Button type="submit" className="w-full">
                  Next
                </Button>
                <Button
                  type="button"
                  variant="link"
                  className="h-auto px-0 text-sm text-muted-foreground"
                  onClick={navigateToForgotPassword}
                >
                  Forgot password?
                </Button>
              </form>
            )}

            {step === STEP.OPTIONS && (
              <div className="grid gap-4">
                <div className="grid gap-3">
                  <Label htmlFor="email-options">Email</Label>
                  <Input
                    id="email-options"
                    type="email"
                    value={email}
                    readOnly
                    className="bg-muted/40"
                  />
                </div>
                <div className="rounded-md border border-dashed border-border bg-muted/40 p-4 text-sm text-muted-foreground">
                  Choose how you&apos;d like to sign in.
                </div>
                <div className="grid gap-3">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={handleMagicLink}
                    disabled={magicLinkLoading}
                  >
                    <Mail className="h-4 w-4" aria-hidden="true" />
                    {magicLinkLoading ? 'Sending magic link...' : 'Send magic link'}
                  </Button>
                  <Button type="button" className="w-full" onClick={() => setStep(STEP.PASSWORD)}>
                    <Lock className="h-4 w-4" aria-hidden="true" />
                    Enter password
                  </Button>
                </div>
                {renderMagicLinkFeedback()}
                <Button
                  type="button"
                  variant="link"
                  className="h-auto px-0 text-sm text-muted-foreground"
                  onClick={handleChangeEmail}
                >
                  Use a different email
                </Button>
                <Button
                  type="button"
                  variant="link"
                  className="h-auto px-0 text-sm text-muted-foreground"
                  onClick={navigateToForgotPassword}
                >
                  Forgot password?
                </Button>
              </div>
            )}

            {step === STEP.PASSWORD && (
              <form className="grid gap-6" onSubmit={handlePasswordSubmit}>
                <div className="grid gap-3">
                  <Label htmlFor="email-readonly">Email</Label>
                  <Input
                    id="email-readonly"
                    type="email"
                    value={email}
                    readOnly
                    className="bg-muted/40"
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="password">Password for {email}</Label>
                  <Input
                    id="password"
                    type="password"
                    autoComplete="current-password"
                    required
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                  />
                </div>
                <div className="flex flex-wrap items-center justify-between gap-2 text-sm">
                  <Button
                    type="button"
                    variant="link"
                    className="h-auto px-0 font-normal text-muted-foreground"
                    onClick={navigateToForgotPassword}
                  >
                    Forgot password?
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="px-2 text-sm"
                    onClick={handleMagicLink}
                    disabled={magicLinkLoading}
                  >
                    <Mail className="h-4 w-4" aria-hidden="true" />
                    {magicLinkLoading ? 'Sending magic link...' : 'Send magic link instead'}
                  </Button>
                </div>
                {renderMagicLinkFeedback()}
                {error && <div className="text-red-500 text-sm">{error}</div>}
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? 'Logging in...' : 'Login'}
                </Button>
                <Button
                  type="button"
                  variant="link"
                  className="h-auto px-0 text-sm text-muted-foreground"
                  onClick={handleChangeEmail}
                >
                  Use a different email
                </Button>
              </form>
            )}

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
            <div className="text-center text-sm">
              Don&apos;t have an account?{' '}
              <Button variant="link" onClick={navigateToSignup}>
                Sign up
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
