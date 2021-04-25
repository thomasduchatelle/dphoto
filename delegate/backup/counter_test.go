package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter_GetFound(t *testing.T) {
	a := assert.New(t)

	// it should init a tracker starting at 0
	counter := tracker{}
	a.Equal(uint32(0), counter.GetFoundCount())
	a.Equal(uint32(0), counter.GetFound(model.MediaTypeImage))
	a.Equal(uint32(0), counter.GetFound(model.MediaTypeVideo))

	// it should increment both total and media type sub-tracker
	counter.incrementFoundCounter(model.MediaTypeImage)
	a.Equal(uint32(1), counter.GetFoundCount())
	a.Equal(uint32(1), counter.GetFound(model.MediaTypeImage))
	a.Equal(uint32(0), counter.GetFound(model.MediaTypeVideo))

	// it should keep count of each media type
	counter.incrementFoundCounter(model.MediaTypeImage)
	counter.incrementFoundCounter(model.MediaTypeVideo)
	counter.incrementFoundCounter("Audio")
	a.Equal(uint32(4), counter.GetFoundCount())
	a.Equal(uint32(2), counter.GetFound(model.MediaTypeImage))
	a.Equal(uint32(1), counter.GetFound(model.MediaTypeVideo))
	a.Equal(uint32(0), counter.GetFound("Audio"))
}
