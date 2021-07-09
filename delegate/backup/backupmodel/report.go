package backupmodel

const numberOfMediaType = 3 // exclude "other", include "total" in position 0

type BackupReport interface {
	NewAlbums() []string
	Skipped() MediaCounter
	CountPerAlbum() map[string]*TypeCounter
	WaitToComplete()
}

type MediaCounter struct {
	Count uint // Count is the number of medias
	Size  uint // Size is the sum of the size of the medias
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

type TypeCounter struct {
	counts [numberOfMediaType]uint
	sizes  [numberOfMediaType]uint
}

func (c *TypeCounter) IncrementFoundCounter(mediaType MediaType, count uint, size uint) {
	c.IncrementCounter(&c.counts, mediaType, count)
	c.IncrementCounter(&c.sizes, mediaType, size)
}

func (c *TypeCounter) IncrementCounter(counter *[numberOfMediaType]uint, mediaType MediaType, delta uint) {
	index := c.GetMediaIndex(mediaType)
	if index > 0 {
		counter[index] = counter[index] + delta
	}

	counter[0] = counter[0] + delta
}

func (c *TypeCounter) GetMediaIndex(mediaType MediaType) int {
	switch mediaType {
	case MediaTypeImage:
		return 1
	case MediaTypeVideo:
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

func (c *TypeCounter) OfType(mediaType MediaType) MediaCounter {
	index := c.GetMediaIndex(mediaType)
	if index < 0 {
		return MediaCounter{}
	}

	return MediaCounter{
		Count: c.counts[index],
		Size:  c.sizes[index],
	}
}
