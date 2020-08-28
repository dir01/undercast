/* eslint-disable @typescript-eslint/no-use-before-define */
import { h, FunctionComponent, Fragment } from "preact";
import { Download, DownloadInput } from "../../api";
import * as style from "./style.css";
import { useState, useCallback } from "preact/hooks";

type DownloadsListProps = { downloads: (Download | DownloadInput)[] };

const DownloadsList: FunctionComponent<DownloadsListProps> = ({
    downloads
}: DownloadsListProps) => {
    return (
        <ul class={style.downloadList}>
            {downloads.map(d => (
                <li class={style.downloadItem} key={d.id}>
                    {isDownloadInput(d) ? (
                        <JustAddedDownloadItem downloadInput={d} />
                    ) : (
                        <DownloadItem download={d} />
                    )}
                </li>
            ))}
        </ul>
    );
};

const DownloadItem: FunctionComponent<{ download: Download }> = ({
    download
}: {
    download: Download;
}) => {
    const [areFilesVisible, setFilesVisible] = useState(false);
    const toggleFiles = useCallback(() => setFilesVisible(!areFilesVisible), [
        areFilesVisible
    ]);
    return (
        <Fragment>
            <span class={style.title} onClick={toggleFiles}>
                {download.name}
            </span>
            {areFilesVisible ? (
                <ul class={style.fileList}>
                    {download.files.map(f => (
                        <li key={f}>{f}</li>
                    ))}
                </ul>
            ) : null}
        </Fragment>
    );
};

const JustAddedDownloadItem: FunctionComponent<{
    downloadInput: DownloadInput;
}> = ({ downloadInput }: { downloadInput: DownloadInput }) => {
    return (
        <Fragment>
            <span class={style.title}>
                {extractDescription(downloadInput.source) ||
                    "Fetching metadata"}
            </span>
        </Fragment>
    );
};

function isDownloadInput(d: Download | DownloadInput): d is DownloadInput {
    return (d as any).name === undefined;
}

function extractDescription(s: string): string {
    try {
        return new URL(s).searchParams.get("dn") || "";
    } catch (e) {
        return "";
    }
}

export default DownloadsList;
