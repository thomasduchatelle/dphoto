package scanner

import (
	"io"
	"os"
	"time"
)

type InmemoryMedia struct {
	filename  string
	size      uint
	mediaDate time.Time
}

type InmemoryMediaWithHash struct {
	*InmemoryMedia
	hash string
}

// NewInmemoryMedia creates a new FoundMedia with content in memory ; for testing purpose
func NewInmemoryMedia(name string, size uint, date time.Time) FoundMedia {
	return &InmemoryMedia{name, size, date}
}

func (i *InmemoryMedia) Filename() string {
	return i.filename
}

func (i *InmemoryMedia) LastModificationDate() time.Time {
	return i.mediaDate
}

func (i *InmemoryMedia) SimpleSignature() *SimpleMediaSignature {
	return &SimpleMediaSignature{
		RelativePath: i.filename,
		Size:         i.size,
	}
}

func (i *InmemoryMedia) ReadMedia() (io.Reader, error) {
	return os.Open("../test_resources/scan/london_skyline_southbank.jpg")
}

func (m *InmemoryMediaWithHash) Sha256Hash() string {
	return m.hash
}
