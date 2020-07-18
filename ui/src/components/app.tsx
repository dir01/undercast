import "milligram/dist/milligram.css";

import { FunctionalComponent, h } from "preact";
import { Route, Router } from "preact-router";

import Home from "../routes/home";
import NotFoundPage from "../routes/notfound";
import Header from "./header";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
if ((module as any).hot) {
    // tslint:disable-next-line:no-var-requires
    require("preact/debug");
}

const App: FunctionalComponent = () => {
    return (
        <div id="app">
            <Header />
            <Router>
                <Route path="/" component={Home} />
                <NotFoundPage default />
            </Router>
        </div>
    );
};

export default App;
