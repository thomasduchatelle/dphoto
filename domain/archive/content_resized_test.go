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
	unreadableReader := io.NopCloser(new(failWhenRead))

	type args struct {
		owner    string
		mediaId  string
		width    int
		maxBytes int
	}
	tests := []struct {
		name        string
		args        args
		initMocks   func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter)
		wantContent []byte
		wantType    string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should resize the image and store the results when the cache is empty",
			args: args{owner, mediaId, 1440, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=1440"+cacheIdSuffix).Once().Return(nil, 0, "", archive.NotFoundError)

				fullContentReader := io.NopCloser(bytes.NewReader(fullContent))
				repository.On("FindById", owner, mediaId).Once().Return("main-store-key-01", nil)
				store.On("Download", "main-store-key-01").Once().Return(fullContentReader, nil)

				resizer.On("ResizeImage", fullContentReader, 1440, false).Once().Return(resizedContent, mediaType, nil)
				cache.On("Put", "w=1440"+cacheIdSuffix, mediaType, mock.Anything).Once().Return(func(id string, mediaType string, reader io.Reader) error {
					content, err := io.ReadAll(reader)
					if assert.NoError(t, err) {
						assert.Equal(t, resizedContent, content)
					}

					return nil
				})

				asyncJob.On("WarmUpCacheByFolder", owner, "main-store-key-01", 1440).Once().Return(nil)
			},
			wantContent: resizedContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should use cached image if on the right size",
			args: args{owner, mediaId, 1440, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=1440"+cacheIdSuffix).Once().Return(io.NopCloser(bytes.NewReader(resizedContent)), 42, mediaType, nil)
			},
			wantContent: resizedContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should store a miniature image in the cache and return a smaller one",
			args: args{owner, mediaId, 180, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
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

				asyncJob.On("WarmUpCacheByFolder", owner, "main-store-key-01", archive.MiniatureCachedWidth).Once().Return(nil)
			},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should get the miniature image from the cache and return a smaller one",
			args: args{owner, mediaId, 180, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				resizedContentReader := io.NopCloser(bytes.NewReader(resizedContent))
				cache.On("Get", "miniatures"+cacheIdSuffix).Once().Return(resizedContentReader, 42, mediaType, nil)

				resizer.On("ResizeImage", resizedContentReader, 180, true).Once().Return(miniContent, mediaType, nil)
			},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should use the appropriate cached width and resize after",
			args: args{owner, mediaId, 1024, 0},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				resizedContentReader := io.NopCloser(bytes.NewReader(resizedContent))
				cache.On("Get", "w=1440"+cacheIdSuffix).Once().Return(resizedContentReader, 42, mediaType, nil)

				resizer.On("ResizeImage", resizedContentReader, 1024, true).Once().Return(miniContent, mediaType, nil)
			},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return an overflow error when the image is too big after having storing it",
			args: args{owner, mediaId, archive.MediumQualityCachedWidth, 8},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				cacheKey := fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix)
				cache.On("Get", cacheKey).Once().Return(nil, 0, "", archive.NotFoundError)

				fullContentReader := io.NopCloser(bytes.NewReader(fullContent))
				repository.On("FindById", owner, mediaId).Once().Return("main-store-key-01", nil)
				store.On("Download", "main-store-key-01").Once().Return(fullContentReader, nil)

				resizer.On("ResizeImage", fullContentReader, archive.MediumQualityCachedWidth, false).Once().Return(resizedContent, mediaType, nil)
				cache.On("Put", cacheKey, mediaType, mock.Anything).Once().Return(nil)

				asyncJob.On("WarmUpCacheByFolder", owner, "main-store-key-01", archive.MediumQualityCachedWidth).Once().Return(nil)
			},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, err, archive.MediaOverflowError, i)
			},
		},
		{
			name: "it should return an overflow error when the cached image is too big",
			args: args{owner, mediaId, archive.MediumQualityCachedWidth, 41},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				cacheKey := fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix)
				cache.On("Get", cacheKey).Once().Return(unreadableReader, 42, mediaType, nil)
			},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.MediaOverflowError, err, i)
			},
		},
		{
			name: "it should return an overflow error when the resized image is too big",
			args: args{owner, mediaId, 1024, 8},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				resizedContentReader := io.NopCloser(bytes.NewReader(resizedContent))
				cache.On("Get", "w=1440"+cacheIdSuffix).Once().Return(resizedContentReader, 40, mediaType, nil)

				resizer.On("ResizeImage", resizedContentReader, 1024, true).Once().Return(miniContent, mediaType, nil)
			},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.MediaOverflowError, err, i)
			},
		},
		{
			name: "it should return the resized image even if the cached version is too big",
			args: args{owner, mediaId, 1024, 16},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				resizedContentReader := io.NopCloser(bytes.NewReader(resizedContent))
				cache.On("Get", "w=1440"+cacheIdSuffix).Once().Return(resizedContentReader, 40, mediaType, nil)

				resizer.On("ResizeImage", resizedContentReader, 1024, true).Once().Return(miniContent, mediaType, nil)
			},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return not found if the image is unknown",
			args: args{owner, mediaId, 1440, 8},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {
				cache.On("Get", "w=1440"+cacheIdSuffix).Once().Return(nil, 0, "", archive.NotFoundError)
				repository.On("FindById", owner, mediaId).Once().Return("", archive.NotFoundError)
			},
			wantContent: nil,
			wantType:    "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.NotFoundError, err, i)
			},
		},
		{
			name: "it should reject width request higher than max cached resolution",
			args: args{owner, mediaId, 151000, 16},
			initMocks: func(t *testing.T, repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter, cache *mocks.CacheAdapter, asyncJob *mocks.AsyncJobAdapter, resizer *mocks.ResizerAdapter) {

			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewARepositoryAdapter(t)
			store := mocks.NewStoreAdapter(t)
			cache := mocks.NewCacheAdapter(t)
			resizer := mocks.NewResizerAdapter(t)
			asyncJob := mocks.NewAsyncJobAdapter(t)
			tt.initMocks(t, repository, store, cache, asyncJob, resizer)
			archive.ResizerPort = resizer
			archive.Init(repository, store, cache, asyncJob)

			archive.CacheableWidths = []int{archive.MediumQualityCachedWidth, 1440, archive.MiniatureCachedWidth}

			gotContent, gotMediaType, err := archive.GetResizedImage(tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)
			if !tt.wantErr(t, err, fmt.Sprintf("GetResizedImage(%v, %v, %v, %v)", tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)) {
				return
			}
			assert.Equal(t, tt.wantContent, gotContent)
			assert.Equal(t, tt.wantType, gotMediaType)
		})
	}
}

func TestGetResizedImageURL(t *testing.T) {
	t.Run("it should pass-through the request to the cache", func(t *testing.T) {
		cacheAdapter := mocks.NewCacheAdapter(t)
		archive.Init(mocks.NewARepositoryAdapter(t), mocks.NewStoreAdapter(t), cacheAdapter, mocks.NewAsyncJobAdapter(t))

		cacheAdapter.On("SignedURL", "miniatures/ironman@avenger.hero/id-01", archive.DownloadUrlValidityDuration).Once().Return("https://id-01.example.com", nil)

		gotUrl, gotErr := archive.GetResizedImageURL("ironman@avenger.hero", "id-01", 200)
		if assert.NoError(t, gotErr) {
			assert.Equal(t, "https://id-01.example.com", gotUrl)
		}
	})
}

type failWhenRead struct {
}

func (f failWhenRead) Read(p []byte) (n int, err error) {
	panic("DO NOT READ ME")
}
