import { FunctionComponent, h } from "preact";
import * as style from "./style.css";

type SpinnerProps = {
    style?: string | { [key: string]: string | number };
    class?: string;
};

const Spinner: FunctionComponent<SpinnerProps> = ({
    class: className,
    style: extraStyle
}: SpinnerProps) => {
    return (
        <div class={`${className} ${style.spinner}`} style={extraStyle}>
            {[1, 2, 3, 4, 5, 6, 7, 8, 9].map(i => (
                <div
                    key={i}
                    style={{ backgroundColor: "var(--fg-primary-color)" }}
                ></div>
            ))}
        </div>
    );
};

export default Spinner;
