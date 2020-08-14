import { FunctionalComponent, h } from "preact";
import { useState, useCallback } from "preact/hooks";

import * as style from "./style.css";

type LoginFormProps = {
    onSubmit: (token: string) => Promise<void>;
    error: string;
};

const LoginForm: FunctionalComponent<LoginFormProps> = ({
    onSubmit,
    error
}: LoginFormProps) => {
    const [password, setPassword] = useState("");
    const onInputPassword = useCallback((event: any) => {
        setPassword(event.target.value);
    }, []);

    return (
        <div class={style.loginPromptContainer}>
            <p>
                {(() => {
                    if (!error) {
                        return (
                            <p>
                                <h1>Hi there</h1>
                                <h2>
                                    This application is intended to be used by
                                    friends and family. If you are friends or
                                    family, then what is the magic word?
                                </h2>
                            </p>
                        );
                    }
                    if (error === "wrong_password") {
                        return (
                            <p>
                                <h1>Suspicious</h1>
                                <h2>
                                    Hmmm, that&apos;s not a magic word. Are you
                                    sure you&apos;ve typed it just right?
                                </h2>
                            </p>
                        );
                    }
                })()}
            </p>
            <form
                onSubmit={event => {
                    event.preventDefault();
                    onSubmit(password);
                    setPassword("");
                    return false;
                }}
            >
                <fieldset>
                    <input
                        placeholder="Magic word"
                        value={password}
                        onInput={onInputPassword}
                    />
                    <input class="button-primary" type="submit" value="Send" />
                </fieldset>
            </form>
        </div>
    );
};

export default LoginForm;
