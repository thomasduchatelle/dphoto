package tracker

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	log "github.com/sirupsen/logrus"
)

const numberOfMediaType = 3 // exclude "other", include "total" in position 0

type TrackScanComplete interface {
	OnScanComplete(total MediaCounter)
}

type TrackDownloaded interface {
	OnDownloaded(done, total MediaCounter)
}

type TrackAnalysed interface {
	OnAnalysed(done, total MediaCounter)
}

// TrackUploaded includes both uploaded and skipped
type TrackUploaded interface {
	OnUploaded(done, total MediaCounter)
}

type MediaCounter struct {
	Count uint // Count is the number of medias
	Size  uint // Size is the sum of the size of the medias
}

type TypeCounter struct {
	counts [numberOfMediaType]uint
	sizes  [numberOfMediaType]uint
}

// Tracker is consuming progress channel, keep a record of counts, and call listeners
type Tracker struct {
	listeners           []interface{} // listeners will receive aggregated and typed updates
	total               MediaCounter
	scanComplete        bool
	skipped             MediaCounter
	downloaded          MediaCounter
	analysed            MediaCounter
	skippedBeforeUpload MediaCounter
	uploaded            MediaCounter
	createdAlbums       []string
	detailedCount       map[string]*TypeCounter
	Done                chan struct{} // Done is closed when all events have been processed.
}

// NewTracker creates the Tracker and start consuming (async)
func NewTracker(progressChannel chan *model.ProgressEvent, listeners []interface{}) *Tracker {
	tracker := &Tracker{
		listeners:     listeners,
		Done:          make(chan struct{}),
		detailedCount: make(map[string]*TypeCounter),
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
	return t.skipped.AddCounter(t.skippedBeforeUpload)
}

func (t *Tracker) CountPerAlbum() map[string]*TypeCounter {
	return t.detailedCount
}

func (t *Tracker) consume(progressChannel chan *model.ProgressEvent) {
	for event := range progressChannel {
		switch event.Type {
		case model.ProgressEventScanComplete:
			t.scanComplete = true
			t.total = MediaCounter{
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
				typeCount = &TypeCounter{}
				t.detailedCount[event.Album] = typeCount
			}
			typeCount.incrementFoundCounter(event.MediaType, event.Count, event.Size)

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

// Add creates a new MediaCounter with the delta applied ; initial MediaCounter is not updated.
func (c MediaCounter) Add(count uint, size uint) MediaCounter {
	return MediaCounter{
		Count: c.Count + count,
		Size:  c.Size + size,
	}
}

// AddCounter creates a new MediaCounter which is the sum of the 2 counters provided.
func (c MediaCounter) AddCounter(counter MediaCounter) MediaCounter {
	return c.Add(counter.Count, counter.Size)
}

// IsZero returns true if it's the default value
func (c MediaCounter) IsZero() bool {
	return c.Size == 0 && c.Count == 0
}

func (c *TypeCounter) incrementFoundCounter(mediaType model.MediaType, count uint, size uint) {
	c.incrementCounter(&c.counts, mediaType, count)
	c.incrementCounter(&c.sizes, mediaType, size)
}

func (c *TypeCounter) incrementCounter(counter *[numberOfMediaType]uint, mediaType model.MediaType, delta uint) {
	index := c.getMediaIndex(mediaType)
	if index > 0 {
		counter[index] = counter[index] + delta
	}

	counter[0] = counter[0] + delta
}

func (c *TypeCounter) getMediaIndex(mediaType model.MediaType) int {
	switch mediaType {
	case model.MediaTypeImage:
		return 1
	case model.MediaTypeVideo:
		return 2
	}

	return -1
}

func (c *TypeCounter) Total() MediaCounter {
	return MediaCounter{
		Count: c.counts[0],
		Size:  c.sizes[0],
	}
}

func (c *TypeCounter) OfType(mediaType model.MediaType) MediaCounter {
	index := c.getMediaIndex(mediaType)
	if index < 0 {
		return MediaCounter{}
	}

	return MediaCounter{
		Count: c.counts[index],
		Size:  c.sizes[index],
	}
}
