package catalog_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
	"time"
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
		avenger1Id: []catalog.MediaId{"media-1", "media-2"},
	}

	type fields struct {
		findAlbumsByOwner         func(t *testing.T) catalog.FindAlbumsByOwnerPort
		countMediasBySelectors    func(t *testing.T) catalog.CountMediasBySelectorsPort
		amendAlbumDateRepository  func(t *testing.T) catalog.AmendAlbumDateRepositoryPort
		transferMedias            func(t *testing.T) catalog.TransferMediasRepositoryPort
		timelineMutationObservers func(t *testing.T) catalog.TimelineMutationObserver
	}
	type args struct {
		albumId catalog.AlbumId
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError assert.ErrorAssertionFunc
	}{
		{
			name: "it should amend the dates of an album, end to end, and call the observers",
			fields: fields{
				findAlbumsByOwner:        stubFindAlbumsByOwnerWith("ironman", &existingAlbum, &allYearAlbum),
				countMediasBySelectors:   stubCountMediasBySelectorsPort(1),
				amendAlbumDateRepository: expectAmendAlbumDateRepositoryCalled(avenger1Id, may24, jan25),
				transferMedias: expectTransferMediasRepositoryPortCalled(catalog.MediaTransferRecords{
					avenger1Id: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{allYearAlbum.AlbumId},
							Start:      jul24,
							End:        jan25,
						},
					},
				}, transferredMedias),
				timelineMutationObservers: expectTimelineMutationObserverCalled(transferredMedias),
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jan25,
			},
			wantError: assert.NoError,
		},
		{
			name: "it should not amend the dates and not call the observer if OrphanMediasError is raised",
			fields: fields{
				findAlbumsByOwner:         stubFindAlbumsByOwnerWith("ironman", &existingAlbum),
				countMediasBySelectors:    stubCountMediasBySelectorsPort(1),
				amendAlbumDateRepository:  expectAmendAlbumDateRepositoryNotCalled(),
				transferMedias:            expectTransferMediasPortNotCalled(),
				timelineMutationObservers: expectTimelineMutationObserverNotCalled(),
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jun24,
			},
			wantError: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amendAlbumDates := catalog.NewAmendAlbumDates(
				tt.fields.findAlbumsByOwner(t),
				tt.fields.countMediasBySelectors(t),
				tt.fields.amendAlbumDateRepository(t),
				tt.fields.transferMedias(t),
				tt.fields.timelineMutationObservers(t),
			)

			err := amendAlbumDates.AmendAlbumDates(context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
			tt.wantError(t, err, fmt.Sprintf("AmendAlbumDates(%v, %v, %v, %v)", context.Background(), tt.args.albumId, tt.args.start, tt.args.end))
		})
	}
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
		FindAlbumsByOwnerPort    func(t *testing.T) catalog.FindAlbumsByOwnerPort
		AmendAlbumDatesObservers func(t *testing.T) catalog.AmendAlbumDatesObserver
	}
	type args struct {
		albumId catalog.AlbumId
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an error if the album is not found",
			fields: fields{
				FindAlbumsByOwnerPort:    stubFindAlbumsByOwnerWith(owner),
				AmendAlbumDatesObservers: expectAmendAlbumDatesObserverNotCalled(),
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jun24,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundError, i...)
			},
		},
		{
			name: "it should return immediately if dates haven't changed",
			fields: fields{
				FindAlbumsByOwnerPort:    stubFindAlbumsByOwnerWith(owner, &existingAlbum),
				AmendAlbumDatesObservers: expectAmendAlbumDatesObserverNotCalled(),
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
				FindAlbumsByOwnerPort: stubFindAlbumsByOwnerWith(owner, &existingAlbum),
				AmendAlbumDatesObservers: expectAmendAlbumDatesObserverCalled([]*catalog.Album{&existingAlbum}, catalog.Album{
					AlbumId: avenger1Id,
					Name:    "Avenger 1",
					Start:   may24,
					End:     jul24,
				}),
			},
			args: args{
				albumId: avenger1Id,
				start:   may24,
				end:     jul24,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &catalog.AmendAlbumDates{
				FindAlbumsByOwnerPort:    tt.fields.FindAlbumsByOwnerPort(t),
				AmendAlbumDatesObservers: []catalog.AmendAlbumDatesObserver{tt.fields.AmendAlbumDatesObservers(t)},
			}
			err := a.AmendAlbumDates(context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
			tt.wantErr(t, err, fmt.Sprintf("AmendAlbumDates(%v, %v, %v, %v)", context.Background(), tt.args.albumId, tt.args.start, tt.args.end))
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

	type fields struct {
		CountMediasBySelectors func(t *testing.T) catalog.CountMediasBySelectorsPort
		MediaTransfer          func(t *testing.T) catalog.MediaTransfer
	}
	type args struct {
		existingTimeline []*catalog.Album
		updatedAlbum     catalog.Album
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should not transfer any media because there were not other albums are present - GROWING BOTH SIDES",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer:          expectMediaTransferNotCalled(),
			},
			args: args{
				existingTimeline: []*catalog.Album{&mayAlbum},
				updatedAlbum:     amendWithDatesOf(mayAlbum.AlbumId, aprToJunAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not transfer any media because there were not other albums are present - SHRINKING BOTH SIDES",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(0),
				MediaTransfer:          expectMediaTransferNotCalled(),
			},
			args: args{
				existingTimeline: []*catalog.Album{&aprToJunAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum.AlbumId, mayAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should transfer medias IN the amended album - GROWING BOTH SIDES",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					mayAlbum.AlbumId: []catalog.MediaSelector{
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
				}),
			},
			args: args{
				existingTimeline: []*catalog.Album{&fullYearAlbum, &mayAlbum},
				updatedAlbum:     amendWithDatesOf(mayAlbum.AlbumId, aprToJunAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should transfer medias OUT the amended album - SHRINKING BOTH SIDES",
			// TODO also check for no other album is present
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					fullYearAlbum.AlbumId: []catalog.MediaSelector{
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
				}),
			},
			args: args{
				existingTimeline: []*catalog.Album{&fullYearAlbum, &aprToJunAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum.AlbumId, mayAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should transfer medias OUT the amended album - SHRINKING BEFORE with another high priority album",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					fullYearAlbum.AlbumId: []catalog.MediaSelector{
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
				}),
			},
			args: args{
				existingTimeline: []*catalog.Album{&fullYearAlbum, &aprToJunAlbum, &mayAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum.AlbumId, fifthJunAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should loose some segments on covered time range when growing",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					junAlbum.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{fifthJunAlbum.AlbumId},
							Start:      fifthJunAlbum.Start,
							End:        junAlbum.End,
						},
					},
				}),
			},
			args: args{
				existingTimeline: []*catalog.Album{&fifthJunAlbum, &junAlbum},
				updatedAlbum:     amendWithDatesOf(fifthJunAlbum.AlbumId, aprToJunAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should win some segments on covered time range when shrinking",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer: expectMediaTransferCalled(catalog.MediaTransferRecords{
					fullYearAlbum.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{aprToJunAlbum.AlbumId},
							Start:      aprToJunAlbum.Start,
							End:        junAlbum.Start,
						},
					},
					aprToJunAlbum.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{junAlbum.AlbumId, fullYearAlbum.AlbumId},
							Start:      fifthJunAlbum.Start,
							End:        fifthJunAlbum.End,
						},
					},
				}),
			},
			args: args{
				existingTimeline: []*catalog.Album{&aprToJunAlbum, &junAlbum, &fullYearAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum.AlbumId, fifthJunAlbum),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should abort if some medias are made orphan - SHRINKING BOTH SIDES",
			fields: fields{
				CountMediasBySelectors: stubCountMediasBySelectorsPort(1),
				MediaTransfer:          expectMediaTransferNotCalled(),
			},
			args: args{
				existingTimeline: []*catalog.Album{&aprToJunAlbum},
				updatedAlbum:     amendWithDatesOf(aprToJunAlbum.AlbumId, junAlbum),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &catalog.AmendAlbumMediaTransfer{
				CountMediasBySelectors: tt.fields.CountMediasBySelectors(t),
				MediaTransfer:          tt.fields.MediaTransfer(t),
			}

			err := a.OnAlbumDatesAmended(context.Background(), tt.args.existingTimeline, tt.args.updatedAlbum)
			tt.wantErr(t, err, fmt.Sprintf("OnAlbumDatesAmended(%v, %v, %v)", context.Background(), tt.args.existingTimeline, tt.args.updatedAlbum))
		})
	}
}

func amendWithDatesOf(id catalog.AlbumId, album catalog.Album) catalog.Album {
	album.AlbumId = id
	return album
}

func expectAmendAlbumDatesObserverCalled(timeline []*catalog.Album, updatedAlbum catalog.Album) func(t *testing.T) catalog.AmendAlbumDatesObserver {
	return func(t *testing.T) catalog.AmendAlbumDatesObserver {
		observer := mocks.NewAmendAlbumDatesObserver(t)
		observer.EXPECT().OnAlbumDatesAmended(mock.Anything, timeline, updatedAlbum).Return(nil)
		return observer
	}
}

func expectAmendAlbumDatesObserverNotCalled() func(t *testing.T) catalog.AmendAlbumDatesObserver {
	return func(t *testing.T) catalog.AmendAlbumDatesObserver {
		return catalog.AmendAlbumDatesObserverFunc(func(ctx context.Context, existingTimeline []*catalog.Album, updatedAlbum catalog.Album) error {
			assert.Failf(t, "AmendAlbumDatesObserver should not be called", "OnAlbumDatesAmended(%v, %v, %+v)", ctx, existingTimeline, updatedAlbum)
			return nil
		})
	}
}

func expectAmendAlbumDateRepositoryNotCalled() func(t *testing.T) catalog.AmendAlbumDateRepositoryPort {
	return func(t *testing.T) catalog.AmendAlbumDateRepositoryPort {
		return mocks.NewAmendAlbumDateRepositoryPort(t)
	}
}

func expectAmendAlbumDateRepositoryCalled(albumId catalog.AlbumId, start time.Time, end time.Time) func(t *testing.T) catalog.AmendAlbumDateRepositoryPort {
	return func(t *testing.T) catalog.AmendAlbumDateRepositoryPort {
		repo := mocks.NewAmendAlbumDateRepositoryPort(t)
		repo.EXPECT().AmendDates(mock.Anything, albumId, start, end).Return(nil).Once()
		return repo
	}
}
