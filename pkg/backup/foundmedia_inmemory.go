package backup

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"time"
)

type InMemoryMedia struct {
	filename string
	date     time.Time
	content  []byte
}

// NewInMemoryMedia creates a new FoundMedia for TESTING PURPOSE ONLY
func NewInMemoryMedia(name string, date time.Time, content []byte) FoundMedia {
	return &InMemoryMedia{filename: name, date: date, content: content}
}

func (i *InMemoryMedia) Size() int {
	return len(i.content)
}

func (i *InMemoryMedia) String() string {
	return fmt.Sprintf("RAM/%s [%d bytes]", i.filename, i.Size())
}

func (i *InMemoryMedia) MediaPath() MediaPath {
	return MediaPath{
		ParentFullPath: path.Join("/ram", path.Dir(i.filename)),
		Root:           "/ram",
		Path:           path.Dir(i.filename),
		Filename:       path.Base(i.filename),
		ParentDir:      path.Base(path.Dir(i.filename)),
	}
}

func (i *InMemoryMedia) LastModification() time.Time {
	return i.date
}

func (i *InMemoryMedia) ReadMedia() (io.ReadCloser, error) {
	return &readerCloserWrapper{bytes.NewReader(i.content)}, nil
}

type readerCloserWrapper struct {
	io.Reader
}

func (r *readerCloserWrapper) Close() error {
	// do nothing
	return nil
}
