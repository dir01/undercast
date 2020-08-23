// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"sync"
	"undercast"
)

// Ensure, that DownloaderMock does implement undercast.Downloader.
// If this is not the case, regenerate this file with moq.
var _ undercast.Downloader = &DownloaderMock{}

// DownloaderMock is a mock implementation of undercast.Downloader.
//
//     func TestSomethingThatUsesDownloader(t *testing.T) {
//
//         // make and configure a mocked undercast.Downloader
//         mockedDownloader := &DownloaderMock{
//             DownloadFunc: func(id string, source string) error {
// 	               panic("mock out the Download method")
//             },
//             IsMatchingFunc: func(source string) bool {
// 	               panic("mock out the IsMatching method")
//             },
//             OnProgressFunc: func(in1 func(id string, downloadInfo *undercast.DownloadInfo))  {
// 	               panic("mock out the OnProgress method")
//             },
//         }
//
//         // use mockedDownloader in code that requires undercast.Downloader
//         // and then make assertions.
//
//     }
type DownloaderMock struct {
	// DownloadFunc mocks the Download method.
	DownloadFunc func(id string, source string) error

	// IsMatchingFunc mocks the IsMatching method.
	IsMatchingFunc func(source string) bool

	// OnProgressFunc mocks the OnProgress method.
	OnProgressFunc func(in1 func(id string, downloadInfo *undercast.DownloadInfo))

	// calls tracks calls to the methods.
	calls struct {
		// Download holds details about calls to the Download method.
		Download []struct {
			// ID is the id argument value.
			ID string
			// Source is the source argument value.
			Source string
		}
		// IsMatching holds details about calls to the IsMatching method.
		IsMatching []struct {
			// Source is the source argument value.
			Source string
		}
		// OnProgress holds details about calls to the OnProgress method.
		OnProgress []struct {
			// In1 is the in1 argument value.
			In1 func(id string, downloadInfo *undercast.DownloadInfo)
		}
	}
	lockDownload   sync.RWMutex
	lockIsMatching sync.RWMutex
	lockOnProgress sync.RWMutex
}

// Download calls DownloadFunc.
func (mock *DownloaderMock) Download(id string, source string) error {
	if mock.DownloadFunc == nil {
		panic("DownloaderMock.DownloadFunc: method is nil but Downloader.Download was just called")
	}
	callInfo := struct {
		ID     string
		Source string
	}{
		ID:     id,
		Source: source,
	}
	mock.lockDownload.Lock()
	mock.calls.Download = append(mock.calls.Download, callInfo)
	mock.lockDownload.Unlock()
	return mock.DownloadFunc(id, source)
}

// DownloadCalls gets all the calls that were made to Download.
// Check the length with:
//     len(mockedDownloader.DownloadCalls())
func (mock *DownloaderMock) DownloadCalls() []struct {
	ID     string
	Source string
} {
	var calls []struct {
		ID     string
		Source string
	}
	mock.lockDownload.RLock()
	calls = mock.calls.Download
	mock.lockDownload.RUnlock()
	return calls
}

// IsMatching calls IsMatchingFunc.
func (mock *DownloaderMock) IsMatching(source string) bool {
	if mock.IsMatchingFunc == nil {
		panic("DownloaderMock.IsMatchingFunc: method is nil but Downloader.IsMatching was just called")
	}
	callInfo := struct {
		Source string
	}{
		Source: source,
	}
	mock.lockIsMatching.Lock()
	mock.calls.IsMatching = append(mock.calls.IsMatching, callInfo)
	mock.lockIsMatching.Unlock()
	return mock.IsMatchingFunc(source)
}

// IsMatchingCalls gets all the calls that were made to IsMatching.
// Check the length with:
//     len(mockedDownloader.IsMatchingCalls())
func (mock *DownloaderMock) IsMatchingCalls() []struct {
	Source string
} {
	var calls []struct {
		Source string
	}
	mock.lockIsMatching.RLock()
	calls = mock.calls.IsMatching
	mock.lockIsMatching.RUnlock()
	return calls
}

// OnProgress calls OnProgressFunc.
func (mock *DownloaderMock) OnProgress(in1 func(id string, downloadInfo *undercast.DownloadInfo)) {
	if mock.OnProgressFunc == nil {
		panic("DownloaderMock.OnProgressFunc: method is nil but Downloader.OnProgress was just called")
	}
	callInfo := struct {
		In1 func(id string, downloadInfo *undercast.DownloadInfo)
	}{
		In1: in1,
	}
	mock.lockOnProgress.Lock()
	mock.calls.OnProgress = append(mock.calls.OnProgress, callInfo)
	mock.lockOnProgress.Unlock()
	mock.OnProgressFunc(in1)
}

// OnProgressCalls gets all the calls that were made to OnProgress.
// Check the length with:
//     len(mockedDownloader.OnProgressCalls())
func (mock *DownloaderMock) OnProgressCalls() []struct {
	In1 func(id string, downloadInfo *undercast.DownloadInfo)
} {
	var calls []struct {
		In1 func(id string, downloadInfo *undercast.DownloadInfo)
	}
	mock.lockOnProgress.RLock()
	calls = mock.calls.OnProgress
	mock.lockOnProgress.RUnlock()
	return calls
}
