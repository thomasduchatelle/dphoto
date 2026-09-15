package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestNewAmendAlbumDatesAcceptance(t *testing.T) {
	const owner = "ironman"
	avenger1Id := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avenger-1")}
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
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
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/all-year")},
		Name:    "All Year",
		Start:   jan24,
		End:     jan25,
	}
	transferredMedias := catalog.TransferredMedias{
		Transfers: map[catalog.AlbumId][]catalog.MediaId{
			avenger1Id: {"media-1", "media-2"},
		},
	}

	countOneMedia := func(_ ownermodel.Owner, _ catalog.MediaSelector) int {
		return 1
	}

	albumRepositoryWithOneMediaPerSelector := func(albums ...catalog.Album) *AlbumRepositoryInMemory {
		copies := make([]*catalog.Album, 0, len(albums))
		for _, album := range albums {
			fresh := album
			copies = append(copies, &fresh)
		}
		repo := NewAlbumRepositoryInMemory(copies...)
		repo.MediasBySelector = countOneMedia
		return repo
	}
	transferMediasReturning := func(transferred catalog.TransferredMedias) *TransferMediasInMemory {
		transfer := NewTransferMediasInMemory()
		transfer.TransferredMedias = transferred
		return transfer
	}

	type fields struct {
		AlbumRepository  *AlbumRepositoryInMemory
		TransferMedias   *TransferMediasInMemory
		TimelineObserver *TimelineMutationObserverInMemory
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
		expectNotifications   []catalog.TransferredMedias
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name: "it should amend the dates of an album, end to end, and call the observers",
			fields: fields{
				AlbumRepository:  albumRepositoryWithOneMediaPerSelector(existingAlbum, allYearAlbum),
				TransferMedias:   transferMediasReturning(transferredMedias),
				TimelineObserver: &TimelineMutationObserverInMemory{},
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jan25,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id:           {start: may24, end: jan25},
				allYearAlbum.AlbumId: {start: allYearAlbum.Start, end: allYearAlbum.End},
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				avenger1Id: {
					{
						FromAlbums: []catalog.AlbumId{allYearAlbum.AlbumId},
						Start:      jul24,
						End:        jan25,
					},
				},
			}},
			expectNotifications: []catalog.TransferredMedias{{
				Transfers:  transferredMedias.Transfers,
				FromAlbums: []catalog.AlbumId{allYearAlbum.AlbumId},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should not amend the dates and not call the observer if OrphanMediasError is raised",
			fields: fields{
				AlbumRepository:  albumRepositoryWithOneMediaPerSelector(existingAlbum),
				TransferMedias:   NewTransferMediasInMemory(),
				TimelineObserver: &TimelineMutationObserverInMemory{},
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jun24,
			},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				avenger1Id: {start: existingAlbum.Start, end: existingAlbum.End},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amendAlbumDates := catalog.NewAmendAlbumDates(
				tt.fields.AlbumRepository,
				tt.fields.AlbumRepository,
				tt.fields.AlbumRepository,
				tt.fields.TransferMedias,
				nil,
				tt.fields.TimelineObserver,
			)

			err := amendAlbumDates.AmendAlbumDates(context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
			if !tt.wantErr(t, err, fmt.Sprintf("AmendAlbumDates(%v, %v, %v)", tt.args.albumId, tt.args.start, tt.args.end)) {
				return
			}

			dates := make(map[catalog.AlbumId]albumDates, len(tt.fields.AlbumRepository.Albums))
			for id, album := range tt.fields.AlbumRepository.Albums {
				dates[id] = albumDates{start: album.Start, end: album.End}
			}
			assert.Equal(t, tt.expectAlbumDates, dates, "album dates in the repository")
			assert.Equal(t, tt.expectTransferRecords, tt.fields.TransferMedias.Records)
			assert.Equal(t, tt.expectNotifications, tt.fields.TimelineObserver.Notifications)
		})
	}
}

type albumDates struct {
	start time.Time
	end   time.Time
}

func TestAmendAlbumDates_AmendAlbumDates(t *testing.T) {
	const owner = "ironman"
	avenger1Id := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avenger-1")}
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	existingAlbum := catalog.Album{
		AlbumId: avenger1Id,
		Name:    "Avenger 1",
		Start:   may24,
		End:     jun24,
	}
	type fields struct {
		Albums []*catalog.Album
	}
	type args struct {
		albumId catalog.AlbumId
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectObserved []catalog.DatesUpdate
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name:   "it should return an error if the album is not found",
			fields: fields{},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jun24,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
		{
			name: "it should return immediately if dates haven't changed",
			fields: fields{
				Albums: []*catalog.Album{&existingAlbum},
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jun24,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should call the observers if dates have changed",
			fields: fields{
				Albums: []*catalog.Album{&existingAlbum},
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jul24,
			},
			expectObserved: []catalog.DatesUpdate{
				{
					UpdatedAlbum: catalog.Album{
						AlbumId: avenger1Id,
						Name:    "Avenger 1",
						Start:   may24,
						End:     jul24,
					},
					PreviousStart: existingAlbum.Start,
					PreviousEnd:   existingAlbum.End,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &AlbumDatesAmendedObserverInMemory{}
			a := &catalog.AmendAlbumDatesStateless{
				Observers: []catalog.AlbumDatesAmendedObserverWithTimeline{
					observer,
				},
			}

			err := a.AmendAlbumDates(context.Background(), catalog.NewLazyTimelineAggregate(tt.fields.Albums), tt.args.albumId, tt.args.start, tt.args.end)
			if !tt.wantErr(t, err, fmt.Sprintf("AmendAlbumDates(%v, %v, %v, %v)", context.Background(), tt.args.albumId, tt.args.start, tt.args.end)) {
				return
			}

			assert.ElementsMatchf(t, observer.DateAmendedAlbums, tt.expectObserved, "AmendAlbumDates(%v, %v, %v, %v)", context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
		})
	}
}

func TestAmendAlbumMediaTransfer_OnAlbumDatesAmended(t *testing.T) {
	fullYearAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/full-year")},
		Name:    "Full Year",
		Start:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	mayAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/may")},
		Name:    "May",
		Start:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	junAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/jun")},
		Name:    "Jun",
		Start:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
	}
	fifthJunAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/jun-fifth")},
		Name:    "Jun Fifth",
		Start:   time.Date(2024, 6, 5, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
	}
	aprToJunAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/apr-to-jun")},
		Name:    "Apr to Jun",
		Start:   time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
	}

	countOneMediaPerSelector := func(_ ownermodel.Owner, _ catalog.MediaSelector) int {
		return 1
	}

	albumRepositoryWithOneMediaPerSelector := func() *AlbumRepositoryInMemory {
		repo := NewAlbumRepositoryInMemory()
		repo.MediasBySelector = countOneMediaPerSelector
		return repo
	}

	type fields struct {
		AlbumRepository *AlbumRepositoryInMemory
	}
	type args struct {
		existingTimeline []*catalog.Album
		updatedAlbum     catalog.DatesUpdate
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectTransferRecords []catalog.MediaTransferRecords
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name:   "it should not transfer any media because there were not other albums are present - GROWING BOTH SIDES",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&mayAlbum},
				updatedAlbum:     amendWithDatesOf(mayAlbum, aprToJunAlbum.Start, aprToJunAlbum.End),
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should not transfer any media because there were not other albums are present - SHRINKING BOTH SIDES",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			args: args{
				existingTimeline: []*catalog.Album{&aprToJunAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum, mayAlbum.Start, mayAlbum.End),
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should transfer medias IN the amended album - GROWING BOTH SIDES",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&fullYearAlbum, &mayAlbum},
				updatedAlbum:     amendWithDatesOf(mayAlbum, aprToJunAlbum.Start, aprToJunAlbum.End),
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				mayAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{fullYearAlbum.AlbumId},
						Start:      aprToJunAlbum.Start,
						End:        mayAlbum.Start,
					},
					{
						FromAlbums: []catalog.AlbumId{fullYearAlbum.AlbumId},
						Start:      mayAlbum.End,
						End:        aprToJunAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should transfer medias OUT the amended album - SHRINKING BOTH SIDES",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&fullYearAlbum, &aprToJunAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum, mayAlbum.Start, mayAlbum.End),
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				fullYearAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
						Start:      aprToJunAlbum.Start,
						End:        mayAlbum.Start,
					},
					{
						FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
						Start:      mayAlbum.End,
						End:        aprToJunAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should transfer medias OUT the amended album - SHRINKING BEFORE with another high priority album",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&fullYearAlbum, &aprToJunAlbum, &mayAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum, fifthJunAlbum.Start, fifthJunAlbum.End),
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				fullYearAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
						Start:      aprToJunAlbum.Start,
						End:        mayAlbum.Start,
					},
					{
						FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
						Start:      mayAlbum.End,
						End:        fifthJunAlbum.Start,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should loose some segments on covered time range when growing",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&fifthJunAlbum, &junAlbum},
				updatedAlbum:     amendWithDatesOf(fifthJunAlbum, aprToJunAlbum.Start, junAlbum.End),
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				junAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{fifthJunAlbum.AlbumId},
						Start:      fifthJunAlbum.Start,
						End:        junAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should win some segments on covered time range when shrinking",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&aprToJunAlbum, &junAlbum, &fullYearAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum, fifthJunAlbum.Start, fifthJunAlbum.End),
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				fullYearAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
						Start:      aprToJunAlbum.Start,
						End:        junAlbum.Start,
					},
				},
				aprToJunAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{junAlbum.AlbumId, fullYearAlbum.AlbumId},
						Start:      fifthJunAlbum.Start,
						End:        fifthJunAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should transfer medias even when timeline has already been updated in case the changes is re-applied [scenario: win some segments on covered time range when shrinking]",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{{
					AlbumId: aprToJunAlbum.AlbumId,
					Name:    aprToJunAlbum.Name,
					Start:   fifthJunAlbum.Start,
					End:     fifthJunAlbum.End,
				}, &junAlbum, &fullYearAlbum},
				updatedAlbum: amendWithDatesOf(aprToJunAlbum, fifthJunAlbum.Start, fifthJunAlbum.End),
			},
			expectTransferRecords: []catalog.MediaTransferRecords{{
				fullYearAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
						Start:      aprToJunAlbum.Start,
						End:        junAlbum.Start,
					},
				},
				aprToJunAlbum.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{junAlbum.AlbumId, fullYearAlbum.AlbumId},
						Start:      fifthJunAlbum.Start,
						End:        fifthJunAlbum.End,
					},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name:   "it should abort if some medias are made orphan - SHRINKING BOTH SIDES",
			fields: fields{AlbumRepository: albumRepositoryWithOneMediaPerSelector()},
			args: args{
				existingTimeline: []*catalog.Album{&aprToJunAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum, junAlbum.Start, junAlbum.End),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediaTransfer := &MediaTransferInMemory{}
			a := &catalog.AmendAlbumMediaTransfer{
				CountMediasBySelectors: tt.fields.AlbumRepository,
				MediaTransfer:          mediaTransfer,
			}

			err := a.OnAlbumDatesAmendedWithTimeline(context.Background(), catalog.NewLazyTimelineAggregate(tt.args.existingTimeline), tt.args.updatedAlbum)
			if !tt.wantErr(t, err, fmt.Sprintf("OnAlbumDatesAmended(%v, %v, %v)", context.Background(), tt.args.existingTimeline, tt.args.updatedAlbum)) {
				return
			}
			assert.Equal(t, tt.expectTransferRecords, mediaTransfer.Records)
		})
	}
}

func amendWithDatesOf(album catalog.Album, start, end time.Time) catalog.DatesUpdate {
	update := catalog.DatesUpdate{
		UpdatedAlbum:  album,
		PreviousStart: album.Start,
		PreviousEnd:   album.End,
	}
	update.UpdatedAlbum.Start = start
	update.UpdatedAlbum.End = end

	return update
}
