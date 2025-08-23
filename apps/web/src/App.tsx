import Header from '@/components/layout/header';
import { AuthProvider } from './contexts/auth-provider';
import { AppRoutes } from './routes';
import { Toaster } from '@/components/ui/sonner';

function App() {
  return (
    <div className="h-screen flex flex-col">
      <AuthProvider>
        <Header />
        <main className="flex-1 overflow-auto">
          <AppRoutes />
        </main>
        <Toaster />
      </AuthProvider>
    </div>
  );
}

export default App;
