package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter_GetFound(t *testing.T) {
	a := assert.New(t)

	// it should init a counter starting at 0
	counter := Counter{}
	a.Equal(uint32(0), counter.GetFoundCount())
	a.Equal(uint32(0), counter.GetFound(IMAGE))
	a.Equal(uint32(0), counter.GetFound(VIDEO))

	// it should increment both total and media type sub-counter
	counter.incrementFoundCounter(IMAGE)
	a.Equal(uint32(1), counter.GetFoundCount())
	a.Equal(uint32(1), counter.GetFound(IMAGE))
	a.Equal(uint32(0), counter.GetFound(VIDEO))

	// it should keep count of each media type
	counter.incrementFoundCounter(IMAGE)
	counter.incrementFoundCounter(VIDEO)
	counter.incrementFoundCounter("Audio")
	a.Equal(uint32(4), counter.GetFoundCount())
	a.Equal(uint32(2), counter.GetFound(IMAGE))
	a.Equal(uint32(1), counter.GetFound(VIDEO))
	a.Equal(uint32(0), counter.GetFound("Audio"))
}
