import { AuthContext } from '@/contexts/auth-context';
import { use } from 'react';

const useAuth = () => {
  const context = use(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export default useAuth;
