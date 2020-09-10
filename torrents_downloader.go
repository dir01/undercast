package undercast

import (
	"fmt"
	anacrolix "github.com/anacrolix/torrent"
	"time"
)

func NewTorrentsDownloader(dataDir string) *TorrentsDownloader {
	cfg := anacrolix.NewDefaultClientConfig()
	cfg.DataDir = dataDir
	anacrolixClient, err := anacrolix.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}
	return &TorrentsDownloader{
		torrentClient: anacrolixClient,
		dataDir:       dataDir,
	}
}

type TorrentsDownloader struct {
	torrentClient *anacrolix.Client
	dataDir       string
	onProgress    func(id string, p *DownloadProgress)
	onInfo        func(id string, di *DownloadInfo)
}

func (td *TorrentsDownloader) IsMatching(source string) bool {
	return true
}

func (td *TorrentsDownloader) Download(id string, source string) error {
	var t *anacrolix.Torrent
	t, err := td.torrentClient.AddMagnet(source)
	if err != nil {
		return err
	}

	go func() {
		<-t.GotInfo()
		t.DownloadAll()

		if td.onInfo != nil {
			go td.onInfo(id, &DownloadInfo{
				Name:    t.Name(),
				Files:   extractFilenames(t.Files()),
				RootDir: td.dataDir,
			})
		}

		isComplete := false
		for !isComplete {
			time.Sleep(1 * time.Second)
			if td.onProgress == nil {
				continue
			}
			bytesCompleted := t.BytesCompleted()
			bytesMissing := t.BytesMissing()
			if bytesMissing == 0 {
				isComplete = true
			}
			p := &DownloadProgress{
				TotalBytes:         bytesMissing + bytesCompleted,
				CompleteBytes:      bytesCompleted,
				IsDownloadComplete: isComplete,
			}
			go td.onProgress(id, p)
		}
	}()

	return nil
}

func (td *TorrentsDownloader) OnProgress(onProgress func(id string, p *DownloadProgress)) {
	td.onProgress = onProgress
}

func (td *TorrentsDownloader) OnInfo(onInfo func(id string, di *DownloadInfo)) {
	td.onInfo = onInfo
}

func extractFilenames(files []*anacrolix.File) []string {
	result := make([]string, 0, len(files))
	for _, f := range files {
		result = append(result, f.Path())
	}
	return result
}
