import { useContext, useEffect } from "preact/hooks";
import ApiContext from "../../contexts/ApiContext";
import { Download } from "../../api";
import usePersistedState from "../../utils/hooks/usePersistedState";

export const useDownloads = () => {
    const [downloads, setDownloads] = usePersistedState<Download[]>(
        "downloads",
        []
    );

    const api = useContext(ApiContext);

    useEffect(() => {
        if (!api) return;
        api.getDownloads().then(result => {
            if (result.isOk()) {
                setDownloads(result.getValue());
            }
        });
    }, [api]);

    const addDownload = (d: Download) => setDownloads([d, ...downloads]);

    return { downloads, addDownload };
};
