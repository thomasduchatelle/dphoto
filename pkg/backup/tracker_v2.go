package backup

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	trackScanComplete           trackEvent = "scan-complete"             // trackScanComplete set the total of files
	trackAnalysedFromCache      trackEvent = "analysed-from-cache"       // trackAnalysedFromCache is sent instead of ProgressEventAnalysed when the analysis has been cached
	trackAnalysisFailed         trackEvent = "analysis-failed"           // trackAnalysisFailed count files skipped because the analysis failed (no date, invalid format, ...)
	trackWrongAlbum             trackEvent = "wrong-album"               // trackWrongAlbum count files in filtered out albums (if filter used), subtracted from trackScanComplete
	trackAlreadyExistsInCatalog trackEvent = "already-exists-in-catalog" // trackAlreadyExistsInCatalog count files already known in catalog, subtracted from trackScanComplete
	trackDuplicatedInVolume     trackEvent = "duplicated-in-volume"      // trackDuplicatedInVolume count files present twice in this backup/scan process, subtracted from trackScanComplete
	trackCatalogued             trackEvent = "catalogued"                // trackCatalogued files remaining after analysis, cataloguing, and filters: trackCatalogued = trackScanComplete - trackDuplicatedInVolume - trackAlreadyExistsInCatalog - trackWrongAlbum
	trackUploaded               trackEvent = "uploaded"                  // trackUploaded files uploaded, is equals to trackCatalogued when complete
	trackAlbumCreated           trackEvent = "album-created"             // trackAlbumCreated notify when a new album is created
)

type trackEvent string

type progressEvent struct {
	Type      trackEvent // Type defines what's count, and size are about ; some might not be used.
	Count     int        // Count is the number of media
	Size      int        // Size is the sum of the size of the media concerned by this event
	Album     string     // Album is the folder name of the medias concerned by this event
	MediaType MediaType  // MediaType is the type of media ; only mandatory with 'uploaded' event
}

type TrackAnalysed interface {
	OnAnalysed(done, total MediaCounter, others ExtraCounts)
}

// TrackUploaded includes both uploaded and skipped
type TrackUploaded interface {
	OnUploaded(done, total MediaCounter)
}

type TrackEvents interface {
	OnEvent(event progressEvent)
}

type TrackScanComplete interface {
	OnScanComplete(total MediaCounter)
}

type ExtraCounts struct {
	Cached   MediaCounter
	Rejected MediaCounter
}

func (c ExtraCounts) String() interface{} {
	var extraDetails []string

	if c.Cached.Count > 0 {
		extraDetails = append(extraDetails, fmt.Sprintf("from cache: %d", c.Cached.Count))
	}
	if c.Rejected.Count > 0 {
		extraDetails = append(extraDetails, fmt.Sprintf("rejected: %d", c.Rejected.Count))
	}

	cachedExplanation := ""
	if len(extraDetails) > 0 {
		cachedExplanation = "[" + strings.Join(extraDetails, " ; ") + "]"
	}
	return cachedExplanation
}

// newTrackerV2 creates the trackerV2 and start consuming (async)
func newTrackerV2(options Options) (*trackerObserver, *trackerV2) {
	tracker := &trackerV2{
		listeners:  panicIfDoesNotImplementsTrackerInterface(options.Listener),
		eventCount: make(map[trackEvent]MediaCounter),
	}

	observer := &trackerObserver{
		channel: make(chan *progressEvent, defaultValue(options.BatchSize, 1)*8),
		done:    make(chan struct{}),
	}
	go func() {
		defer close(observer.done)
		tracker.consume(observer.channel)
	}()
	return observer, tracker
}

func panicIfDoesNotImplementsTrackerInterface(listener interface{}) (listeners []interface{}) {
	if listener == nil {
		return
	}

	listeners = append(listeners, listener)

	if _, ok := listener.(TrackEvents); ok {
		return
	}

	if _, ok := listener.(TrackScanComplete); ok {
		return
	}

	if _, ok := listener.(TrackAnalysed); ok {
		return
	}

	if _, ok := listener.(TrackUploaded); ok {
		return
	}

	panic("listener must implement at least one of the tracker interfaces")
}

// trackerV2 is simplifying the consumption of events from scans and backups to implement progress bars.
type trackerV2 struct {
	listeners  []interface{} // listeners will receive aggregated and typed updates
	eventCount map[trackEvent]MediaCounter
}

func (t *trackerV2) consume(progressChannel chan *progressEvent) {
	for event := range progressChannel {
		t.fireRawEvent(event)

		current, _ := t.eventCount[event.Type]
		t.eventCount[event.Type] = current.Add(event.Count, event.Size)

		switch event.Type {
		case trackScanComplete:
			t.fireScanComplete()

		case trackAnalysisFailed,
			trackAlreadyExistsInCatalog,
			trackDuplicatedInVolume,
			trackWrongAlbum,
			trackCatalogued:
			t.fireAnalysedEvent()

		case trackUploaded:
			t.fireUploadedEvent()

		case trackAlbumCreated:
		case trackAnalysedFromCache:
			// nothing

		default:
			log.Warnf("Progress type '%s' is not supported", event.Type)
		}
	}
}

func (t *trackerV2) fireScanComplete() {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackScanComplete); ok {
			dispatch.OnScanComplete(t.eventCount[trackScanComplete])
		}
	}
}

func (t *trackerV2) fireAnalysedEvent() {
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

func (t *trackerV2) fireUploadedEvent() {
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

func (t *trackerV2) fireRawEvent(event *progressEvent) {
	for _, listener := range t.listeners {
		if dispatch, ok := listener.(TrackEvents); ok {
			dispatch.OnEvent(*event)
		}
	}

}

type trackerObserver struct {
	channel chan *progressEvent
	done    chan struct{}
}

func (p *trackerObserver) NoMoreEvents() {
	close(p.channel)
	<-p.done
}

func (p *trackerObserver) OnScanComplete(ctx context.Context, count, size int) error {
	p.channel <- &progressEvent{Type: trackScanComplete, Count: count, Size: size}
	return nil
}

func (p *trackerObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	p.channel <- &progressEvent{Type: trackAnalysisFailed, Count: 1, Size: found.Size()}
	return nil
}

func (p *trackerObserver) OnSkipDelegateAnalyser(ctx context.Context, found FoundMedia) error {
	p.channel <- &progressEvent{Type: trackAnalysedFromCache, Count: 1, Size: found.Size()}

	return nil
}

func (p *trackerObserver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	count := MediaCounterZero
	for _, request := range requests {
		if request.CatalogReference.AlbumCreated() {
			p.channel <- &progressEvent{Type: trackAlbumCreated, Count: 1, Album: request.CatalogReference.AlbumFolderName()}
		}

		count = count.Add(1, request.AnalysedMedia.FoundMedia.Size())
	}

	p.channel <- &progressEvent{Type: trackCatalogued, Count: count.Count, Size: count.Size}

	return nil
}

func (p *trackerObserver) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	switch {
	case errors.Is(cause, ErrCatalogerFilterMustBeInAlbum):
		p.channel <- &progressEvent{Type: trackWrongAlbum, Count: 1, Size: media.FoundMedia.Size()}
		return nil
	case errors.Is(cause, ErrCatalogerFilterMustNotAlreadyExists):
		p.channel <- &progressEvent{Type: trackAlreadyExistsInCatalog, Count: 1, Size: media.FoundMedia.Size()}
		return nil
	case errors.Is(cause, ErrMediaMustNotBeDuplicated):
		p.channel <- &progressEvent{Type: trackDuplicatedInVolume, Count: 1, Size: media.FoundMedia.Size()}
		return nil

	default:
		return errors.Wrapf(cause, "filter error is not supported. Media: %s", media.FoundMedia)
	}
}

func (p *trackerObserver) OnBackingUpMediaRequestUploaded(ctx context.Context, request BackingUpMediaRequest) error {
	p.channel <- &progressEvent{
		Type:      trackUploaded,
		Count:     1,
		Size:      request.AnalysedMedia.FoundMedia.Size(),
		Album:     request.CatalogReference.AlbumFolderName(),
		MediaType: request.AnalysedMedia.Type,
	}

	return nil
}
