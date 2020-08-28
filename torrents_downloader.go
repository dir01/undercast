package undercast

import (
	"fmt"
	anacrolix "github.com/anacrolix/torrent"
	"time"
)

func NewTorrentsDownloader() *TorrentsDownloader {
	cfg := anacrolix.NewDefaultClientConfig()
	cfg.DataDir = "./data"
	anacrolixClient, err := anacrolix.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}
	return &TorrentsDownloader{
		torrentClient: anacrolixClient,
	}
}

type TorrentsDownloader struct {
	torrentClient *anacrolix.Client
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
			td.onInfo(id, &DownloadInfo{
				Name:  t.Name(),
				Files: extractFilenames(t.Files()),
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
			td.onProgress(id, p)
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
