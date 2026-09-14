package archive_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestGetResizedImage(t *testing.T) {
	const resizedOwner = "ironman@avenger.hero"
	const mediaId = "id-01"
	const cacheIdSuffix = "/ironman@avenger.hero/id-01"
	const mediaType = "image/jpeg"
	const storeKey = "main-store-key-01"
	fullContent := []byte("full-content-01")
	resizedContent := []byte("resized-content-01")
	miniContent := []byte("mini-content-01")

	repositoryWithMedia := func() *ARepositoryInMemory {
		repository := NewARepositoryInMemory()
		_ = repository.AddLocation(resizedOwner, mediaId, storeKey)
		return repository
	}
	storeWithMedia := func() *StoreInMemory {
		store := NewStoreInMemory()
		store.Content[storeKey] = fullContent
		return store
	}
	cacheWith := func(key string, content []byte) *CacheInMemory {
		cache := NewCacheInMemory()
		cache.Content[key] = CacheEntry{MediaType: mediaType, Content: content}
		return cache
	}
	resizerReturning := func(byWidth map[int][]byte) *ResizerInMemory {
		resizer := NewResizerInMemory()
		resizer.MediaType = mediaType
		for w, c := range byWidth {
			resizer.ByWidth[w] = c
		}
		return resizer
	}

	type fields struct {
		Repository *ARepositoryInMemory
		Store      *StoreInMemory
		Cache      *CacheInMemory
		AsyncJob   *AsyncJobInMemory
		Resizer    *ResizerInMemory
	}
	type args struct {
		owner    string
		mediaId  string
		width    int
		maxBytes int
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantContent            []byte
		wantType               string
		wantErr                assert.ErrorAssertionFunc
		expectCachedKey        string
		expectCachedBytes      []byte
		expectPendingWarmUpJob *WarmUpJob
	}{
		{
			name: "it should resize the image and store the result when the cache is empty",
			fields: fields{
				Repository: repositoryWithMedia(),
				Store:      storeWithMedia(),
				Cache:      NewCacheInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(map[int][]byte{1440: resizedContent}),
			},
			args:                   args{resizedOwner, mediaId, 1440, 0},
			wantContent:            resizedContent,
			wantType:               mediaType,
			wantErr:                assert.NoError,
			expectCachedKey:        "w=1440" + cacheIdSuffix,
			expectCachedBytes:      resizedContent,
			expectPendingWarmUpJob: &WarmUpJob{Owner: resizedOwner, MissedKey: storeKey, Width: 1440},
		},
		{
			name: "it should use the cached image when it exists at the requested cacheable width",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      cacheWith("w=1440"+cacheIdSuffix, resizedContent),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(nil),
			},
			args:        args{resizedOwner, mediaId, 1440, 0},
			wantContent: resizedContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should store a miniature in the cache and return a smaller image resized on the fly",
			fields: fields{
				Repository: repositoryWithMedia(),
				Store:      storeWithMedia(),
				Cache:      NewCacheInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer: resizerReturning(map[int][]byte{
					archive.MiniatureCachedWidth: resizedContent,
					180:                          miniContent,
				}),
			},
			args:                   args{resizedOwner, mediaId, 180, 0},
			wantContent:            miniContent,
			wantType:               mediaType,
			wantErr:                assert.NoError,
			expectCachedKey:        "miniatures" + cacheIdSuffix,
			expectCachedBytes:      resizedContent,
			expectPendingWarmUpJob: &WarmUpJob{Owner: resizedOwner, MissedKey: storeKey, Width: archive.MiniatureCachedWidth},
		},
		{
			name: "it should get the miniature image from the cache and return a smaller one resized on the fly",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      cacheWith("miniatures"+cacheIdSuffix, resizedContent),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(map[int][]byte{180: miniContent}),
			},
			args:        args{resizedOwner, mediaId, 180, 0},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should use the appropriate cached width and resize down after",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      cacheWith("w=1440"+cacheIdSuffix, resizedContent),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(map[int][]byte{1024: miniContent}),
			},
			args:        args{resizedOwner, mediaId, 1024, 0},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return an overflow error when the freshly cached image is too big",
			fields: fields{
				Repository: repositoryWithMedia(),
				Store:      storeWithMedia(),
				Cache:      NewCacheInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(map[int][]byte{archive.MediumQualityCachedWidth: resizedContent}),
			},
			args: args{resizedOwner, mediaId, archive.MediumQualityCachedWidth, 8},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.MediaOverflowError)
			},
			wantType:               mediaType,
			expectCachedKey:        fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix),
			expectCachedBytes:      resizedContent,
			expectPendingWarmUpJob: &WarmUpJob{Owner: resizedOwner, MissedKey: storeKey, Width: archive.MediumQualityCachedWidth},
		},
		{
			name: "it should return an overflow error when the cached image is too big",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      cacheWith(fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix), []byte("this-is-a-large-cached-blob")),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(nil),
			},
			args: args{resizedOwner, mediaId, archive.MediumQualityCachedWidth, 8},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.MediaOverflowError)
			},
			wantType: mediaType,
		},
		{
			name: "it should return an overflow error when the on-the-fly resized image is too big",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      cacheWith("w=1440"+cacheIdSuffix, []byte("cached-fits")),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(map[int][]byte{1024: []byte("way-too-big-resized-payload")}),
			},
			args: args{resizedOwner, mediaId, 1024, 8},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.MediaOverflowError)
			},
			wantType: mediaType,
		},
		{
			name: "it should return the resized image even if the cached version is too big for the consumer",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      cacheWith("w=1440"+cacheIdSuffix, []byte("cached-too-big-for-consumer-but-not-returned")),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(map[int][]byte{1024: miniContent}),
			},
			args:        args{resizedOwner, mediaId, 1024, 16},
			wantContent: miniContent,
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name: "it should return NotFoundError when the image id is unknown",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      NewCacheInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(nil),
			},
			args: args{resizedOwner, mediaId, 1440, 8},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.NotFoundError)
			},
		},
		{
			name: "it should reject widths higher than the maximum cacheable resolution",
			fields: fields{
				Repository: NewARepositoryInMemory(),
				Store:      NewStoreInMemory(),
				Cache:      NewCacheInMemory(),
				AsyncJob:   NewAsyncJobInMemory(),
				Resizer:    resizerReturning(nil),
			},
			args: args{resizedOwner, mediaId, 151000, 16},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive.ResizerPort = tt.fields.Resizer
			archive.Init(tt.fields.Repository, tt.fields.Store, tt.fields.Cache, tt.fields.AsyncJob)
			archive.CacheableWidths = []int{archive.MediumQualityCachedWidth, 1440, archive.MiniatureCachedWidth}

			gotContent, gotMediaType, err := archive.GetResizedImage(tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)
			if !tt.wantErr(t, err, fmt.Sprintf("GetResizedImage(%v, %v, %v, %v)", tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)) {
				return
			}
			assert.Equal(t, tt.wantContent, gotContent)
			assert.Equal(t, tt.wantType, gotMediaType)

			if tt.expectCachedKey != "" {
				entry, ok := tt.fields.Cache.Content[tt.expectCachedKey]
				if assert.Truef(t, ok, "expected cache entry %s to be present", tt.expectCachedKey) {
					assert.Equal(t, tt.expectCachedBytes, entry.Content)
				}
			}
			if tt.expectPendingWarmUpJob != nil {
				assert.Equal(t, []WarmUpJob{*tt.expectPendingWarmUpJob}, tt.fields.AsyncJob.PendingWarmUpJobs)
			} else {
				assert.Empty(t, tt.fields.AsyncJob.PendingWarmUpJobs)
			}
		})
	}
}

func TestGetResizedImageURL(t *testing.T) {
	t.Run("it should return a signed URL from the cache adapter for a miniature request", func(t *testing.T) {
		cache := NewCacheInMemory()
		archive.Init(NewARepositoryInMemory(), NewStoreInMemory(), cache, NewAsyncJobInMemory())
		archive.CacheableWidths = []int{archive.MediumQualityCachedWidth, 1440, archive.MiniatureCachedWidth}

		gotUrl, gotErr := archive.GetResizedImageURL("ironman@avenger.hero", "id-01", 200)

		if assert.NoError(t, gotErr) {
			assert.Equal(t, "signed://miniatures/ironman@avenger.hero/id-01", gotUrl)
		}
	})
}
