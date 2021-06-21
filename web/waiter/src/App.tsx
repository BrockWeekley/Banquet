import React from 'react';
import './App.css';
import {
    BrowserRouter as Router,
    Switch,
    Route
} from 'react-router-dom'
import Home from "./pages/Home/Home";
import Create from './pages/Create/Create';

const App: React.FC = () => {
  return (
    <Router>
        <Switch>
            <Route path="/" exact={true}>
                <Home />
            </Route>
            <Route path="/home">
                <Home />
            </Route>
            <Route path="/create">
                <Create />
            </Route>
        </Switch>
    </Router>
  );
};

export default App;
