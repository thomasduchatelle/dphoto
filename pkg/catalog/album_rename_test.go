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

	t.Run("it should create a new album end to end", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(existingAlbum)
		transferMedias := NewTransferMediasInMemory()
		transferMedias.TransferredMedias = transferredMedias
		timelineObserver := &TimelineMutationObserverInMemory{}

		renameAlbum := catalog.NewRenameAlbum(
			albumRepository,
			albumRepository,
			albumRepository,
			albumRepository,
			transferMedias,
			albumRepository,
			timelineObserver,
		)

		err := renameAlbum.RenameAlbum(context.Background(), renameRequest)
		if !assert.NoError(t, err) {
			return
		}

		assert.NotContains(t, albumRepository.Albums, existingAlbum.AlbumId, "the old album should be deleted")
		assert.Contains(t, albumRepository.Albums, newAlbum.AlbumId, "the new album should be created")
		assert.Equal(t, []catalog.MediaTransferRecords{{
			newAlbum.AlbumId: {
				{
					FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
					Start:      existingAlbum.Start,
					End:        existingAlbum.End,
				},
			},
		}}, transferMedias.Records)
		assert.Equal(t, []catalog.TransferredMedias{{
			Transfers:  transferredMedias.Transfers,
			FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
		}}, timelineObserver.Notifications)
	})

	// E6: uses testify/mock inline for the InsertAlbumPort to force insert to fail. Verify by
	// state that the old album is still present, the transfer never happened, and the observer
	// was not notified (rollback semantics).
	t.Run("it should interrupt the transfer if the album insertion fails", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(existingAlbum)
		failingInsert := new(insertAlbumPortMock)
		failingInsert.On("InsertAlbum", mock.Anything, newAlbum).Return(testError).Once()
		transferMedias := NewTransferMediasInMemory()
		transferMedias.TransferredMedias = transferredMedias
		timelineObserver := &TimelineMutationObserverInMemory{}

		renameAlbum := catalog.NewRenameAlbum(
			albumRepository,
			albumRepository,
			failingInsert,
			albumRepository,
			transferMedias,
			albumRepository,
			timelineObserver,
		)

		err := renameAlbum.RenameAlbum(context.Background(), renameRequest)
		assert.ErrorIs(t, err, testError)
		failingInsert.AssertExpectations(t)
		assert.Contains(t, albumRepository.Albums, existingAlbum.AlbumId, "the old album must NOT be deleted when insert fails")
		assert.Empty(t, transferMedias.Records, "medias must NOT be transferred when insert fails")
		assert.Empty(t, timelineObserver.Notifications, "observer must NOT be notified when insert fails")
	})
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
