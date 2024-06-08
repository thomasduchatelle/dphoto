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
		Transfers: map[catalog.AlbumId][]catalog.MediaId{
			createAlbum.AlbumId: {"media-1", "media-2"},
		},
	}
	testErrorInsertingAlbum := errors.Errorf("TEST error insering album")
	testErrorFindingAlbums := errors.New("TEST error finding albums")

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
				FindAlbumsByOwnerPort: stubFindAlbumsByOwnerWith(owner, lifetimeAlbum),
				InsertAlbumPort:       expectAlbumInserted(createAlbum),
				TransferMediasPort:    stubTransferMediaPort(transferredMedias),
				TimelineMutationObserver: expectTimelineMutationObserverCalled(catalog.TransferredMedias{
					Transfers:  transferredMedias.Transfers,
					FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
				}),
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
		{
			name: "it should list the existing albums before creating the new one (otherwise there are duplicates in the timeline)",
			fields: fields{
				FindAlbumsByOwnerPort:    stubFindAlbumsByOwnerPortWithError(testErrorFindingAlbums),
				InsertAlbumPort:          stubInsertAlbumPortWithError(testErrorInsertingAlbum),
				TransferMediasPort:       stubTransferMediaPort(transferredMedias),
				TimelineMutationObserver: expectTimelineMutationObserverNotCalled(),
			},
			args: args{
				request: standardRequest,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErrorFindingAlbums)
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

func TestCreateAlbumStateless_Create(t *testing.T) {
	const owner = "tonystark"
	ironmanOneAlbum := catalog.Album{
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
		Name:  ironmanOneAlbum.Name,
		Start: ironmanOneAlbum.Start,
		End:   ironmanOneAlbum.End,
	}

	type fields struct {
		Albums []*catalog.Album
	}
	type args struct {
		request catalog.CreateAlbumRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantObserved []catalog.Album
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "it should NOT create the album without owner",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: "",
					Name:  "foobar",
					Start: ironmanOneAlbum.Start,
					End:   ironmanOneAlbum.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ownermodel.EmptyOwnerError)
			},
		},
		{
			name: "it should NOT create the album without name",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "",
					Start: ironmanOneAlbum.Start,
					End:   ironmanOneAlbum.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album without start date",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					End:   ironmanOneAlbum.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumStartAndEndDateMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album without end date",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					Start: ironmanOneAlbum.Start,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumStartAndEndDateMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album with start and end reversed",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					Start: ironmanOneAlbum.End,
					End:   ironmanOneAlbum.Start,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumEndDateMustBeAfterStartErr)
			},
		},
		{
			name: "it should create the album with a generated name",
			args: args{
				request: standardRequest,
			},
			wantObserved: []catalog.Album{ironmanOneAlbum},
			wantErr:      assert.NoError,
		},
		{
			name: "it should create the album with a forced name",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             "Avenger 1",
					Start:            ironmanOneAlbum.Start,
					End:              ironmanOneAlbum.End,
					ForcedFolderName: "Phase_1_Avenger",
				},
			},
			wantObserved: []catalog.Album{{
				AlbumId: catalog.AlbumId{
					Owner:      owner,
					FolderName: "/Phase_1_Avenger",
				},
				Name:  "Avenger 1",
				Start: ironmanOneAlbum.Start,
				End:   ironmanOneAlbum.End,
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should NOT create the album if the forced name is already taken",
			fields: fields{
				Albums: []*catalog.Album{&ironmanOneAlbum},
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             "A different name",
					Start:            ironmanOneAlbum.Start,
					End:              ironmanOneAlbum.End,
					ForcedFolderName: ironmanOneAlbum.AlbumId.FolderName.String(),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumFolderNameAlreadyTakenErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := new(CreateAlbumObserverFake)
			c := &catalog.CreateAlbumStateless{
				Observers: []catalog.CreateAlbumObserverWithTimeline{&catalog.CreateAlbumObserverWrapper{CreateAlbumObserver: observer}},
			}

			_, err := c.Create(context.Background(), catalog.NewLazyTimelineAggregate(tt.fields.Albums), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.args.request)) {
				return
			}

			assert.Equal(t, tt.wantObserved, observer.CreatedAlbums)
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
		MediaTransfer func(t *testing.T) catalog.MediaTransfer
	}
	type args struct {
		createdAlbum   catalog.Album
		existingAlbums []*catalog.Album
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
				MediaTransfer: expectMediaTransferCalled(nil),
			},
			args: args{
				createdAlbum: album,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from a lower priority album",
			fields: fields{
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
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 2 lower priority albums ; selector still in one single block",
			fields: fields{
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
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum, remainingLifetimeAlbum},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 1 lower priority albums, avoiding 1 high priority (selectors in two blocks)",
			fields: fields{
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
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum, highPriorityAlbum},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &catalog.CreateAlbumMediaTransfer{
				MediaTransfer: tt.fields.MediaTransfer(t),
			}
			tt.wantErr(t, c.ObserveCreateAlbum(context.Background(), catalog.NewLazyTimelineAggregate(tt.args.existingAlbums), tt.args.createdAlbum), fmt.Sprintf("ObserveCreateAlbum(%v)", tt.args.createdAlbum))
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

type FindAlbumsByOwnerFake map[ownermodel.Owner][]*catalog.Album

func (f FindAlbumsByOwnerFake) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	albums, _ := f[owner]
	return albums, nil
}

type CreateAlbumObserverFake struct {
	CreatedAlbums []catalog.Album
}

func (c *CreateAlbumObserverFake) ObserveCreateAlbum(ctx context.Context, createdAlbum catalog.Album) error {
	c.CreatedAlbums = append(c.CreatedAlbums, createdAlbum)
	return nil
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
			assert.Failf(t, "TimelineMutationObserverFunc", "should not be called", "OnTransferredMedias(%+v)", transfers)
			return nil
		})
	}
}
