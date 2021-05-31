package model

import (
	"io"
	"os"
	"time"
)

type InmemoryMedia struct {
	filename string
	size     uint
	date     time.Time
}

type InmemoryMediaWithHash struct {
	*InmemoryMedia
	hash string
}

// NewInmemoryMedia creates a new InmemoryMedia. TESTING PURPOSE ONLY
func NewInmemoryMedia(name string, size uint, date time.Time) FoundMedia {
	return &InmemoryMedia{filename: name, size: size, date: date}
}

// NewInmemoryMediaWithHash creates a new InmemoryMediaWithHash. TESTING PURPOSE ONLY
func NewInmemoryMediaWithHash(name string, size uint, date time.Time, hash string) FoundMedia {
	return &InmemoryMediaWithHash{&InmemoryMedia{filename: name, size: size, date: date}, hash}
}

func (i *InmemoryMedia) Filename() string {
	return i.filename
}

func (i *InmemoryMedia) LastModificationDate() time.Time {
	return i.date
}

func (i *InmemoryMedia) SimpleSignature() *SimpleMediaSignature {
	return &SimpleMediaSignature{
		RelativePath: i.filename,
		Size:         i.size,
	}
}

func (i *InmemoryMedia) ReadMedia() (io.Reader, error) {
	return os.Open("../../../test_resources/scan/london_skyline_southbank.jpg")
}

func (m *InmemoryMediaWithHash) Sha256Hash() string {
	return m.hash
}
