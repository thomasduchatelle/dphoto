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

func TestCreateAlbum_Create(t *testing.T) {
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
	standardRequest := catalog.CreateAlbumRequest{
		Owner: owner,
		Name:  album.Name,
		Start: album.Start,
		End:   album.End,
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
		FindAlbumsByOwnerPort func(t *testing.T) catalog.FindAlbumsByOwnerPort
		Observer              func(t *testing.T) catalog.CreateAlbumObserver
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
			name: "it should NOT create the album without owner",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer:              expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: "",
					Name:  "foobar",
					Start: album.Start,
					End:   album.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.EmptyOwnerError)
			},
		},
		{
			name: "it should NOT create the album without name",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer:              expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "",
					Start: album.Start,
					End:   album.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album without start date",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer:              expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					End:   album.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumStartAndEndDateMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album without end date",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer:              expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					Start: album.Start,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumStartAndEndDateMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album with start and end reversed",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer:              expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					Start: album.End,
					End:   album.Start,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumEndDateMustBeAfterStartErr)
			},
		},
		{
			name: "it should create the album with a generated name",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer:              expectCreateAlbumObserved(album, nil),
			},
			args: args{
				request: standardRequest,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should create the album with a forced name",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				Observer: expectCreateAlbumObserved(catalog.Album{
					AlbumId: catalog.AlbumId{
						Owner:      owner,
						FolderName: "/Phase_1_Avenger",
					},
					Name:  "Avenger 1",
					Start: album.Start,
					End:   album.End,
				}, nil),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             "Avenger 1",
					Start:            album.Start,
					End:              album.End,
					ForcedFolderName: "Phase_1_Avenger",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from a lower priority album",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner, lifetimeAlbum),
				Observer: expectCreateAlbumObserved(album, catalog.MediaTransferRecords{
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
				request: standardRequest,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 2 lower priority albums ; selector still in one single block",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner, lifetimeAlbum, remainingLifetimeAlbum),
				Observer: expectCreateAlbumObserved(album, catalog.MediaTransferRecords{
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
				request: standardRequest,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 1 lower priority albums, avoiding 1 high priority (selectors in two blocks)",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner, lifetimeAlbum, highPriorityAlbum),
				Observer: expectCreateAlbumObserved(album, catalog.MediaTransferRecords{
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
				request: standardRequest,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var observers []catalog.CreateAlbumObserver
			{
			}
			if tt.fields.Observer != nil {
				observers = append(observers, tt.fields.Observer(t))
			}
			c := &catalog.CreateAlbum{
				FindAlbumsByOwnerPort: tt.fields.FindAlbumsByOwnerPort(t),
				Observers:             observers,
			}

			err := c.Create(context.TODO(), tt.args.request)
			tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.args.request))
		})
	}
}

func expectCreateAlbumObserveNotCalled() func(t *testing.T) catalog.CreateAlbumObserver {
	return func(t *testing.T) catalog.CreateAlbumObserver {
		return catalog.CreateAlbumObserverFunc(func(ctx context.Context, album catalog.Album, records catalog.MediaTransferRecords) error {
			assert.Fail(t, "CreateAlbumObserverFunc", "should not be called")
			return nil
		})
	}
}

func returnListOfAlbums(expectedOwner catalog.Owner, albums ...*catalog.Album) func(t *testing.T) catalog.FindAlbumsByOwnerPort {
	return func(t *testing.T) catalog.FindAlbumsByOwnerPort {
		return catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner catalog.Owner) ([]*catalog.Album, error) {
			if owner == expectedOwner && len(albums) > 0 {
				return albums, nil
			}

			return nil, nil
		})
	}
}

func expectCreateAlbumObserved(album catalog.Album, records catalog.MediaTransferRecords) func(t *testing.T) catalog.CreateAlbumObserver {
	return func(t *testing.T) catalog.CreateAlbumObserver {
		observer := mocks.NewCreateAlbumObserver(t)
		observer.EXPECT().
			ObserveCreateAlbum(mock.Anything, album, records).
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

func TestNewAlbumCreate(t *testing.T) {
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
		TransferMediasPort       func(t *testing.T) catalog.TransferMediasPort
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
				FindAlbumsByOwnerPort:    returnListOfAlbums(owner, lifetimeAlbum),
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
				FindAlbumsByOwnerPort:    returnListOfAlbums(owner, lifetimeAlbum),
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

			err := albumCreate.Create(context.Background(), tt.args.request)
			tt.wantErr(t, err)
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
