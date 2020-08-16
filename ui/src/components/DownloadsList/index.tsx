import { h, FunctionComponent } from "preact";
import { Download } from "../../api";

type DownloadsListProps = { downloads: Download[] };

const DownloadsList: FunctionComponent<DownloadsListProps> = ({
    downloads
}: DownloadsListProps) => {
    return (
        <ul>
            {downloads.map(d => (
                <li key={d.id}>{d.source}</li>
            ))}
        </ul>
    );
};

export default DownloadsList;
