package archive_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	mocks2 "github.com/thomasduchatelle/dphoto/internal/mocks"
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
		initMocks func(repository *mocks2.ARepositoryAdapter, store *mocks2.StoreAdapter)
		want      string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should return the requested media",
			args: args{owner, "id-01"},
			initMocks: func(repository *mocks2.ARepositoryAdapter, store *mocks2.StoreAdapter) {
				repository.On("FindById", owner, "id-01").Once().Return("key-01", nil)
				store.On("SignedURL", "key-01", archive.DownloadUrlValidityDuration).Once().Return("/a/url?signed", nil)
			},
			want:    "/a/url?signed",
			wantErr: assert.NoError,
		},
		{
			name: "it should return not found if the media id doesn't exists",
			args: args{owner, "id-01"},
			initMocks: func(repository *mocks2.ARepositoryAdapter, store *mocks2.StoreAdapter) {
				repository.On("FindById", owner, "id-01").Once().Return("", archive.NotFoundError)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.NotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks2.NewARepositoryAdapter(t)
			store := mocks2.NewStoreAdapter(t)
			tt.initMocks(repository, store)
			archive.Init(repository, store, mocks2.NewCacheAdapter(t), mocks2.NewAsyncJobAdapter(t))

			got, err := archive.GetMediaOriginalURL(tt.args.owner, tt.args.mediaId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)
		})
	}
}
