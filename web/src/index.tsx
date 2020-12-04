import "index.scss"
import "favicon.ico"

import * as React from "react"
import { render } from "react-dom"

import { Provider } from "react-redux"
import { store } from "data/store"
import NavbarComponent from "components/Navbar"

import App from "components/App"
import Home from "components/Home"
import { Route, Router } from "react-router"
import { createHashHistory } from "history"

const history = createHashHistory()

render(
    <Provider store={store}>
        <Router history={history}>
            <div className="app-container">
                <NavbarComponent />
                <main className="app-content">
                    <Route exact path="/" component={Home}/>
                    <Route path="/:service" component={App}/>
                </main>
            </div>
        </Router>
    </Provider>,
    document.getElementById("app")
);
