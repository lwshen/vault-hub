import HeroSection from '@/components/hero/hero-section';
import { ProtectedRoute } from '@/components/protected-route';
import { PATH } from '@/const/path';
import ForgotPassword from '@/pages/auth/forgot-password';
import Login from '@/pages/auth/login';
import MagicLink from '@/pages/auth/magic-link';
import ResetPassword from '@/pages/auth/reset-password';
import Signup from '@/pages/auth/signup';
import ApiKeys from '@/pages/dashboard/api-keys';
import AuditLog from '@/pages/dashboard/audit-log';
import Dashboard from '@/pages/dashboard/dashboard';
import VaultDetail from '@/pages/dashboard/vault-detail';
import Vaults from '@/pages/dashboard/vaults';
import Documentation from '@/pages/documentation';
import Features from '@/pages/features';
import Mock from '@/pages/mock';
import { Route, Switch } from 'wouter';

export const AppRoutes = () => (
  <Switch>
    <Route path={PATH.HOME}>
      <HeroSection />
    </Route>
    <Route path={PATH.FEATURES}>
      <Features />
    </Route>
    <Route path={PATH.DOCS}>
      <Documentation />
    </Route>
    <Route path={PATH.MOCK}>
      <Mock />
    </Route>
    <Route path="/users/:name">{(params) => <>Hello, {params.name}!</>}</Route>
    <Route path={PATH.LOGIN}>
      <Login />
    </Route>
    <Route path={PATH.SIGNUP}>
      <Signup />
    </Route>
    <Route path={PATH.FORGOT_PASSWORD}>
      <ForgotPassword />
    </Route>
    <Route path={PATH.RESET_PASSWORD}>
      <ResetPassword />
    </Route>
    <Route path={PATH.MAGIC_LINK_LOGIN}>
      <MagicLink />
    </Route>
    <Route path={PATH.DASHBOARD}>
      <ProtectedRoute>
        <Dashboard />
      </ProtectedRoute>
    </Route>
    <Route path={PATH.VAULTS}>
      <ProtectedRoute>
        <Vaults />
      </ProtectedRoute>
    </Route>
    <Route path={PATH.VAULT_DETAIL}>
      {(params: { vaultId: string; }) => (
        <ProtectedRoute>
          <VaultDetail vaultId={params.vaultId} />
        </ProtectedRoute>
      )}
    </Route>
    <Route path={PATH.API_KEYS}>
      <ProtectedRoute>
        <ApiKeys />
      </ProtectedRoute>
    </Route>
    <Route path={PATH.AUDIT_LOG}>
      <ProtectedRoute>
        <AuditLog />
      </ProtectedRoute>
    </Route>
    <Route>404: No such page!</Route>
  </Switch>
);
