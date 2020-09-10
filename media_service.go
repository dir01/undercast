package undercast

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
)

//go:generate moq -stub -out ./mocks/MediaDownloadsService.go -pkg mocks . MediaDownloadsService
type MediaDownloadsService interface {
	OnDownloadCreated(func(download *Download))
	OnDownloadComplete(func(download *Download))
	GetById(ctx context.Context, downloadId string) (download *Download, err error)
}

//go:generate moq -stub -out ./mocks/MediaRepository.go -pkg mocks . MediaRepository
type MediaRepository interface {
	Save(ctx context.Context, media *Media) error
	GetMedia(ctx context.Context, mediaId string) (*Media, error)
	ListByDownloadId(ctx context.Context, downloadId string) ([]Media, error)
}

//go:generate moq -stub -out ./mocks/MediaStorage.go -pkg mocks . MediaStorage
type MediaStorage interface {
	Store(ctx context.Context, filepath, filename string) (url string, err error)
}

//go:generate moq -stub -out ./mocks/MediaConverter.go -pkg mocks . MediaConverter
type MediaConverter interface {
	Concatenate(filepaths []string, filename string, format string) (resultFilePath string, err error)
}

type mediaState string

const waiting mediaState = "waiting"

const uploaded mediaState = "uploaded"

type Media struct {
	ID         string   `json:"id"`
	DownloadId string   `json:"downloadId"`
	Files      []string `json:"files"`
	Url        string   `json:"url"`
	State      mediaState
}

type MediaService struct {
	downloadsService MediaDownloadsService
	repository       MediaRepository
	converter        MediaConverter
	storage          MediaStorage
	onMediaUploaded  func(media *Media)
}

type CreateMediaRequest struct {
	ID         string   `json:"id"`
	DownloadId string   `json:"downloadId"`
	Files      []string `json:"files"`
}

func (srv *MediaService) Run() {
	srv.downloadsService.OnDownloadCreated(func(d *Download) {
		srv.createDefaultMedias(d)
	})
	srv.downloadsService.OnDownloadComplete(func(d *Download) {
		medias, err := srv.repository.ListByDownloadId(context.TODO(), d.ID)
		if err != nil {
			log.Printf("failed to list media of download %s: %s", d.ID, err)
		}
		for _, m := range medias {
			go srv.convertAndUpload(m, d)
		}
	})
}

func (srv *MediaService) Create(ctx context.Context, req CreateMediaRequest) (*Media, error) {
	download, err := srv.downloadsService.GetById(ctx, req.DownloadId)
	if err != nil {
		return nil, err
	}
	media := &Media{ID: req.ID, DownloadId: req.DownloadId, Files: req.Files}
	if err := srv.repository.Save(ctx, media); err != nil {
		return nil, err
	}
	if download.IsDownloadComplete {
		go srv.convertAndUpload(*media, download)
	}
	return media, err
}

func (srv *MediaService) GetMedia(ctx context.Context, mediaId string) (*Media, error) {
	return srv.repository.GetMedia(ctx, mediaId)
}

func (srv *MediaService) OnMediaUploaded(onMediaUploaded func(media *Media)) {
	srv.onMediaUploaded = onMediaUploaded
}

func (srv *MediaService) createDefaultMedias(d *Download) {
	// Currently we just create one media containing all of download's files
	// In recent future, we need to become a bit smarter about this
	media := &Media{ID: uuid.NewV4().String(), DownloadId: d.ID, Files: d.Files, State: waiting}
	if err := srv.repository.Save(context.TODO(), media); err != nil {
		fmt.Printf("Failed to save default media: %s", err)
	}
}

func (srv *MediaService) convertAndUpload(media Media, download *Download) {
	filename := fmt.Sprintf("%s.mp3", media.ID)
	filepaths := download.AbsPaths(media.Files)
	resultFilepath, err := srv.converter.Concatenate(filepaths, filename, "mp3")
	if err != nil {
		log.Printf("Error while converting media %s: %s", media.ID, err)
		return
	}
	log.Printf("Uploading %s", filename)
	url, err := srv.storage.Store(context.TODO(), resultFilepath, filename)
	log.Printf("Finished uploading %s as %s", resultFilepath, url)
	media.Url = url
	media.State = uploaded
	if err = srv.repository.Save(context.TODO(), &media); err != nil {
		log.Printf("Error while saving media: %s", err)
		return
	}
	if srv.onMediaUploaded != nil {
		srv.onMediaUploaded(&media)
	}
}
