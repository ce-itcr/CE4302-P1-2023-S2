import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter, Route, Switch, Redirect } from "react-router-dom";
import App from './routes/App';
import "./assets/styles/tailwind.css";
import "@fortawesome/fontawesome-free/css/all.min.css";



const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <BrowserRouter>
    <Switch>
      <Route path="/app" component={App} />
      <Redirect from="*" to="/app" />
    </Switch>
  </BrowserRouter>
);