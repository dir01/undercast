package undercast

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"regexp"
	"time"
)

type Download struct {
	ID                 string    `json:"id"`
	Source             string    `json:"source"`
	Name               string    `json:"name"`
	CreatedAt          time.Time `json:"createdAt"`
	TotalBytes         int64     `json:"totalBytes"`
	CompleteBytes      int64     `json:"completeBytes"`
	IsDownloadComplete bool
}

//go:generate moq -out ./mocks/DownloadsRepository.go -pkg mocks . DownloadsRepository
type DownloadsRepository interface {
	Save(ctx context.Context, download *Download) error
	List(ctx context.Context) ([]Download, error)
	GetById(ctx context.Context, id string) (*Download, error)
	ListIncomplete(ctx context.Context) ([]Download, error)
}

type DownloadInfo struct {
	Name               string
	TotalBytes         int64
	CompleteBytes      int64
	IsDownloadComplete bool
}

//go:generate moq -out ./mocks/Downloader.go -pkg mocks . Downloader
type Downloader interface {
	IsMatching(source string) bool
	Download(id string, source string) error
	OnProgress(func(id string, downloadInfo *DownloadInfo))
}

func NewDownloadsService(repository DownloadsRepository, downloader Downloader) *DownloadsService {
	return &DownloadsService{repository: repository, downloader: downloader}
}

type DownloadsService struct {
	repository DownloadsRepository
	downloader Downloader
}

func (s *DownloadsService) Run() {
	ctx := context.Background()

	s.downloader.OnProgress(func(id string, di *DownloadInfo) {
		download, err := s.repository.GetById(ctx, id)
		if err != nil {
			log.Printf("Repository failed to load download with id %s:\n%s\n", id, err.Error())
			return
		}
		download.CompleteBytes = di.CompleteBytes
		download.TotalBytes = di.TotalBytes
		download.IsDownloadComplete = di.IsDownloadComplete
		download.Name = di.Name
		err = s.repository.Save(context.Background(), download)
		if err != nil {
			log.Printf("Repository failed to save download with id %s:\n%s\n", id, err.Error())
		}
	})

	incomplete, err := s.repository.ListIncomplete(ctx)
	if err != nil {
		log.Printf("Error while fetching incomplete downloads:\n%s\n", err.Error())
	}
	log.Printf("Resuming %d incomplete downloads", len(incomplete))

	for _, d := range incomplete {
		err = s.downloader.Download(d.ID, d.Source)
		if err != nil {
			log.Printf("Error while resuming download %s:\n%s\n", d.ID, err.Error())
		}
	}
}

func (s *DownloadsService) Add(ctx context.Context, source string) (*Download, error) {
	match, _ := regexp.MatchString("magnet:\\?xt=urn:btih:[A-Z0-9]{20,50}", source)
	if !match {
		return nil, fmt.Errorf("Bad download source: %s", source)
	}

	d := &Download{
		ID:        uuid.NewV4().String(),
		Source:    source,
		CreatedAt: time.Now(),
	}

	err := s.repository.Save(ctx, d)
	if err != nil {
		return nil, err
	}

	err = s.downloader.Download(d.ID, d.Source)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *DownloadsService) List(ctx context.Context) ([]Download, error) {
	return s.repository.List(ctx)
}
