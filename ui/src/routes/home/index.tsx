import { FunctionalComponent, h } from "preact";
import { useContext } from "preact/hooks";
import { v4 as uuidv4 } from "uuid";

import DownloadForm from "../../components/DownloadForm";
import DownloadsList from "../../components/DownloadsList";
import { DownloadInput, Download } from "../../api";
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
            <DownloadsList
                downloads={downloads}
                createMedia={async (d: Download) => {
                    if (!api) return;
                    await api.createMedia({
                        id: uuidv4(),
                        downloadId: d.id,
                        files: d.files
                    });
                }}
            />
        </div>
    );
};

export default Home;
