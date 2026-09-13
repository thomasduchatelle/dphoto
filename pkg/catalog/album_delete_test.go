package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestDeleteAlbum_DeleteAlbum(t *testing.T) {
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	const owner = "ironman"
	toDeleteAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := catalog.Album{
		AlbumId: toDeleteAlbumId,
		Name:    "Avenger 1",
		Start:   mar24,
		End:     may24,
	}
	existingAllYearAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")},
		Name:    "lifetime",
		Start:   jan24,
		End:     jan25,
	}
	existingQ1Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q1")},
		Name:    "q1",
		Start:   jan24,
		End:     apr24,
	}
	existingQ2Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q2")},
		Name:    "q2",
		Start:   apr24,
		End:     jul24,
	}

	type fields struct {
		albums           []*catalog.Album
		orphanMediaCount int
	}
	type args struct {
		albumId catalog.AlbumId
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantDeleted   []catalog.AlbumId
		wantTransfers []catalog.MediaTransferRecords
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "it should delete album if all segments can be transferred to 1 other album",
			fields: fields{
				albums: []*catalog.Album{&existingAllYearAlbum, &toDeleteAlbum},
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantDeleted: []catalog.AlbumId{toDeleteAlbumId},
			wantTransfers: []catalog.MediaTransferRecords{{
				existingAllYearAlbum.AlbumId: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
						Start:      toDeleteAlbum.Start,
						End:        toDeleteAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete album if all segments can be transferred to several other albums",
			fields: fields{
				albums: []*catalog.Album{&existingQ1Album, &existingQ2Album, &toDeleteAlbum},
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantDeleted: []catalog.AlbumId{toDeleteAlbumId},
			wantTransfers: []catalog.MediaTransferRecords{{
				existingQ1Album.AlbumId: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
						Start:      toDeleteAlbum.Start,
						End:        existingQ1Album.End,
					},
				},
				existingQ2Album.AlbumId: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
						Start:      existingQ2Album.Start,
						End:        toDeleteAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete album even if the segments are not covered by other albums as long as there is no medias to become orphaned",
			fields: fields{
				albums:           []*catalog.Album{&existingQ1Album, &toDeleteAlbum},
				orphanMediaCount: 0,
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantDeleted: []catalog.AlbumId{toDeleteAlbumId},
			wantTransfers: []catalog.MediaTransferRecords{{
				existingQ1Album.AlbumId: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
						Start:      toDeleteAlbum.Start,
						End:        existingQ1Album.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should raise an error if medias are about to be orphaned",
			fields: fields{
				albums:           []*catalog.Album{&existingQ1Album, &toDeleteAlbum},
				orphanMediaCount: 1,
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewAlbumRepositoryInMemory(tt.fields.albums...)
			repository.MediaCountsBySelect[owner] = tt.fields.orphanMediaCount
			observer := &DeleteAlbumObserverInMemory{}

			d := &catalog.DeleteAlbum{
				FindAlbumsByOwner:      repository,
				CountMediasBySelectors: repository,
				Observers:              []catalog.DeleteAlbumObserver{observer},
			}
			err := d.DeleteAlbum(context.Background(), tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteAlbum(%v)", tt.args.albumId)) {
				return
			}
			assert.Equal(t, tt.wantDeleted, observer.Deleted)
			assert.Equal(t, tt.wantTransfers, observer.Transfers)
		})
	}
}

type ExternalTimelineMutationObserver struct {
	Transfers catalog.TransferredMedias
}

func (e *ExternalTimelineMutationObserver) OnTransferredMedias(ctx context.Context, transfers catalog.TransferredMedias) error {
	e.Transfers = transfers
	return nil
}

func TestNewDeleteAlbum(t *testing.T) {
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	const owner = "ironman"
	toDeleteAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := catalog.Album{
		AlbumId: toDeleteAlbumId,
		Name:    "Avenger 1",
		Start:   mar24,
		End:     may24,
	}
	existingAllYearAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")},
		Name:    "lifetime",
		Start:   jan24,
		End:     jan25,
	}

	externalObserver := new(ExternalTimelineMutationObserver)

	transferredMedias := catalog.TransferredMedias{
		Transfers: map[catalog.AlbumId][]catalog.MediaId{
			existingAllYearAlbum.AlbumId: {"media-1", "media-2"},
		},
	}

	repository := NewAlbumRepositoryInMemory(&existingAllYearAlbum, &toDeleteAlbum)
	transfer := &MediaTransferInMemory{TransferredMedias: transferredMedias}

	deleteAlbum := catalog.NewDeleteAlbum(
		repository,
		repository,
		transfer,
		repository,
		externalObserver,
	)

	err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
	if assert.NoError(t, err) {
		assert.Equal(t, externalObserver.Transfers, catalog.TransferredMedias{
			Transfers:  transferredMedias.Transfers,
			FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
		})
		_, stillExists := repository.Albums[toDeleteAlbumId]
		assert.False(t, stillExists, "the album should be deleted from the repository")
	}
}
