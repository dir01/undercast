import { FunctionalComponent, h } from "preact";
import { useContext } from "preact/hooks";

import DownloadForm from "../../components/DownloadForm";
import DownloadsList from "../../components/DownloadsList";
import { Download, DownloadInput } from "../../api";
import ApiContext from "../../contexts/ApiContext";

import { useDownloads } from "./hooks";
import * as style from "./style.css";

const Home: FunctionalComponent = () => {
    const api = useContext(ApiContext);
    const { downloads, addDownload } = useDownloads();

    return (
        <div class={style.home}>
            <DownloadForm
                onSubmitDownload={async (d: DownloadInput) => {
                    if (!api) return;
                    addDownload(d);
                    await api.saveDownload(d);
                }}
            />
            <DownloadsList downloads={downloads} />
        </div>
    );
};

export default Home;
