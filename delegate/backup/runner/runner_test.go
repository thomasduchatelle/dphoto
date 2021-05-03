package runner

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
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

func (u *unitMedia) Filename() string {
	return path.Base(u.filename)
}

func (u *unitMedia) LastModificationDate() time.Time {
	return time.Now()
}

func (u *unitMedia) SimpleSignature() *model.SimpleMediaSignature {
	return &model.SimpleMediaSignature{
		RelativePath: u.filename,
		Size:         u.size,
	}
}

func (u *unitMedia) ReadMedia() (io.Reader, error) {
	return strings.NewReader(u.filename), nil
}

func newMedia(url string, size int) *unitMedia {
	return &unitMedia{url, size}
}

func TestRunner(t *testing.T) {
	a := assert.New(t)

	// given
	filter := new(MockFilter)
	downloader := new(MockDownloader)
	analyser := new(MockAnalyser)
	uploader := new(MockUploader)
	preCompletion := new(MockPreCompletion)

	r := Runner{
		MDC: log.WithField("Unit", "Test"),
		Source: func(medias chan model.FoundMedia) error {
			log.Infoln("Starting to push medias")

			medias <- newMedia("foo", 1)
			medias <- newMedia("bar", 2)
			medias <- newMedia("baz", 3)

			log.Infoln("Pushed 3 scanned medias")

			return nil
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
	filter.On("Execute", mock.Anything).Times(3).Return(func(m model.FoundMedia) bool {
		log.Infof("Filter %+v", m)
		return m.SimpleSignature().Size >= 2
	})

	downloader.On("Execute", mock.Anything).Times(2).Return(func(m model.FoundMedia) model.FoundMedia {
		log.Infof("Download %+v", m)
		return m
	}, nil)

	analyser.On("Execute", mock.Anything).Times(2).Return(func(m model.FoundMedia) *model.AnalysedMedia {
		log.Infof("Analyse %+v", m)
		return &model.AnalysedMedia{
			FoundMedia: m,
			Type:       model.MediaTypeImage,
			Signature:  nil,
			Details:    nil,
		}
	}, nil)

	uploader.On("Execute", []*model.AnalysedMedia{
		{
			FoundMedia: newMedia("bar", 2),
			Type:       model.MediaTypeImage,
		},
		{
			FoundMedia: newMedia("baz", 3),
			Type:       model.MediaTypeImage,
		},
	}).Once().Return(nil)

	preCompletion.On("Execute").Return(nil)

	// when
	completionChannel := Start(r)
	report := <-completionChannel
	a.Empty(report.Errors)

	for _, m := range []HasExpectations{filter, downloader, analyser, uploader, preCompletion} {
		m.AssertExpectations(t)
	}
}

type HasExpectations interface {
	AssertExpectations(t mock.TestingT) bool
}
