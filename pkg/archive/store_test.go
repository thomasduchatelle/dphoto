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
	content := io.NopCloser(bytes.NewReader([]byte("foobar")))
	opener := func() (io.ReadCloser, error) {
		return content, nil
	}

	baseRequest := func(id, filename string) *archive.StoreRequest {
		return &archive.StoreRequest{
			DateTime:         time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
			FolderName:       "/folder-1",
			Id:               id,
			Open:             opener,
			OriginalFilename: filename,
			Owner:            owner,
			SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
		}
	}

	t.Run("it should store a media online with the right names", func(t *testing.T) {
		repository := NewARepositoryInMemory()
		store := NewStoreInMemory()
		cache := NewCacheInMemory()
		asyncJob := NewAsyncJobInMemory()
		archive.ResizerPort = NewResizerInMemory()
		archive.Init(repository, store, cache, asyncJob)

		got, err := archive.Store(baseRequest("media-1", "randomName.photo.JPG"))

		expectedKey := owner + "/folder-1/2022-06-26_15-48-42_qwertyui.jpg"
		if assert.NoError(t, err) {
			assert.Equal(t, "2022-06-26_15-48-42_qwertyui.jpg", got)

			storedKey, err := repository.FindById(owner, "media-1")
			if assert.NoError(t, err) {
				assert.Equal(t, expectedKey, storedKey)
			}
			assert.True(t, store.Has(expectedKey))

			if assert.Len(t, asyncJob.LoadedImages, 1) && assert.Len(t, asyncJob.LoadedImages[0], 1) {
				loaded := asyncJob.LoadedImages[0][0]
				assert.Equal(t, owner, loaded.Owner)
				assert.Equal(t, "media-1", loaded.MediaId)
				assert.Equal(t, archive.CacheableWidths, loaded.Widths)
				assert.NotNil(t, loaded.Open)
			}
		}
	})

	t.Run("it should not store anything if the media is already present", func(t *testing.T) {
		repository := NewARepositoryInMemory()
		store := NewStoreInMemory()
		cache := NewCacheInMemory()
		asyncJob := NewAsyncJobInMemory()
		archive.ResizerPort = NewResizerInMemory()
		archive.Init(repository, store, cache, asyncJob)

		existingKey := owner + "/folder-1/previous_id.jpg"
		_ = repository.AddLocation(owner, "media-1", existingKey)

		got, err := archive.Store(baseRequest("media-1", "randomName.photo.JPG"))

		if assert.NoError(t, err) {
			assert.Equal(t, "previous_id.jpg", got)
			assert.False(t, store.Has(existingKey), "no upload should have happened")
			assert.Empty(t, asyncJob.LoadedImages)
		}
	})

	t.Run("it should store a media online without caching it if its extension is not supported", func(t *testing.T) {
		repository := NewARepositoryInMemory()
		store := NewStoreInMemory()
		cache := NewCacheInMemory()
		asyncJob := NewAsyncJobInMemory()
		archive.ResizerPort = NewResizerInMemory()
		archive.Init(repository, store, cache, asyncJob)

		got, err := archive.Store(baseRequest("video-1", "randomName.photo.Mpeg"))

		expectedKey := owner + "/folder-1/2022-06-26_15-48-42_qwertyui.mpeg"
		if assert.NoError(t, err) {
			assert.Equal(t, "2022-06-26_15-48-42_qwertyui.mpeg", got)

			storedKey, err := repository.FindById(owner, "video-1")
			if assert.NoError(t, err) {
				assert.Equal(t, expectedKey, storedKey)
			}
			assert.True(t, store.Has(expectedKey))
			assert.Empty(t, asyncJob.LoadedImages, "unsupported extensions must not queue a resize")
		}
	})

	// E2 — data-loss safety: no orphaned index entry when upload fails.
	// Uses testify/mock inline because a Fake cannot inject an Upload failure.
	t.Run("it should not index the new location if the upload failed", func(t *testing.T) {
		repository := NewARepositoryInMemory()
		failingStore := &failingUploadStore{StoreInMemory: NewStoreInMemory()}
		failingStore.On("Upload", mock.Anything, mock.Anything).Return("", errors.Errorf("TEST - simulate failure while uploading"))
		cache := NewCacheInMemory()
		asyncJob := NewAsyncJobInMemory()
		archive.ResizerPort = NewResizerInMemory()
		archive.Init(repository, failingStore, cache, asyncJob)

		_, err := archive.Store(baseRequest("media-1", "randomName.photo.JPG"))

		assert.Error(t, err)
		_, findErr := repository.FindById(owner, "media-1")
		assert.ErrorIs(t, findErr, archive.NotFoundError, "no index entry should have been recorded")
	})
}

// failingUploadStore embeds a real Fake but overrides Upload with a testify/mock
// stub so we can inject an error without letting the Fake record the bytes.
type failingUploadStore struct {
	*StoreInMemory
	mock.Mock
}

func (f *failingUploadStore) Upload(values archive.DestructuredKey, content io.Reader) (string, error) {
	args := f.Called(values, content)
	return args.String(0), args.Error(1)
}
