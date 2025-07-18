import { Button } from '@/components/ui/button';
import useAuth from '@/hooks/use-auth';

export default function Mock() {
  const { signup, login } = useAuth();

  const handleRegister = async () => {
    await signup('mock@demo.com', 'Test1234!', 'demo');
  };

  const handleLogin = async () => {
    await login('mock@demo.com', 'Test1234!');
  };

  return (
    <div className="flex flex-col gap-4 p-4">
      <h1>Mock</h1>
      <Button variant="outline" onClick={handleRegister}>
        Register Demo User
      </Button>
      <Button variant="outline" onClick={handleLogin}>
        Login Demo User
      </Button>
    </div>
  );
}
