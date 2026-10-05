package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

const deleteOwner = "ironman"

var (
	deleteJan24 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	deleteMar24 = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	deleteApr24 = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	deleteMay24 = time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	deleteNov24 = time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	deleteDec24 = time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	deleteJan25 = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	toDeleteAlbumId = catalog.AlbumId{Owner: deleteOwner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum   = &catalog.Album{AlbumId: toDeleteAlbumId, Name: "Avenger 1", Start: deleteMar24, End: deleteMay24}
	allYearAlbumId  = catalog.AlbumId{Owner: deleteOwner, FolderName: catalog.NewFolderName("/lifetime")}
	allYearAlbum    = &catalog.Album{AlbumId: allYearAlbumId, Name: "lifetime", Start: deleteJan24, End: deleteJan25}
	isolatedAlbumId = catalog.AlbumId{Owner: deleteOwner, FolderName: catalog.NewFolderName("/isolated")}
	isolatedAlbum   = &catalog.Album{AlbumId: isolatedAlbumId, Name: "Isolated", Start: deleteNov24, End: deleteDec24}
	q1AlbumId       = catalog.AlbumId{Owner: deleteOwner, FolderName: catalog.NewFolderName("/q1")}
	q1Album         = &catalog.Album{AlbumId: q1AlbumId, Name: "q1", Start: deleteJan24, End: deleteApr24}

	movingMedia      = deleteMedia("movie-1", deleteMar24.AddDate(0, 0, 10))
	otherMovingMedia = deleteMedia("movie-2", deleteMar24.AddDate(0, 0, 20))
	orphanMedia      = deleteMedia("orphan-1", deleteApr24.AddDate(0, 0, 10))

	movingCherryPicked     = catalog.Cover{MediaId: movingMedia.Id, Filename: movingMedia.Filename, Origin: catalog.CoverOriginCherryPicked}
	otherMovingRandom      = catalog.Cover{MediaId: otherMovingMedia.Id, Filename: otherMovingMedia.Filename, Origin: catalog.CoverOriginRandom}
	ghostCherryPickedCover = catalog.Cover{MediaId: "ghost", Filename: "ghost.jpg", Origin: catalog.CoverOriginCherryPicked}

	deleteTestError = errors.Errorf("TEST error throwing")
)

func deleteMedia(id string, when time.Time) *catalog.MediaMeta {
	return &catalog.MediaMeta{
		Id:       catalog.MediaId(id),
		Filename: id + ".jpg",
		Type:     catalog.MediaTypeImage,
		Details:  catalog.MediaDetails{DateTime: when},
	}
}

func TestDeleteAlbum_DeleteAlbum(t *testing.T) {
	isOrphanedMediasErr := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
	}
	isAlbumNotFoundErr := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
	}
	isTestError := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, deleteTestError, i...)
	}

	type fields struct {
		TimelineRepository catalog.TimelineRepository
		MediaStore         *MediaReadRepositoryInMemory
		CoverRepository    *CoverRepositoryInMemory
	}
	type args struct {
		albumId catalog.AlbumId
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectAlbumIds      []catalog.AlbumId
		expectMediasByAlbum map[catalog.AlbumId][]*catalog.MediaMeta
		expectCoversByAlbum map[catalog.AlbumId][]catalog.Cover
		expectDeletedEvents []catalog.AlbumDeleted
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should transfer medias to the surrounding album and inherit the deleted album's CHERRY_PICKED cover",
			fields: fields{
				TimelineRepository: NewAlbumRepositoryInMemory(allYearAlbum, toDeleteAlbum),
				MediaStore: &MediaReadRepositoryInMemory{Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
					toDeleteAlbumId: {movingMedia, otherMovingMedia},
				}},
				CoverRepository: NewCoverRepositoryInMemory(coversFor(toDeleteAlbumId, movingCherryPicked, otherMovingRandom)),
			},
			args:           args{albumId: toDeleteAlbumId},
			expectAlbumIds: []catalog.AlbumId{allYearAlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				toDeleteAlbumId: {},
				allYearAlbumId:  {movingMedia, otherMovingMedia},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{
				allYearAlbumId: {movingCherryPicked, otherMovingRandom},
			},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId: toDeleteAlbumId,
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{allYearAlbumId: {movingMedia.Id, otherMovingMedia.Id}},
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					allYearAlbumId:  {movingCherryPicked, otherMovingRandom},
					toDeleteAlbumId: nil,
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete an empty album, and fire an empty-transfer event when no surrounding album covers it",
			fields: fields{
				TimelineRepository: NewAlbumRepositoryInMemory(toDeleteAlbum),
				MediaStore:         &MediaReadRepositoryInMemory{Medias: map[catalog.AlbumId][]*catalog.MediaMeta{}},
				CoverRepository:    NewCoverRepositoryInMemory(),
			},
			args:                args{albumId: toDeleteAlbumId},
			expectAlbumIds:      []catalog.AlbumId{},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId:    toDeleteAlbumId,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
				Covers:            nil,
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should return OrphanedMediasErr without touching any album when the deletion would orphan medias",
			fields: fields{
				TimelineRepository: NewAlbumRepositoryInMemory(q1Album, toDeleteAlbum),
				MediaStore: &MediaReadRepositoryInMemory{Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
					toDeleteAlbumId: {orphanMedia},
				}},
				CoverRepository: NewCoverRepositoryInMemory(),
			},
			args:           args{albumId: toDeleteAlbumId},
			expectAlbumIds: []catalog.AlbumId{q1AlbumId, toDeleteAlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				toDeleteAlbumId: {orphanMedia},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:             isOrphanedMediasErr,
		},
		{
			name: "it should return AlbumNotFoundErr when the album does not exist",
			fields: fields{
				TimelineRepository: NewAlbumRepositoryInMemory(allYearAlbum),
				MediaStore:         &MediaReadRepositoryInMemory{},
				CoverRepository:    NewCoverRepositoryInMemory(),
			},
			args:                args{albumId: toDeleteAlbumId},
			expectAlbumIds:      []catalog.AlbumId{allYearAlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta(nil),
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:             isAlbumNotFoundErr,
		},
		{
			name: "it should not fire the event when the repository fails to delete the album after the transfer happened",
			fields: fields{
				TimelineRepository: &failingDeleteAlbumInterceptor{
					AlbumRepositoryInMemory: NewAlbumRepositoryInMemory(allYearAlbum, toDeleteAlbum),
					err:                     deleteTestError,
				},
				MediaStore: &MediaReadRepositoryInMemory{Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
					toDeleteAlbumId: {movingMedia},
				}},
				CoverRepository: NewCoverRepositoryInMemory(),
			},
			args:           args{albumId: toDeleteAlbumId},
			expectAlbumIds: []catalog.AlbumId{allYearAlbumId, toDeleteAlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				toDeleteAlbumId: {},
				allYearAlbumId:  {movingMedia},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:             isTestError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &AlbumDeletedObserverInMemory{}
			deleteAlbum := catalog.NewDeleteAlbum(
				tt.fields.TimelineRepository,
				tt.fields.MediaStore,
				&catalog.TransferMediasFromRepository{TransferMediasRepository: tt.fields.MediaStore},
				&catalog.CoverService{
					CoverRepository:     tt.fields.CoverRepository,
					MediaReadRepository: tt.fields.MediaStore,
					Randomiser:          deterministicRandomiser,
				},
				observer,
			)

			err := deleteAlbum.DeleteAlbum(context.Background(), tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteAlbum(%v)", tt.args.albumId)) {
				return
			}

			albumRepository := underlyingDeleteRepository(tt.fields.TimelineRepository)
			var albumIds []catalog.AlbumId
			for id := range albumRepository.Albums {
				albumIds = append(albumIds, id)
			}
			assert.ElementsMatch(t, tt.expectAlbumIds, albumIds, "albums remaining in the repository")
			assert.Equal(t, tt.expectMediasByAlbum, tt.fields.MediaStore.Medias, "medias remaining per album")
			assert.Equal(t, tt.expectCoversByAlbum, tt.fields.CoverRepository.Covers, "canonical covers stored")
			assert.Equal(t, tt.expectDeletedEvents, observer.Events, "AlbumDeleted events fired")
		})
	}
}

type failingDeleteAlbumInterceptor struct {
	*AlbumRepositoryInMemory
	err error
}

func (i *failingDeleteAlbumInterceptor) DeleteAlbum(_ context.Context, _ catalog.AlbumId) error {
	return i.err
}

func underlyingDeleteRepository(port catalog.TimelineRepository) *AlbumRepositoryInMemory {
	if repository, ok := port.(*AlbumRepositoryInMemory); ok {
		return repository
	}
	return port.(*failingDeleteAlbumInterceptor).AlbumRepositoryInMemory
}
