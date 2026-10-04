package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type albumDates struct {
	start time.Time
	end   time.Time
}

func TestAmendAlbumDates_AmendAlbumDates(t *testing.T) {
	const owner = "ironman"
	avenger1Id := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avenger-1")}
	allYearId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/all-year")}
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	may01 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	may03 := time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC)
	may04 := time.Date(2024, 5, 4, 0, 0, 0, 0, time.UTC)
	may05 := time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	existingAvenger1 := catalog.Album{
		AlbumId: avenger1Id,
		Name:    "Avenger 1",
		Start:   may01,
		End:     may05,
	}
	existingAllYear := catalog.Album{
		AlbumId: allYearId,
		Name:    "All Year",
		Start:   jan24,
		End:     jan25,
	}
	testError := errors.Errorf("TEST error throwing")

	repositoryWithOneMediaPerSelector := func(albums ...catalog.Album) *AlbumRepositoryInMemory {
		copies := make([]*catalog.Album, 0, len(albums))
		for _, album := range albums {
			fresh := album
			copies = append(copies, &fresh)
		}
		repo := NewAlbumRepositoryInMemory(copies...)
		repo.MediasBySelector = func(_ ownermodel.Owner, _ catalog.MediaSelector) int {
			return 1
		}
		return repo
	}

	shrunkAvenger := catalog.Album{
		AlbumId: avenger1Id,
		Name:    "Avenger 1",
		Start:   may01,
		End:     may04,
	}
	movedMediaId := fakeMediaId(avenger1Id, may04)
	shrinkTransferred := catalog.TransferredMedias{
		Transfers:  map[catalog.AlbumId][]catalog.MediaId{allYearId: {movedMediaId}},
		FromAlbums: []catalog.AlbumId{avenger1Id},
	}
	shrinkRecords := catalog.MediaTransferRecords{
		allYearId: []catalog.MediaSelector{
			{
				FromAlbums: []catalog.AlbumId{avenger1Id},
				Start:      may04,
				End:        may05,
			},
		},
	}

	avengerRemain1 := &catalog.MediaMeta{Id: "a-remain-1", Filename: "a-remain-1.jpg", Type: catalog.MediaTypeImage}
	avengerRemain2 := &catalog.MediaMeta{Id: "a-remain-2", Filename: "a-remain-2.jpg", Type: catalog.MediaTypeImage}
	avengerRemain3 := &catalog.MediaMeta{Id: "a-remain-3", Filename: "a-remain-3.jpg", Type: catalog.MediaTypeImage}
	bRandom1 := &catalog.MediaMeta{Id: "b-random-1", Filename: "b-random-1.jpg", Type: catalog.MediaTypeImage}
	bRandom2 := &catalog.MediaMeta{Id: "b-random-2", Filename: "b-random-2.jpg", Type: catalog.MediaTypeImage}
	bRandom3 := &catalog.MediaMeta{Id: "b-random-3", Filename: "b-random-3.jpg", Type: catalog.MediaTypeImage}
	bRandom4 := &catalog.MediaMeta{Id: "b-random-4", Filename: "b-random-4.jpg", Type: catalog.MediaTypeImage}
	movedMediaMeta := &catalog.MediaMeta{Id: movedMediaId, Filename: "moved-may04.jpg", Type: catalog.MediaTypeImage}

	postTransferMedias := func() *MediaReadRepositoryInMemory {
		return &MediaReadRepositoryInMemory{
			Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
				avenger1Id: {avengerRemain1, avengerRemain2, avengerRemain3},
				allYearId:  {bRandom1, bRandom2, bRandom3, bRandom4, movedMediaMeta},
			},
		}
	}
	seededCoverRepository := func() *CoverRepositoryInMemory {
		return NewCoverRepositoryInMemory(
			coversFor(avenger1Id, catalog.Cover{MediaId: movedMediaId, Filename: "moved-may04.jpg", Origin: catalog.CoverOriginCherryPicked}),
			coversFor(allYearId,
				catalog.Cover{MediaId: bRandom1.Id, Filename: bRandom1.Filename, Origin: catalog.CoverOriginRandom},
				catalog.Cover{MediaId: bRandom2.Id, Filename: bRandom2.Filename, Origin: catalog.CoverOriginRandom},
				catalog.Cover{MediaId: bRandom3.Id, Filename: bRandom3.Filename, Origin: catalog.CoverOriginRandom},
				catalog.Cover{MediaId: bRandom4.Id, Filename: bRandom4.Filename, Origin: catalog.CoverOriginRandom},
			),
		)
	}

	expectedAvenger1Covers := []catalog.Cover{
		{MediaId: avengerRemain1.Id, Filename: avengerRemain1.Filename, Origin: catalog.CoverOriginRandom},
		{MediaId: avengerRemain2.Id, Filename: avengerRemain2.Filename, Origin: catalog.CoverOriginRandom},
		{MediaId: avengerRemain3.Id, Filename: avengerRemain3.Filename, Origin: catalog.CoverOriginRandom},
	}
	expectedAllYearCovers := []catalog.Cover{
		{MediaId: movedMediaId, Filename: "moved-may04.jpg", Origin: catalog.CoverOriginCherryPicked},
		{MediaId: bRandom2.Id, Filename: bRandom2.Filename, Origin: catalog.CoverOriginRandom},
		{MediaId: bRandom3.Id, Filename: bRandom3.Filename, Origin: catalog.CoverOriginRandom},
		{MediaId: bRandom4.Id, Filename: bRandom4.Filename, Origin: catalog.CoverOriginRandom},
	}
	expectedReconciledCovers := map[catalog.AlbumId][]catalog.Cover{
		avenger1Id: expectedAvenger1Covers,
		allYearId:  expectedAllYearCovers,
	}
	shrinkEvent := catalog.AlbumDatesAmended{
		DatesUpdate: catalog.DatesUpdate{
			UpdatedAlbum:  shrunkAvenger,
			PreviousStart: may01,
			PreviousEnd:   may05,
		},
		TransferredMedias: shrinkTransferred,
		Covers:            expectedReconciledCovers,
	}

	type fields struct {
		AlbumRepository     catalog.TimelineRepository
		TransferErr         error
		CoverRepository     catalog.CoverRepository
		MediaReadRepository catalog.MediaReadRepository
	}
	type args struct {
		albumId catalog.AlbumId
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectAlbumDates      map[catalog.AlbumId]albumDates
		expectTransferRecords []catalog.MediaTransferRecords
		expectAmendedEvents   []catalog.AlbumDatesAmended
		expectStoredCovers    map[catalog.AlbumId][]catalog.Cover
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name: "it should amend the dates, transfer medias, reconcile covers across both albums (CHERRY_PICKED cover moves to destination, source is backfilled) and fire the AlbumDatesAmended event carrying both cover sets",
			fields: fields{
				AlbumRepository:     repositoryWithOneMediaPerSelector(existingAvenger1, existingAllYear),
				CoverRepository:     seededCoverRepository(),
				MediaReadRepository: postTransferMedias(),
			},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may04,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may01, end: may04},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{shrinkRecords},
			expectAmendedEvents:   []catalog.AlbumDatesAmended{shrinkEvent},
			expectStoredCovers: map[catalog.AlbumId][]catalog.Cover{
				avenger1Id: expectedAvenger1Covers,
				allYearId:  expectedAllYearCovers,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the error and not fire the event when the cover repository fails to save the reconciled covers",
			fields: fields{
				AlbumRepository:     repositoryWithOneMediaPerSelector(existingAvenger1, existingAllYear),
				CoverRepository:     &failingSaveCoverRepository{CoverRepositoryInMemory: seededCoverRepository(), err: testError},
				MediaReadRepository: postTransferMedias(),
			},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may04,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may01, end: may04},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{shrinkRecords},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
		{
			name:   "it should return without firing any event when the dates have not changed",
			fields: fields{AlbumRepository: repositoryWithOneMediaPerSelector(existingAvenger1)},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may05,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may01, end: may05},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should return OrphanedMediasErr and not persist nor transfer when medias would be orphaned",
			fields: fields{AlbumRepository: repositoryWithOneMediaPerSelector(existingAvenger1)},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may03,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may01, end: may05},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
		{
			name:   "it should return AlbumNotFoundErr when the album does not exist",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may05,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
		{
			name: "it should not persist the new dates nor fire the event when the media transfer fails (persistence happens after the transfer)",
			fields: fields{
				AlbumRepository: repositoryWithOneMediaPerSelector(existingAvenger1, existingAllYear),
				TransferErr:     testError,
			},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may04,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may01, end: may05},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{shrinkRecords},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
		{
			name:   "it should not fire the event when persisting the new dates fails after a successful media transfer",
			fields: fields{AlbumRepository: failingAmendDates(repositoryWithOneMediaPerSelector(existingAvenger1, existingAllYear), testError)},
			args: args{
				albumId: avenger1Id,
				start:   may01,
				end:     may04,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may01, end: may05},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{shrinkRecords},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferService := &TransferMediasServiceFake{Err: tt.fields.TransferErr}
			observer := &AlbumDatesAmendedObserverInMemory{}

			coverRepository := tt.fields.CoverRepository
			if coverRepository == nil {
				coverRepository = NewCoverRepositoryInMemory()
			}
			mediaReadRepository := tt.fields.MediaReadRepository
			if mediaReadRepository == nil {
				mediaReadRepository = &MediaReadRepositoryInMemory{}
			}
			coverService := &catalog.CoverService{
				CoverRepository:     coverRepository,
				MediaReadRepository: mediaReadRepository,
				Randomiser:          deterministicRandomiser,
			}

			repository := underlyingAmendRepository(tt.fields.AlbumRepository)
			amendAlbumDates := catalog.NewAmendAlbumDates(
				tt.fields.AlbumRepository,
				repository,
				transferService,
				coverService,
				observer,
			)

			err := amendAlbumDates.AmendAlbumDates(context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
			if !tt.wantErr(t, err, fmt.Sprintf("AmendAlbumDates(%v, %v, %v)", tt.args.albumId, tt.args.start, tt.args.end)) {
				return
			}

			dates := make(map[catalog.AlbumId]albumDates, len(repository.Albums))
			for id, album := range repository.Albums {
				dates[id] = albumDates{start: album.Start, end: album.End}
			}
			assert.Equal(t, tt.expectAlbumDates, dates, "album dates in the repository")
			assert.Equal(t, tt.expectTransferRecords, transferService.Records, "records passed to TransferMedias")
			assert.Equal(t, tt.expectAmendedEvents, observer.Events, "AlbumDatesAmended events fired")
			if tt.expectStoredCovers != nil {
				stored, ok := coverRepository.(*CoverRepositoryInMemory)
				if assert.True(t, ok, "cover repository should be the in-memory fake to be able to assert its content") {
					assert.Equal(t, tt.expectStoredCovers, stored.Covers, "covers stored in the repository")
				}
			}
		})
	}
}

func failingAmendDates(memory *AlbumRepositoryInMemory, err error) catalog.TimelineRepository {
	return &failingAmendDatesInterceptor{
		AlbumRepositoryInMemory: memory,
		err:                     err,
	}
}

type failingAmendDatesInterceptor struct {
	*AlbumRepositoryInMemory
	err error
}

func (i *failingAmendDatesInterceptor) AmendDates(_ context.Context, _ catalog.AlbumId, _, _ time.Time) error {
	return i.err
}

func underlyingAmendRepository(port catalog.TimelineRepository) *AlbumRepositoryInMemory {
	if repository, ok := port.(*AlbumRepositoryInMemory); ok {
		return repository
	}
	return port.(*failingAmendDatesInterceptor).AlbumRepositoryInMemory
}

type failingSaveCoverRepository struct {
	*CoverRepositoryInMemory
	err error
}

func (f *failingSaveCoverRepository) SaveCovers(_ context.Context, _ catalog.AlbumId, _ []catalog.Cover) error {
	return f.err
}
