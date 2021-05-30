package backup

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTracker(t *testing.T) {
	a := assert.New(t)

	channel := make(chan *scanner.ProgressEvent, 256)
	tracker := NewTracker(channel, nil)

	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventScanComplete, Count: 3, Size: 42}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventSkipped, Count: 1, Size: 20}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventDownloaded, Count: 1, Size: 10}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventAnalysed, Count: 1, Size: 10}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventSkipped, Count: 1, Size: 10}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventDownloaded, Count: 1, Size: 12}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventAnalysed, Count: 1, Size: 12}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventUploaded, Count: 1, Size: 12, Album: "albums/007", MediaType: scanner.MediaTypeImage}
	channel <- &scanner.ProgressEvent{Type: scanner.ProgressEventAlbumCreated, Album: "albums/007"}
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
//	a.Equal(uint32(0), counter.GetFound(scanner.MediaTypeImage))
//	a.Equal(uint32(0), counter.GetFound(scanner.MediaTypeVideo))
//
//	// it should increment both total and media type sub-tracker
//	counter.incrementFoundCounter(scanner.MediaTypeImage)
//	a.Equal(uint32(1), counter.GetFoundCount())
//	a.Equal(uint32(1), counter.GetFound(scanner.MediaTypeImage))
//	a.Equal(uint32(0), counter.GetFound(scanner.MediaTypeVideo))
//
//	// it should keep count of each media type
//	counter.incrementFoundCounter(scanner.MediaTypeImage)
//	counter.incrementFoundCounter(scanner.MediaTypeVideo)
//	counter.incrementFoundCounter("Audio")
//	a.Equal(uint32(4), counter.GetFoundCount())
//	a.Equal(uint32(2), counter.GetFound(scanner.MediaTypeImage))
//	a.Equal(uint32(1), counter.GetFound(scanner.MediaTypeVideo))
//	a.Equal(uint32(0), counter.GetFound("Audio"))
//}
