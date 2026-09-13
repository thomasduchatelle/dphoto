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

func TestNewRenameAlbumAcceptance_HappyPath(t *testing.T) {
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

	repository := NewAlbumRepositoryInMemory(existingAlbum)
	transfer := &MediaTransferInMemory{TransferredMedias: transferredMedias}
	observer := &TimelineMutationObserverInMemory{}

	renameAlbum := catalog.NewRenameAlbum(
		repository,
		repository,
		repository,
		repository,
		transfer,
		repository,
		observer,
	)

	err := renameAlbum.RenameAlbum(context.Background(), catalog.RenameAlbumRequest{
		CurrentId:        existingAlbum.AlbumId,
		NewName:          newName,
		RenameFolder:     true,
		ForcedFolderName: "",
	})
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, &newAlbum, repository.Albums[newAlbum.AlbumId], "new album is inserted")
	_, oldStillExists := repository.Albums[existingAlbum.AlbumId]
	assert.False(t, oldStillExists, "old album is removed from the repository")
	assert.Equal(t, []catalog.MediaTransferRecords{{
		newAlbum.AlbumId: {
			{
				FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
				Start:      existingAlbum.Start,
				End:        existingAlbum.End,
			},
		},
	}}, transfer.TransferRecords, "medias are transferred from the old album to the new album")
	assert.Equal(t, []catalog.TransferredMedias{{
		Transfers:  transferredMedias.Transfers,
		FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
	}}, observer.Notifications, "timeline mutation observer is notified")
}

// E6: failure-injection test kept with inline testify/mock — rollback semantics ensuring the
// old album is not deleted and no medias are transferred when insertion of the new album fails.
func TestNewRenameAlbumAcceptance_shouldInterruptTransferIfAlbumInsertionFails(t *testing.T) {
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

	testError := errors.New("TEST error inserting album")

	repository := NewAlbumRepositoryInMemory(existingAlbum)
	insertAlbumPort := new(insertAlbumPortMock)
	insertAlbumPort.On("InsertAlbum", mock.Anything, newAlbum).Return(testError).Once()
	defer insertAlbumPort.AssertExpectations(t)

	transfer := &MediaTransferInMemory{}
	observer := &TimelineMutationObserverInMemory{}

	renameAlbum := catalog.NewRenameAlbum(
		repository,
		repository,
		insertAlbumPort,
		repository,
		transfer,
		repository,
		observer,
	)

	err := renameAlbum.RenameAlbum(context.Background(), catalog.RenameAlbumRequest{
		CurrentId:        existingAlbum.AlbumId,
		NewName:          newName,
		RenameFolder:     true,
		ForcedFolderName: "",
	})
	assert.ErrorIs(t, err, testError)
	assert.Contains(t, repository.Albums, existingAlbum.AlbumId, "old album should still exist (not deleted)")
	assert.Empty(t, transfer.TransferRecords, "no medias should have been transferred")
	assert.Empty(t, observer.Notifications, "no timeline mutation observer should have been notified")
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

	type fields struct {
		albums []*catalog.Album
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantAlbumName        string // expected name of existingAlbum after the call (only meaningful when it should still exist and be renamed in-place)
		wantRenamedFrom      []catalog.AlbumId
		wantCreationRequests []catalog.CreateAlbumRequest
		wantErr              assert.ErrorAssertionFunc
	}{
		{
			name:   "it should get an error if the new name is empty",
			fields: fields{albums: []*catalog.Album{existingAlbum}},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          "",
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantAlbumName: existingAlbum.Name,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name:   "it should get an error if the album doesn't exists",
			fields: fields{},
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
			fields: fields{albums: []*catalog.Album{existingAlbum}},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantAlbumName: newName,
			wantErr:       assert.NoError,
		},
		{
			name:   "it should create a new album if the album is found and folder name is changed",
			fields: fields{albums: []*catalog.Album{existingAlbum}},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			wantAlbumName:   existingAlbum.Name,
			wantRenamedFrom: []catalog.AlbumId{existingAlbum.AlbumId},
			wantCreationRequests: []catalog.CreateAlbumRequest{{
				Owner:            owner,
				Name:             newName,
				Start:            may24,
				End:              jun24,
				ForcedFolderName: "",
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should create a new album if the album is found and folder name is forced to a certain value",
			fields: fields{albums: []*catalog.Album{existingAlbum}},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "Avengers vs Loki",
				},
			},
			wantAlbumName:   existingAlbum.Name,
			wantRenamedFrom: []catalog.AlbumId{existingAlbum.AlbumId},
			wantCreationRequests: []catalog.CreateAlbumRequest{{
				Owner:            owner,
				Name:             newName,
				Start:            may24,
				End:              jun24,
				ForcedFolderName: "Avengers vs Loki",
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumsCopy := make([]*catalog.Album, len(tt.fields.albums))
			for i, a := range tt.fields.albums {
				copy := *a
				albumsCopy[i] = &copy
			}
			repository := NewAlbumRepositoryInMemory(albumsCopy...)
			observer := &RenameAlbumObserverInMemory{}

			r := &catalog.RenameAlbum{
				FindAlbumById:        repository,
				UpdateAlbumName:      repository,
				RenameAlbumObservers: []catalog.RenameAlbumObserver{observer},
			}
			err := r.RenameAlbum(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v)", tt.args.request)) {
				return
			}

			assert.Equal(t, tt.wantRenamedFrom, observer.RenamedFrom, "renamed-from observer records")
			assert.Equal(t, tt.wantCreationRequests, observer.CreationRequests, "creation-request observer records")

			if len(tt.fields.albums) > 0 && tt.wantAlbumName != "" {
				assert.Equal(t, tt.wantAlbumName, repository.Albums[existingAlbum.AlbumId].Name, "album name in the repository")
			}
		})
	}
}
