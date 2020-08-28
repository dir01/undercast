/* eslint-disable @typescript-eslint/no-explicit-any */
import { FunctionalComponent, h } from "preact";
import { useState, useCallback } from "preact/hooks";

import { DownloadInput } from "../../api";

type DownloadFormProps = {
    onSubmitDownload: (d: DownloadInput) => Promise<void>;
};

const DownloadForm: FunctionalComponent<DownloadFormProps> = ({
    onSubmitDownload
}: DownloadFormProps) => {
    const [source, setSource] = useState("");

    const onInputSource = useCallback((event: any) => {
        setSource(event.target.value);
    }, []);

    const onSubmit = useCallback(
        (event: any) => {
            const d = new DownloadInput({ source });
            onSubmitDownload(d);
            setSource("");
            event.preventDefault();
            return false;
        },
        [source, onSubmitDownload]
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
