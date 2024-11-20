package backup

import (
	"context"
	"sync"
)

var (
	MediaCounterZero = MediaCounter{}
)

type Report interface {
	Skipped() MediaCounter
	CountPerAlbum() map[string]*AlbumReport
}

type MediaCounter struct {
	Count int // Count is the number of medias
	Size  int // Size is the sum of the size of the medias
}

func NewMediaCounter(count int, size int) MediaCounter {
	return MediaCounter{
		Count: count,
		Size:  size,
	}
}

// Add creates a new MediaCounter with the delta applied ; initial MediaCounter is not updated.
func (c MediaCounter) Add(count int, size int) MediaCounter {
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

// NewAlbumReport is a convenience method for testing or mocking 'backup' domain
func NewAlbumReport(mediaType MediaType, count int, size int, isNew bool) *AlbumReport {
	counter := NewMediaCounter(count, size)

	report := &AlbumReport{isNew: isNew}
	switch mediaType {
	case MediaTypeImage:
		report.image = counter
	case MediaTypeVideo:
		report.video = counter
	default:
		report.other = counter
	}
	return report
}

func newBackupReportBuilder() *backupReportBuilder {
	return &backupReportBuilder{
		lock:          sync.Mutex{},
		countPerAlbum: make(map[string]*AlbumReport),
	}
}

type backupReportBuilder struct {
	lock          sync.Mutex
	skipped       MediaCounter
	countPerAlbum map[string]*AlbumReport
}

func (r *backupReportBuilder) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.skipped = r.skipped.Add(1, found.Size())

	return nil
}

func (r *backupReportBuilder) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.skipped = r.skipped.Add(1, media.FoundMedia.Size())

	return nil
}

func (r *backupReportBuilder) OnBackingUpMediaRequestUploaded(ctx context.Context, request BackingUpMediaRequest) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	albumName := request.CatalogReference.AlbumFolderName()
	count, ok := r.countPerAlbum[albumName]
	if !ok {
		count = &AlbumReport{}
		r.countPerAlbum[albumName] = count
	}

	if request.CatalogReference.AlbumCreated() {
		count.isNew = true
	}

	switch request.AnalysedMedia.Type {
	case MediaTypeImage:
		count.image = count.image.Add(1, request.AnalysedMedia.FoundMedia.Size())
	case MediaTypeVideo:
		count.video = count.video.Add(1, request.AnalysedMedia.FoundMedia.Size())
	default:
		count.other = count.other.Add(1, request.AnalysedMedia.FoundMedia.Size())
	}

	return nil
}

func (r *backupReportBuilder) Skipped() MediaCounter {
	return r.skipped
}

func (r *backupReportBuilder) CountPerAlbum() map[string]*AlbumReport {
	return r.countPerAlbum
}

type AlbumReport struct {
	isNew bool
	image MediaCounter
	video MediaCounter
	other MediaCounter
}

func (c *AlbumReport) IsNew() bool {
	return c.isNew
}

func (c *AlbumReport) Total() MediaCounter {
	return c.image.AddCounter(c.video).AddCounter(c.other)
}

func (c *AlbumReport) OfType(mediaType MediaType) MediaCounter {
	switch mediaType {
	case MediaTypeImage:
		return c.image
	case MediaTypeVideo:
		return c.video
	default:
		return c.other
	}
}
