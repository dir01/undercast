package undercast

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type DownloadPathsProvider interface {
	AbsPaths(ctx context.Context, downloadId string, relativeFilenames []string) ([]string, error)
}

type MediaRepository interface {
	Save(ctx context.Context, media *Media) error
}

type Media struct {
	ID         string   `json:"id"`
	DownloadId string   `json:"downloadId"`
	Files      []string `json:"files"`
}

type MediaService struct {
	downloadPathsProvider DownloadPathsProvider
	repository            MediaRepository
}

type CreateMediaRequest struct {
	ID         string   `json:"id"`
	DownloadId string   `json:"downloadId"`
	Files      []string `json:"files"`
}

func (service *MediaService) Create(ctx context.Context, req CreateMediaRequest) (*Media, error) {
	absPaths, err := service.downloadPathsProvider.AbsPaths(ctx, req.DownloadId, req.Files)
	if err != nil {
		return nil, err
	}
	resultFilepath, err := glueAudioFiles(absPaths)
	if err != nil {
		return nil, err
	}
	_ = resultFilepath
	return &Media{ID: req.ID, DownloadId: req.DownloadId, Files: req.Files}, nil

}

func glueAudioFiles(filepaths []string) (string, error) {
	output := path.Join(os.TempDir(), uuid.NewV4().String()+".mp3")
	args := []string{"-i", "concat:" + strings.Join(filepaths, "|"), "-acodec", "copy", output}
	cmd := exec.Command("ffmpeg", args...)
	log.Printf("Running %s\n", cmd)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return output, nil
}
