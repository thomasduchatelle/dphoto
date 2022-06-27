package backup

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"time"
)

type InmemoryMedia struct {
	filename string
	date     time.Time
	content  []byte
}

// NewInmemoryMedia creates a new FoundMedia for TESTING PURPOSE ONLY
func NewInmemoryMedia(name string, date time.Time, content []byte) FoundMedia {
	return &InmemoryMedia{filename: name, date: date, content: content}
}

func (i *InmemoryMedia) Size() int {
	return len(i.content)
}

func (i *InmemoryMedia) String() string {
	return fmt.Sprintf("RAM/%s [%d bytes]", i.filename, i.Size())
}

func (i *InmemoryMedia) MediaPath() MediaPath {
	return MediaPath{
		ParentFullPath: path.Join("/ram", path.Dir(i.filename)),
		Root:           "/ram",
		Path:           path.Dir(i.filename),
		Filename:       path.Base(i.filename),
		ParentDir:      path.Base(path.Dir(i.filename)),
	}
}

func (i *InmemoryMedia) ReadMedia() (io.ReadCloser, error) {
	return &readerCloserWrapper{bytes.NewReader(i.content)}, nil
}

type readerCloserWrapper struct {
	io.Reader
}

func (r *readerCloserWrapper) Close() error {
	// do nothing
	return nil
}
