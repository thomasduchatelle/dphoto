package backup_test

import (
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/mocks"
	"io"
	"testing"
	"time"
)

func mockDetailsReaderAdapter(t *testing.T) *mocks.DetailsReaderAdapter {
	readerAdapter := mocks.NewDetailsReaderAdapter(t)
	readerAdapter.On("Supports", mock.Anything, backup.MediaTypeImage).Maybe().Return(true)
	readerAdapter.On("Supports", mock.Anything, backup.MediaTypeVideo).Maybe().Return(true)
	readerAdapter.On("ReadDetails", mock.Anything, mock.Anything).Maybe().Return(func(reader io.Reader, options backup.DetailsReaderOptions) *backup.MediaDetails {
		content, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		datetime, err := time.Parse("2006-01-02", string(content)[:10])
		if err != nil {
			panic(err)
		}

		return &backup.MediaDetails{
			DateTime: datetime,
		}
	}, nil)
	return readerAdapter
}

type CapturedMedia struct {
	id         string
	folderName string
	filename   string
}
