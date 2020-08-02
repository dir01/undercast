/* eslint-disable @typescript-eslint/no-use-before-define */
import "milligram/dist/milligram.css";

import { FunctionalComponent, h } from "preact";
import { Route, Router } from "preact-router";
import { useState } from "preact/hooks";

import API, { Profile } from "../API";
import ApiContext from "../contexts/ApiContext";
import Home from "../routes/home";
import NotFoundPage from "../routes/notfound";
import Header from "./Header";
import Auth from "./Auth";
import { AuthContainer } from "./Auth/hooks";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
if ((module as any).hot) {
    // tslint:disable-next-line:no-var-requires
    require("preact/debug");
}

const App: FunctionalComponent = () => {
    const { profile, logout } = AuthContainer.useContainer();
    return (
        <Auth>
            <div id="app">
                <Header profile={profile as Profile} onLogout={logout} />
                <Router>
                    <Route path="/" component={Home} />
                    <NotFoundPage default />
                </Router>
            </div>
        </Auth>
    );
};

const Wrapped: FunctionalComponent = () => {
    const [api] = useState(new API("http://localhost:4242"));
    return (
        <ApiContext.Provider value={api}>
            <AuthContainer.Provider initialState={api}>
                <App />
            </AuthContainer.Provider>
        </ApiContext.Provider>
    );
};

export default Wrapped;
