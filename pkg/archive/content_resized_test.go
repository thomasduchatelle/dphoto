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

	happyRepository := func() *ARepositoryInMemory {
		r := NewARepositoryInMemory()
		_ = r.AddLocation(owner, mediaId, storeKey)
		return r
	}
	happyStore := func() *StoreInMemory {
		s := NewStoreInMemory()
		s.Content[storeKey] = fullContent
		return s
	}
	seededCacheAt := func(key string, content []byte) *CacheInMemory {
		c := NewCacheInMemory()
		c.Content[key] = content
		c.MediaType[key] = mediaType
		return c
	}

	type fields struct {
		repository *ARepositoryInMemory
		store      *StoreInMemory
		cache      *CacheInMemory
	}
	type args struct {
		owner    string
		mediaId  string
		width    int
		maxBytes int
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantContent       []byte
		wantType          string
		wantErr           assert.ErrorAssertionFunc
		wantCacheHasKey   string
		wantCacheContent  []byte
		wantWarmUpTouched bool
	}{
		{
			name:              "it should resize the image and store the results when the cache is empty",
			fields:            fields{repository: happyRepository(), store: happyStore(), cache: NewCacheInMemory()},
			args:              args{owner, mediaId, 1440, 0},
			wantContent:       resizedAt(1440),
			wantType:          mediaType,
			wantErr:           assert.NoError,
			wantCacheHasKey:   "w=1440" + cacheIdSuffix,
			wantCacheContent:  resizedAt(1440),
			wantWarmUpTouched: true,
		},
		{
			name:        "it should use cached image if on the right size",
			fields:      fields{repository: happyRepository(), store: happyStore(), cache: seededCacheAt("w=1440"+cacheIdSuffix, []byte("pre-cached-1440"))},
			args:        args{owner, mediaId, 1440, 0},
			wantContent: []byte("pre-cached-1440"),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name:              "it should store a miniature image in the cache and return a smaller one",
			fields:            fields{repository: happyRepository(), store: happyStore(), cache: NewCacheInMemory()},
			args:              args{owner, mediaId, 180, 0},
			wantContent:       resizedAt(180),
			wantType:          mediaType,
			wantErr:           assert.NoError,
			wantCacheHasKey:   "miniatures" + cacheIdSuffix,
			wantCacheContent:  resizedAt(archive.MiniatureCachedWidth),
			wantWarmUpTouched: true,
		},
		{
			name:        "it should get the miniature image from the cache and return a smaller one",
			fields:      fields{repository: happyRepository(), store: happyStore(), cache: seededCacheAt("miniatures"+cacheIdSuffix, []byte("pre-cached-mini"))},
			args:        args{owner, mediaId, 180, 0},
			wantContent: resizedAt(180),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name:        "it should use the appropriate cached width and resize after",
			fields:      fields{repository: happyRepository(), store: happyStore(), cache: seededCacheAt("w=1440"+cacheIdSuffix, []byte("pre-cached-1440"))},
			args:        args{owner, mediaId, 1024, 0},
			wantContent: resizedAt(1024),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name:        "it should return an overflow error when the image is too big after having storing it",
			fields:      fields{repository: happyRepository(), store: happyStore(), cache: NewCacheInMemory()},
			args:        args{owner, mediaId, archive.MediumQualityCachedWidth, 8},
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
			fields: fields{
				repository: happyRepository(),
				store:      happyStore(),
				cache: func() *CacheInMemory {
					key := fmt.Sprintf("w=%d%s", archive.MediumQualityCachedWidth, cacheIdSuffix)
					c := NewCacheInMemory()
					c.Content[key] = make([]byte, 42)
					c.MediaType[key] = mediaType
					return c
				}(),
			},
			args:        args{owner, mediaId, archive.MediumQualityCachedWidth, 41},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.MediaOverflowError, err, i)
			},
		},
		{
			name:        "it should return an overflow error when the resized image is too big",
			fields:      fields{repository: happyRepository(), store: happyStore(), cache: seededCacheAt("w=1440"+cacheIdSuffix, []byte("pre-cached-1440"))},
			args:        args{owner, mediaId, 1024, 8},
			wantContent: nil,
			wantType:    mediaType,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.MediaOverflowError, err, i)
			},
		},
		{
			name: "it should return the resized image even if the cached version is too big",
			fields: fields{
				repository: happyRepository(),
				store:      happyStore(),
				cache: func() *CacheInMemory {
					c := NewCacheInMemory()
					c.Content["w=1440"+cacheIdSuffix] = make([]byte, 40)
					c.MediaType["w=1440"+cacheIdSuffix] = mediaType
					return c
				}(),
			},
			args:        args{owner, mediaId, 1024, 16},
			wantContent: resizedAt(1024),
			wantType:    mediaType,
			wantErr:     assert.NoError,
		},
		{
			name:        "it should return not found if the image is unknown",
			fields:      fields{repository: NewARepositoryInMemory(), store: NewStoreInMemory(), cache: NewCacheInMemory()},
			args:        args{owner, mediaId, 1440, 8},
			wantContent: nil,
			wantType:    "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.NotFoundError, err, i)
			},
		},
		{
			name:   "it should reject width request higher than max cached resolution",
			fields: fields{repository: happyRepository(), store: happyStore(), cache: NewCacheInMemory()},
			args:   args{owner, mediaId, 151000, 16},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asyncJob := NewAsyncJobInMemory()
			archive.ResizerPort = NewResizerInMemory()
			archive.Init(tt.fields.repository, tt.fields.store, tt.fields.cache, asyncJob)
			archive.CacheableWidths = []int{archive.MediumQualityCachedWidth, 1440, archive.MiniatureCachedWidth}

			gotContent, gotMediaType, err := archive.GetResizedImage(tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)
			if !tt.wantErr(t, err, fmt.Sprintf("GetResizedImage(%v, %v, %v, %v)", tt.args.owner, tt.args.mediaId, tt.args.width, tt.args.maxBytes)) {
				return
			}
			assert.Equal(t, tt.wantContent, gotContent)
			assert.Equal(t, tt.wantType, gotMediaType)
			if tt.wantCacheHasKey != "" {
				assert.Equal(t, tt.wantCacheContent, tt.fields.cache.Content[tt.wantCacheHasKey], "cached content at %s", tt.wantCacheHasKey)
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
