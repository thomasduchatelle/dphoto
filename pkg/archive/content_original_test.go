package archive_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"testing"
)

func TestGetMediaOriginalURL(t *testing.T) {
	const owner = "ironman"

	type args struct {
		owner   string
		mediaId string
	}
	tests := []struct {
		name      string
		args      args
		initMocks func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter)
		want      string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should return the requested media",
			args: args{owner, "id-01"},
			initMocks: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				repository.On("FindById", owner, "id-01").Once().Return("key-01", nil)
				store.On("SignedURL", "key-01", archive.DownloadUrlValidityDuration).Once().Return("/a/url?signed", nil)
			},
			want:    "/a/url?signed",
			wantErr: assert.NoError,
		},
		{
			name: "it should return not found if the media id doesn't exists",
			args: args{owner, "id-01"},
			initMocks: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				repository.On("FindById", owner, "id-01").Once().Return("", archive.NotFoundError)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.NotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewARepositoryAdapter(t)
			store := mocks.NewStoreAdapter(t)
			tt.initMocks(repository, store)
			archive.Init(repository, store, mocks.NewCacheAdapter(t), mocks.NewAsyncJobAdapter(t))

			got, err := archive.GetMediaOriginalURL(tt.args.owner, tt.args.mediaId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)
		})
	}
}
