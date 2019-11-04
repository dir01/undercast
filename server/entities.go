package server

type state string

const added state = "ADDED"
const downloaded state = "DOWNLOADED"

type Torrent struct {
	ID             int       `json:"id"`
	State          state     `json:"state"`
	Name           string    `json:"name"`
	Source         string    `json:"source"`
	FileNames      []string  `json:"filenames"`
	BytesCompleted int64     `json:"bytesCompleted"`
	BytesMissing   int64     `json:"bytesMissing"`
	Episodes       []Episode `json:"episodes"`
}

func (t *Torrent) markAsAdded() {
	t.State = added
}

func (t *Torrent) markAsDownloaded() {
	t.State = downloaded
}

func (t *Torrent) maybeSetDefaultEpisodes() {
	if len(t.Episodes) > 0 || len(t.FileNames) == 0 {
		return
	}
	t.Episodes = suggestEpisodes(t.Name, t.FileNames)
}

type Episode struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	FileNames []string `json:"filenames"`
}
