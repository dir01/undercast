import { FunctionalComponent, h } from "preact";
import * as style from "./style.css";
import { Profile } from "../../API";

type HeaderProps = { onLogout: () => Promise<void>; profile: Profile };

const Header: FunctionalComponent<HeaderProps> = ({
    onLogout
}: HeaderProps) => {
    return (
        <header class={style.header}>
            <h1>Undercast</h1>
            <a onClick={onLogout}>logout</a>
        </header>
    );
};

export default Header;
