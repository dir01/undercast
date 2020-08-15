package undercast

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type Download struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type downloadsService struct {
	repository *downloadsRepository
}

func (s *downloadsService) Add(ctx context.Context, source string) (*Download, error) {
	d := &Download{
		ID:     uuid.NewV4().String(),
		Source: source,
	}
	err := s.repository.Save(ctx, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *downloadsService) List(ctx context.Context) ([]Download, error) {
	return s.repository.List(ctx)
}
