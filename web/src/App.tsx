import { Route, Switch } from "wouter";
import Home from "./components/home";

function App() {
  return (
    <>
      <Switch>
        <Route path="/" component={Home} />
        <Route path="/about">
          <div>About</div>
        </Route>
        <Route path="/users/:name">
          {(params) => <>Hello, {params.name}!</>}
        </Route>
        <Route>404: No such page!</Route>
      </Switch>
    </>
  );
}

export default App;
