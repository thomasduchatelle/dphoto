package archive_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestGetMediaOriginalURL(t *testing.T) {
	const owner = "ironman"

	type args struct {
		owner   string
		mediaId string
	}
	tests := []struct {
		name    string
		args    args
		seed    func(repository *ARepositoryInMemory, store *StoreInMemory)
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return the requested media",
			args: args{owner, "id-01"},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory) {
				_ = repository.AddLocation(owner, "id-01", "key-01")
				store.Content["key-01"] = []byte("original")
			},
			want:    "signed://key-01",
			wantErr: assert.NoError,
		},
		{
			name: "it should return not found if the media id doesn't exists",
			args: args{owner, "id-01"},
			seed: func(repository *ARepositoryInMemory, store *StoreInMemory) {},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.NotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewARepositoryInMemory()
			store := NewStoreInMemory()
			tt.seed(repository, store)
			archive.Init(repository, store, NewCacheInMemory(), NewAsyncJobInMemory())

			got, err := archive.GetMediaOriginalURL(tt.args.owner, tt.args.mediaId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)
		})
	}
}
