package tracker

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTracker(t *testing.T) {
	a := assert.New(t)

	channel := make(chan *backupmodel.ProgressEvent, 256)
	tracker := NewTracker(channel, nil)

	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventScanComplete, Count: 3, Size: 42}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventSkipped, Count: 1, Size: 20}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventDownloaded, Count: 1, Size: 10}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventAnalysed, Count: 1, Size: 10}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventSkipped, Count: 1, Size: 10}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventDownloaded, Count: 1, Size: 12}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventAnalysed, Count: 1, Size: 12}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventUploaded, Count: 1, Size: 12, Album: "albums/007", MediaType: backupmodel.MediaTypeImage}
	channel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventAlbumCreated, Album: "albums/007"}
	close(channel)

	<-tracker.done

	a.Equal(backupmodel.MediaCounter{Count: 2, Size: 30}, tracker.Skipped())
	c := &backupmodel.TypeCounter{}
	c.IncrementFoundCounter(backupmodel.MediaTypeImage, 1, 12)
	a.Equal(map[string]*backupmodel.TypeCounter{
		"albums/007": c,
	}, tracker.CountPerAlbum())
	a.Equal([]string{"albums/007"}, tracker.NewAlbums())
}
