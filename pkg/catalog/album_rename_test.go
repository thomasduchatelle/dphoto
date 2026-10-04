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

func TestRenameAlbum_RenameAlbum(t *testing.T) {
	const owner = "ironman"
	may24 := time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)
	newName := "Avenger 1"

	existingAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/avenger"),
		},
		Name:  "Avenger",
		Start: may24,
		End:   jun24,
	}
	generatedRenamedId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/2024-05_Avenger_1")}
	forcedRenamedId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/Avengers_vs_Loki")}
	testError := errors.Errorf("TEST error throwing")

	renameFolderRequest := catalog.RenameAlbumRequest{
		CurrentId:        existingAlbum.AlbumId,
		NewName:          newName,
		RenameFolder:     true,
		ForcedFolderName: "",
	}

	repositoryWithAvenger := func() *AlbumRepositoryInMemory {
		freshAlbum := *existingAlbum
		return NewAlbumRepositoryInMemory(&freshAlbum)
	}

	transferFromAvengerTo := func(newId catalog.AlbumId) catalog.MediaTransferRecords {
		return catalog.MediaTransferRecords{
			newId: []catalog.MediaSelector{{
				FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
				Start:      may24,
				End:        jun24,
			}},
		}
	}

	transferredFromAvengerTo := func(newId catalog.AlbumId) catalog.TransferredMedias {
		var ids []catalog.MediaId
		for day := may24; day.Before(jun24); day = day.AddDate(0, 0, 1) {
			ids = append(ids, fakeMediaId(existingAlbum.AlbumId, day))
		}
		return catalog.TransferredMedias{
			Transfers:  map[catalog.AlbumId][]catalog.MediaId{newId: ids},
			FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
		}
	}

	renamedEvent := func(newId catalog.AlbumId) catalog.AlbumRenamed {
		return catalog.AlbumRenamed{
			ExistingAlbum: *existingAlbum,
			RenamedAlbum: catalog.Album{
				AlbumId: newId,
				Name:    newName,
				Start:   may24,
				End:     jun24,
			},
			TransferredMedias: transferredFromAvengerTo(newId),
		}
	}

	cherryPickedCover := catalog.Cover{MediaId: "media-cherry", Filename: "cherry.jpg", Origin: catalog.CoverOriginCherryPicked}
	randomCover := catalog.Cover{MediaId: "media-random", Filename: "random.jpg", Origin: catalog.CoverOriginRandom}
	coversOnAvenger := func() *CoverRepositoryInMemory {
		return NewCoverRepositoryInMemory(CoverRepositorySeed{
			AlbumId: existingAlbum.AlbumId,
			Covers:  []catalog.Cover{cherryPickedCover, randomCover},
		})
	}

	type fields struct {
		AlbumRepository catalog.TimelineRepository
		CoverRepository *CoverRepositoryInMemory
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectAlbumsByIds     map[catalog.AlbumId]string
		expectTransferRecords []catalog.MediaTransferRecords
		expectRenamedEvents   []catalog.AlbumRenamed
		expectCoversByAlbum   map[catalog.AlbumId][]catalog.Cover
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name:   "it should get an error if the new name is empty",
			fields: fields{AlbumRepository: repositoryWithAvenger(), CoverRepository: coversOnAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          "",
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{existingAlbum.AlbumId: existingAlbum.Name},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{existingAlbum.AlbumId: {cherryPickedCover, randomCover}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name:   "it should get an error if the album doesn't exists (name-only path)",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(), CoverRepository: NewCoverRepositoryInMemory()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
			},
		},
		{
			name:   "it should get an error if the album doesn't exists (replace path)",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(), CoverRepository: NewCoverRepositoryInMemory()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
			},
		},
		{
			name:   "it should update the name in place and fire the AlbumRenamed event (folder unchanged, no transfer)",
			fields: fields{AlbumRepository: repositoryWithAvenger(), CoverRepository: coversOnAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds: map[catalog.AlbumId]string{existingAlbum.AlbumId: newName},
			expectRenamedEvents: []catalog.AlbumRenamed{
				{
					ExistingAlbum: *existingAlbum,
					RenamedAlbum: catalog.Album{
						AlbumId: existingAlbum.AlbumId,
						Name:    newName,
						Start:   existingAlbum.Start,
						End:     existingAlbum.End,
					},
				},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{existingAlbum.AlbumId: {cherryPickedCover, randomCover}},
			wantErr:             assert.NoError,
		},
		{
			name:                  "it should replace the album with a new folder name generated from the new name, transferring medias and firing the AlbumRenamed event",
			fields:                fields{AlbumRepository: repositoryWithAvenger(), CoverRepository: coversOnAvenger()},
			args:                  args{request: renameFolderRequest},
			expectAlbumsByIds:     map[catalog.AlbumId]string{generatedRenamedId: newName},
			expectTransferRecords: []catalog.MediaTransferRecords{transferFromAvengerTo(generatedRenamedId)},
			expectRenamedEvents:   []catalog.AlbumRenamed{renamedEvent(generatedRenamedId)},
			expectCoversByAlbum:   map[catalog.AlbumId][]catalog.Cover{generatedRenamedId: {cherryPickedCover, randomCover}},
			wantErr:               assert.NoError,
		},
		{
			name:   "it should replace the album with a forced folder name, transferring medias and firing the AlbumRenamed event",
			fields: fields{AlbumRepository: repositoryWithAvenger(), CoverRepository: coversOnAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "Avengers_vs_Loki",
				},
			},
			expectAlbumsByIds:     map[catalog.AlbumId]string{forcedRenamedId: newName},
			expectTransferRecords: []catalog.MediaTransferRecords{transferFromAvengerTo(forcedRenamedId)},
			expectRenamedEvents:   []catalog.AlbumRenamed{renamedEvent(forcedRenamedId)},
			expectCoversByAlbum:   map[catalog.AlbumId][]catalog.Cover{forcedRenamedId: {cherryPickedCover, randomCover}},
			wantErr:               assert.NoError,
		},
		{
			name:                  "it should migrate no covers when the renamed album had none, leaving both identities empty",
			fields:                fields{AlbumRepository: repositoryWithAvenger(), CoverRepository: NewCoverRepositoryInMemory()},
			args:                  args{request: renameFolderRequest},
			expectAlbumsByIds:     map[catalog.AlbumId]string{generatedRenamedId: newName},
			expectTransferRecords: []catalog.MediaTransferRecords{transferFromAvengerTo(generatedRenamedId)},
			expectRenamedEvents:   []catalog.AlbumRenamed{renamedEvent(generatedRenamedId)},
			expectCoversByAlbum:   map[catalog.AlbumId][]catalog.Cover{},
			wantErr:               assert.NoError,
		},
		{
			name:                  "it should interrupt the transfer and skip firing the event if the album insertion fails",
			fields:                fields{AlbumRepository: failingInsertAlbum(repositoryWithAvenger(), testError), CoverRepository: coversOnAvenger()},
			args:                  args{request: renameFolderRequest},
			expectAlbumsByIds:     map[catalog.AlbumId]string{existingAlbum.AlbumId: existingAlbum.Name},
			expectTransferRecords: nil,
			expectRenamedEvents:   nil,
			expectCoversByAlbum:   map[catalog.AlbumId][]catalog.Cover{existingAlbum.AlbumId: {cherryPickedCover, randomCover}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferService := &TransferMediasServiceFake{}
			observer := &AlbumRenamedObserverInMemory{}

			renameAlbum := catalog.NewRenameAlbum(
				tt.fields.AlbumRepository,
				transferService,
				tt.fields.CoverRepository,
				observer,
			)

			err := renameAlbum.RenameAlbum(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v)", tt.args.request)) {
				return
			}

			repository := underlyingRepository(tt.fields.AlbumRepository)
			names := make(map[catalog.AlbumId]string, len(repository.Albums))
			for id, album := range repository.Albums {
				names[id] = album.Name
			}
			assert.Equal(t, tt.expectAlbumsByIds, names, "album names in the repository")
			assert.Equal(t, tt.expectTransferRecords, transferService.Records, "records passed to TransferMedias")
			assert.Equal(t, tt.expectRenamedEvents, observer.Events, "AlbumRenamed events fired")
			assert.Equal(t, tt.expectCoversByAlbum, tt.fields.CoverRepository.Covers, "covers in the repository")
		})
	}
}

func failingInsertAlbum(memory *AlbumRepositoryInMemory, err error) catalog.TimelineRepository {
	return &failingInsertAlbumInterceptor{
		AlbumRepositoryInMemory: memory,
		err:                     err,
	}
}

type failingInsertAlbumInterceptor struct {
	*AlbumRepositoryInMemory
	err error
}

func (i *failingInsertAlbumInterceptor) InsertAlbum(_ context.Context, _ catalog.Album) error {
	return i.err
}

func underlyingRepository(port catalog.TimelineRepository) *AlbumRepositoryInMemory {
	if repository, ok := port.(*AlbumRepositoryInMemory); ok {
		return repository
	}
	return port.(*failingInsertAlbumInterceptor).AlbumRepositoryInMemory
}
