package archive_test

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/mocks"
	"io"
	"testing"
	"time"
)

const owner = "ironman"

func TestStore(t *testing.T) {
	content := io.NopCloser(bytes.NewReader([]byte("foobar")))
	miniatureContent := []byte("miniature_foobar")

	tests := []struct {
		name             string
		mocksExpectation func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter)
		request          *archive.StoreRequest
		want             string
		wantErr          bool
	}{
		{
			name: "it should store a media online with the right names",
			mocksExpectation: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				repository.On("FindById", owner, "media-1").Once().Return("", archive.NotFoundError)
				repository.On("AddLocation", owner, "media-1", owner+"/folder-1/my_choice.jpg").Once().Return(nil)
				store.On("Upload", archive.DestructuredKey{Prefix: owner + "/folder-1/2022-06-26_15-48-42_qwertyui", Suffix: ".jpg"}, mock.Anything).Once().Return(owner+"/folder-1/my_choice.jpg", nil)

				resizer.On("ResizeImage", mock.Anything, archive.MiniatureCachedWidth, false).Once().Return(miniatureContent, "image/jpg", nil)
				cache.On("Put", fmt.Sprintf("miniatures/%s/media-1", owner), "image/jpg", mock.Anything).Once().Return(nil)
			},
			request: &archive.StoreRequest{
				DateTime:   time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
				FolderName: "/folder-1",
				Id:         "media-1",
				Open: func() (io.ReadCloser, error) {
					return content, nil
				},
				OriginalFilename: "randomName.photo.JPG",
				Owner:            owner,
				SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
			},
			want: "my_choice.jpg",
		},
		{
			name: "it should not store anything is the media is already present",
			mocksExpectation: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				repository.On("FindById", owner, "media-1").Once().Return(owner+"/folder-1/previous_id.jpg", nil)
			},
			request: &archive.StoreRequest{
				DateTime:   time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
				FolderName: "/folder-1",
				Id:         "media-1",
				Open: func() (io.ReadCloser, error) {
					return content, nil
				},
				OriginalFilename: "randomName.photo.JPG",
				Owner:            owner,
				SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
			},
			want: "previous_id.jpg",
		},
		{
			name: "it should not index the new location if the upload failed",
			mocksExpectation: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				repository.On("FindById", owner, "media-1").Once().Return("", archive.NotFoundError)
				store.On("Upload", mock.Anything, mock.Anything).Once().Return("", errors.Errorf("TEST - simulate failure while uploading"))
			},
			request: &archive.StoreRequest{
				DateTime:   time.Date(2022, 6, 26, 15, 48, 42, 0, time.UTC),
				FolderName: "/folder-1",
				Id:         "media-1",
				Open: func() (io.ReadCloser, error) {
					return content, nil
				},
				OriginalFilename: "randomName.photo.JPG",
				Owner:            owner,
				SignatureSha256:  "qwertyuiopasdfghjklzxcvbnm",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			repository := mocks.NewARepositoryAdapter(t)
			store := mocks.NewStoreAdapter(t)
			cache := mocks.NewCacheAdapter(t)
			resizer := mocks.NewResizerAdapter(t)

			tt.mocksExpectation(repository, store, cache, resizer)

			archive.ResizerPort = resizer
			archive.Init(repository, store, cache)

			got, err := archive.Store(tt.request)
			if !tt.wantErr && a.NoError(err, tt.name) {
				a.Equal(tt.want, got, tt.name)
			} else if tt.wantErr {
				a.Error(err, tt.name)
			}
		})
	}
}
