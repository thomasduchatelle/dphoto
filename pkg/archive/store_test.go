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

	repositoryWithMedia := func(id, key string) *ARepositoryInMemory {
		r := NewARepositoryInMemory()
		_ = r.AddLocation(owner, id, key)
		return r
	}

	type fields struct {
		repository *ARepositoryInMemory
		store      *StoreInMemory
	}
	type args struct {
		request *archive.StoreRequest
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		want             string
		wantErr          assert.ErrorAssertionFunc
		wantStoredAtKey  string
		wantNoUpload     bool
		wantResizeQueued bool
	}{
		{
			name:             "it should store a media online with the right names",
			fields:           fields{repository: NewARepositoryInMemory(), store: NewStoreInMemory()},
			args:             args{request: baseRequest("media-1", "randomName.photo.JPG")},
			want:             "2022-06-26_15-48-42_qwertyui.jpg",
			wantErr:          assert.NoError,
			wantStoredAtKey:  owner + "/folder-1/2022-06-26_15-48-42_qwertyui.jpg",
			wantResizeQueued: true,
		},
		{
			name:         "it should not store anything if the media is already present",
			fields:       fields{repository: repositoryWithMedia("media-1", owner+"/folder-1/previous_id.jpg"), store: NewStoreInMemory()},
			args:         args{request: baseRequest("media-1", "randomName.photo.JPG")},
			want:         "previous_id.jpg",
			wantErr:      assert.NoError,
			wantNoUpload: true,
		},
		{
			name:            "it should store a media online without caching it if its extension is not supported",
			fields:          fields{repository: NewARepositoryInMemory(), store: NewStoreInMemory()},
			args:            args{request: baseRequest("video-1", "randomName.photo.Mpeg")},
			want:            "2022-06-26_15-48-42_qwertyui.mpeg",
			wantErr:         assert.NoError,
			wantStoredAtKey: owner + "/folder-1/2022-06-26_15-48-42_qwertyui.mpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asyncJob := NewAsyncJobInMemory()
			archive.ResizerPort = NewResizerInMemory()
			archive.Init(tt.fields.repository, tt.fields.store, NewCacheInMemory(), asyncJob)

			got, err := archive.Store(tt.args.request)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)

			if tt.wantStoredAtKey != "" {
				storedKey, findErr := tt.fields.repository.FindById(owner, tt.args.request.Id)
				if assert.NoError(t, findErr) {
					assert.Equal(t, tt.wantStoredAtKey, storedKey)
				}
				assert.True(t, tt.fields.store.Has(tt.wantStoredAtKey))
			}
			if tt.wantNoUpload {
				assert.Empty(t, tt.fields.store.Content, "no upload should have happened")
				assert.Empty(t, asyncJob.LoadedImages)
			}
			if tt.wantResizeQueued {
				if assert.Len(t, asyncJob.LoadedImages, 1) && assert.Len(t, asyncJob.LoadedImages[0], 1) {
					loaded := asyncJob.LoadedImages[0][0]
					assert.Equal(t, owner, loaded.Owner)
					assert.Equal(t, tt.args.request.Id, loaded.MediaId)
					assert.Equal(t, archive.CacheableWidths, loaded.Widths)
					assert.NotNil(t, loaded.Open)
				}
			} else {
				assert.Empty(t, asyncJob.LoadedImages, "unsupported extensions must not queue a resize")
			}
		})
	}

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
