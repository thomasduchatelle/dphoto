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
		Repository archive.ARepositoryAdapter
		Store      archive.StoreAdapter
		AsyncJob   *AsyncJobInMemory
	}
	tests := []struct {
		name              string
		fields            fields
		request           *archive.StoreRequest
		want              string
		wantErr           assert.ErrorAssertionFunc
		expectStoredKey   string
		expectStoredBytes []byte
		expectLoadedImage *archive.ImageToResize
	}{
		{
			name: "it should upload the media, index it, and queue a cache warm-up for the image",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
			},
			request:           baseRequest,
			want:              "2022-06-26_15-48-42_qwertyui.jpg",
			wantErr:           assert.NoError,
			expectStoredKey:   owner + "/folder-1/2022-06-26_15-48-42_qwertyui.jpg",
			expectStoredBytes: []byte("foobar"),
			expectLoadedImage: &archive.ImageToResize{
				Owner:   owner,
				MediaId: "media-1",
				Widths:  archive.CacheableWidths,
			},
		},
		{
			name: "it should skip upload when the media is already indexed",
			fields: fields{
				Repository: repositoryWithExistingMedia(),
				Store:      NewStoreInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
			},
			request: baseRequest,
			want:    "previous_id.jpg",
			wantErr: assert.NoError,
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
			want:              "2022-06-26_15-48-42_qwertyui.mpeg",
			wantErr:           assert.NoError,
			expectStoredKey:   owner + "/folder-1/2022-06-26_15-48-42_qwertyui.mpeg",
			expectStoredBytes: []byte("foobar"),
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

			if tt.expectStoredKey != "" {
				if store, ok := tt.fields.Store.(*StoreInMemory); assert.True(t, ok) {
					assert.Equal(t, tt.expectStoredBytes, store.Content[tt.expectStoredKey])
				}
				if repository, ok := tt.fields.Repository.(*ARepositoryInMemory); assert.True(t, ok) {
					key, err := repository.FindById(tt.request.Owner, tt.request.Id)
					if assert.NoError(t, err) {
						assert.Equal(t, tt.expectStoredKey, key)
					}
				}
			}

			if tt.expectLoadedImage != nil {
				if assert.Len(t, tt.fields.AsyncJob.LoadedImages, 1) && assert.Len(t, tt.fields.AsyncJob.LoadedImages[0], 1) {
					got := tt.fields.AsyncJob.LoadedImages[0][0]
					assert.Equal(t, tt.expectLoadedImage.Owner, got.Owner)
					assert.Equal(t, tt.expectLoadedImage.MediaId, got.MediaId)
					assert.Equal(t, tt.expectLoadedImage.Widths, got.Widths)
					assert.NotNil(t, got.Open)
				}
			} else {
				assert.Empty(t, tt.fields.AsyncJob.LoadedImages)
			}
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
