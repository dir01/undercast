export class Torrent {
  id: number;
  source: string;
  name: string;
  files: string[];
  bytesComplete: number;
  bytesMissing: number;
  episods: Object[];
  done: boolean;
}
