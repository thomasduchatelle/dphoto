package backup

import "sync/atomic"

type Counter struct {
	found [numberOfMediaType]uint32
}

func (c *Counter) GetFoundCount() uint32 {
	return c.found[0]
}

func (c *Counter) GetFound(mediaType MediaType) uint32 {
	index := c.getMediaIndex(mediaType)
	if index > 0 {
		return c.found[index]
	}

	return 0
}

// atomically increments found counter per type
func (c *Counter) incrementFoundCounter(mediaType MediaType) uint32 {
	return c.incrementCounter(&c.found, mediaType)
}

// atomically increments counter per type
func (c *Counter) incrementCounter(counter *[numberOfMediaType]uint32, mediaType MediaType) uint32 {
	index := c.getMediaIndex(mediaType)
	if index > 0 {
		atomic.AddUint32(&(counter[index]), 1)
	}

	return atomic.AddUint32(&(counter[0]), 1)
}

// return -1 if media type is unknown
func (c *Counter) getMediaIndex(mediaType MediaType) int {
	switch mediaType {
	case IMAGE:
		return 1
	case VIDEO:
		return 2
	}

	return -1
}
