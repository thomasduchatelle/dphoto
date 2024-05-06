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

func TestDeleteAlbum_DeleteAlbum(t *testing.T) {
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	const owner = "ironman"
	toDeleteAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := catalog.Album{
		AlbumId: toDeleteAlbumId,
		Name:    "Avenger 1",
		Start:   mar24,
		End:     may24,
	}
	existingAllYearAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")},
		Name:    "lifetime",
		Start:   jan24,
		End:     jan25,
	}
	existingQ1Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q1")},
		Name:    "q1",
		Start:   jan24,
		End:     apr24,
	}
	existingQ2Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q2")},
		Name:    "q2",
		Start:   apr24,
		End:     jul24,
	}
	anExpectedError := errors.Errorf("TEST error")

	type fields struct {
		FindAlbumsByOwner         func(t *testing.T) catalog.FindAlbumsByOwnerPort
		CountMediasBySelectors    func(t *testing.T) catalog.CountMediasBySelectorsPort
		AlbumCanBeDeletedObserver func(t *testing.T) catalog.AlbumCanBeDeletedObserver
	}
	type args struct {
		albumId catalog.AlbumId
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should delete album if all segments can be transferred to 1 other album",
			fields: fields{
				FindAlbumsByOwner:      returnListOfAlbums(toDeleteAlbum.Owner, &existingAllYearAlbum, &toDeleteAlbum),
				CountMediasBySelectors: expectCountMediasBySelectorsPortNotCalled(),
				AlbumCanBeDeletedObserver: expectAlbumCanBeDeletedObservedCalled(toDeleteAlbumId, catalog.MediaTransferRecords{
					existingAllYearAlbum.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
							Start:      toDeleteAlbum.Start,
							End:        toDeleteAlbum.End,
						},
					},
				}),
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete album if all segments can be transferred to several other albums",
			fields: fields{
				FindAlbumsByOwner:      returnListOfAlbums(toDeleteAlbum.Owner, &existingQ1Album, &existingQ2Album, &toDeleteAlbum),
				CountMediasBySelectors: expectCountMediasBySelectorsPortNotCalled(),
				AlbumCanBeDeletedObserver: expectAlbumCanBeDeletedObservedCalled(toDeleteAlbumId, catalog.MediaTransferRecords{
					existingQ1Album.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
							Start:      toDeleteAlbum.Start,
							End:        existingQ1Album.End,
						},
					},
					existingQ2Album.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
							Start:      existingQ2Album.Start,
							End:        toDeleteAlbum.End,
						},
					},
				}),
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete album even if the segments are not covered by other albums as long as there is no medias to become orphaned",
			fields: fields{
				FindAlbumsByOwner: returnListOfAlbums(toDeleteAlbum.Owner, &existingQ1Album, &toDeleteAlbum),
				CountMediasBySelectors: expectCountMediasBySelectorsPortCalled(0, owner, catalog.MediaSelector{
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
					Start:      apr24,
					End:        may24,
				}),
				AlbumCanBeDeletedObserver: expectAlbumCanBeDeletedObservedCalled(toDeleteAlbumId, catalog.MediaTransferRecords{
					existingQ1Album.AlbumId: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
							Start:      toDeleteAlbum.Start,
							End:        existingQ1Album.End,
						},
					},
				}),
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should raise an error if medias are about to be orphaned",
			fields: fields{
				FindAlbumsByOwner: returnListOfAlbums(toDeleteAlbum.Owner, &existingQ1Album, &toDeleteAlbum),
				CountMediasBySelectors: expectCountMediasBySelectorsPortCalled(1, owner, catalog.MediaSelector{
					FromAlbums: []catalog.AlbumId{toDeleteAlbumId},
					Start:      apr24,
					End:        may24,
				}),
				AlbumCanBeDeletedObserver: expectAlbumCanBeDeletedObserverNotCalled(),
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasError, i...)
			},
		},
		{
			name: "it should fails if listing albums raises an error",
			fields: fields{
				FindAlbumsByOwner:         stubFindAlbumsByOwnerPortWithError(anExpectedError),
				CountMediasBySelectors:    expectCountMediasBySelectorsPortNotCalled(),
				AlbumCanBeDeletedObserver: expectAlbumCanBeDeletedObserverNotCalled(),
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, anExpectedError, i...)
			},
		},
		{
			name: "it should fails if listing albums raises an error",
			fields: fields{
				FindAlbumsByOwner:         returnListOfAlbums(toDeleteAlbum.Owner, &toDeleteAlbum),
				CountMediasBySelectors:    stubCountMediasBySelectorsPortWithError(anExpectedError),
				AlbumCanBeDeletedObserver: expectAlbumCanBeDeletedObserverNotCalled(),
			},
			args: args{
				albumId: toDeleteAlbumId,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, anExpectedError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var observers []catalog.AlbumCanBeDeletedObserver
			if tt.fields.AlbumCanBeDeletedObserver != nil {
				observers = append(observers, tt.fields.AlbumCanBeDeletedObserver(t))
			}
			d := &catalog.DeleteAlbum{
				FindAlbumsByOwner:         tt.fields.FindAlbumsByOwner(t),
				CountMediasBySelectors:    tt.fields.CountMediasBySelectors(t),
				AlbumCanBeDeletedObserver: observers,
			}
			err := d.DeleteAlbum(context.Background(), tt.args.albumId)
			tt.wantErr(t, err, fmt.Sprintf("DeleteAlbum(%v)", tt.args.albumId))
		})
	}
}

func TestDeleteAlbumMediaTransfer_Observe(t *testing.T) {
	avenger1Id := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers-1")}
	ironman1Id := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/ironman-1")}
	records := catalog.MediaTransferRecords{
		avenger1Id: {
			{
				FromAlbums: []catalog.AlbumId{ironman1Id},
				Start:      time.Time{},
				End:        time.Time{},
			},
		},
	}
	transfers := catalog.TransferredMedias{
		avenger1Id: []catalog.MediaId{"media-1", "media-2"},
	}
	emptyTransfers := catalog.TransferredMedias{
		avenger1Id: []catalog.MediaId{},
		ironman1Id: nil,
	}

	type fields struct {
		TransferMedias           func(t *testing.T) catalog.TransferMediasPort
		TimelineMutationObserver func(t *testing.T) catalog.TimelineMutationObserver
	}
	type args struct {
		deletedAlbum catalog.AlbumId
		transfers    catalog.MediaTransferRecords
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should perform the media transfer on the catalog DB",
			fields: fields{
				TransferMedias: expectTransferMediasPortCalled(records, transfers),
			},
			args: args{
				deletedAlbum: ironman1Id,
				transfers:    records,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should notify observers that medias should be transferred",
			fields: fields{
				TransferMedias:           expectTransferMediasPortCalled(records, transfers),
				TimelineMutationObserver: expectTimelineMutationObserverCalled(transfers),
			},
			args: args{
				deletedAlbum: ironman1Id,
				transfers:    records,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not notify observers when no medias should be transferred",
			fields: fields{
				TransferMedias:           expectTransferMediasPortCalled(records, emptyTransfers),
				TimelineMutationObserver: timelineMutationObserverNotCalled(),
			},
			args: args{
				deletedAlbum: ironman1Id,
				transfers:    records,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var observers []catalog.TimelineMutationObserver
			if tt.fields.TimelineMutationObserver != nil {
				observers = append(observers, tt.fields.TimelineMutationObserver(t))
			}
			d := &catalog.DeleteAlbumMediaTransfer{
				TransferMedias:            tt.fields.TransferMedias(t),
				TimelineMutationObservers: observers,
			}
			err := d.Observe(context.Background(), tt.args.deletedAlbum, tt.args.transfers)
			tt.wantErr(t, err, fmt.Sprintf("Observe(%v, %v)", tt.args.deletedAlbum, tt.args.transfers))
		})
	}
}

func expectAlbumCanBeDeletedObserverNotCalled() func(t *testing.T) catalog.AlbumCanBeDeletedObserver {
	return func(t *testing.T) catalog.AlbumCanBeDeletedObserver {
		return mocks.NewAlbumCanBeDeletedObserver(t)
	}
}

func expectAlbumCanBeDeletedObservedCalled(toDeleteAlbumId catalog.AlbumId, records catalog.MediaTransferRecords) func(t *testing.T) catalog.AlbumCanBeDeletedObserver {
	return func(t *testing.T) catalog.AlbumCanBeDeletedObserver {
		observer := mocks.NewAlbumCanBeDeletedObserver(t)
		observer.EXPECT().Observe(mock.Anything, toDeleteAlbumId, records).Return(nil).Once()
		return observer
	}
}

func timelineMutationObserverNotCalled() func(t *testing.T) catalog.TimelineMutationObserver {
	return func(t *testing.T) catalog.TimelineMutationObserver {
		return mocks.NewTimelineMutationObserver(t)
	}
}

func expectTimelineMutationObserverCalled(transfers catalog.TransferredMedias) func(t *testing.T) catalog.TimelineMutationObserver {
	return func(t *testing.T) catalog.TimelineMutationObserver {
		observer := mocks.NewTimelineMutationObserver(t)
		observer.EXPECT().Observe(mock.Anything, transfers).Return(nil).Once()
		return observer
	}
}

// TODO use expectTransferMediasPortCalled on create album as well
func expectTransferMediasPortCalled(expectedRecords catalog.MediaTransferRecords, returnedTransfers catalog.TransferredMedias) func(t *testing.T) catalog.TransferMediasPort {
	return func(t *testing.T) catalog.TransferMediasPort {
		port := mocks.NewTransferMediasPort(t)
		port.EXPECT().
			TransferMediasFromRecords(mock.Anything, expectedRecords).
			Return(returnedTransfers, nil).
			Once()
		return port
	}
}

func stubFindAlbumsByOwnerPortWithError(err error) func(t *testing.T) catalog.FindAlbumsByOwnerPort {
	return func(t *testing.T) catalog.FindAlbumsByOwnerPort {
		return catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner catalog.Owner) ([]*catalog.Album, error) {
			return nil, err
		})
	}
}

func expectCountMediasBySelectorsPortNotCalled() func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector) (int, error) {
			assert.Failf(t, "unexpected call", "CountMediasBySelectors(%+v)", selectors)
			return 0, nil
		})
	}
}

func expectCountMediasBySelectorsPortCalled(count int, expectedOwner catalog.Owner, expectedSelectors ...catalog.MediaSelector) func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector) (int, error) {
			assert.Equalf(t, expectedOwner, owner, "unexpected call CountMediasBySelectors(%v, %+v)", owner, selectors)
			assert.Equalf(t, expectedSelectors, selectors, "unexpected call CountMediasBySelectors(%v, %+v)", owner, selectors)

			return count, nil
		})
	}
}

func stubCountMediasBySelectorsPortWithError(err error) func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector) (int, error) {
			return 0, err
		})
	}
}

type ExternalTimelineMutationObserver struct {
	Transfers catalog.TransferredMedias
}

func (e *ExternalTimelineMutationObserver) Observe(ctx context.Context, transfers catalog.TransferredMedias) error {
	e.Transfers = transfers
	return nil
}

func TestNewDeleteAlbum(t *testing.T) {
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	const owner = "ironman"
	toDeleteAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := catalog.Album{
		AlbumId: toDeleteAlbumId,
		Name:    "Avenger 1",
		Start:   mar24,
		End:     may24,
	}
	existingAllYearAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")},
		Name:    "lifetime",
		Start:   jan24,
		End:     jan25,
	}

	externalObserver := new(ExternalTimelineMutationObserver)

	transferredMedias := catalog.TransferredMedias{
		existingAllYearAlbum.AlbumId: []catalog.MediaId{"media-1", "media-2"},
	}

	deleteAlbum := catalog.NewDeleteAlbum(
		catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner catalog.Owner) ([]*catalog.Album, error) {
			return []*catalog.Album{&existingAllYearAlbum, &toDeleteAlbum}, nil
		}),
		catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector) (int, error) {
			return 0, nil
		}),
		catalog.TransferMediasFunc(func(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
			return transferredMedias, nil
		}),
		catalog.DeleteAlbumRepositoryFunc(func(ctx context.Context, albumId catalog.AlbumId) error {
			return nil
		}),
		externalObserver,
	)

	err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
	if assert.NoError(t, err) {
		assert.Equal(t, externalObserver.Transfers, transferredMedias)
	}
}
