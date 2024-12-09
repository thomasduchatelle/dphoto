package backup

const numberOfMediaType = 3 // exclude "other", include "total" in position 0
var (
	MediaCounterZero = MediaCounter{}
)

type CompletionReport interface {
	Skipped() MediaCounter
	CountPerAlbum() map[string]IAlbumReport
}

type IAlbumReport interface {
	IsNew() bool
	Total() MediaCounter
	OfType(mediaType MediaType) MediaCounter
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
func NewAlbumReport(mediaType MediaType, count int, size int, isNew bool) IAlbumReport {
	counter := NewMediaCounter(count, size)

	report := &albumReport{isNew: isNew}
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
