package backup

import (
	log "github.com/sirupsen/logrus"
)

type TrackEvents interface {
	OnEvent(event ProgressEvent)
}

type TrackScanComplete interface {
	OnScanComplete(total MediaCounter)
}

type TrackAnalysed interface {
	OnAnalysed(done, total MediaCounter)
}

// TrackUploaded includes both uploaded and skipped
type TrackUploaded interface {
	OnUploaded(done, total MediaCounter)
}

// Tracker is consuming progress channel, keep a record of counts, and call listeners
type Tracker struct {
	listeners     []interface{} // listeners will receive aggregated and typed updates
	scanComplete  bool
	eventCount    map[ProgressEventType]MediaCounter
	createdAlbums []string
	detailedCount map[string]*TypeCounter
	Done          chan struct{} // Done is closed when all events have been processed.
}

// NewTracker creates the Tracker and start consuming (async)
func NewTracker(progressChannel chan *ProgressEvent, listeners ...interface{}) *Tracker {
	tracker := &Tracker{
		listeners:     listeners,
		eventCount:    make(map[ProgressEventType]MediaCounter),
		detailedCount: make(map[string]*TypeCounter),
		Done:          make(chan struct{}),
	}
	go func() {
		defer close(tracker.Done)
		tracker.consume(progressChannel)
	}()
	return tracker
}

func (t *Tracker) NewAlbums() []string {
	return t.createdAlbums
}

func (t *Tracker) Skipped() MediaCounter {
	exists, _ := t.eventCount[ProgressEventAlreadyExists]
	duplicates, _ := t.eventCount[ProgressEventDuplicate]
	wrongAlbum, _ := t.eventCount[ProgressEventWrongAlbum]
	return exists.AddCounter(duplicates).AddCounter(wrongAlbum)
}

func (t *Tracker) CountPerAlbum() map[string]*TypeCounter {
	return t.detailedCount
}

func (t *Tracker) WaitToComplete() {
	<-t.Done
}

func (t *Tracker) consume(progressChannel chan *ProgressEvent) {
	for event := range progressChannel {
		t.fireRawEvent(event)

		current, _ := t.eventCount[event.Type]
		t.eventCount[event.Type] = current.Add(event.Count, event.Size)

		switch event.Type {
		case ProgressEventScanComplete:
			t.scanComplete = true
			t.fireScanComplete()

		case ProgressEventAnalysed:
		case ProgressEventCatalogued:
			// nothing

		case ProgressEventAlreadyExists,
			ProgressEventDuplicate,
			ProgressEventWrongAlbum,
			ProgressEventReadyForUpload:
			t.fireAnalysedEvent()

		case ProgressEventUploaded:
			typeCount, ok := t.detailedCount[event.Album]
			if !ok {
				typeCount = &TypeCounter{}
				t.detailedCount[event.Album] = typeCount
			}
			typeCount.IncrementFoundCounter(event.MediaType, event.Count, event.Size)

			t.fireUploadedEvent()

		case ProgressEventAlbumCreated:
			t.createdAlbums = append(t.createdAlbums, event.Album)

		default:
			log.Warnf("Progress type '%s' is not supported", event.Type)
		}
	}
}

func (t *Tracker) fireScanComplete() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackScanComplete); ok {
			dispatch.OnScanComplete(t.eventCount[ProgressEventScanComplete])
		}
	}
}

func (t *Tracker) fireAnalysedEvent() {
	passed, _ := t.eventCount[ProgressEventReadyForUpload]
	exists, _ := t.eventCount[ProgressEventAlreadyExists]
	duplicates, _ := t.eventCount[ProgressEventDuplicate]
	wrongAlbum, _ := t.eventCount[ProgressEventWrongAlbum]

	done := passed.AddCounter(exists).AddCounter(duplicates).AddCounter(wrongAlbum)

	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackAnalysed); ok {
			dispatch.OnAnalysed(done, t.eventCount[ProgressEventScanComplete])
		}
	}
}

func (t *Tracker) fireUploadedEvent() {
	scanned, _ := t.eventCount[ProgressEventScanComplete]
	exists, _ := t.eventCount[ProgressEventAlreadyExists]
	duplicates, _ := t.eventCount[ProgressEventDuplicate]
	wrongAlbum, _ := t.eventCount[ProgressEventWrongAlbum]
	ready, _ := t.eventCount[ProgressEventReadyForUpload]
	uploaded, _ := t.eventCount[ProgressEventUploaded]

	total := MediaCounterZero
	if ready.AddCounter(duplicates).AddCounter(exists).AddCounter(wrongAlbum).Count == scanned.Count {
		total = ready
	}

	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackUploaded); ok {
			dispatch.OnUploaded(uploaded, total)
		}
	}
}

func (t *Tracker) fireRawEvent(event *ProgressEvent) {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackEvents); ok {
			dispatch.OnEvent(*event)
		}
	}

}
