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
	apr24 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	existingAlbum := catalog.Album{
		AlbumId: avenger1Id,
		Name:    "Avenger 1",
		Start:   may24,
		End:     jul24,
	}
	allYearAlbum := catalog.Album{
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

	extendedAvenger := catalog.Album{
		AlbumId: avenger1Id,
		Name:    "Avenger 1",
		Start:   may24,
		End:     jan25,
	}
	extendedTransferredIds := func() []catalog.MediaId {
		var ids []catalog.MediaId
		for day := jul24; day.Before(jan25); day = day.AddDate(0, 0, 1) {
			ids = append(ids, fakeMediaId(allYearId, day))
		}
		return ids
	}
	extendedEvent := catalog.AlbumDatesAmended{
		DatesUpdate: catalog.DatesUpdate{
			UpdatedAlbum:  extendedAvenger,
			PreviousStart: may24,
			PreviousEnd:   jul24,
		},
		TransferredMedias: catalog.TransferredMedias{
			Transfers:  map[catalog.AlbumId][]catalog.MediaId{avenger1Id: extendedTransferredIds()},
			FromAlbums: []catalog.AlbumId{allYearId},
		},
	}
	extendedRecords := catalog.MediaTransferRecords{
		avenger1Id: []catalog.MediaSelector{
			{
				FromAlbums: []catalog.AlbumId{allYearId},
				Start:      jul24,
				End:        jan25,
			},
		},
	}

	grownAvenger := catalog.Album{
		AlbumId: avenger1Id,
		Name:    "Avenger 1",
		Start:   apr24,
		End:     jan25,
	}
	grownAloneEvent := catalog.AlbumDatesAmended{
		DatesUpdate: catalog.DatesUpdate{
			UpdatedAlbum:  grownAvenger,
			PreviousStart: may24,
			PreviousEnd:   jul24,
		},
		TransferredMedias: catalog.NewTransferredMedias(),
	}

	type fields struct {
		AlbumRepository catalog.TimelineRepository
		TransferErr     error
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
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name:   "it should amend the dates end to end, transfer medias, and fire the AlbumDatesAmended event",
			fields: fields{AlbumRepository: repositoryWithOneMediaPerSelector(existingAlbum, allYearAlbum)},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jan25,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may24, end: jan25},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{extendedRecords},
			expectAmendedEvents:   []catalog.AlbumDatesAmended{extendedEvent},
			wantErr:               assert.NoError,
		},
		{
			name:   "it should return without firing any event when the dates have not changed",
			fields: fields{AlbumRepository: repositoryWithOneMediaPerSelector(existingAlbum)},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jul24,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may24, end: jul24},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should return OrphanedMediasErr and not persist nor transfer when medias would be orphaned",
			fields: fields{AlbumRepository: repositoryWithOneMediaPerSelector(existingAlbum)},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jun24,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may24, end: jul24},
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
				start:   may24,
				end:     jul24,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
		{
			name:   "it should not persist the new dates nor fire the event when the media transfer fails (persistence happens after the transfer)",
			fields: fields{
				AlbumRepository: repositoryWithOneMediaPerSelector(existingAlbum, allYearAlbum),
				TransferErr:     testError,
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jan25,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may24, end: jul24},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{extendedRecords},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
		{
			name:   "it should not fire the event when persisting the new dates fails after a successful media transfer",
			fields: fields{AlbumRepository: failingAmendDates(repositoryWithOneMediaPerSelector(existingAlbum, allYearAlbum), testError)},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jan25,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: may24, end: jul24},
				allYearId:  {start: jan24, end: jan25},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{extendedRecords},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
		{
			name:   "it should amend the dates and fire the event with empty TransferredMedias when no records need transferring",
			fields: fields{AlbumRepository: repositoryWithOneMediaPerSelector(existingAlbum)},
			args: args{
				albumId: avenger1Id,
				start:   apr24,
				end:     jan25,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: apr24, end: jan25},
			},
			expectAmendedEvents: []catalog.AlbumDatesAmended{grownAloneEvent},
			wantErr:             assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferService := &TransferMediasServiceFake{Err: tt.fields.TransferErr}
			observer := &AlbumDatesAmendedObserverInMemory{}

			repository := underlyingAmendRepository(tt.fields.AlbumRepository)
			amendAlbumDates := catalog.NewAmendAlbumDates(
				tt.fields.AlbumRepository,
				repository,
				transferService,
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
