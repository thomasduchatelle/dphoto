package archive_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestGetResizedImage(t *testing.T) {
	const owner = "ironman@avenger.hero"
	const mediaId = "id-01"
	const cacheIdSuffix = "/ironman@avenger.hero/id-01"
	const mediaType = "image/jpeg"
	const storeKey = "main-store-key-01"
	fullContent := []byte("full-content-01")

	resizedAt := func(width int) []byte { return []byte(fmt.Sprintf("resized-w=%d", width)) }

	type args struct {
		owner    string
		mediaId  string
		width    int
		maxBytes int
	}
	tests := []struct {
		name              string
		args              args
		seed              func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory)
		wantContent       []byte
		wantType          string
		wantErr           assert.ErrorAssertionFunc
		wantCacheHasKey   string
		wantCacheContent  []byte
		wantWarmUpTouched bool
	}{
		{
			name: "it should resize the image and store the results when the cache is empty",
			args: args{owner, mediaId, 1440, 0},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				_ = repository.AddLocation(owner, mediaId, storeKey)
				store.Content[storeKey] = fullContent
			},
			wantContent:       resizedAt(1440),
			wantType:          mediaType,
			wantErr:           assert.NoError,
			wantCacheHasKey:   "w=1440" + cacheIdSuffix,
			wantCacheContent:  resizedAt(1440),
			wantWarmUpTouched: true,
		},
		{
			name: "it should use cached image if on the right size",
			args: args{owner, mediaId, 1440, 0},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				cache.Content["w=1440"+cacheIdSuffix] = []byte("pre-cached-1440")
				cache.MediaType["w=1440"+cacheIdSuffix] = mediaType
			},
			wantContent: []byte("pre-cached-1440"),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should store a miniature image in the cache and return a smaller one",
			args: args{owner, mediaId, 180, 0},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				_ = repository.AddLocation(owner, mediaId, storeKey)
				store.Content[storeKey] = fullContent
			},
			wantContent:       resizedAt(180),
			wantType:          mediaType,
			wantErr:           assert.NoError,
			wantCacheHasKey:   "miniatures" + cacheIdSuffix,
			wantCacheContent:  resizedAt(archive.MiniatureCachedWidth),
			wantWarmUpTouched: true,
		},
		{
			name: "it should get the miniature image from the cache and return a smaller one",
			args: args{owner, mediaId, 180, 0},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				cache.Content["miniatures"+cacheIdSuffix] = []byte("pre-cached-mini")
				cache.MediaType["miniatures"+cacheIdSuffix] = mediaType
			},
			wantContent: resizedAt(180),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should use the appropriate cached width and resize after",
			args: args{owner, mediaId, 1024, 0},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				cache.Content["w=1440"+cacheIdSuffix] = []byte("pre-cached-1440")
				cache.MediaType["w=1440"+cacheIdSuffix] = mediaType
			},
			wantContent: resizedAt(1024),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return an overflow error when the image is too big after having storing it",
			args: args{owner, mediaId, archive.MediumQualityCachedWidth, 8},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				_ = repository.AddLocation(owner, mediaId, storeKey)
				store.Content[storeKey] = fullContent
			},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.MediaOverflowError, err, i)
			},
			wantCacheHasKey:   fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix),
			wantCacheContent:  resizedAt(archive.MediumQualityCachedWidth),
			wantWarmUpTouched: true,
		},
		{
			name: "it should return an overflow error when the cached image is too big",
			args: args{owner, mediaId, archive.MediumQualityCachedWidth, 41},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				key := fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix)
				cache.Content[key] = make([]byte, 42)
				cache.MediaType[key] = mediaType
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
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				cache.Content["w=1440"+cacheIdSuffix] = []byte("pre-cached-1440")
				cache.MediaType["w=1440"+cacheIdSuffix] = mediaType
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
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {
				cache.Content["w=1440"+cacheIdSuffix] = make([]byte, 40)
				cache.MediaType["w=1440"+cacheIdSuffix] = mediaType
			},
			wantContent: resizedAt(1024),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return not found if the image is unknown",
			args: args{owner, mediaId, 1440, 8},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {},
			wantContent: nil,
			wantType:    "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.NotFoundError, err, i)
			},
		},
		{
			name: "it should reject width request higher than max cached resolution",
			args: args{owner, mediaId, 151000, 16},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory, cache *CacheInMemory) {},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewARepositoryInMemory()
			store := NewStoreInMemory()
			cache := NewCacheInMemory()
			asyncJob := NewAsyncJobInMemory()
			tt.seed(repository, store, cache)
			archive.ResizerPort = NewResizerInMemory()
			archive.Init(repository, store, cache, asyncJob)
			archive.CacheableWidths = []int{archive.MediumQualityCachedWidth, 1440, archive.MiniatureCachedWidth}

			gotContent, gotMediaType, err := archive.GetResizedImage(tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)
			if !tt.wantErr(t, err, fmt.Sprintf("GetResizedImage(%v, %v, %v, %v)", tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)) {
				return
			}
			assert.Equal(t, tt.wantContent, gotContent)
			assert.Equal(t, tt.wantType, gotMediaType)
			if tt.wantCacheHasKey != "" {
				assert.Equal(t, tt.wantCacheContent, cache.Content[tt.wantCacheHasKey], "cached content at %s", tt.wantCacheHasKey)
			}
			if tt.wantWarmUpTouched {
				assert.NotEmpty(t, asyncJob.WarmUpCalls)
			} else {
				assert.Empty(t, asyncJob.WarmUpCalls)
			}
		})
	}
}

func TestGetResizedImageURL(t *testing.T) {
	t.Run("it should pass-through the request to the cache", func(t *testing.T) {
		repository := NewARepositoryInMemory()
		store := NewStoreInMemory()
		cache := NewCacheInMemory()
		archive.Init(repository, store, cache, NewAsyncJobInMemory())

		gotUrl, gotErr := archive.GetResizedImageURL("ironman@avenger.hero", "id-01", 200)
		if assert.NoError(t, gotErr) {
			assert.Equal(t, "signed-cache://miniatures/ironman@avenger.hero/id-01", gotUrl)
		}
	})
}
