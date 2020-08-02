import { FunctionComponent, h, ComponentChildren, Fragment } from "preact";

import Spinner from "../Spinner";
import * as style from "./style.css";
import { AuthContainer } from "./hooks";
import LoginForm from "./LoginForm";

type AuthProps = {
    children: ComponentChildren;
};

const Auth: FunctionComponent<AuthProps> = ({ children }: AuthProps) => {
    const { isLoggedIn, isLoading, login } = AuthContainer.useContainer();

    if (isLoading) {
        return (
            <div class={style.spinnerContainer}>
                <Spinner style={{ width: "120px", height: "120px" }} />
            </div>
        );
    }

    if (isLoggedIn) {
        return <Fragment>{children}</Fragment>;
    }

    return <LoginForm onSubmit={login} />;
};

export default Auth;
