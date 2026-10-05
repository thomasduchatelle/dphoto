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

func newAlbumViewForEventTest(repo AlbumSummaryRepository, counter MediaCounterPort, covers ...FindCoversByAlbumPort) *AlbumView {
	findCovers := FindCoversByAlbumPort(FindCoversByAlbumPortFake(nil))
	if len(covers) > 0 {
		findCovers = covers[0]
	}
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
		findCovers,
	)
}

func TestAlbumView_AlbumCreated(t *testing.T) {
	coverA := catalog.Cover{MediaId: "media-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}
	coverB := catalog.Cover{MediaId: "media-b", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked}
	staleCover := catalog.Cover{MediaId: "stale", Filename: "stale.jpg", Origin: catalog.CoverOriginRandom}

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
		{
			name: "it should propagate the covers of the new album and of every source album carried by the event to every viewer",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 10, Covers: []catalog.Cover{staleCover}},
							Availability: OwnerAvailability(ownerUserId),
						},
						{
							AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 10, Covers: []catalog.Cover{staleCover}},
							Availability: VisitorAvailability(visitorUserId),
						},
					},
				},
				MediaCounterPort: MediaCounterPortFake{albumBeta: 7},
			},
			event: catalog.AlbumCreated{
				CreatedAlbum: catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{albumAlpha: {"media-a", "m2"}},
					FromAlbums: []catalog.AlbumId{albumBeta},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					albumAlpha: {coverA, coverB},
					albumBeta:  {},
				},
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 7},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 7},
					Availability: VisitorAvailability(visitorUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 2, Covers: []catalog.Cover{coverA, coverB}},
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
	coverA := catalog.Cover{MediaId: "media-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}
	coverB := catalog.Cover{MediaId: "media-b", Filename: "b.jpg", Origin: catalog.CoverOriginRandom}
	coverC := catalog.Cover{MediaId: "media-c", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked}

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
	seededOwnerAndVisitorWithCovers := func() *AlbumSummaryInMemoryRepository {
		repo := seededOwnerAndVisitorRepo()
		for i := range repo.Summaries {
			repo.Summaries[i].AlbumSummary.Covers = []catalog.Cover{coverA}
		}
		return repo
	}

	type fields struct {
		Repository *AlbumSummaryInMemoryRepository
	}
	tests := []struct {
		name            string
		fields          fields
		event           catalog.MediasInserted
		expectSummaries []UserAlbumSummary
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name:   "it should increment the count on every viewer row",
			fields: fields{Repository: seededOwnerAndVisitorRepo()},
			event: catalog.MediasInserted{
				Inserted: map[catalog.AlbumId][]catalog.MediaId{albumAlpha: {"a", "b", "c"}},
			},
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
			name:            "it should be a no-op when the event is empty",
			fields:          fields{Repository: seededOwnerAndVisitorRepo()},
			event:           catalog.MediasInserted{},
			expectSummaries: seededOwnerAndVisitorRepo().Summaries,
			wantErr:         assert.NoError,
		},
		{
			name:   "it should write the covers of every album carried by the event alongside the count update",
			fields: fields{Repository: seededOwnerAndVisitorRepo()},
			event: catalog.MediasInserted{
				Inserted: map[catalog.AlbumId][]catalog.MediaId{albumAlpha: {"a", "b"}},
				Covers:   map[catalog.AlbumId][]catalog.Cover{albumAlpha: {coverA, coverB, coverC}},
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3, Covers: []catalog.Cover{coverA, coverB, coverC}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3, Covers: []catalog.Cover{coverA, coverB, coverC}},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should clear the covers when the event carries an empty cover list for an album",
			fields: fields{Repository: seededOwnerAndVisitorWithCovers()},
			event: catalog.MediasInserted{
				Covers: map[catalog.AlbumId][]catalog.Cover{albumAlpha: {}},
			},
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
			err := view.OnMediasInserted(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumRenamed(t *testing.T) {
	cherryCover := catalog.Cover{MediaId: "media-cherry", Filename: "cherry.jpg", Origin: catalog.CoverOriginCherryPicked}
	randomCover := catalog.Cover{MediaId: "media-random", Filename: "random.jpg", Origin: catalog.CoverOriginRandom}

	seededOwnerAndVisitorOnOld := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
		}
	}
	seededOldAndSource := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "Old Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
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
			name: "it should update the name on all viewer rows and leave the covers untouched when folder unchanged",
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
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumOld, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should replace the rows on all viewers when folder changes, carrying the original covers to the new identity verbatim",
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
					AlbumSummary: AlbumSummary{AlbumId: albumNew, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumNew, Name: "New Name", Start: jan24, End: feb24, MediaCount: 4, Covers: []catalog.Cover{cherryCover, randomCover}},
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
	coverA := catalog.Cover{MediaId: "media-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}
	coverB := catalog.Cover{MediaId: "media-b", Filename: "b.jpg", Origin: catalog.CoverOriginRandom}
	coverPreexisting := catalog.Cover{MediaId: "media-pre", Filename: "pre.jpg", Origin: catalog.CoverOriginCherryPicked}

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
	seededOwnerAndVisitorWithCovers := func() *AlbumSummaryInMemoryRepository {
		repo := seededOwnerAndVisitor()
		for i := range repo.Summaries {
			repo.Summaries[i].AlbumSummary.Covers = []catalog.Cover{coverPreexisting}
		}
		return repo
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
		{
			name: "it should write covers for every album carried by the event alongside the dates update",
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
				Covers:            map[catalog.AlbumId][]catalog.Cover{albumAlpha: {coverA, coverB}},
			},
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24, MediaCount: 3, Covers: []catalog.Cover{coverA, coverB}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24, MediaCount: 3, Covers: []catalog.Cover{coverA, coverB}},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should clear covers for an album whose event entry is explicitly empty",
			fields: fields{
				Repository:       seededOwnerAndVisitorWithCovers(),
				MediaCounterPort: MediaCounterPortFake(nil),
			},
			event: catalog.AlbumDatesAmended{
				DatesUpdate: catalog.DatesUpdate{
					UpdatedAlbum:  catalog.Album{AlbumId: albumAlpha, Name: "Alpha", Start: feb24, End: mar24},
					PreviousStart: jan24,
					PreviousEnd:   feb24,
				},
				TransferredMedias: catalog.NewTransferredMedias(),
				Covers:            map[catalog.AlbumId][]catalog.Cover{albumAlpha: {}},
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
		{
			name: "it should write the destination album covers carried by the event onto every viewer row",
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
					Transfers: map[catalog.AlbumId][]catalog.MediaId{albumGamma: {"m1"}},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					albumAlpha: nil,
					albumGamma: {{MediaId: "m1", Filename: "m1.jpg", Origin: catalog.CoverOriginCherryPicked}},
				},
			},
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{
						AlbumId:    albumGamma,
						Name:       "Gamma",
						Start:      feb24,
						End:        mar24,
						MediaCount: 5,
						Covers:     []catalog.Cover{{MediaId: "m1", Filename: "m1.jpg", Origin: catalog.CoverOriginCherryPicked}},
					},
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
	coverA := catalog.Cover{MediaId: "media-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}
	coverB := catalog.Cover{MediaId: "media-b", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked}

	type fields struct {
		Repository            *AlbumSummaryInMemoryRepository
		MediaCounterPort      MediaCounterPort
		FindCoversByAlbumPort FindCoversByAlbumPort
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
				MediaCounterPort:      MediaCounterPortFake{albumAlpha: 7},
				FindCoversByAlbumPort: FindCoversByAlbumPortFake(nil),
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
		{
			name: "it should copy the owner's current covers onto the visitor row",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
						Availability: OwnerAvailability(ownerUserId),
					},
				}},
				MediaCounterPort:      MediaCounterPortFake{albumAlpha: 7},
				FindCoversByAlbumPort: FindCoversByAlbumPortFake{albumAlpha: {coverA, coverB}},
			},
			album:  weddingsAlbum,
			userId: visitorUserId,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should leave the visitor row without covers when the owner has none",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7},
						Availability: OwnerAvailability(ownerUserId),
					},
				}},
				MediaCounterPort:      MediaCounterPortFake{albumAlpha: 7},
				FindCoversByAlbumPort: FindCoversByAlbumPortFake{},
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
		{
			name: "it should be idempotent when the visitor was already shared with the same covers",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
						Availability: OwnerAvailability(ownerUserId),
					},
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
						Availability: VisitorAvailability(visitorUserId),
					},
				}},
				MediaCounterPort:      MediaCounterPortFake{albumAlpha: 7},
				FindCoversByAlbumPort: FindCoversByAlbumPortFake{albumAlpha: {coverA, coverB}},
			},
			album:  weddingsAlbum,
			userId: visitorUserId,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA, coverB}},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newAlbumViewForEventTest(tt.fields.Repository, tt.fields.MediaCounterPort, tt.fields.FindCoversByAlbumPort)
			err := view.AlbumShared(context.Background(), tt.album, tt.userId)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectRepo, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumUnShared(t *testing.T) {
	coverA := catalog.Cover{MediaId: "media-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}

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
		{
			name: "it should remove the visitor row along with its covers while leaving the owner covers intact",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{Summaries: []UserAlbumSummary{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA}},
						Availability: OwnerAvailability(ownerUserId),
					},
					{
						AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA}},
						Availability: VisitorAvailability(visitorUserId),
					},
				}},
			},
			albumId: albumAlpha,
			userId:  visitorUserId,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 7, Covers: []catalog.Cover{coverA}},
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
