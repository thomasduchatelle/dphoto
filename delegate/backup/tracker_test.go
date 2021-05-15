package backup

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

	<-tracker.Done

	a.Equal(MediaCounter{Count: 2, Size: 30}, tracker.Skipped())
	a.Equal(map[string]*TypeCounter{
		"albums/007": {
			counts: [3]uint{1, 1, 0},
			sizes:  [3]uint{12, 12, 0},
		},
	}, tracker.CountPerAlbum())
	a.Equal([]string{"albums/007"}, tracker.NewAlbums())
}

//func TestCounter_GetFound(t *testing.T) {
//	a := assert.New(t)
//
//	// it should init a tracker starting at 0
//	counter := {}
//	a.Equal(uint32(0), counter.GetFoundCount())
//	a.Equal(uint32(0), counter.GetFound(model.MediaTypeImage))
//	a.Equal(uint32(0), counter.GetFound(model.MediaTypeVideo))
//
//	// it should increment both total and media type sub-tracker
//	counter.incrementFoundCounter(model.MediaTypeImage)
//	a.Equal(uint32(1), counter.GetFoundCount())
//	a.Equal(uint32(1), counter.GetFound(model.MediaTypeImage))
//	a.Equal(uint32(0), counter.GetFound(model.MediaTypeVideo))
//
//	// it should keep count of each media type
//	counter.incrementFoundCounter(model.MediaTypeImage)
//	counter.incrementFoundCounter(model.MediaTypeVideo)
//	counter.incrementFoundCounter("Audio")
//	a.Equal(uint32(4), counter.GetFoundCount())
//	a.Equal(uint32(2), counter.GetFound(model.MediaTypeImage))
//	a.Equal(uint32(1), counter.GetFound(model.MediaTypeVideo))
//	a.Equal(uint32(0), counter.GetFound("Audio"))
//}
