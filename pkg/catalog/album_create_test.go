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
	testErrorInsertingAlbum := errors.Errorf("TEST error inserting album")
	testErrorFindingAlbums := errors.New("TEST error finding albums")

	transferMediasWithMedias := func() *TransferMediasInMemory {
		transfer := NewTransferMediasInMemory()
		transfer.TransferredMedias = transferredMedias
		return transfer
	}

	t.Run("it should create a happy path full album create process", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(lifetimeAlbum)
		transferMedias := transferMediasWithMedias()
		timelineObserver := &TimelineMutationObserverInMemory{}

		albumCreate := catalog.NewAlbumCreate(
			albumRepository,
			albumRepository,
			transferMedias,
			timelineObserver,
		)

		_, err := albumCreate.Create(context.Background(), standardRequest)
		if !assert.NoError(t, err) {
			return
		}

		assert.Contains(t, albumRepository.Albums, createAlbum.AlbumId, "album should be inserted")
		assert.Equal(t, []catalog.TransferredMedias{{
			Transfers:  transferredMedias.Transfers,
			FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
		}}, timelineObserver.Notifications)
	})

	// E3: uses testify/mock inline for the InsertAlbumPort to force the insert to fail and
	// verify by state that neither the transfer nor the observer was called (order guarantee).
	t.Run("it should not call transfer observer if album insert fails (verify the order)", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(lifetimeAlbum)
		failingInsert := new(insertAlbumPortMock)
		failingInsert.On("InsertAlbum", mock.Anything, createAlbum).Return(testErrorInsertingAlbum).Once()
		transferMedias := transferMediasWithMedias()
		timelineObserver := &TimelineMutationObserverInMemory{}

		albumCreate := catalog.NewAlbumCreate(
			albumRepository,
			failingInsert,
			transferMedias,
			timelineObserver,
		)

		_, err := albumCreate.Create(context.Background(), standardRequest)
		assert.ErrorIs(t, err, testErrorInsertingAlbum)
		failingInsert.AssertExpectations(t)
		assert.Empty(t, transferMedias.Records, "transfer must not run when insert fails")
		assert.Empty(t, timelineObserver.Notifications, "observer must not be notified when insert fails")
	})

	// E4: uses testify/mock inline for FindAlbumsByOwnerPort to force the list to fail and
	// verify by state that no album is inserted (ordering guarantee: list before insert).
	t.Run("it should list the existing albums before creating the new one (otherwise there are duplicates in the timeline)", func(t *testing.T) {
		findAlbumsByOwner := new(findAlbumsByOwnerMock)
		findAlbumsByOwner.On("FindAlbumsByOwner", mock.Anything, ownermodel.Owner(owner)).Return([]*catalog.Album(nil), testErrorFindingAlbums).Once()

		albumRepository := NewAlbumRepositoryInMemory(lifetimeAlbum)
		transferMedias := transferMediasWithMedias()
		timelineObserver := &TimelineMutationObserverInMemory{}

		albumCreate := catalog.NewAlbumCreate(
			findAlbumsByOwner,
			albumRepository,
			transferMedias,
			timelineObserver,
		)

		_, err := albumCreate.Create(context.Background(), standardRequest)
		assert.ErrorIs(t, err, testErrorFindingAlbums)
		findAlbumsByOwner.AssertExpectations(t)
		assert.NotContains(t, albumRepository.Albums, createAlbum.AlbumId, "album must not be inserted when list fails")
		assert.Empty(t, transferMedias.Records)
		assert.Empty(t, timelineObserver.Notifications)
	})
}

// insertAlbumPortMock is a testify/mock used only inline by the failure-injection tests (E3
// in this file, E6 in album_rename_test.go) to force InsertAlbum to return an error. Every
// other test uses the AlbumRepositoryInMemory Fake.
type insertAlbumPortMock struct {
	mock.Mock
}

func (m *insertAlbumPortMock) InsertAlbum(ctx context.Context, album catalog.Album) error {
	return m.Called(ctx, album).Error(0)
}

// findAlbumsByOwnerMock is a testify/mock used only inline by E4 to force FindAlbumsByOwner
// to return an error. Every other test in this file uses the AlbumRepositoryInMemory Fake.
type findAlbumsByOwnerMock struct {
	mock.Mock
}

func (m *findAlbumsByOwnerMock) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
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
			observer := &CreateAlbumObserverInMemory{}
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
		name                string
		args                args
		wantTransferRecords []catalog.MediaTransferRecords
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the album with a generated name",
			args: args{
				createdAlbum: album,
			},
			wantTransferRecords: []catalog.MediaTransferRecords{nil},
			wantErr:             assert.NoError,
		},
		{
			name: "it should re-allocate medias from a lower priority album",
			args: args{
				createdAlbum:   album,
				existingAlbums: []*catalog.Album{lifetimeAlbum},
			},
			wantTransferRecords: []catalog.MediaTransferRecords{{
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
			wantTransferRecords: []catalog.MediaTransferRecords{{
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
			wantTransferRecords: []catalog.MediaTransferRecords{{
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
			mediaTransfer := &MediaTransferInMemory{}
			c := &catalog.CreateAlbumMediaTransfer{
				MediaTransfer: mediaTransfer,
			}
			err := c.ObserveCreateAlbum(context.Background(), catalog.NewLazyTimelineAggregate(tt.args.existingAlbums), tt.args.createdAlbum)
			if !tt.wantErr(t, err, fmt.Sprintf("ObserveCreateAlbum(%v)", tt.args.createdAlbum)) {
				return
			}
			assert.Equal(t, tt.wantTransferRecords, mediaTransfer.Records)
		})
	}
}
