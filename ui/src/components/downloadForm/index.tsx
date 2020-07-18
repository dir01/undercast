import { FunctionalComponent, h } from "preact";
import { useState, useCallback } from "preact/hooks";
import API, { Download } from "../../API";

const api = new API("");

const DownloadForm: FunctionalComponent = () => {
    const [source, setSource] = useState("");

    const onInputSource = useCallback((event: any) => {
        setSource(event.target.value);
    }, [source]);

    const onSubmit = useCallback((event: any) => {
        const d = new Download({ source });
        api.saveDownload(d);
        setSource("");
        event.preventDefault();
        return false;
    }, [source]);

    return (
        <form onSubmit={onSubmit}>
            <fieldset>
                <label for="source">Source</label>
                <input type="text" onInput={onInputSource} value={source} />
                <input class="button" type="submit" value="Send" disabled={!source} />
            </fieldset>
        </form>
    );
};

export default DownloadForm