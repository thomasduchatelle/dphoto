package backup

const numberOfMediaType = 3 // exclude "other", include "total" in position 0
var (
	MediaCounterZero = MediaCounter{}
)

type CompletionReport interface {
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

type AlbumReport struct {
	New    bool // New is true if the album have been created during the backup
	counts [numberOfMediaType]int
	sizes  [numberOfMediaType]int
}

// NewTypeCounter is a convenience method for testing or mocking 'backup' domain
func NewTypeCounter(mediaType MediaType, count int, size int, isNew bool) *AlbumReport {
	counter := new(AlbumReport)
	counter.IncrementFoundCounter(mediaType, count, size)
	counter.New = isNew
	return counter
}

func (c *AlbumReport) IncrementFoundCounter(mediaType MediaType, count int, size int) *AlbumReport {
	c.IncrementCounter(&c.counts, mediaType, count)
	c.IncrementCounter(&c.sizes, mediaType, size)
	return c
}

func (c *AlbumReport) IncrementCounter(counter *[numberOfMediaType]int, mediaType MediaType, delta int) {
	index := c.GetMediaIndex(mediaType)
	if index > 0 {
		counter[index] = counter[index] + delta
	}

	counter[0] = counter[0] + delta
}

func (c *AlbumReport) GetMediaIndex(mediaType MediaType) int {
	switch mediaType {
	case MediaTypeImage:
		return 1
	case MediaTypeVideo:
		return 2
	}

	return -1
}

func (c *AlbumReport) Total() MediaCounter {
	return MediaCounter{
		Count: c.counts[0],
		Size:  c.sizes[0],
	}
}

func (c *AlbumReport) OfType(mediaType MediaType) MediaCounter {
	index := c.GetMediaIndex(mediaType)
	if index < 0 {
		return MediaCounter{}
	}

	return MediaCounter{
		Count: c.counts[index],
		Size:  c.sizes[index],
	}
}
