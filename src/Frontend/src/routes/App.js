import React from "react";
import { Switch, Route, Redirect } from "react-router-dom";
import Selector from "../pages/selector/Selector";
import Dashboard from "../pages/dashboard/Dashboard";
import Summary from "../pages/summary/Summary";

export default function App() {
    return (
        <>
            <Switch>
                <Route path="/app/selector" exact component={Selector} />
                <Route path="/app/dashboard" exact component={Dashboard} />
                <Route path="/app/summary" exact component={Summary} />
                <Redirect from="/app" to="/app/selector" />
            </Switch>
        </>
    );
}
