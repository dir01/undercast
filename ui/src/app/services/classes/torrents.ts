export class Torrent {
  id: number;
  name: string;
  files: string[];
  bytesComplete: number;
  bytesMissing: number;
  done: boolean;
}
