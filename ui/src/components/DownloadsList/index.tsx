/* eslint-disable @typescript-eslint/no-use-before-define */
import { h, FunctionComponent } from "preact";
import { Download, DownloadInput } from "../../api";

type DownloadsListProps = { downloads: (Download | DownloadInput)[] };

const DownloadsList: FunctionComponent<DownloadsListProps> = ({
    downloads
}: DownloadsListProps) => {
    return (
        <ul>
            {downloads.map(d => (
                <li key={d.id}>
                    {d?.name ||
                        extractDescription(d.source) ||
                        "Fetching metadata"}
                    {d.percentDone ? ` (${d.percentDone}%)` : ""}
                </li>
            ))}
        </ul>
    );
};

function extractDescription(s: string): string {
    try {
        return new URL(s).searchParams.get("dn") || "";
    } catch (e) {
        return "";
    }
}

export default DownloadsList;
