package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

const (
	layout                  = "2006-01-02T15"
	owner  ownermodel.Owner = "ironman"
)

var (
	albumId1 = catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/MyAlbum")}
)

func TestAlbumQueries_FindAlbum(t *testing.T) {
	album1 := &catalog.Album{
		AlbumId: albumId1,
		Name:    "My Album 1",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	type fields struct {
		Repository catalog.RepositoryAdapter
	}
	type args struct {
		ctx     context.Context
		albumId catalog.AlbumId
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *catalog.Album
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "it should return the album that has been found",
			fields: fields{Repository: NewAlbumRepositoryInMemory(album1)},
			args: args{
				ctx:     context.TODO(),
				albumId: albumId1,
			},
			want:    album1,
			wantErr: assert.NoError,
		},
		{
			name:   "it should return a not found error if no album has been found",
			fields: fields{Repository: NewAlbumRepositoryInMemory()},
			args: args{
				ctx:     context.TODO(),
				albumId: albumId1,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &catalog.AlbumQueries{
				Repository: tt.fields.Repository,
			}
			got, err := a.FindAlbum(tt.args.ctx, tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("FindAlbum(%v, %v)", tt.args.ctx, tt.args.albumId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindAlbum(%v, %v)", tt.args.ctx, tt.args.albumId)
		})
	}
}
