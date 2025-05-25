import { Route, Switch } from 'wouter';
import HeroSection from '@/components/hero/hero-section';
import Header from '@/components/layout/header';
import Login from '@/pages/auth/login';
import Signup from '@/pages/auth/signup';
import { PATH } from '@/const/path';

function App() {
  return (
    <div className="h-screen flex flex-col">
      <Header />
      <main className="flex-1 overflow-hidden">
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
      </main>
    </div>
  );
}

export default App;
