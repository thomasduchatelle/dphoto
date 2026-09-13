package archive_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestGetMediaOriginalURL(t *testing.T) {
	const owner = "ironman"

	happyRepository := NewARepositoryInMemory()
	_ = happyRepository.AddLocation(owner, "id-01", "key-01")
	happyStore := NewStoreInMemory()
	happyStore.Content["key-01"] = []byte("original")

	type fields struct {
		repository *ARepositoryInMemory
		store      *StoreInMemory
	}
	type args struct {
		owner   string
		mediaId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should return the requested media",
			fields:  fields{repository: happyRepository, store: happyStore},
			args:    args{owner, "id-01"},
			want:    "signed://key-01",
			wantErr: assert.NoError,
		},
		{
			name:   "it should return not found if the media id doesn't exists",
			fields: fields{repository: NewARepositoryInMemory(), store: NewStoreInMemory()},
			args:   args{owner, "id-01"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.NotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive.Init(tt.fields.repository, tt.fields.store, NewCacheInMemory(), NewAsyncJobInMemory())

			got, err := archive.GetMediaOriginalURL(tt.args.owner, tt.args.mediaId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)
		})
	}
}
