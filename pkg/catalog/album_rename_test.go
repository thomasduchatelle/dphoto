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

	renamedEvent := func(newId catalog.AlbumId) catalog.AlbumRenamed {
		return catalog.AlbumRenamed{
			ExistingAlbum: *existingAlbum,
			RenamedAlbum: catalog.Album{
				AlbumId: newId,
				Name:    newName,
				Start:   may24,
				End:     jun24,
			},
			MediaTransfer: transferFromAvengerTo(newId),
		}
	}

	type fields struct {
		AlbumRepository catalog.FindAndRenameAlbumPort
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
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name:   "it should get an error if the new name is empty",
			fields: fields{AlbumRepository: repositoryWithAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          "",
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds: map[catalog.AlbumId]string{existingAlbum.AlbumId: existingAlbum.Name},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name:   "it should get an error if the album doesn't exists (name-only path)",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds: map[catalog.AlbumId]string{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
			},
		},
		{
			name:   "it should get an error if the album doesn't exists (replace path)",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds: map[catalog.AlbumId]string{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
			},
		},
		{
			name:   "it should update the name in place if the folder name is unchanged (no transfer, no event fired)",
			fields: fields{AlbumRepository: repositoryWithAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			expectAlbumsByIds: map[catalog.AlbumId]string{existingAlbum.AlbumId: newName},
			wantErr:           assert.NoError,
		},
		{
			name:                  "it should replace the album with a new folder name generated from the new name, transferring medias and firing the AlbumRenamed event",
			fields:                fields{AlbumRepository: repositoryWithAvenger()},
			args:                  args{request: renameFolderRequest},
			expectAlbumsByIds:     map[catalog.AlbumId]string{generatedRenamedId: newName},
			expectTransferRecords: []catalog.MediaTransferRecords{transferFromAvengerTo(generatedRenamedId)},
			expectRenamedEvents:   []catalog.AlbumRenamed{renamedEvent(generatedRenamedId)},
			wantErr:               assert.NoError,
		},
		{
			name:   "it should replace the album with a forced folder name, transferring medias and firing the AlbumRenamed event",
			fields: fields{AlbumRepository: repositoryWithAvenger()},
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
			wantErr:               assert.NoError,
		},
		{
			name:                  "it should interrupt the transfer and skip firing the event if the album insertion fails",
			fields:                fields{AlbumRepository: failingInsertAlbum(repositoryWithAvenger(), testError)},
			args:                  args{request: renameFolderRequest},
			expectAlbumsByIds:     map[catalog.AlbumId]string{existingAlbum.AlbumId: existingAlbum.Name},
			expectTransferRecords: nil,
			expectRenamedEvents:   nil,
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
		})
	}
}

func failingInsertAlbum(memory *AlbumRepositoryInMemory, err error) catalog.FindAndRenameAlbumPort {
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

func underlyingRepository(port catalog.FindAndRenameAlbumPort) *AlbumRepositoryInMemory {
	if repository, ok := port.(*AlbumRepositoryInMemory); ok {
		return repository
	}
	return port.(*failingInsertAlbumInterceptor).AlbumRepositoryInMemory
}
