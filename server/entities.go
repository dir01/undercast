package server

type state string

const downloading state = "DOWNLOADING"
const enconding state = "ENCODING"
const published state = "PUBLISHED"

// Torrent represents a lifecycle of a single torrent download
type Torrent struct {
	ID             int       `json:"id"`
	State          state     `json:"state"`
	Name           string    `json:"name"`
	Source         string    `json:"source"`
	FilePaths      []string  `json:"filepaths"`
	BytesCompleted int64     `json:"bytesCompleted"`
	BytesMissing   int64     `json:"bytesMissing"`
	Episodes       []Episode `json:"episodes"`
}

// NewTorrent creates new Torrent instance
func NewTorrent() *Torrent {
	return &Torrent{State: downloading}
}

// UpdateFromTorrentState updates Torrent based on data in TorrentState
func (t *Torrent) UpdateFromTorrentState(state TorrentState) {
	t.Name = state.Name
	t.FilePaths = state.FilePaths
	t.BytesCompleted = state.BytesCompleted
	t.BytesMissing = state.BytesMissing
	if state.Done {
		t.State = enconding
	}
	t.maybeSetDefaultEpisodes()
}

func (t *Torrent) maybeSetDefaultEpisodes() {
	if len(t.Episodes) > 0 || len(t.FilePaths) == 0 {
		return
	}
	t.Episodes = suggestEpisodes(t.Name, t.FilePaths)
}

// TorrentState is intended for progress reporting
type TorrentState struct {
	Name           string   `json:"name"`
	FilePaths      []string `json:"filepaths"`
	BytesCompleted int64    `json:"bytesCompleted"`
	BytesMissing   int64    `json:"bytesMissing"`
	Done           bool     `json:"done"`
}

type Episode struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	FilePaths []string `json:"filepaths"`
}
