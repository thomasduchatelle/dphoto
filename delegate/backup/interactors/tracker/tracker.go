package tracker

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	log "github.com/sirupsen/logrus"
)

type TrackScanComplete interface {
	OnScanComplete(total model.MediaCounter)
}

type TrackDownloaded interface {
	OnDownloaded(done, total model.MediaCounter)
}

type TrackAnalysed interface {
	OnAnalysed(done, total model.MediaCounter)
}

// TrackUploaded includes both uploaded and skipped
type TrackUploaded interface {
	OnUploaded(done, total model.MediaCounter)
}

// Tracker is consuming progress channel, keep a record of counts, and call listeners
type Tracker struct {
	listeners           []interface{} // listeners will receive aggregated and typed updates
	total               model.MediaCounter
	scanComplete        bool
	skipped             model.MediaCounter
	downloaded          model.MediaCounter
	analysed            model.MediaCounter
	skippedBeforeUpload model.MediaCounter
	uploaded            model.MediaCounter
	createdAlbums       []string
	detailedCount       map[string]*model.TypeCounter
	done                chan struct{} // done is closed when all events have been processed.
}

// NewTracker creates the Tracker and start consuming (async)
func NewTracker(progressChannel chan *model.ProgressEvent, listeners []interface{}) *Tracker {
	tracker := &Tracker{
		listeners:     listeners,
		done:          make(chan struct{}),
		detailedCount: make(map[string]*model.TypeCounter),
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

func (t *Tracker) Skipped() model.MediaCounter {
	return t.skipped.AddCounter(t.skippedBeforeUpload)
}

func (t *Tracker) CountPerAlbum() map[string]*model.TypeCounter {
	return t.detailedCount
}

func (t *Tracker) WaitToComplete() {
	<-t.done
}

func (t *Tracker) consume(progressChannel chan *model.ProgressEvent) {
	for event := range progressChannel {
		switch event.Type {
		case model.ProgressEventScanComplete:
			t.scanComplete = true
			t.total = model.MediaCounter{
				Count: event.Count,
				Size:  event.Size,
			}

			for _, listener := range t.listeners {
				if dispatch, ok := listener.(TrackScanComplete); ok {
					dispatch.OnScanComplete(t.total)
				}
			}

		case model.ProgressEventSkipped:
			t.skipped = t.skipped.Add(event.Count, event.Size)

			t.fireDownloadedEvent()
			t.fireAnalysedEvent()
			t.fireUploadedEvent()

		case model.ProgressEventDownloaded:
			t.downloaded = t.downloaded.Add(event.Count, event.Size)
			t.fireDownloadedEvent()

		case model.ProgressEventAnalysed:
			t.analysed = t.analysed.Add(event.Count, event.Size)
			t.fireAnalysedEvent()

		case model.ProgressEventSkippedAfterAnalyse:
			t.skippedBeforeUpload.Add(event.Count, event.Size)
			t.fireUploadedEvent()

		case model.ProgressEventUploaded:
			t.uploaded = t.uploaded.Add(event.Count, event.Size)

			typeCount, ok := t.detailedCount[event.Album]
			if !ok {
				typeCount = &model.TypeCounter{}
				t.detailedCount[event.Album] = typeCount
			}
			typeCount.IncrementFoundCounter(event.MediaType, event.Count, event.Size)

			t.fireUploadedEvent()

		case model.ProgressEventAlbumCreated:
			t.createdAlbums = append(t.createdAlbums, event.Album)

		default:
			log.Warnf("Progress type '%s' is not supported", event.Type)
		}
	}
}

func (t *Tracker) fireDownloadedEvent() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackDownloaded); ok {
			dispatch.OnDownloaded(t.downloaded.AddCounter(t.skipped), t.total)
		}
	}
}

func (t *Tracker) fireAnalysedEvent() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackAnalysed); ok {
			dispatch.OnAnalysed(t.analysed.AddCounter(t.skipped), t.total)
		}
	}
}

func (t *Tracker) fireUploadedEvent() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackUploaded); ok {
			dispatch.OnUploaded(t.uploaded.AddCounter(t.skipped).AddCounter(t.skippedBeforeUpload), t.total)
		}
	}
}
