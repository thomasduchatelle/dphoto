package backup

import (
	log "github.com/sirupsen/logrus"
)

const (
	ProgressEventAnalysed   trackEvent = "analysed"      // ProgressEventAnalysed is not useful for progress, it will be fined grained before upload
	ProgressEventCatalogued trackEvent = "catalogued-v1" // ProgressEventCatalogued is not useful for progress, it will be fined grained before upload
)

// Tracker is simplifying the consumption of events from scans and backups to implement progress bars.
type Tracker struct {
	listeners     []interface{} // listeners will receive aggregated and typed updates
	scanComplete  bool
	eventCount    map[trackEvent]MediaCounter
	createdAlbums []string
	detailedCount map[string]*AlbumReport
	Done          chan struct{} // Done is closed when all events have been processed.
}

// NewTracker creates the Tracker and start consuming (async)
func NewTracker(progressChannel chan *progressEvent, listeners ...interface{}) *Tracker {
	tracker := &Tracker{
		listeners:     listeners,
		eventCount:    make(map[trackEvent]MediaCounter),
		detailedCount: make(map[string]*AlbumReport),
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
	exists, _ := t.eventCount[trackAlreadyExistsInCatalog]
	duplicates, _ := t.eventCount[trackDuplicatedInVolume]
	wrongAlbum, _ := t.eventCount[trackWrongAlbum]
	return exists.AddCounter(duplicates).AddCounter(wrongAlbum)
}

func (t *Tracker) CountPerAlbum() map[string]*AlbumReport {
	newAlbums := make(map[string]interface{})
	for _, album := range t.createdAlbums {
		newAlbums[album] = nil
	}

	for folderName, counts := range t.detailedCount {
		if _, isNew := newAlbums[folderName]; isNew {
			counts.New = true
		}
	}

	return t.detailedCount
}

func (t *Tracker) WaitToComplete() {
	<-t.Done
}

func (t *Tracker) consume(progressChannel chan *progressEvent) {
	for event := range progressChannel {
		t.fireRawEvent(event)

		current, _ := t.eventCount[event.Type]
		t.eventCount[event.Type] = current.Add(event.Count, event.Size)

		switch event.Type {
		case trackScanComplete:
			t.scanComplete = true
			t.fireScanComplete()

		case ProgressEventAnalysed:
		case trackAnalysedFromCache:
		case ProgressEventCatalogued:
			// nothing

		case trackAnalysisFailed,
			trackAlreadyExistsInCatalog,
			trackDuplicatedInVolume,
			trackWrongAlbum,
			trackCatalogued:
			t.fireAnalysedEvent()

		case trackUploaded:
			typeCount, ok := t.detailedCount[event.Album]
			if !ok {
				typeCount = &AlbumReport{}
				t.detailedCount[event.Album] = typeCount
			}
			typeCount.IncrementFoundCounter(event.MediaType, event.Count, event.Size)

			t.fireUploadedEvent()

		case trackAlbumCreated:
			t.createdAlbums = append(t.createdAlbums, event.Album)

		default:
			log.Warnf("Progress type '%s' is not supported", event.Type)
		}
	}
}

func (t *Tracker) fireScanComplete() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackScanComplete); ok {
			dispatch.OnScanComplete(t.eventCount[trackScanComplete])
		}
	}
}

func (t *Tracker) fireAnalysedEvent() {
	passed, _ := t.eventCount[trackCatalogued]
	exists, _ := t.eventCount[trackAlreadyExistsInCatalog]
	duplicates, _ := t.eventCount[trackDuplicatedInVolume]
	wrongAlbum, _ := t.eventCount[trackWrongAlbum]
	rejected, _ := t.eventCount[trackAnalysisFailed]
	analysedFromCache, _ := t.eventCount[trackAnalysedFromCache]

	done := passed.AddCounter(exists).AddCounter(duplicates).AddCounter(wrongAlbum).AddCounter(rejected)

	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackAnalysed); ok {
			dispatch.OnAnalysed(done, t.eventCount[trackScanComplete], ExtraCounts{
				Cached:   analysedFromCache,
				Rejected: rejected,
			})
		}
	}
}

func (t *Tracker) fireUploadedEvent() {
	scanned, _ := t.eventCount[trackScanComplete]
	rejected, _ := t.eventCount[trackAnalysisFailed]
	exists, _ := t.eventCount[trackAlreadyExistsInCatalog]
	duplicates, _ := t.eventCount[trackDuplicatedInVolume]
	wrongAlbum, _ := t.eventCount[trackWrongAlbum]
	ready, _ := t.eventCount[trackCatalogued]
	uploaded, _ := t.eventCount[trackUploaded]

	total := MediaCounterZero
	if ready.AddCounter(duplicates).AddCounter(exists).AddCounter(wrongAlbum).AddCounter(rejected).Count == scanned.Count {
		// total-to-upload is confirmed
		total = ready
	}

	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackUploaded); ok {
			dispatch.OnUploaded(uploaded, total)
		}
	}
}

func (t *Tracker) fireRawEvent(event *progressEvent) {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackEvents); ok {
			dispatch.OnEvent(*event)
		}
	}

}

func (t *Tracker) MediaCount() int {
	count := 0
	for _, counter := range t.detailedCount {
		count += counter.Total().Count
	}

	return count
}
