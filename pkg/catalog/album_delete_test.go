package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestDeleteAlbum_DeleteAlbum(t *testing.T) {
	const owner = "ironman"
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	nov24 := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	dec24 := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	toDeleteAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := catalog.Album{
		AlbumId: toDeleteAlbumId,
		Name:    "Avenger 1",
		Start:   mar24,
		End:     may24,
	}
	isolatedAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/isolated")}
	isolatedAlbum := catalog.Album{
		AlbumId: isolatedAlbumId,
		Name:    "Isolated",
		Start:   nov24,
		End:     dec24,
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
	testError := errors.Errorf("TEST error throwing")

	mediaIdsFromToDelete := func(start, end time.Time) []catalog.MediaId {
		var ids []catalog.MediaId
		for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
			ids = append(ids, fakeMediaId(toDeleteAlbumId, day))
		}
		return ids
	}

	type fields struct {
		AlbumRepository *AlbumRepositoryInMemory
	}
	type args struct {
		albumId catalog.AlbumId
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectStoredAlbumIds  []catalog.AlbumId
		expectTransferRecords []catalog.MediaTransferRecords
		expectDeletedEvents   []catalog.AlbumDeleted
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name:                 "it should delete the album and transfer every media to the single surrounding album",
			fields:               fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingAllYearAlbum, &toDeleteAlbum)},
			args:                 args{albumId: toDeleteAlbumId},
			expectStoredAlbumIds: []catalog.AlbumId{existingAllYearAlbum.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				existingAllYearAlbum.AlbumId: {{
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
					Start:      mar24,
					End:        may24,
				}},
			}},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId: toDeleteAlbumId,
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{existingAllYearAlbum.AlbumId: mediaIdsFromToDelete(mar24, may24)},
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:                 "it should delete the album and split medias between several surrounding albums",
			fields:               fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingQ1Album, &existingQ2Album, &toDeleteAlbum)},
			args:                 args{albumId: toDeleteAlbumId},
			expectStoredAlbumIds: []catalog.AlbumId{existingQ1Album.AlbumId, existingQ2Album.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				existingQ1Album.AlbumId: {{
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
					Start:      mar24,
					End:        apr24,
				}},
				existingQ2Album.AlbumId: {{
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
					Start:      apr24,
					End:        may24,
				}},
			}},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId: toDeleteAlbumId,
				TransferredMedias: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{
						existingQ1Album.AlbumId: mediaIdsFromToDelete(mar24, apr24),
						existingQ2Album.AlbumId: mediaIdsFromToDelete(apr24, may24),
					},
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:                 "it should delete the album even when part of its range is uncovered as long as no media would be orphaned",
			fields:               fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingQ1Album, &toDeleteAlbum)},
			args:                 args{albumId: toDeleteAlbumId},
			expectStoredAlbumIds: []catalog.AlbumId{existingQ1Album.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				existingQ1Album.AlbumId: {{
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
					Start:      mar24,
					End:        apr24,
				}},
			}},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId: toDeleteAlbumId,
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{existingQ1Album.AlbumId: mediaIdsFromToDelete(mar24, apr24)},
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:                  "it should delete the album and fire the event with an empty TransferredMedias when no surrounding album covers it",
			fields:                fields{AlbumRepository: NewAlbumRepositoryInMemory(&isolatedAlbum, &toDeleteAlbum)},
			args:                  args{albumId: toDeleteAlbumId},
			expectStoredAlbumIds:  []catalog.AlbumId{isolatedAlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{{}},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId:    toDeleteAlbumId,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should return OrphanedMediasErr and change nothing if the deletion would orphan medias",
			fields: fields{AlbumRepository: func() *AlbumRepositoryInMemory {
				repo := NewAlbumRepositoryInMemory(&existingQ1Album, &toDeleteAlbum)
				repo.MediasBySelector = countOneMediaForOrphanSelector
				return repo
			}()},
			args:                 args{albumId: toDeleteAlbumId},
			expectStoredAlbumIds: []catalog.AlbumId{existingQ1Album.AlbumId, toDeleteAlbumId},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
		{
			name:                 "it should return AlbumNotFoundErr when the album does not exist",
			fields:               fields{AlbumRepository: NewAlbumRepositoryInMemory(&existingAllYearAlbum)},
			args:                 args{albumId: toDeleteAlbumId},
			expectStoredAlbumIds: []catalog.AlbumId{existingAllYearAlbum.AlbumId},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferService := &TransferMediasServiceFake{}
			observer := &AlbumDeletedObserverInMemory{}

			deleteAlbum := catalog.NewDeleteAlbum(
				tt.fields.AlbumRepository,
				tt.fields.AlbumRepository,
				transferService,
				tt.fields.AlbumRepository,
				observer,
			)

			err := deleteAlbum.DeleteAlbum(context.Background(), tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteAlbum(%v)", tt.args.albumId)) {
				return
			}

			var storedIds []catalog.AlbumId
			for id := range tt.fields.AlbumRepository.Albums {
				storedIds = append(storedIds, id)
			}
			assert.ElementsMatch(t, tt.expectStoredAlbumIds, storedIds, "albums remaining in the repository")
			assert.Equal(t, tt.expectTransferRecords, transferService.Records, "records passed to TransferMedias")
			assert.Equal(t, tt.expectDeletedEvents, observer.Events, "AlbumDeleted events fired")
		})
	}

	t.Run("it should skip firing the event when the repository fails to delete the album (transfer already happened, no rollback)", func(t *testing.T) {
		repository := NewAlbumRepositoryInMemory(&existingAllYearAlbum, &toDeleteAlbum)
		transferService := &TransferMediasServiceFake{}
		observer := &AlbumDeletedObserverInMemory{}

		deleteAlbum := catalog.NewDeleteAlbum(
			repository,
			repository,
			transferService,
			catalog.DeleteAlbumRepositoryFunc(func(_ context.Context, _ catalog.AlbumId) error { return testError }),
			observer,
		)

		err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
		assert.ErrorIs(t, err, testError)
		assert.Contains(t, repository.Albums, toDeleteAlbumId, "album should still be in the fake repository since the failing port is a different implementation")
		assert.Len(t, transferService.Records, 1, "the transfer should have happened before the deletion attempt")
		assert.Empty(t, observer.Events, "no AlbumDeleted event should be fired when the repository fails")
	})
}
