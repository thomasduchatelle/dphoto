package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"sync/atomic"
)

// Tracker is an object provided by consumer to track a backup progress
type Tracker interface {
	// IncrementFound is used during the scan each time a file is found
	IncrementFound(by int, mediaType model.MediaType)
	// ScanComplete is called when volume has been fully scanned
	ScanComplete()

	IncrementToBackup(by int, mediaType model.MediaType)
	IncrementBackedUp(by int, mediaType model.MediaType)
}

// count is an array of 1 + number of media types: [<total>,<number of images>,<number of video>]
type count [numberOfMediaType]uint32

type tracker struct {
	total count
	found [numberOfMediaType]uint32
}

func (c *tracker) GetFoundCount() uint32 {
	return c.found[0]
}

func (c *tracker) GetFound(mediaType model.MediaType) uint32 {
	index := c.getMediaIndex(mediaType)
	if index > 0 {
		return c.found[index]
	}

	return 0
}

// atomically increments found tracker per type
func (c *tracker) incrementFoundCounter(mediaType model.MediaType) uint32 {
	return c.incrementCounter(&c.found, mediaType)
}

// atomically increments tracker per type
func (c *tracker) incrementCounter(counter *[numberOfMediaType]uint32, mediaType model.MediaType) uint32 {
	index := c.getMediaIndex(mediaType)
	if index > 0 {
		atomic.AddUint32(&(counter[index]), 1)
	}

	return atomic.AddUint32(&(counter[0]), 1)
}

// return -1 if media type is unknown
func (c *tracker) getMediaIndex(mediaType model.MediaType) int {
	switch mediaType {
	case model.MediaTypeImage:
		return 1
	case model.MediaTypeVideo:
		return 2
	}

	return -1
}
