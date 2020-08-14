import { FunctionalComponent, h } from "preact";
import DownloadForm from "../../components/DownloadForm";
import * as style from "./style.css";

const Home: FunctionalComponent = () => {
    return (
        <div class={style.home}>
            <DownloadForm />
        </div>
    );
};

export default Home;
