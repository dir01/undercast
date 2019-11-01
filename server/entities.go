package server

type state string

const added state = "ADDED"
const downloaded state = "DOWNLOADED"

type Torrent struct {
	ID             int      `json:"id"`
	State          state    `json:"state"`
	Name           string   `json:"name"`
	Source         string   `json:"source"`
	FileNames      []string `json:"filenames"`
	BytesCompleted int64    `json:"bytesCompleted"`
	BytesMissing   int64    `json:"bytesMissing"`
}

func (t *Torrent) markAsAdded() {
	t.State = added
}

func (t *Torrent) markAsDownloaded() {
	t.State = downloaded
}

type Episode struct {
	Name      string
	Filenames []string
}
