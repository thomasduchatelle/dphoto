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

func TestNewAlbumCreateAcceptance_HappyPath(t *testing.T) {
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

	repository := NewAlbumRepositoryInMemory(lifetimeAlbum)
	transfer := &MediaTransferInMemory{TransferredMedias: transferredMedias}
	observer := &TimelineMutationObserverInMemory{}

	albumCreate := catalog.NewAlbumCreate(repository, repository, transfer, observer)

	_, err := albumCreate.Create(context.Background(), standardRequest)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, &createAlbum, repository.Albums[createAlbum.AlbumId], "album is inserted")
	assert.Equal(t, []catalog.MediaTransferRecords{{
		createAlbum.AlbumId: {
			{
				FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
				Start:      createAlbum.Start,
				End:        createAlbum.End,
			},
		},
	}}, transfer.TransferRecords, "medias are transferred from the lifetime album")
	assert.Equal(t, []catalog.TransferredMedias{{
		Transfers:  transferredMedias.Transfers,
		FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
	}}, observer.Notifications, "timeline mutation observer is notified")
}

// E3: failure-injection test kept with inline testify/mock — event-ordering guarantee that no
// transfer/timeline events are emitted when the album insert fails.
func TestNewAlbumCreateAcceptance_shouldNotCallTransferObserverIfAlbumInsertFails(t *testing.T) {
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
	testErrorInsertingAlbum := errors.New("TEST error inserting album")

	repository := NewAlbumRepositoryInMemory(lifetimeAlbum)
	insertAlbumPort := new(insertAlbumPortMock)
	insertAlbumPort.On("InsertAlbum", mock.Anything, createAlbum).Return(testErrorInsertingAlbum).Once()
	defer insertAlbumPort.AssertExpectations(t)

	transfer := &MediaTransferInMemory{}
	observer := &TimelineMutationObserverInMemory{}

	albumCreate := catalog.NewAlbumCreate(repository, insertAlbumPort, transfer, observer)

	_, err := albumCreate.Create(context.Background(), standardRequest)
	assert.ErrorIs(t, err, testErrorInsertingAlbum)
	assert.Empty(t, transfer.TransferRecords, "no media transfer should have been attempted")
	assert.Empty(t, observer.Notifications, "no timeline mutation observer should have been notified")
}

// E4: failure-injection test kept with inline testify/mock — ordering guarantee that the
// existing albums are listed before the new album is inserted.
func TestNewAlbumCreateAcceptance_shouldListExistingAlbumsBeforeCreatingTheNewOne(t *testing.T) {
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
	testErrorFindingAlbums := errors.New("TEST error finding albums")

	findAlbumsByOwner := new(findAlbumsByOwnerPortMock)
	findAlbumsByOwner.On("FindAlbumsByOwner", mock.Anything, ownermodel.Owner(owner)).Return(([]*catalog.Album)(nil), testErrorFindingAlbums).Once()
	defer findAlbumsByOwner.AssertExpectations(t)

	repository := NewAlbumRepositoryInMemory()
	transfer := &MediaTransferInMemory{}
	observer := &TimelineMutationObserverInMemory{}

	albumCreate := catalog.NewAlbumCreate(findAlbumsByOwner, repository, transfer, observer)

	_, err := albumCreate.Create(context.Background(), standardRequest)
	assert.ErrorIs(t, err, testErrorFindingAlbums)
	assert.Empty(t, repository.Albums, "no album should have been inserted")
	assert.Empty(t, transfer.TransferRecords, "no media transfer should have been attempted")
	assert.Empty(t, observer.Notifications, "no timeline mutation observer should have been notified")
}

type insertAlbumPortMock struct {
	mock.Mock
}

func (m *insertAlbumPortMock) InsertAlbum(ctx context.Context, album catalog.Album) error {
	args := m.Called(ctx, album)
	return args.Error(0)
}

type findAlbumsByOwnerPortMock struct {
	mock.Mock
}

func (m *findAlbumsByOwnerPortMock) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	args := m.Called(ctx, owner)
	albums, _ := args.Get(0).([]*catalog.Album)
	return albums, args.Error(1)
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
			name: "it should create the album with a generated name when the folderName is just a slash '/'",
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            standardRequest.Owner,
					Name:             standardRequest.Name,
					Start:            standardRequest.Start,
					End:              standardRequest.End,
					ForcedFolderName: "/",
				},
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
			observer := new(CreateAlbumObserverInMemory)
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

	type args struct {
		createdAlbum   catalog.Album
		existingAlbums []*catalog.Album
	}
	tests := []struct {
		name        string
		args        args
		wantRecords []catalog.MediaTransferRecords
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the album with a generated name",
			args: args{
				createdAlbum: album,
			},
			wantRecords: []catalog.MediaTransferRecords{nil},
			wantErr:     assert.NoError,
		},
		{
			name: "it should re-allocate medias from a lower priority album",
			args: args{
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum},
			},
			wantRecords: []catalog.MediaTransferRecords{{
				album.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
						Start:      album.Start,
						End:        album.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 2 lower priority albums ; selector still in one single block",
			args: args{
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum, remainingLifetimeAlbum},
			},
			wantRecords: []catalog.MediaTransferRecords{{
				album.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{remainingLifetimeAlbum.AlbumId, lifetimeAlbum.AlbumId},
						Start:      album.Start,
						End:        album.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 1 lower priority albums, avoiding 1 high priority (selectors in two blocks)",
			args: args{
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum, highPriorityAlbum},
			},
			wantRecords: []catalog.MediaTransferRecords{{
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
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := &MediaTransferInMemory{}
			c := &catalog.CreateAlbumMediaTransfer{
				MediaTransfer: transfer,
			}
			err := c.ObserveCreateAlbum(context.Background(), catalog.NewLazyTimelineAggregate(tt.args.existingAlbums), tt.args.createdAlbum)
			if !tt.wantErr(t, err, fmt.Sprintf("ObserveCreateAlbum(%v)", tt.args.createdAlbum)) {
				return
			}
			assert.Equal(t, tt.wantRecords, transfer.TransferRecords)
		})
	}
}
