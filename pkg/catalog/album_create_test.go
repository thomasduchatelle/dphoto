package catalog_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"testing"
	"time"
)

func TestNewAlbumCreateAcceptance(t *testing.T) {
	const owner = "tonystark"
	createAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.FolderName("/2024-04_Ironman_1"),
		},
		Name:  "Ironman 1",
		Start: time.Date(2024, 04, 28, 8, 33, 42, 0, time.UTC),
		End:   time.Date(2024, 05, 1, 0, 0, 0, 0, time.UTC),
	}
	standardRequest := catalog.CreateAlbumRequest{
		Owner: owner,
		Name:  createAlbum.Name,
		Start: createAlbum.Start,
		End:   createAlbum.End,
	}

	lifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/lifetime"),
		},
		Name:  "lifetime",
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	transferredMedias := catalog.TransferredMedias{
		createAlbum.AlbumId: []catalog.MediaId{"media-1", "media-2"},
	}
	testErrorInsertingAlbum := errors.Errorf("TEST error insering album")

	type fields struct {
		FindAlbumsByOwnerPort    func(t *testing.T) catalog.FindAlbumsByOwnerPort
		InsertAlbumPort          func(t *testing.T) catalog.InsertAlbumPort
		TransferMediasPort       func(t *testing.T) catalog.TransferMediasRepositoryPort
		TimelineMutationObserver func(t *testing.T) catalog.TimelineMutationObserver
	}
	type args struct {
		request catalog.CreateAlbumRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should create a happy path full album create process",
			fields: fields{
				FindAlbumsByOwnerPort:    stubFindAlbumsByOwnerWith(owner, lifetimeAlbum),
				InsertAlbumPort:          expectAlbumInserted(createAlbum),
				TransferMediasPort:       stubTransferMediaPort(transferredMedias),
				TimelineMutationObserver: expectTimelineMutationObserverCalled(transferredMedias),
			},
			args: args{
				request: standardRequest,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not call transfer observer if album insert fails (verify the order)",
			fields: fields{
				FindAlbumsByOwnerPort:    stubFindAlbumsByOwnerWith(owner, lifetimeAlbum),
				InsertAlbumPort:          stubInsertAlbumPortWithError(testErrorInsertingAlbum),
				TransferMediasPort:       stubTransferMediaPort(transferredMedias),
				TimelineMutationObserver: expectTimelineMutationObserverNotCalled(),
			},
			args: args{
				request: standardRequest,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErrorInsertingAlbum)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumCreate := catalog.NewAlbumCreate(
				tt.fields.FindAlbumsByOwnerPort(t),
				tt.fields.InsertAlbumPort(t),
				tt.fields.TransferMediasPort(t),
				tt.fields.TimelineMutationObserver(t),
			)

			_, err := albumCreate.Create(context.Background(), tt.args.request)
			tt.wantErr(t, err)
		})
	}
}

func TestCreateAlbumMediaTransfer_ObserveCreateAlbum(t *testing.T) {
	const owner = "tonystark"
	album := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.FolderName("/2024-04_Ironman_1"),
		},
		Name:  "Ironman 1",
		Start: time.Date(2024, 04, 28, 8, 33, 42, 0, time.UTC),
		End:   time.Date(2024, 05, 1, 0, 0, 0, 0, time.UTC),
	}

	lifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/lifetime"),
		},
		Name:  "lifetime",
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	remainingLifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/remaining-lifetime"),
		},
		Name:  "remaining-lifetime",
		Start: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	highPriorityAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/high-priority"),
		},
		Name:  "remaining-lifetime",
		Start: time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
	}

	type fields struct {
		MediaTransfer         func(t *testing.T) catalog.MediaTransfer
		FindAlbumsByOwnerPort func(t *testing.T) catalog.FindAlbumsByOwnerPort
	}
	type args struct {
		createdAlbum catalog.Album
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the album with a generated name",
			fields: fields{
				FindAlbumsByOwnerPort: stubFindAlbumsByOwnerWith(owner),
				MediaTransfer:         expectMediaTransferCalled(nil),
			},
			args: args{
				createdAlbum: album,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from a lower priority album",
			fields: fields{
				FindAlbumsByOwnerPort: stubFindAlbumsByOwnerWith(owner, lifetimeAlbum),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					album.AlbumId: {
						{
							FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
							Start:      album.Start,
							End:        album.End,
						},
					},
				}),
			},
			args: args{
				createdAlbum: album,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 2 lower priority albums ; selector still in one single block",
			fields: fields{
				FindAlbumsByOwnerPort: stubFindAlbumsByOwnerWith(owner, lifetimeAlbum, remainingLifetimeAlbum),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					album.AlbumId: {
						{
							FromAlbums: []catalog.AlbumId{remainingLifetimeAlbum.AlbumId, lifetimeAlbum.AlbumId},
							Start:      album.Start,
							End:        album.End,
						},
					},
				}),
			},
			args: args{
				createdAlbum: album,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 1 lower priority albums, avoiding 1 high priority (selectors in two blocks)",
			fields: fields{
				FindAlbumsByOwnerPort: stubFindAlbumsByOwnerWith(owner, lifetimeAlbum, highPriorityAlbum),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					album.AlbumId: {
						{
							FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
							Start:      album.Start,
							End:        highPriorityAlbum.Start,
						},
						{
							FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
							Start:      highPriorityAlbum.End,
							End:        album.End,
						},
					},
				}),
			},
			args: args{
				createdAlbum: album,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &catalog.CreateAlbumMediaTransfer{
				MediaTransfer:         tt.fields.MediaTransfer(t),
				FindAlbumsByOwnerPort: tt.fields.FindAlbumsByOwnerPort(t),
			}
			tt.wantErr(t, c.ObserveCreateAlbum(context.Background(), tt.args.createdAlbum), fmt.Sprintf("ObserveCreateAlbum(%v)", tt.args.createdAlbum))
		})
	}
}

func expectMediaTransferCalled(records catalog.MediaTransferRecords) func(t *testing.T) catalog.MediaTransfer {
	return func(t *testing.T) catalog.MediaTransfer {
		transfer := mocks.NewMediaTransfer(t)
		transfer.EXPECT().Transfer(mock.Anything, records).Return(nil).Once()
		return transfer
	}
}

func expectCreateAlbumObserveNotCalled() func(t *testing.T) catalog.CreateAlbumObserver {
	return func(t *testing.T) catalog.CreateAlbumObserver {
		return catalog.CreateAlbumObserverFunc(func(ctx context.Context, album catalog.Album) error {
			assert.Fail(t, "CreateAlbumObserverFunc", "should not be called")
			return nil
		})
	}
}

func stubFindAlbumsByOwnerWith(expectedOwner ownermodel.Owner, albums ...*catalog.Album) func(t *testing.T) catalog.FindAlbumsByOwnerPort {
	return func(t *testing.T) catalog.FindAlbumsByOwnerPort {
		return catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
			if owner == expectedOwner && len(albums) > 0 {
				return albums, nil
			}

			return nil, nil
		})
	}
}

func expectCreateAlbumObserved(album catalog.Album) func(t *testing.T) catalog.CreateAlbumObserver {
	return func(t *testing.T) catalog.CreateAlbumObserver {
		observer := mocks.NewCreateAlbumObserver(t)
		observer.EXPECT().
			ObserveCreateAlbum(mock.Anything, album).
			Return(nil).
			Once()
		return observer
	}
}

func expectAlbumInserted(album catalog.Album) func(t *testing.T) catalog.InsertAlbumPort {
	return func(t *testing.T) catalog.InsertAlbumPort {
		observer := mocks.NewInsertAlbumPort(t)
		observer.EXPECT().
			InsertAlbum(mock.Anything, album).
			Return(nil).
			Once()
		return observer
	}
}

func stubInsertAlbumPortWithError(err error) func(t *testing.T) catalog.InsertAlbumPort {
	return func(t *testing.T) catalog.InsertAlbumPort {
		return catalog.InsertAlbumPortFunc(func(ctx context.Context, album catalog.Album) error {
			return err
		})
	}
}

func expectTimelineMutationObserverNotCalled() func(t *testing.T) catalog.TimelineMutationObserver {
	return func(t *testing.T) catalog.TimelineMutationObserver {
		return catalog.TimelineMutationObserverFunc(func(ctx context.Context, transfers catalog.TransferredMedias) error {
			assert.Failf(t, "TimelineMutationObserverFunc", "should not be called", "Observe(%+v)", transfers)
			return nil
		})
	}
}
