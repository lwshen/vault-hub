import { Route, Switch } from 'wouter';
import HeroSection from '@/components/hero/hero-section';
import Header from '@/components/layout/header';
import Login from '@/pages/auth/login';
import Signup from '@/pages/auth/signup';
import { PATH } from '@/const/path';

function App() {
  return (
    <>
      <Header />
      <Switch>
        <Route path={PATH.HOME}>
          <HeroSection />
        </Route>
        <Route path={PATH.FEATURES}>
          <div className="mt-16">Features</div>
        </Route>
        <Route path={PATH.PRICING}>
          <div className="mt-16">Pricing</div>
        </Route>
        <Route path={PATH.DOCS}>
          <div className="mt-16">Docs</div>
        </Route>
        <Route path={PATH.ABOUT}>
          <div className="mt-16">About</div>
        </Route>
        <Route path="/users/:name">
          {(params) => <>Hello, {params.name}!</>}
        </Route>
        <Route path={PATH.LOGIN}>
          <Login />
        </Route>
        <Route path={PATH.SIGNUP}>
          <Signup />
        </Route>
        <Route>404: No such page!</Route>
      </Switch>
    </>
  );
}

export default App;
