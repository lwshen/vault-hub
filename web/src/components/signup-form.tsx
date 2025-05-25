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
import { FaOpenid, FaGoogle, FaApple } from 'react-icons/fa';
import { useLocation } from 'wouter';
import { PATH } from '@/const/path';

export function SignupForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const [, navigate] = useLocation();

  const navigateToLogin = () => {
    navigate(PATH.LOGIN);
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
          <form>
            <div className="grid gap-6">
              <div className="grid gap-6">
                <div className="grid gap-3">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    type="text"
                    required
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    required
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="password">Password</Label>
                  <Input id="password" type="password" required />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="confirm-password">Confirm password</Label>
                  <Input id="confirm-password" type="password" required />
                </div>
                <Button type="submit" className="w-full">
                  Create account
                </Button>
              </div>
              <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
                <span className="bg-card text-muted-foreground relative z-10 px-2">
                  Or continue with
                </span>
              </div>

              <div className="flex flex-col gap-4">
                <Button variant="outline" className="w-full">
                  <FaApple />
                  Sign up with Apple
                </Button>
                <Button variant="outline" className="w-full">
                  <FaGoogle />
                  Sign up with Google
                </Button>
                <Button variant="outline" className="w-full">
                  <FaOpenid />
                  Sign up with OIDC
                </Button>
              </div>
              <div className="text-center text-sm">
                Already have an account?{' '}
                <Button variant="link" onClick={navigateToLogin}>
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
