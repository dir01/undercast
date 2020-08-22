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
	onProgress    func(id string, di *DownloadInfo)
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
			di := &DownloadInfo{
				TotalBytes:         bytesMissing + bytesCompleted,
				CompleteBytes:      bytesCompleted,
				IsDownloadComplete: isComplete,
				Name:               t.Name(),
			}
			td.onProgress(id, di)
		}
	}()

	return nil
}

func (td *TorrentsDownloader) OnProgress(onProgress func(id string, di *DownloadInfo)) {
	td.onProgress = onProgress
}
