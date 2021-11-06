package runner

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/delegate/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"path"
	"strings"
	"testing"
	"time"
)

type unitMedia struct {
	filename string
	size     int
}

func (u *unitMedia) MediaPath() backupmodel.MediaPath {
	return backupmodel.MediaPath{
		ParentFullPath: path.Dir(u.filename),
		Root:           path.Dir(u.filename),
		Path:           "",
		Filename:       path.Base(u.filename),
		ParentDir:      path.Dir(u.filename),
	}
}

func (u *unitMedia) LastModificationDate() time.Time {
	return time.Now()
}

func (u *unitMedia) SimpleSignature() *backupmodel.SimpleMediaSignature {
	return &backupmodel.SimpleMediaSignature{
		RelativePath: u.filename,
		Size:         uint(u.size),
	}
}

func (u *unitMedia) ReadMedia() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(u.filename)), nil
}

func newMedia(url string, size int) *unitMedia {
	return &unitMedia{url, size}
}

func TestRunner(t *testing.T) {
	a := assert.New(t)

	// given
	filter := new(mocks.Filter)
	downloader := new(mocks.Downloader)
	analyser := new(mocks.Analyser)
	uploader := new(mocks.Uploader)
	preCompletion := new(mocks.PreCompletion)

	r := Runner{
		MDC: log.WithField("Unit", "Test"),
		Source: func(medias chan backupmodel.FoundMedia) (uint, uint, error) {
			log.Infoln("Starting to push medias")

			medias <- newMedia("foo", 1)
			medias <- newMedia("bar", 2)
			medias <- newMedia("baz", 3)

			log.Infoln("Pushed 3 scanned medias")

			return 3, 6, nil
		},
		Filter:               filter.Execute,
		Downloader:           downloader.Execute,
		Analyser:             analyser.Execute,
		Uploader:             uploader.Execute,
		PreCompletion:        preCompletion.Execute,
		BufferSize:           10,
		ConcurrentDownloader: 1,
		ConcurrentAnalyser:   1,
		ConcurrentUploader:   1,
		UploadBatchSize:      10,
	}

	// pre-then
	filter.On("Execute", mock.Anything).Times(3).Return(func(m backupmodel.FoundMedia) bool {
		log.Infof("Filter %+v", m)
		return m.SimpleSignature().Size >= 2
	})

	downloader.On("Execute", mock.Anything).Times(2).Return(func(m backupmodel.FoundMedia) backupmodel.FoundMedia {
		log.Infof("Download %+v", m)
		return m
	}, nil)

	analyser.On("Execute", mock.Anything).Times(2).Return(func(m backupmodel.FoundMedia) *backupmodel.AnalysedMedia {
		log.Infof("Analyse %+v", m)
		return &backupmodel.AnalysedMedia{
			FoundMedia: m,
			Type:       backupmodel.MediaTypeImage,
			Signature:  nil,
			Details:    nil,
		}
	}, nil)

	uploader.On("Execute", []*backupmodel.AnalysedMedia{
		{
			FoundMedia: newMedia("bar", 2),
			Type:       backupmodel.MediaTypeImage,
		},
		{
			FoundMedia: newMedia("baz", 3),
			Type:       backupmodel.MediaTypeImage,
		},
	}, mock.Anything).Once().Return(nil)

	preCompletion.On("Execute").Return(nil)

	// when
	completionChannel, progressChannel := Start(r)
	DummyProgressListener(progressChannel)

	report := <-completionChannel
	a.Empty(report.Errors)

	for _, m := range []HasExpectations{filter, downloader, analyser, uploader, preCompletion} {
		m.AssertExpectations(t)
	}
}

type HasExpectations interface {
	AssertExpectations(t mock.TestingT) bool
}
