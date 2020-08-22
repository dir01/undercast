package undercast_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestDownloadsService(t *testing.T) {
	repoMock := &mocks.DownloadsRepository{}
	downloaderMock := &mocks.Downloader{}
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
	repoMock       *mocks.DownloadsRepository
	downloaderMock *mocks.Downloader
	service        *undercast.DownloadsService
}

func (suite *DownloadsServiceSuite) TestAddMagnet() {
	magnetUrl := "magnet:?xt=urn:btih:980E4184AEE6F326A9F9E2EE3E9D40ACAA90BC40"
	savedDownloads := make([]undercast.Download, 0, 0)

	suite.repoMock.On(
		"Save",
		mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("*undercast.Download"),
	).Run(func(args mock.Arguments) {
		download := args[1].(*undercast.Download)
		savedDownloads = append(savedDownloads, *download)
	}).Return(nil)

	suite.downloaderMock.On(
		"Download",
		mock.AnythingOfType("string"),
		magnetUrl,
	).Run(func(args mock.Arguments) {
		id := args[0].(string)
		source := args[1].(string)
		suite.Assert().Equal(savedDownloads[0].ID, id)
		suite.Assert().Equal(savedDownloads[0].Source, source)
	}).Return(nil)

	d, err := suite.service.Add(context.Background(), magnetUrl)

	suite.Require().NoError(err)

	suite.Assert().Equal(magnetUrl, d.Source)

	suite.Assert().Equal(magnetUrl, savedDownloads[0].Source)
	suite.Assert().Equal(int64(0), savedDownloads[0].TotalBytes)
	suite.Assert().Equal(int64(0), savedDownloads[0].CompleteBytes)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *DownloadsServiceSuite) TestAddInvalidSource() {
	ctx := context.Background()
	_, err := suite.service.Add(ctx, "foo")
	suite.Assert().Equal("Bad download source: foo", err.Error())
}
