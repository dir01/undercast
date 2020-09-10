package undercast

import (
	"context"
	"fmt"
	"log"
	"path"
	"regexp"
	"sort"
	"time"
)

type Download struct {
	ID                 string    `json:"id"`
	Source             string    `json:"source"`
	Name               string    `json:"name"`
	CreatedAt          time.Time `json:"createdAt"`
	TotalBytes         int64     `json:"totalBytes"`
	CompleteBytes      int64     `json:"completeBytes"`
	Files              []string  `json:"files"`
	RootDir            string
	IsDownloadComplete bool
}

func (d *Download) AbsPaths(relFilePaths []string) (absFilePaths []string) {
	absFilePaths = make([]string, 0, len(relFilePaths))
	for _, p := range relFilePaths {
		absFilePaths = append(absFilePaths, path.Join(d.RootDir, p))
	}
	return absFilePaths
}

//go:generate moq -out ./mocks/DownloadsRepository.go -pkg mocks . DownloadsRepository
type DownloadsRepository interface {
	Save(ctx context.Context, download *Download) error
	List(ctx context.Context) ([]Download, error)
	GetById(ctx context.Context, id string) (*Download, error)
	ListIncomplete(ctx context.Context) ([]Download, error)
}

type DownloadInfo struct {
	Name    string
	RootDir string
	Files   []string
}

type DownloadProgress struct {
	TotalBytes         int64
	CompleteBytes      int64
	IsDownloadComplete bool
}

//go:generate moq -stub -out ./mocks/Downloader.go -pkg mocks . Downloader
type Downloader interface {
	IsMatching(source string) bool
	Download(id string, source string) error
	OnInfo(onInfo func(id string, info *DownloadInfo))
	OnProgress(onProgress func(id string, progress *DownloadProgress))
}

func NewDownloadsService(repository DownloadsRepository, downloader Downloader) *DownloadsService {
	return &DownloadsService{repository: repository, downloader: downloader}
}

type DownloadsService struct {
	repository         DownloadsRepository
	downloader         Downloader
	onDownloadCreated  func(download *Download)
	onDownloadComplete func(download *Download)
}

func (srv *DownloadsService) Run() {
	ctx := context.Background()

	srv.downloader.OnInfo(func(id string, i *DownloadInfo) {
		d, err := srv.repository.GetById(ctx, id)
		if err != nil {
			log.Printf("Repository failed to load d with id %s:\n%s\n", id, err.Error())
			return
		}
		if d.Name != "" && len(d.Files) != 0 && d.RootDir != "" {
			log.Printf("Download %s already got info, bailing\n", d.ID)
			return
		}
		d.Name = i.Name
		d.Files = i.Files
		d.RootDir = i.RootDir
		sort.Strings(d.Files)
		err = srv.repository.Save(context.Background(), d)
		if err != nil {
			log.Printf("Repository failed to save d with id %s:\n%s\n", id, err.Error())
		}
		if srv.onDownloadCreated != nil {
			srv.onDownloadCreated(d)
		}
	})

	srv.downloader.OnProgress(func(id string, di *DownloadProgress) {
		download, err := srv.repository.GetById(ctx, id)
		if err != nil {
			log.Printf("Repository failed to load download with id %s:\n%s\n", id, err.Error())
			return
		}
		download.CompleteBytes = di.CompleteBytes
		download.TotalBytes = di.TotalBytes
		download.IsDownloadComplete = di.IsDownloadComplete
		err = srv.repository.Save(context.Background(), download)
		if err != nil {
			log.Printf("Repository failed to save download with id %s:\n%s\n", id, err.Error())
		}
		if download.IsDownloadComplete && srv.onDownloadComplete != nil {
			srv.onDownloadComplete(download)
		}
	})

	incomplete, err := srv.repository.ListIncomplete(ctx)
	if err != nil {
		log.Printf("Error while fetching incomplete downloads:\n%s\n", err.Error())
	}
	log.Printf("Found %d incomplete downloads", len(incomplete))

	for _, d := range incomplete {
		log.Printf("Resuming download %s (%s)", d.ID, d.Name)
		err = srv.downloader.Download(d.ID, d.Source)
		if err != nil {
			log.Printf("Error while resuming download %s:\n%s\n", d.ID, err.Error())
		}
	}
}

type AddDownloadRequest struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

func (srv *DownloadsService) Add(ctx context.Context, req AddDownloadRequest) (*Download, error) {
	match, _ := regexp.MatchString("magnet:\\?xt=urn:btih:[A-Z0-9]{20,50}", req.Source)
	if !match {
		return nil, fmt.Errorf("Bad download source: %s", req.Source)
	}

	d := &Download{ID: req.ID, Source: req.Source, CreatedAt: time.Now()}

	err := srv.repository.Save(ctx, d)
	if err != nil {
		return nil, err
	}

	err = srv.downloader.Download(d.ID, d.Source)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (srv *DownloadsService) GetById(ctx context.Context, id string) (*Download, error) {
	return srv.repository.GetById(ctx, id)
}

func (srv *DownloadsService) List(ctx context.Context) ([]Download, error) {
	return srv.repository.List(ctx)
}

func (srv *DownloadsService) OnDownloadCreated(onDownloadCreated func(download *Download)) {
	srv.onDownloadCreated = onDownloadCreated
}

func (srv *DownloadsService) OnDownloadComplete(onDownloadComplete func(download *Download)) {
	srv.onDownloadComplete = onDownloadComplete
}
