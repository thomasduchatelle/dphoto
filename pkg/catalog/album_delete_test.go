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

	orphanSelector := catalog.MediaSelector{
		FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
		Start:      apr24,
		End:        may24,
	}
	countOneMediaForOrphanSelector := func(_ ownermodel.Owner, selector catalog.MediaSelector) int {
		if selector.Start.Equal(orphanSelector.Start) && selector.End.Equal(orphanSelector.End) {
			return 1
		}
		return 0
	}

	type fields struct {
		AlbumRepository *AlbumRepositoryInMemory
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
			name:   "it should delete album if all segments can be transferred to 1 other album",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingAllYearAlbum, &toDeleteAlbum)},
			args:   args{albumId: toDeleteAlbumId},
			wantDeleted: []catalog.AlbumId{toDeleteAlbumId},
			wantTransfers: []catalog.MediaTransferRecords{{
				existingAllYearAlbum.AlbumId: {
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
			name:   "it should delete album if all segments can be transferred to several other albums",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingQ1Album, &existingQ2Album, &toDeleteAlbum)},
			args:   args{albumId: toDeleteAlbumId},
			wantDeleted: []catalog.AlbumId{toDeleteAlbumId},
			wantTransfers: []catalog.MediaTransferRecords{{
				existingQ1Album.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
						Start:      toDeleteAlbum.Start,
						End:        existingQ1Album.End,
					},
				},
				existingQ2Album.AlbumId: {
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
			name:   "it should delete album even if the segments are not covered by other albums as long as there is no medias to become orphaned",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingQ1Album, &toDeleteAlbum)},
			args:   args{albumId: toDeleteAlbumId},
			wantDeleted: []catalog.AlbumId{toDeleteAlbumId},
			wantTransfers: []catalog.MediaTransferRecords{{
				existingQ1Album.AlbumId: {
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
			fields: fields{AlbumRepository: func() *AlbumRepositoryInMemory {
				repo := NewAlbumRepositoryInMemory(&existingQ1Album, &toDeleteAlbum)
				repo.MediasBySelector = countOneMediaForOrphanSelector
				return repo
			}()},
			args: args{albumId: toDeleteAlbumId},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &DeleteAlbumObserverInMemory{}
			d := &catalog.DeleteAlbum{
				FindAlbumsByOwner:      tt.fields.AlbumRepository,
				CountMediasBySelectors: tt.fields.AlbumRepository,
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

	transferredMedias := catalog.TransferredMedias{
		Transfers: map[catalog.AlbumId][]catalog.MediaId{
			existingAllYearAlbum.AlbumId: {"media-1", "media-2"},
		},
	}

	albumRepository := NewAlbumRepositoryInMemory(&existingAllYearAlbum, &toDeleteAlbum)
	transferMedias := NewTransferMediasInMemory()
	transferMedias.TransferredMedias = transferredMedias
	timelineObserver := &TimelineMutationObserverInMemory{}

	deleteAlbum := catalog.NewDeleteAlbum(
		albumRepository,
		albumRepository,
		transferMedias,
		albumRepository,
		timelineObserver,
	)

	err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
	if assert.NoError(t, err) {
		assert.NotContains(t, albumRepository.Albums, toDeleteAlbumId, "album should be removed from the repository")
		assert.Equal(t, []catalog.TransferredMedias{{
			Transfers:  transferredMedias.Transfers,
			FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
		}}, timelineObserver.Notifications)
	}
}
