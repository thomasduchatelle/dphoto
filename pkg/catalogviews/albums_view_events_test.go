package catalogviews

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

var (
	tonyOwner       = ownermodel.Owner("tony")
	ownerUserId     = usermodel.UserId("ironman@avenger.hero")
	visitorUserId   = usermodel.UserId("pepper@stark.com")
	visitor2UserId  = usermodel.UserId("wanda@avenger.hero")
	jan24           = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24           = time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	mar24           = time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	apr24           = time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	albumAlpha      = catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("alpha")}
	albumBeta       = catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("beta")}
	albumGamma      = catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("gamma")}
	albumOld        = catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("old-folder")}
	albumNew        = catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("new-folder")}
	ownerUserIdPort = stubOwnerUserIdPort(tonyOwner, ownerUserId)
)

func newAlbumViewForEventTest(repo AlbumSummaryRepository, counter MediaCounterPort) *AlbumView {
	return NewAlbumView(
		repo,
		GetAlbumSharingGridFunc(func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
			return nil, nil
		}),
		counter,
		FindAlbumsByIdsFunc(func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) {
			return nil, nil
		}),
		ownerUserIdPort,
	)
}

func TestAlbumView_AlbumCreated(t *testing.T) {
	type fields struct {
		Repository       *AlbumSummaryInMemoryRepository
		MediaCounterPort MediaCounterPort
	}
	tests := []struct {
		name            string
		fields          fields
		event           catalog.AlbumCreated
		expectSummaries []UserAlbumSummary
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "it should make the album visible to the owner",
			fields: fields{
				Repository:       &AlbumSummaryInMemoryRepository{},
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event: catalog.AlbumCreated{
				CreatedAlbum:      catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24},
				TransferredMedias: catalog.NewTransferredMedias(),
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 0},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should overwrite an existing owner row with the transfer-derived count",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "stale", Start: mar24, End: apr24, MediaCount: 5},
							Availability: OwnerAvailability(ownerUserId),
						},
					},
				},
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event: catalog.AlbumCreated{
				CreatedAlbum:      catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24},
				TransferredMedias: catalog.NewTransferredMedias(),
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 0},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should recount transferred source albums and leave their display fields untouched",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 10},
							Availability: OwnerAvailability(ownerUserId),
						},
					},
				},
				MediaCounterPort: MediaCounterPortFake{albumBeta: 7},
			},
			event: catalog.AlbumCreated{
				CreatedAlbum: catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{albumAlpha: {"m1", "m2", "m3"}},
					FromAlbums: []catalog.AlbumId{albumBeta},
				},
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 7},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, tt.fields.MediaCounterPort)
			err := view.OnAlbumCreated(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_MediasInserted(t *testing.T) {
	seededOwnerAndVisitorRepo := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 1},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 1},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
		}
	}

	type fields struct {
		Repository *AlbumSummaryInMemoryRepository
	}
	tests := []struct {
		name            string
		fields          fields
		medias          map[catalog.AlbumId][]catalog.MediaId
		expectSummaries []UserAlbumSummary
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name:   "it should increment the count on every viewer row",
			fields: fields{Repository: seededOwnerAndVisitorRepo()},
			medias: map[catalog.AlbumId][]catalog.MediaId{albumAlpha: {"a", "b", "c"}},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 4},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 4},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should be a no-op when the event is empty",
			fields: fields{Repository: seededOwnerAndVisitorRepo()},
			medias: map[catalog.AlbumId][]catalog.MediaId{},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 1},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 1},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, MediaCounterPortFake(nil))
			err := view.OnMediasInserted(context.Background(), tt.medias)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumRenamed(t *testing.T) {
	seededOwnerAndVisitorOnOld := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
		}
	}
	seededOldAndSource := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: VisitorAvailability(visitorUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Source", Start: feb24, End: mar24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
		}
	}

	type fields struct {
		Repository       *AlbumSummaryInMemoryRepository
		MediaCounterPort MediaCounterPort
	}
	tests := []struct {
		name            string
		fields          fields
		event           catalog.AlbumRenamed
		expectSummaries []UserAlbumSummary
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "it should update the name on all viewer rows when folder unchanged",
			fields: fields{
				Repository:       seededOwnerAndVisitorOnOld(),
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event: catalog.AlbumRenamed{
				ExistingAlbum:     catalog.Album{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24},
				RenamedAlbum:      catalog.Album{AlbumId: albumOld, Name: "New Name", Start: jan24, End: feb24},
				TransferredMedias: catalog.NewTransferredMedias(),
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should replace the rows on all viewers when folder changes, inheriting count and dates from the owner projection",
			fields: fields{
				Repository:       seededOldAndSource(),
				MediaCounterPort: MediaCounterPortFake{},
			},
			event: catalog.AlbumRenamed{
				ExistingAlbum: catalog.Album{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24},
				RenamedAlbum:  catalog.Album{AlbumId: albumNew, Name: "New Name", Start: jan24, End: feb24},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{albumNew: {"m1", "m2", "m3", "m4"}},
					FromAlbums: []catalog.AlbumId{albumOld},
				},
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Source", Start: feb24, End: mar24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumNew, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumNew, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, tt.fields.MediaCounterPort)
			err := view.OnAlbumRenamed(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumDatesAmended(t *testing.T) {
	seededOwnerAndVisitor := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
		}
	}
	seededSourceAndAmended := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Source", Start: jan24, End: feb24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24, MediaCount: 0},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
		}
	}

	type fields struct {
		Repository       *AlbumSummaryInMemoryRepository
		MediaCounterPort MediaCounterPort
	}
	tests := []struct {
		name       string
		fields     fields
		event      catalog.AlbumDatesAmended
		expectRepo []UserAlbumSummary
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should update start and end on all viewer rows",
			fields: fields{
				Repository:       seededOwnerAndVisitor(),
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event: catalog.AlbumDatesAmended{
				DatesUpdate: catalog.DatesUpdate{
					UpdatedAlbum:  catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24},
					PreviousStart: jan24,
					PreviousEnd:   feb24,
				},
				TransferredMedias: catalog.NewTransferredMedias(),
			},
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24, MediaCount: 3},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should recount transferred source albums and leave their display fields",
			fields: fields{
				Repository:       seededSourceAndAmended(),
				MediaCounterPort: MediaCounterPortFake{albumBeta: 3, albumAlpha: 2},
			},
			event: catalog.AlbumDatesAmended{
				DatesUpdate: catalog.DatesUpdate{
					UpdatedAlbum:  catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: apr24},
					PreviousStart: feb24,
					PreviousEnd:   mar24,
				},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{albumAlpha: {"m1", "m2"}},
					FromAlbums: []catalog.AlbumId{albumBeta},
				},
			},
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Source", Start: jan24, End: feb24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: apr24, MediaCount: 0},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, tt.fields.MediaCounterPort)
			err := view.OnAlbumDatesAmended(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectRepo, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumDeleted(t *testing.T) {
	deletedRow := func(availability Availability) UserAlbumSummary {
		return UserAlbumSummary{
			AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 1},
			Availability: availability,
		}
	}
	survivingRow := UserAlbumSummary{
		AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: mar24, End: apr24, MediaCount: 5},
		Availability: OwnerAvailability(ownerUserId),
	}
	destinationRow := UserAlbumSummary{
		AlbumSummary: AlbumSummary{AlbumId: albumGamma, Name: "Gamma", Start: feb24, End: mar24, MediaCount: 2},
		Availability: OwnerAvailability(ownerUserId),
	}

	type fields struct {
		Repository       *AlbumSummaryInMemoryRepository
		MediaCounterPort MediaCounterPort
	}
	tests := []struct {
		name       string
		fields     fields
		event      catalog.AlbumDeleted
		expectRepo []UserAlbumSummary
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should remove the owner row",
			fields: fields{
				Repository:       &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{deletedRow(OwnerAvailability(ownerUserId))}},
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event:      catalog.AlbumDeleted{DeletedAlbumId: albumAlpha, TransferredMedias: catalog.NewTransferredMedias()},
			expectRepo: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should remove all visitor rows",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					deletedRow(OwnerAvailability(ownerUserId)),
					deletedRow(VisitorAvailability(visitorUserId)),
					deletedRow(VisitorAvailability(visitor2UserId)),
				}},
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event:      catalog.AlbumDeleted{DeletedAlbumId: albumAlpha, TransferredMedias: catalog.NewTransferredMedias()},
			expectRepo: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should leave other albums untouched",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					deletedRow(OwnerAvailability(ownerUserId)),
					survivingRow,
				}},
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event:      catalog.AlbumDeleted{DeletedAlbumId: albumAlpha, TransferredMedias: catalog.NewTransferredMedias()},
			expectRepo: []UserAlbumSummary{survivingRow},
			wantErr:    assert.NoError,
		},
		{
			name: "it should recount transferred destination albums",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					deletedRow(OwnerAvailability(ownerUserId)),
					destinationRow,
				}},
				MediaCounterPort: MediaCounterPortFake{albumGamma: 5},
			},
			event: catalog.AlbumDeleted{
				DeletedAlbumId: albumAlpha,
				TransferredMedias: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{albumGamma: {"m1", "m2", "m3"}},
				},
			},
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumGamma, Name: "Gamma", Start: feb24, End: mar24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, tt.fields.MediaCounterPort)
			err := view.OnAlbumDeleted(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectRepo, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumShared(t *testing.T) {
	weddingsAlbum := catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24}

	type fields struct {
		Repository       *AlbumSummaryInMemoryRepository
		MediaCounterPort MediaCounterPort
	}
	tests := []struct {
		name       string
		fields     fields
		album      catalog.Album
		userId     usermodel.UserId
		expectRepo []UserAlbumSummary
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should add a visitor row with display fields and count",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
						Availability: OwnerAvailability(ownerUserId),
					},
				}},
				MediaCounterPort: MediaCounterPortFake{albumAlpha: 7},
			},
			album:  weddingsAlbum,
			userId: visitorUserId,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, tt.fields.MediaCounterPort)
			err := view.AlbumShared(context.Background(), tt.album, tt.userId)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectRepo, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumUnShared(t *testing.T) {
	type fields struct {
		Repository *AlbumSummaryInMemoryRepository
	}
	tests := []struct {
		name       string
		fields     fields
		albumId    catalog.AlbumId
		userId     usermodel.UserId
		expectRepo []UserAlbumSummary
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should remove the visitor row only",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
						Availability: OwnerAvailability(ownerUserId),
					},
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
						Availability: VisitorAvailability(visitorUserId),
					},
				}},
			},
			albumId: albumAlpha,
			userId:  visitorUserId,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, MediaCounterPortFake(nil))
			err := view.AlbumUnShared(context.Background(), tt.albumId, tt.userId)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectRepo, tt.fields.Repository.Summaries)
			}
		})
	}
}
