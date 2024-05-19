package catalog_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
	"testing"
	"time"
)

func mockAdapters(t *testing.T) *mocks.RepositoryAdapter {
	mockRepository := mocks.NewRepositoryAdapter(t)
	catalog.Init(mockRepository)

	return mockRepository
}

const (
	layout                  = "2006-01-02T15"
	owner  ownermodel.Owner = "ironman"
)

var (
	albumId1 = catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/MyAlbum")}
)

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockRepository := mockAdapters(t)

	album := &catalog.Album{
		AlbumId: albumId1,
		Name:    "My Album",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbumsByOwner", mock.Anything, owner).Return([]*catalog.Album{album}, nil)

	got, err := catalog.FindAllAlbums(owner)
	if a.NoError(err) {
		a.Equal([]*catalog.Album{album}, got)
	}
}

func TestAlbumQueries_FindAlbum(t *testing.T) {
	albumId1 = catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/MyAlbum")}
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
			name: "it should return the album that has been found",
			fields: fields{
				Repository: RepositoryAdapterFake{album1},
			},
			args: args{
				ctx:     context.TODO(),
				albumId: albumId1,
			},
			want:    album1,
			wantErr: assert.NoError,
		},
		{
			name: "it should return a not found error if no album has been found",
			fields: fields{
				Repository: RepositoryAdapterFake{},
			},
			args: args{
				ctx:     context.TODO(),
				albumId: albumId1,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundError, i...)
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

type RepositoryAdapterFake []*catalog.Album

func (r RepositoryAdapterFake) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	panic("implement me")
}

func (r RepositoryAdapterFake) FindAlbumByIds(ctx context.Context, ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, album := range r {
		if slices.Contains(ids, album.AlbumId) {
			albums = append(albums, album)
		}
	}
	return albums, nil
}

func (r RepositoryAdapterFake) FindMedias(ctx context.Context, request *catalog.FindMediaRequest) (medias []*catalog.MediaMeta, err error) {
	panic("implement me")
}

func (r RepositoryAdapterFake) FindMediaCurrentAlbum(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (id *catalog.AlbumId, err error) {
	panic("implement me")
}

func (r RepositoryAdapterFake) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	panic("implement me")
}
