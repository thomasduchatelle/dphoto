package archive

import (
	"io"
	"time"
)

type StoreRequest struct {
	Content          io.ReadCloser
	DateTime         time.Time
	FolderName       string
	Id               string
	OriginalFilename string
	Owner            string
	Size             uint
}
