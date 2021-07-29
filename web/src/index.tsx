import "index.scss"
import "favicon.ico"

import * as React from "react"
import { render } from "react-dom"

import { Provider } from "react-redux"
import { store } from "data/store"
import NavbarComponent from "components/Navbar"

import App from "components/App"
import Home from "components/Home"
import TargetsList from "components/TargetsList"
import { Route, Router, Switch } from "react-router"
import { createHashHistory } from "history"

const history = createHashHistory()

render(
    <Provider store={store}>
        <Router history={history}>
            <div className="app-container">
                <NavbarComponent />
                <main className="app-content">
                    <Switch>
                        <Route exact path="/" component={Home}/>
                        <Route exact path="/targets" component={TargetsList}/>
                        <Route path={["/:service/errors/:errorKey", "/:service"]} component={App}/>
                    </Switch>
                </main>
            </div>
        </Router>
    </Provider>,
    document.getElementById("app")
);
