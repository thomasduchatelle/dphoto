package archive_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestGetMediaOriginalURL(t *testing.T) {
	const owner = "ironman"

	repositoryWithMedia := func() *ARepositoryInMemory {
		repository := NewARepositoryInMemory()
		_ = repository.AddLocation(owner, "id-01", "key-01")
		return repository
	}
	storeWithArchivedMedia := func() *StoreInMemory {
		store := NewStoreInMemory()
		store.Content["key-01"] = []byte("original-content-01")
		return store
	}

	type fields struct {
		Repository *ARepositoryInMemory
		Store      *StoreInMemory
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
			name:    "it should return the signed URL for the archived media",
			fields:  fields{Repository: repositoryWithMedia(), Store: storeWithArchivedMedia()},
			args:    args{owner: owner, mediaId: "id-01"},
			want:    "signed://key-01",
			wantErr: assert.NoError,
		},
		{
			name:   "it should return NotFoundError when the media id is unknown",
			fields: fields{Repository: NewARepositoryInMemory(), Store: NewStoreInMemory()},
			args:   args{owner: owner, mediaId: "id-01"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, archive.NotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive.Init(tt.fields.Repository, tt.fields.Store, NewCacheInMemory(), NewAsyncJobInMemory())

			got, err := archive.GetMediaOriginalURL(tt.args.owner, tt.args.mediaId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMediaOriginalURL(%v, %v)", tt.args.owner, tt.args.mediaId)
		})
	}
}
