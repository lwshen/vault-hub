import { Route, Switch } from 'wouter';
import HeroSection from '@/components/hero/hero-section';
import Features from '@/pages/features';
import Login from '@/pages/auth/login';
import Signup from '@/pages/auth/signup';
import Dashboard from '@/pages/dashboard/dashboard';
import Vaults from '@/pages/dashboard/vaults';
import AuditLog from '@/pages/dashboard/audit-log';
import ApiKeys from '@/pages/dashboard/api-keys';
import Mock from '@/pages/mock';
import { ProtectedRoute } from '@/components/protected-route';
import { PATH } from '@/const/path';

export const AppRoutes = () => (
  <Switch>
    <Route path={PATH.HOME}>
      <HeroSection />
    </Route>
    <Route path={PATH.FEATURES}>
      <Features />
    </Route>
    <Route path={PATH.DOCS}>
      <div>Docs</div>
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
