package archive_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

const owner = "ironman"

func TestStore(t *testing.T) {
	opener := func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte("foobar"))), nil
	}

	baseRequest := &archive.StoreRequest{
		DateTime:         time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
		FolderName:       "/folder-1",
		Id:               "media-1",
		Open:             opener,
		OriginalFilename: "randomName.photo.JPG",
		Owner:            owner,
		SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
	}

	repositoryWithExistingMedia := func() *ARepositoryInMemory {
		repository := NewARepositoryInMemory()
		_ = repository.AddLocation(owner, "media-1", owner+"/folder-1/previous_id.jpg")
		return repository
	}

	type fields struct {
		Repository *ARepositoryInMemory
		Store      *StoreInMemory
		AsyncJob   *AsyncJobInMemory
	}
	tests := []struct {
		name                 string
		fields               fields
		request              *archive.StoreRequest
		want                 string
		wantErr              assert.ErrorAssertionFunc
		expectStore          map[string][]byte
		expectMediasForOwner map[string]string
		expectLoadedImages   [][]*archive.ImageToResize
	}{
		{
			name: "it should upload the media, index it, and queue a cache warm-up for the image",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
			},
			request: baseRequest,
			want:    "2022-06-26_15-48-42_qwertyui.jpg",
			wantErr: assert.NoError,
			expectStore: map[string][]byte{
				owner + "/folder-1/2022-06-26_15-48-42_qwertyui.jpg": []byte("foobar"),
			},
			expectMediasForOwner: map[string]string{
				"media-1": owner + "/folder-1/2022-06-26_15-48-42_qwertyui.jpg",
			},
			expectLoadedImages: [][]*archive.ImageToResize{
				{{Owner: owner, MediaId: "media-1", Widths: archive.CacheableWidths}},
			},
		},
		{
			name: "it should skip upload when the media is already indexed",
			fields: fields{
				Repository: repositoryWithExistingMedia(),
				Store:      NewStoreInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
			},
			request:     baseRequest,
			want:        "previous_id.jpg",
			wantErr:     assert.NoError,
			expectStore: map[string][]byte{},
			expectMediasForOwner: map[string]string{
				"media-1": owner + "/folder-1/previous_id.jpg",
			},
		},
		{
			name: "it should upload a non-resizable media without queuing a cache warm-up",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
			},
			request: &archive.StoreRequest{
				DateTime:         time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
				FolderName:       "/folder-1",
				Id:               "video-1",
				Open:             opener,
				OriginalFilename: "randomName.photo.Mpeg",
				Owner:            owner,
				SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
			},
			want:    "2022-06-26_15-48-42_qwertyui.mpeg",
			wantErr: assert.NoError,
			expectStore: map[string][]byte{
				owner + "/folder-1/2022-06-26_15-48-42_qwertyui.mpeg": []byte("foobar"),
			},
			expectMediasForOwner: map[string]string{
				"video-1": owner + "/folder-1/2022-06-26_15-48-42_qwertyui.mpeg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive.Init(tt.fields.Repository, tt.fields.Store, NewCacheInMemory(), tt.fields.AsyncJob)

			got, err := archive.Store(tt.request)
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.expectStore, tt.fields.Store.Content)
			assert.Equal(t, tt.expectMediasForOwner, tt.fields.Repository.Media[tt.request.Owner])
			assert.Equal(t, tt.expectLoadedImages, tt.fields.AsyncJob.LoadedImages)
		})
	}
}

// TestStore_shouldNotIndexIfUploadFails is a data-loss safety test (E2). It uses testify/mock
// inline to inject an Upload failure on the StoreAdapter — the Fake would otherwise succeed
// unconditionally and we could not verify the "no index update after failed upload" invariant.
func TestStore_shouldNotIndexIfUploadFails(t *testing.T) {
	opener := func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte("foobar"))), nil
	}

	repository := NewARepositoryInMemory()

	failingStore := new(failingUploadStore)
	failingStore.On("Upload", mock.Anything, mock.Anything).Return("", errors.Errorf("TEST - simulate failure while uploading"))

	archive.Init(repository, failingStore, NewCacheInMemory(), NewAsyncJobInMemory())

	_, err := archive.Store(&archive.StoreRequest{
		DateTime:         time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
		FolderName:       "/folder-1",
		Id:               "media-1",
		Open:             opener,
		OriginalFilename: "randomName.photo.JPG",
		Owner:            owner,
		SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
	})

	assert.Error(t, err)
	_, err = repository.FindById(owner, "media-1")
	assert.ErrorIs(t, err, archive.NotFoundError, "repository must not be updated when upload fails")
}

type failingUploadStore struct {
	StoreInMemory
	mock.Mock
}

func (f *failingUploadStore) Upload(values archive.DestructuredKey, content io.Reader) (string, error) {
	args := f.Called(values, content)
	return args.String(0), args.Error(1)
}
