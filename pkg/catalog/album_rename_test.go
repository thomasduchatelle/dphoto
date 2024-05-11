package catalog_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
	"time"
)

func TestNewRenameAlbumAcceptance(t *testing.T) {
	const owner = "ironman"
	may24 := time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)
	newName := "Avenger 1"

	existingAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      catalog.Owner(owner),
			FolderName: catalog.NewFolderName("/avenger"),
		},
		Name:  "Avenger",
		Start: may24,
		End:   jun24,
	}
	newAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      catalog.Owner(owner),
			FolderName: catalog.NewFolderName("/2024-05_Avenger_1"),
		},
		Name:  newName,
		Start: may24,
		End:   jun24,
	}
	transferredMedias := catalog.TransferredMedias{
		newAlbum.AlbumId: []catalog.MediaId{"media-1", "media-2"},
	}

	testError := errors.Errorf("TEST error throwing")

	type fields struct {
		FindAlbumById             func(t *testing.T) catalog.FindAlbumByIdPort
		UpdateAlbumName           func(t *testing.T) catalog.UpdateAlbumNamePort
		InsertAlbumPort           func(t *testing.T) catalog.InsertAlbumPort
		DeleteAlbumRepositoryPort func(t *testing.T) catalog.DeleteAlbumRepositoryPort
		TransferMedias            func(t *testing.T) catalog.TransferMediasRepositoryPort
		TimelineMutationObservers func(t *testing.T) catalog.TimelineMutationObserver
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should create a new album end to end",
			fields: fields{
				FindAlbumById:             stubFindAlbumByIdWith(existingAlbum),
				UpdateAlbumName:           expectUpdateAlbumNameNotCalled(),
				InsertAlbumPort:           expectAlbumInserted(newAlbum),
				DeleteAlbumRepositoryPort: expectDeleteAlbumRepositoryPortCalled(existingAlbum.AlbumId),
				TransferMedias: expectTransferMediasRepositoryPortCalled(catalog.MediaTransferRecords{
					newAlbum.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
							Start:      existingAlbum.Start,
							End:        existingAlbum.End,
						},
					},
				}, transferredMedias),
				TimelineMutationObservers: expectTimelineMutationObserverCalled(transferredMedias),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avenger")},
					NewName:          "Avenger 1",
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should interrupt the transfer if the album insertion fails",
			fields: fields{
				FindAlbumById:             stubFindAlbumByIdWith(existingAlbum),
				UpdateAlbumName:           expectUpdateAlbumNameNotCalled(),
				InsertAlbumPort:           stubInsertAlbumPortWithError(testError),
				DeleteAlbumRepositoryPort: expectDeleteAlbumRepositoryPortNotCalled(),
				TransferMedias:            expectTransferMediasPortNotCalled(),
				TimelineMutationObservers: expectTimelineMutationObserverNotCalled(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avenger")},
					NewName:          "Avenger 1",
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renameAlbum := catalog.NewRenameAlbum(
				tt.fields.FindAlbumById(t),
				tt.fields.UpdateAlbumName(t),
				tt.fields.InsertAlbumPort(t),
				tt.fields.DeleteAlbumRepositoryPort(t),
				tt.fields.TransferMedias(t),
				tt.fields.TimelineMutationObservers(t),
			)

			err := renameAlbum.RenameAlbum(context.Background(), tt.args.request)
			tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v)", tt.args.request))
		})
	}
}

func expectTransferMediasPortNotCalled() func(t *testing.T) catalog.TransferMediasRepositoryPort {
	return func(t *testing.T) catalog.TransferMediasRepositoryPort {
		return catalog.TransferMediasFunc(func(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
			assert.Failf(t, "TransferMediasRepository should not be called", "TransferMediasFunc(%v, %v)", ctx, records)
			return nil, nil
		})
	}
}

func expectMediaTransferNotCalled() func(t *testing.T) catalog.MediaTransfer {
	return func(t *testing.T) catalog.MediaTransfer {
		return catalog.MediaTransferFunc(func(ctx context.Context, records catalog.MediaTransferRecords) error {
			assert.Failf(t, "MediaTransfer should not be called", "MediaTransfer(%v, %v)", ctx, records)
			return nil
		})
	}
}

func TestRenameAlbum_RenameAlbum(t *testing.T) {
	const owner = "ironman"
	may24 := time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	existingAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      catalog.Owner(owner),
			FolderName: catalog.NewFolderName("/avenger"),
		},
		Name:  "Avenger",
		Start: may24,
		End:   jun24,
	}
	newName := "Avenger 1"

	type fields struct {
		FindAlbumById       func(t *testing.T) catalog.FindAlbumByIdPort
		UpdateAlbumName     func(t *testing.T) catalog.UpdateAlbumNamePort
		RenameAlbumObserver func(t *testing.T) catalog.RenameAlbumObserver
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should get an error if the new name is empty",
			fields: fields{
				FindAlbumById:       stubFindAlbumByIdWith(existingAlbum),
				UpdateAlbumName:     expectUpdateAlbumNameNotCalled(),
				RenameAlbumObserver: expectRenameAlbumObserverNotCalled(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          "",
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name: "it should get an error if the album doesn't exists",
			fields: fields{
				FindAlbumById:       stubFindAlbumByIdWith(&catalog.Album{}),
				UpdateAlbumName:     expectUpdateAlbumNameNotCalled(),
				RenameAlbumObserver: expectRenameAlbumObserverNotCalled(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundError)
			},
		},
		{
			name: "it should update the name if the album is found and folder name is unchanged",
			fields: fields{
				FindAlbumById:       stubFindAlbumByIdWith(existingAlbum),
				UpdateAlbumName:     expectUpdateAlbumNameCalledWith(existingAlbum.AlbumId, newName),
				RenameAlbumObserver: expectRenameAlbumObserverNotCalled(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should create a new album if the album is found and folder name is changed",
			fields: fields{
				FindAlbumById:   stubFindAlbumByIdWith(existingAlbum),
				UpdateAlbumName: expectUpdateAlbumNameNotCalled(),
				RenameAlbumObserver: expectRenameAlbumObserverCalledWith(existingAlbum.AlbumId, catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             newName,
					Start:            may24,
					End:              jun24,
					ForcedFolderName: "",
				}),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     true,
					ForcedFolderName: "",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should create a new album if the album is found and folder name is forced to a certain value",
			fields: fields{
				FindAlbumById:   stubFindAlbumByIdWith(existingAlbum),
				UpdateAlbumName: expectUpdateAlbumNameNotCalled(),
				RenameAlbumObserver: expectRenameAlbumObserverCalledWith(existingAlbum.AlbumId, catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             newName,
					Start:            may24,
					End:              jun24,
					ForcedFolderName: "Avengers vs Loki",
				}),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "Avengers vs Loki",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &catalog.RenameAlbum{
				FindAlbumById:        tt.fields.FindAlbumById(t),
				UpdateAlbumName:      tt.fields.UpdateAlbumName(t),
				RenameAlbumObservers: []catalog.RenameAlbumObserver{tt.fields.RenameAlbumObserver(t)},
			}
			tt.wantErr(t, r.RenameAlbum(context.Background(), tt.args.request), fmt.Sprintf("RenameAlbum(%v)", tt.args.request))
		})
	}
}

func expectRenameAlbumObserverNotCalled() func(t *testing.T) catalog.RenameAlbumObserver {
	return func(t *testing.T) catalog.RenameAlbumObserver {
		return mocks.NewRenameAlbumObserver(t)
	}
}

func expectRenameAlbumObserverCalledWith(currentId catalog.AlbumId, request catalog.CreateAlbumRequest) func(t *testing.T) catalog.RenameAlbumObserver {
	return func(t *testing.T) catalog.RenameAlbumObserver {
		observer := mocks.NewRenameAlbumObserver(t)
		observer.EXPECT().OnRenameAlbum(mock.Anything, currentId, request).Return(nil).Once()
		return observer
	}
}

func expectUpdateAlbumNameCalledWith(albumId catalog.AlbumId, newName string) func(t *testing.T) catalog.UpdateAlbumNamePort {
	return func(t *testing.T) catalog.UpdateAlbumNamePort {
		port := mocks.NewUpdateAlbumNamePort(t)
		port.EXPECT().UpdateAlbumName(mock.Anything, albumId, newName).Return(nil).Once()
		return port
	}
}

func expectUpdateAlbumNameNotCalled() func(t *testing.T) catalog.UpdateAlbumNamePort {
	return func(t *testing.T) catalog.UpdateAlbumNamePort {
		return mocks.NewUpdateAlbumNamePort(t)
	}
}

func stubFindAlbumByIdWith(existingAlbum *catalog.Album) func(t *testing.T) catalog.FindAlbumByIdPort {
	return func(t *testing.T) catalog.FindAlbumByIdPort {
		return catalog.FindAlbumByIdFunc(func(ctx context.Context, id catalog.AlbumId) (*catalog.Album, error) {
			if existingAlbum.AlbumId.IsEqual(id) {
				return existingAlbum, nil
			}
			return nil, catalog.AlbumNotFoundError
		})
	}
}

func expectDeleteAlbumRepositoryPortCalled(id catalog.AlbumId) func(t *testing.T) catalog.DeleteAlbumRepositoryPort {
	return func(t *testing.T) catalog.DeleteAlbumRepositoryPort {
		port := mocks.NewDeleteAlbumRepositoryPort(t)
		port.EXPECT().DeleteAlbum(mock.Anything, id).Return(nil).Once()
		return port
	}
}
func expectDeleteAlbumRepositoryPortNotCalled() func(t *testing.T) catalog.DeleteAlbumRepositoryPort {
	return func(t *testing.T) catalog.DeleteAlbumRepositoryPort {
		return catalog.DeleteAlbumRepositoryFunc(func(ctx context.Context, id catalog.AlbumId) error {
			t.Error("DeleteAlbumRepositoryPort should not be called")
			return nil
		})
	}
}
