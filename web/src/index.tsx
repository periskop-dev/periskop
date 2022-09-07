import "./index.scss"

import ReactDOM from "react-dom/client";
import { Provider } from "react-redux"
import { store } from "data/store"
import NavbarComponent from "components/Navbar"

import App from "components/App"
import Home from "components/Home"
import TargetsList from "components/TargetsList"
import { Route, Router, Switch } from "react-router"
import { createHashHistory } from "history"
import reportWebVitals from './reportWebVitals';

const history = createHashHistory()

const app = ReactDOM.createRoot(document.getElementById("app"));
app.render(
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
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
