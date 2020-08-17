package undercast_test

import (
	"context"
	"gotest.tools/assert"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestDownloadsService_Add(t *testing.T) {
	service := undercast.NewDownloadsService(&mocks.DownloadsRepository{})
	ctx := context.Background()

	// Bad source format
	_, err := service.Add(ctx, "foo")
	assert.Equal(t, "Bad download source: foo", err.Error())
}
