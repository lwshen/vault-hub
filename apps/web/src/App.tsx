import Header from '@/components/layout/header';
import { AuthProvider } from './contexts/auth-provider';
import { AppRoutes } from './routes';

function App() {
  return (
    <div className="h-screen flex flex-col">
      <AuthProvider>
        <Header />
        <main className="flex-1 overflow-auto">
          <AppRoutes />
        </main>
      </AuthProvider>
    </div>
  );
}

export default App;
