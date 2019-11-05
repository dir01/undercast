package server

type state string

const added state = "ADDED"
const downloaded state = "DOWNLOADED"

// Torrent represents a lifecycle of a single torrent download
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

// NewTorrent creates new Torrent instance
func NewTorrent() *Torrent {
	return &Torrent{State: added}
}

// UpdateFromTorrentState updates Torrent based on data in TorrentState
func (t *Torrent) UpdateFromTorrentState(state TorrentState) {
	t.Name = state.Name
	t.FileNames = state.FileNames
	t.BytesCompleted = state.BytesCompleted
	t.BytesMissing = state.BytesMissing
	if state.Done {
		t.State = downloaded
	}
	t.maybeSetDefaultEpisodes()
}

func (t *Torrent) maybeSetDefaultEpisodes() {
	if len(t.Episodes) > 0 || len(t.FileNames) == 0 {
		return
	}
	t.Episodes = suggestEpisodes(t.Name, t.FileNames)
}

// TorrentState is intended for progress reporting
type TorrentState struct {
	Name           string   `json:"name"`
	FileNames      []string `json:"filenames"`
	BytesCompleted int64    `json:"bytesCompleted"`
	BytesMissing   int64    `json:"bytesMissing"`
	Done           bool     `json:"done"`
}

type Episode struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	FileNames []string `json:"filenames"`
}
