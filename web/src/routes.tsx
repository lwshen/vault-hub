import { Route, Switch } from 'wouter';
import HeroSection from '@/components/hero/hero-section';
import Login from '@/pages/auth/login';
import Signup from '@/pages/auth/signup';
import Dashboard from '@/pages/dashboard';
import Vaults from '@/pages/vaults';
import { PATH } from '@/const/path';

export const AppRoutes = () => (
  <Switch>
    <Route path={PATH.HOME}>
      <HeroSection />
    </Route>
    <Route path={PATH.FEATURES}>
      <div>Features</div>
    </Route>
    <Route path={PATH.PRICING}>
      <div>Pricing</div>
    </Route>
    <Route path={PATH.DOCS}>
      <div>Docs</div>
    </Route>
    <Route path={PATH.ABOUT}>
      <div>About</div>
    </Route>
    <Route path="/users/:name">{(params) => <>Hello, {params.name}!</>}</Route>
    <Route path={PATH.LOGIN}>
      <Login />
    </Route>
    <Route path={PATH.SIGNUP}>
      <Signup />
    </Route>
    <Route path={PATH.DASHBOARD}>
      <Dashboard />
    </Route>
    <Route path={PATH.VAULTS}>
      <Vaults />
    </Route>
    <Route>404: No such page!</Route>
  </Switch>
);
