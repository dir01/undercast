/* eslint-disable @typescript-eslint/no-explicit-any */
import { FunctionalComponent, h } from "preact";
import { useState, useCallback, useContext } from "preact/hooks";
import API, { Download } from "../../api";
import ApiContext from "../../contexts/ApiContext";

const DownloadForm: FunctionalComponent = () => {
    const api = useContext(ApiContext) as API;

    const [source, setSource] = useState("");

    const onInputSource = useCallback((event: any) => {
        setSource(event.target.value);
    }, []);

    const onSubmit = useCallback(
        (event: any) => {
            const d = new Download({ source });
            api.saveDownload(d);
            setSource("");
            event.preventDefault();
            return false;
        },
        [source, api]
    );

    return (
        <form onSubmit={onSubmit}>
            <fieldset>
                <label htmlFor="source">Source</label>
                <input type="text" onInput={onInputSource} value={source} />
                <input type="submit" value="Send" disabled={!source} />
            </fieldset>
        </form>
    );
};

export default DownloadForm;
