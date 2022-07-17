package archive_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/mocks"
	"io"
	"testing"
)

func TestGetResizedImage(t *testing.T) {
	const owner = "ironman@avenger.hero"
	const mediaId = "id-01"
	const cacheIdSuffix = "/ironman@avenger.hero/id-01"
	const mediaType = "image/jpeg"
	fullContent := []byte("full-content-01")
	resizedContent := []byte("resized-content-01")
	miniContent := []byte("mini-content-01")

	type args struct {
		owner    string
		mediaId  string
		width    int
		maxBytes int
	}
	tests := []struct {
		name        string
		args        args
		initMocks   func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter)
		wantContent []byte
		wantType    string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should resize the image and store the results when the cache is empty",
			args: args{owner, mediaId, 1024, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=1024"+cacheIdSuffix).Once().Return(nil, 0, "", archive.NotFoundError)

				fullContentReader := io.NopCloser(bytes.NewReader(fullContent))
				repository.On("FindById", owner, mediaId).Once().Return("main-store-key-01", nil)
				store.On("Download", "main-store-key-01").Once().Return(fullContentReader, nil)

				resizer.On("ResizeImage", fullContentReader, 1024, false).Once().Return(resizedContent, mediaType, nil)
				cache.On("Put", "w=1024"+cacheIdSuffix, mediaType, mock.Anything).Once().Return(func(id string, mediaType string, reader io.Reader) error {
					content, err := io.ReadAll(reader)
					if assert.NoError(t, err) {
						assert.Equal(t, resizedContent, content)
					}

					return nil
				})
			},
			wantContent: resizedContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should use cached image if on the right size",
			args: args{owner, mediaId, 1024, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=1024"+cacheIdSuffix).Once().Return(io.NopCloser(bytes.NewReader(resizedContent)), 42, mediaType, nil)
			},
			wantContent: resizedContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should store a miniature image in the cache and return a smaller one",
			args: args{owner, mediaId, 180, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "miniatures"+cacheIdSuffix).Once().Return(nil, 0, "", archive.NotFoundError)

				fullContentReader := io.NopCloser(bytes.NewReader(fullContent))
				repository.On("FindById", owner, mediaId).Once().Return("main-store-key-01", nil)
				store.On("Download", "main-store-key-01").Once().Return(fullContentReader, nil)

				resizer.On("ResizeImage", fullContentReader, archive.MiniatureCachedWidth, false).Once().Return(resizedContent, mediaType, nil)
				cache.On("Put", "miniatures"+cacheIdSuffix, mediaType, mock.Anything).Once().Return(nil)

				resizer.On("ResizeImage", mock.Anything, 180, true).Once().Return(miniContent, mediaType, func(reader io.Reader, width int, fast bool) error {
					content, err := io.ReadAll(reader)
					if assert.NoError(t, err) {
						assert.Equal(t, resizedContent, content)
					}

					return nil
				})

			},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should get the miniature image from the cache and return a smaller one",
			args: args{owner, mediaId, 180, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				resizedContentReader := io.NopCloser(bytes.NewReader(resizedContent))
				cache.On("Get", "miniatures"+cacheIdSuffix).Once().Return(resizedContentReader, 42, mediaType, nil)

				resizer.On("ResizeImage", resizedContentReader, 180, true).Once().Return(miniContent, mediaType, nil)
			},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return an overflow error when the image is too big after having storing it",
			args: args{owner, mediaId, 4096, 8},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=4096"+cacheIdSuffix).Once().Return(nil, 0, "", archive.NotFoundError)

				fullContentReader := io.NopCloser(bytes.NewReader(fullContent))
				repository.On("FindById", owner, mediaId).Once().Return("main-store-key-01", nil)
				store.On("Download", "main-store-key-01").Once().Return(fullContentReader, nil)

				resizer.On("ResizeImage", fullContentReader, 4096, false).Once().Return(resizedContent, mediaType, nil)
				cache.On("Put", "w=4096"+cacheIdSuffix, mediaType, mock.Anything).Once().Return(nil)
			},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, err, archive.MediaOverflowError, i)
			},
		},
		{
			name: "it should return an overflow error when the cached image is too big",
			args: args{owner, mediaId, 4096, 41},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=4096"+cacheIdSuffix).Once().Return(io.NopCloser(bytes.NewReader(miniContent)), 42, mediaType, nil)
			},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.MediaOverflowError, err, i)
			},
		},
		{
			name: "it should return not found if the image is unknown",
			args: args{owner, mediaId, 1024, 8},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=1024"+cacheIdSuffix).Once().Return(nil, 0, "", archive.NotFoundError)
				repository.On("FindById", owner, mediaId).Once().Return("", archive.NotFoundError)
			},
			wantContent: nil,
			wantType:    "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.NotFoundError, err, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewARepositoryAdapter(t)
			store := mocks.NewStoreAdapter(t)
			cache := mocks.NewCacheAdapter(t)
			resizer := mocks.NewResizerAdapter(t)
			tt.initMocks(t, repository, store, cache, resizer)
			archive.ResizerPort = resizer
			archive.Init(repository, store, cache)

			gotContent, gotMediaType, err := archive.GetResizedImage(tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)
			if !tt.wantErr(t, err, fmt.Sprintf("GetResizedImage(%v, %v, %v, %v)", tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)) {
				return
			}
			assert.Equal(t, tt.wantContent, gotContent)
			assert.Equal(t, tt.wantType, gotMediaType)
		})
	}
}
