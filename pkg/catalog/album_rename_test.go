package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestNewRenameAlbumAcceptance(t *testing.T) {
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
	newAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/2024-05_Avenger_1"),
		},
		Name:  newName,
		Start: may24,
		End:   jun24,
	}
	transferredMedias := catalog.TransferredMedias{
		Transfers: map[catalog.AlbumId][]catalog.MediaId{
			newAlbum.AlbumId: {"media-1", "media-2"},
		},
	}
	testError := errors.Errorf("TEST error throwing")

	renameRequest := catalog.RenameAlbumRequest{
		CurrentId:        existingAlbum.AlbumId,
		NewName:          newName,
		RenameFolder:     true,
		ForcedFolderName: "",
	}

	transferMediasWithMedias := func() *TransferMediasInMemory {
		transfer := NewTransferMediasInMemory()
		transfer.TransferredMedias = transferredMedias
		return transfer
	}

	failingInsertAlbum := func() *insertAlbumPortMock {
		m := new(insertAlbumPortMock)
		m.On("InsertAlbum", mock.Anything, newAlbum).Return(testError).Once()
		return m
	}

	type fields struct {
		AlbumRepository  *AlbumRepositoryInMemory
		InsertAlbum      catalog.InsertAlbumPort
		TransferMedias   *TransferMediasInMemory
		TimelineObserver *TimelineMutationObserverInMemory
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantOldAlbumDeleted bool
		wantNewAlbumCreated bool
		wantTransferRecords []catalog.MediaTransferRecords
		wantNotifications   []catalog.TransferredMedias
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should create a new album end to end",
			fields: fields{
				AlbumRepository:  NewAlbumRepositoryInMemory(existingAlbum),
				TransferMedias:   transferMediasWithMedias(),
				TimelineObserver: &TimelineMutationObserverInMemory{},
			},
			args:                args{request: renameRequest},
			wantOldAlbumDeleted: true,
			wantNewAlbumCreated: true,
			wantTransferRecords: []catalog.MediaTransferRecords{{
				newAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
						Start:      existingAlbum.Start,
						End:        existingAlbum.End,
					},
				},
			}},
			wantNotifications: []catalog.TransferredMedias{{
				Transfers:  transferredMedias.Transfers,
				FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
			}},
			wantErr: assert.NoError,
		},
		{
			// E6: uses testify/mock inline for the InsertAlbumPort to force insert to fail. Verify by
			// state that the old album is still present, the transfer never happened, and the observer
			// was not notified (rollback semantics).
			name: "it should interrupt the transfer if the album insertion fails",
			fields: fields{
				AlbumRepository:  NewAlbumRepositoryInMemory(existingAlbum),
				InsertAlbum:      failingInsertAlbum(),
				TransferMedias:   transferMediasWithMedias(),
				TimelineObserver: &TimelineMutationObserverInMemory{},
			},
			args: args{request: renameRequest},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insertAlbum := catalog.InsertAlbumPort(tt.fields.AlbumRepository)
			if tt.fields.InsertAlbum != nil {
				insertAlbum = tt.fields.InsertAlbum
			}

			renameAlbum := catalog.NewRenameAlbum(
				tt.fields.AlbumRepository,
				tt.fields.AlbumRepository,
				insertAlbum,
				tt.fields.AlbumRepository,
				tt.fields.TransferMedias,
				tt.fields.AlbumRepository,
				tt.fields.TimelineObserver,
			)

			err := renameAlbum.RenameAlbum(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v)", tt.args.request)) {
				return
			}

			if tt.wantOldAlbumDeleted {
				assert.NotContains(t, tt.fields.AlbumRepository.Albums, existingAlbum.AlbumId, "the old album should be deleted")
			} else {
				assert.Contains(t, tt.fields.AlbumRepository.Albums, existingAlbum.AlbumId, "the old album must NOT be deleted")
			}
			if tt.wantNewAlbumCreated {
				assert.Contains(t, tt.fields.AlbumRepository.Albums, newAlbum.AlbumId, "the new album should be created")
			} else {
				assert.NotContains(t, tt.fields.AlbumRepository.Albums, newAlbum.AlbumId, "the new album must NOT be created")
			}
			assert.Equal(t, tt.wantTransferRecords, tt.fields.TransferMedias.Records)
			assert.Equal(t, tt.wantNotifications, tt.fields.TimelineObserver.Notifications)
		})
	}
}

func TestRenameAlbum_RenameAlbum(t *testing.T) {
	const owner = "ironman"
	may24 := time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	existingAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/avenger"),
		},
		Name:  "Avenger",
		Start: may24,
		End:   jun24,
	}
	newName := "Avenger 1"

	albumRepositoryWithAvenger := func() *AlbumRepositoryInMemory {
		freshAlbum := *existingAlbum
		return NewAlbumRepositoryInMemory(&freshAlbum)
	}

	type fields struct {
		AlbumRepository *AlbumRepositoryInMemory
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRenamed []RenameAlbumCall
		wantName    string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:   "it should get an error if the new name is empty",
			fields: fields{AlbumRepository: albumRepositoryWithAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          "",
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantName: existingAlbum.Name,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name:   "it should get an error if the album doesn't exists",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
			},
		},
		{
			name:   "it should update the name if the album is found and folder name is unchanged",
			fields: fields{AlbumRepository: albumRepositoryWithAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantName: newName,
			wantErr:  assert.NoError,
		},
		{
			name:   "it should create a new album if the album is found and folder name is changed",
			fields: fields{AlbumRepository: albumRepositoryWithAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			wantName: existingAlbum.Name,
			wantRenamed: []RenameAlbumCall{{
				Current: existingAlbum.AlbumId,
				CreationRequest: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             newName,
					Start:            may24,
					End:              jun24,
					ForcedFolderName: "",
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should create a new album if the album is found and folder name is forced to a certain value",
			fields: fields{AlbumRepository: albumRepositoryWithAvenger()},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "Avengers vs Loki",
				},
			},
			wantName: existingAlbum.Name,
			wantRenamed: []RenameAlbumCall{{
				Current: existingAlbum.AlbumId,
				CreationRequest: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             newName,
					Start:            may24,
					End:              jun24,
					ForcedFolderName: "Avengers vs Loki",
				},
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &RenameAlbumObserverInMemory{}
			r := &catalog.RenameAlbum{
				FindAlbumById:        tt.fields.AlbumRepository,
				UpdateAlbumName:      tt.fields.AlbumRepository,
				RenameAlbumObservers: []catalog.RenameAlbumObserver{observer},
			}
			err := r.RenameAlbum(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v)", tt.args.request)) {
				return
			}
			assert.Equal(t, tt.wantRenamed, observer.Renamed)
			if tt.wantName != "" {
				if stored, ok := tt.fields.AlbumRepository.Albums[existingAlbum.AlbumId]; ok {
					assert.Equal(t, tt.wantName, stored.Name, "album Name in the repository")
				}
			}
		})
	}
}
