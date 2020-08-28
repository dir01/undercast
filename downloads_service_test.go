package undercast_test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestDownloadsService(t *testing.T) {
	repoMock := &mocks.DownloadsRepositoryMock{}
	downloaderMock := &mocks.DownloaderMock{}
	service := undercast.NewDownloadsService(repoMock, downloaderMock)
	downloadsServiceSuite := &DownloadsServiceSuite{
		repoMock:       repoMock,
		downloaderMock: downloaderMock,
		service:        service,
	}
	suite.Run(t, downloadsServiceSuite)
}

type DownloadsServiceSuite struct {
	suite.Suite
	repoMock       *mocks.DownloadsRepositoryMock
	downloaderMock *mocks.DownloaderMock
	service        *undercast.DownloadsService
}

func (suite *DownloadsServiceSuite) TestAddMagnet() {
	magnetUrl := "magnet:?xt=urn:btih:980E4184AEE6F326A9F9E2EE3E9D40ACAA90BC40"

	savedDownloads := make([]undercast.Download, 0, 0)
	suite.repoMock.SaveFunc = func(ctx context.Context, download *undercast.Download) error {
		savedDownloads = append(savedDownloads, *download)
		return nil
	}

	suite.downloaderMock.DownloadFunc = func(id, source string) error {
		suite.Assert().Equal(savedDownloads[0].ID, id)
		suite.Assert().Equal(savedDownloads[0].Source, source)
		return nil
	}

	d, err := suite.service.Add(context.Background(), undercast.AddDownloadRequest{
		ID:     "some-id",
		Source: magnetUrl,
	})

	suite.Require().NoError(err)

	suite.Assert().Equal("some-id", d.ID)
	suite.Assert().Equal(magnetUrl, d.Source)

	suite.Assert().Equal(magnetUrl, savedDownloads[0].Source)
	suite.Assert().Equal(int64(0), savedDownloads[0].TotalBytes)
	suite.Assert().Equal(int64(0), savedDownloads[0].CompleteBytes)

}

func (suite *DownloadsServiceSuite) TestAddInvalidSource() {
	ctx := context.Background()
	_, err := suite.service.Add(ctx, undercast.AddDownloadRequest{Source: "foo"})
	suite.Assert().Equal("Bad download source: foo", err.Error())
}
