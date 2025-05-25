import { Route, Switch } from 'wouter';
import HeroSection from '@/components/hero/hero-section';
import Header from '@/components/layout/header';
import Login from '@/pages/auth/login';

function App() {
  return (
    <>
      <Header />
      <Switch>
        <Route path="/">
          <HeroSection />
        </Route>
        <Route path="/features">
          <div className="mt-16">Features</div>
        </Route>
        <Route path="/pricing">
          <div className="mt-16">Pricing</div>
        </Route>
        <Route path="/docs">
          <div className="mt-16">Docs</div>
        </Route>
        <Route path="/about">
          <div className="mt-16">About</div>
        </Route>
        <Route path="/users/:name">
          {(params) => <>Hello, {params.name}!</>}
        </Route>
        <Route path="/login">
          <Login />
        </Route>
        <Route>404: No such page!</Route>
      </Switch>
    </>
  );
}

export default App;
