package undercast

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"regexp"
	"time"
)

type Download struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

type DownloadsRepository interface {
	Save(ctx context.Context, download *Download) error
	List(ctx context.Context) ([]Download, error)
}

func NewDownloadsService(repository DownloadsRepository) *DownloadsService {
	return &DownloadsService{repository: repository}
}

type DownloadsService struct {
	repository DownloadsRepository
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
	return d, nil
}

func (s *DownloadsService) List(ctx context.Context) ([]Download, error) {
	return s.repository.List(ctx)
}
