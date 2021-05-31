package tracker

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTracker(t *testing.T) {
	a := assert.New(t)

	channel := make(chan *model.ProgressEvent, 256)
	tracker := NewTracker(channel, nil)

	channel <- &model.ProgressEvent{Type: model.ProgressEventScanComplete, Count: 3, Size: 42}
	channel <- &model.ProgressEvent{Type: model.ProgressEventSkipped, Count: 1, Size: 20}
	channel <- &model.ProgressEvent{Type: model.ProgressEventDownloaded, Count: 1, Size: 10}
	channel <- &model.ProgressEvent{Type: model.ProgressEventAnalysed, Count: 1, Size: 10}
	channel <- &model.ProgressEvent{Type: model.ProgressEventSkipped, Count: 1, Size: 10}
	channel <- &model.ProgressEvent{Type: model.ProgressEventDownloaded, Count: 1, Size: 12}
	channel <- &model.ProgressEvent{Type: model.ProgressEventAnalysed, Count: 1, Size: 12}
	channel <- &model.ProgressEvent{Type: model.ProgressEventUploaded, Count: 1, Size: 12, Album: "albums/007", MediaType: model.MediaTypeImage}
	channel <- &model.ProgressEvent{Type: model.ProgressEventAlbumCreated, Album: "albums/007"}
	close(channel)

	<-tracker.done

	a.Equal(model.MediaCounter{Count: 2, Size: 30}, tracker.Skipped())
	c := &model.TypeCounter{}
	c.IncrementFoundCounter(model.MediaTypeImage, 1, 12)
	a.Equal(map[string]*model.TypeCounter{
		"albums/007": c,
	}, tracker.CountPerAlbum())
	a.Equal([]string{"albums/007"}, tracker.NewAlbums())
}
