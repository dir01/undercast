/* eslint-disable @typescript-eslint/no-use-before-define */
import { h, FunctionComponent, Fragment } from "preact";
import { useState, useCallback } from "preact/hooks";

import { Download, DownloadInput } from "../../api";

import * as style from "./style.css";

type DownloadsListProps = {
    downloads: (Download | DownloadInput)[];
    createMedia: (download: Download) => Promise<void>;
};

const DownloadsList: FunctionComponent<DownloadsListProps> = ({
    downloads,
    createMedia
}: DownloadsListProps) => {
    return (
        <ul class={style.downloadList}>
            {downloads.map(d => (
                <li class={style.downloadItem} key={d.id}>
                    {isDownloadInput(d) ? (
                        <JustAddedDownloadItem downloadInput={d} />
                    ) : (
                        <DownloadItem download={d} createMedia={createMedia} />
                    )}
                </li>
            ))}
        </ul>
    );
};

type DownloadItemProps = {
    download: Download;
    createMedia: (download: Download) => Promise<void>;
};
const DownloadItem: FunctionComponent<DownloadItemProps> = ({
    download,
    createMedia
}: DownloadItemProps) => {
    const [areFilesVisible, setFilesVisible] = useState(false);
    const toggleFiles = useCallback(() => setFilesVisible(!areFilesVisible), [
        areFilesVisible
    ]);
    return (
        <Fragment>
            <span class={style.title} onClick={toggleFiles}>
                {download.name} ({download.percentDone})
            </span>
            <button onClick={() => createMedia(download)}>publish</button>
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
