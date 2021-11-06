package tracker

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	log "github.com/sirupsen/logrus"
)

type TrackScanComplete interface {
	OnScanComplete(total backupmodel.MediaCounter)
}

type TrackDownloaded interface {
	OnDownloaded(done, total backupmodel.MediaCounter)
}

type TrackAnalysed interface {
	OnAnalysed(done, total backupmodel.MediaCounter)
}

// TrackUploaded includes both uploaded and skipped
type TrackUploaded interface {
	OnUploaded(done, total backupmodel.MediaCounter)
}

// Tracker is consuming progress channel, keep a record of counts, and call listeners
type Tracker struct {
	listeners           []interface{} // listeners will receive aggregated and typed updates
	total               backupmodel.MediaCounter
	scanComplete        bool
	skipped             backupmodel.MediaCounter
	downloaded          backupmodel.MediaCounter
	analysed            backupmodel.MediaCounter
	skippedBeforeUpload backupmodel.MediaCounter
	uploaded            backupmodel.MediaCounter
	createdAlbums       []string
	detailedCount       map[string]*backupmodel.TypeCounter
	done                chan struct{} // done is closed when all events have been processed.
}

// NewTracker creates the Tracker and start consuming (async)
func NewTracker(progressChannel chan *backupmodel.ProgressEvent, listeners []interface{}) *Tracker {
	tracker := &Tracker{
		listeners:     listeners,
		done:          make(chan struct{}),
		detailedCount: make(map[string]*backupmodel.TypeCounter),
	}
	go func() {
		defer close(tracker.done)
		tracker.consume(progressChannel)
	}()
	return tracker
}

func (t *Tracker) NewAlbums() []string {
	return t.createdAlbums
}

func (t *Tracker) Skipped() backupmodel.MediaCounter {
	return t.skipped.AddCounter(t.skippedBeforeUpload)
}

func (t *Tracker) CountPerAlbum() map[string]*backupmodel.TypeCounter {
	return t.detailedCount
}

func (t *Tracker) WaitToComplete() {
	<-t.done
}

func (t *Tracker) consume(progressChannel chan *backupmodel.ProgressEvent) {
	for event := range progressChannel {
		switch event.Type {
		case backupmodel.ProgressEventScanComplete:
			t.scanComplete = true
			t.total = backupmodel.MediaCounter{
				Count: event.Count,
				Size:  event.Size,
			}

			for _, listener := range t.listeners {
				if dispatch, ok := listener.(TrackScanComplete); ok {
					dispatch.OnScanComplete(t.total)
				}
			}

		case backupmodel.ProgressEventSkipped:
			t.skipped = t.skipped.Add(event.Count, event.Size)

		case backupmodel.ProgressEventDownloaded:
			t.downloaded = t.downloaded.Add(event.Count, event.Size)
			t.fireDownloadedEvent()

		case backupmodel.ProgressEventAnalysed:
			t.analysed = t.analysed.Add(event.Count, event.Size)
			t.fireAnalysedEvent()

		case backupmodel.ProgressEventSkippedAfterAnalyse:
			t.skippedBeforeUpload = t.skippedBeforeUpload.Add(event.Count, event.Size)
			t.fireUploadedEvent()

		case backupmodel.ProgressEventUploaded:
			t.uploaded = t.uploaded.Add(event.Count, event.Size)

			typeCount, ok := t.detailedCount[event.Album]
			if !ok {
				typeCount = &backupmodel.TypeCounter{}
				t.detailedCount[event.Album] = typeCount
			}
			typeCount.IncrementFoundCounter(event.MediaType, event.Count, event.Size)

			t.fireUploadedEvent()

		case backupmodel.ProgressEventAlbumCreated:
			t.createdAlbums = append(t.createdAlbums, event.Album)

		default:
			log.Warnf("Progress type '%s' is not supported", event.Type)
		}
	}
}

func (t *Tracker) fireDownloadedEvent() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackDownloaded); ok {
			dispatch.OnDownloaded(t.downloaded, t.total)
		}
	}
}

func (t *Tracker) fireAnalysedEvent() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackAnalysed); ok {
			dispatch.OnAnalysed(t.analysed, t.total)
		}
	}
}

func (t *Tracker) fireUploadedEvent() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackUploaded); ok {
			dispatch.OnUploaded(t.uploaded.AddCounter(t.skippedBeforeUpload), t.total)
		}
	}
}
