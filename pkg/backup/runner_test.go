package backup

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldStopAtFirstError(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(run *runner)
		sizeHint int
		want     [][]string
		wantErr  string
	}{
		{
			name:     "it should run through with channel buffer enabled",
			modifier: func(run *runner) {},
			sizeHint: 42,
			want:     [][]string{{"file_1.jpg", "file_2.jpg"}, {"file_3.jpg", "file_4.jpg"}, {"file_5.jpg"}},
			wantErr:  "",
		},
		{
			name:     "it should run through with channel buffer disabled",
			modifier: func(run *runner) {},
			want:     [][]string{{"file_1.jpg", "file_2.jpg"}, {"file_3.jpg", "file_4.jpg"}, {"file_5.jpg"}},
			wantErr:  "",
		},
		{
			name: "it should re-buffer before assigning ids and before uploading (to not have half empty batch)",
			modifier: func(run *runner) {
				run.Cataloger = RunnerCatalogerFunc(func(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
					var requests []*BackingUpMediaRequest
					for _, media := range medias {
						if media.FoundMedia.MediaPath().Filename != "file_2.jpg" && media.FoundMedia.MediaPath().Filename != "file_3.jpg" {
							requests = append(requests, &BackingUpMediaRequest{
								AnalysedMedia:    media,
								CatalogReference: &CatalogReferenceStub{MediaIdValue: media.Sha256Hash, AlbumFolderNameValue: "/album1"},
							})
						}
					}

					return requests, nil
				})
				run.UniqueFilter = func(media *BackingUpMediaRequest, progressChannel chan *ProgressEvent) bool {
					return media.CatalogReference.UniqueIdentifier() != "file_6.jpg"
				}
			},
			want:    [][]string{{"file_1.jpg", "file_4.jpg"}, {"file_5.jpg"}},
			wantErr: "",
		},
		{
			name: "it should stop the process after an error on the analyser: buffers are NOT flushed",
			modifier: func(run *runner) {
				run.Analyser = &AnalyserFake{
					ErroredFilename: map[string]error{
						"file_4.jpg": errors.Errorf("[TEST] analyser error"),
					},
				}
			},
			want:    [][]string{{"file_1.jpg", "file_2.jpg"}},
			wantErr: "error in analyser: [TEST] analyser error",
		},
		{
			name: "it should stop the process after an error on the cataloguer",
			modifier: func(run *runner) {
				original := run.Cataloger
				run.Cataloger = RunnerCatalogerFunc(func(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
					if medias[0].FoundMedia.MediaPath().Filename == "file_3.jpg" {
						return nil, errors.Errorf("TEST")
					}

					return original.Catalog(context.TODO(), medias, progressChannel)
				})
			},
			want:    [][]string{{"file_1.jpg", "file_2.jpg"}},
			wantErr: "error in cataloguer: TEST",
		},
		{
			name: "it should stop the process after an error on the uploader",
			modifier: func(run *runner) {
				original := run.Uploader
				run.Uploader = RunnerUploaderFunc(func(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error {
					if buffer[0].AnalysedMedia.FoundMedia.MediaPath().Filename == "file_1.jpg" {
						return errors.Errorf("TEST")
					}

					return original.Upload(buffer, progressChannel)
				})
			},
			want:    nil,
			wantErr: "error in uploader: TEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			run, capture := newMockedRun(mockPublisher(
				NewInMemoryMedia("file_1.jpg", time.Now(), []byte("3.14")),
				NewInMemoryMedia("file_2.jpg", time.Now(), []byte("3.1415")),
				NewInMemoryMedia("file_3.jpg", time.Now(), []byte("3.141592")),
				NewInMemoryMedia("file_4.jpg", time.Now(), []byte("3.141592")),
				NewInMemoryMedia("file_5.jpg", time.Now(), []byte("3.141592")),
			))

			tt.modifier(run)

			eventsChannel, completion := run.start(context.Background(), tt.sizeHint)
			_ = collectEvents(eventsChannel)

			caughtErrors := <-completion

			if tt.wantErr == "" {
				a.Empty(caughtErrors, tt.name)

				a.Equal(tt.want, capture.requests, tt.name)
			} else {
				if a.Len(caughtErrors, 1, tt.name) {
					a.Equal(tt.wantErr, caughtErrors[0].Error(), tt.name)
				}
			}
		})
	}
}

// ** UTILS

func mockPublisher(medias ...FoundMedia) runnerPublisher {
	return func(channel chan FoundMedia, events chan *ProgressEvent) error {
		for _, media := range medias {
			channel <- media
		}

		return nil
	}
}

type captureStruct struct {
	requests [][]string
}

func newMockedRun(publisher runnerPublisher) (*runner, *captureStruct) {
	uploadedCapture := new(captureStruct)

	run := &runner{
		MDC:       log.WithField("Testing", "Test"),
		Publisher: publisher,
		Analyser:  new(AnalyserFake),
		Cataloger: RunnerCatalogerFunc(func(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
			var requests []*BackingUpMediaRequest
			for _, media := range medias {
				requests = append(requests, &BackingUpMediaRequest{
					AnalysedMedia:    media,
					CatalogReference: &CatalogReferenceStub{MediaIdValue: media.Sha256Hash, AlbumFolderNameValue: "/album1"},
				})
			}

			return requests, nil
		}),
		UniqueFilter: func(media *BackingUpMediaRequest, progressChannel chan *ProgressEvent) bool {
			return true
		},
		Uploader: RunnerUploaderFunc(func(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error {
			var names []string
			for _, request := range buffer {
				names = append(names, request.AnalysedMedia.FoundMedia.MediaPath().Filename)
			}

			uploadedCapture.requests = append(uploadedCapture.requests, names)
			return nil
		}),
		Options: Options{
			ConcurrentAnalyser:   1,
			ConcurrentCataloguer: 1,
			ConcurrentUploader:   1,
			BatchSize:            2,
		},
	}

	return run, uploadedCapture
}

func collectEvents(channel chan *ProgressEvent) map[ProgressEventType]eventSummary {
	eventCapture := newEventCapture()

	for event := range channel {
		eventCapture.OnEvent(*event)
	}

	return eventCapture.Captured
}

type AnalyserFake struct {
	ErroredFilename map[string]error
}

func (a *AnalyserFake) Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver) error {
	if err, errored := a.ErroredFilename[found.MediaPath().Filename]; errored {
		return err
	}

	return analysedMediaObserver.OnAnalysedMedia(ctx, &AnalysedMedia{
		FoundMedia: found,
		Type:       MediaTypeImage,
		Sha256Hash: found.MediaPath().Filename,
		Details: &MediaDetails{
			DateTime: time.Date(2022, 6, 18, 10, 42, 0, 0, time.UTC),
		},
	})
}
