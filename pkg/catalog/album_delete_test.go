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
		AlbumCanBeDeletedObserver func(t *testing.T) catalog.DeleteAlbumObserver
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
				FindAlbumsByOwner:      stubFindAlbumsByOwnerWith(toDeleteAlbum.Owner, &existingAllYearAlbum, &toDeleteAlbum),
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
				FindAlbumsByOwner:      stubFindAlbumsByOwnerWith(toDeleteAlbum.Owner, &existingQ1Album, &existingQ2Album, &toDeleteAlbum),
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
				FindAlbumsByOwner: stubFindAlbumsByOwnerWith(toDeleteAlbum.Owner, &existingQ1Album, &toDeleteAlbum),
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
				FindAlbumsByOwner: stubFindAlbumsByOwnerWith(toDeleteAlbum.Owner, &existingQ1Album, &toDeleteAlbum),
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
				FindAlbumsByOwner:         stubFindAlbumsByOwnerWith(toDeleteAlbum.Owner, &toDeleteAlbum),
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
			var observers []catalog.DeleteAlbumObserver
			if tt.fields.AlbumCanBeDeletedObserver != nil {
				observers = append(observers, tt.fields.AlbumCanBeDeletedObserver(t))
			}
			d := &catalog.DeleteAlbum{
				FindAlbumsByOwner:      tt.fields.FindAlbumsByOwner(t),
				CountMediasBySelectors: tt.fields.CountMediasBySelectors(t),
				Observers:              observers,
			}
			err := d.DeleteAlbum(context.Background(), tt.args.albumId)
			tt.wantErr(t, err, fmt.Sprintf("DeleteAlbum(%v)", tt.args.albumId))
		})
	}
}

func expectAlbumCanBeDeletedObserverNotCalled() func(t *testing.T) catalog.DeleteAlbumObserver {
	return func(t *testing.T) catalog.DeleteAlbumObserver {
		return mocks.NewDeleteAlbumObserver(t)
	}
}

func expectAlbumCanBeDeletedObservedCalled(toDeleteAlbumId catalog.AlbumId, records catalog.MediaTransferRecords) func(t *testing.T) catalog.DeleteAlbumObserver {
	return func(t *testing.T) catalog.DeleteAlbumObserver {
		observer := mocks.NewDeleteAlbumObserver(t)
		observer.EXPECT().OnDeleteAlbum(mock.Anything, toDeleteAlbumId, records).Return(nil).Once()
		return observer
	}
}

func stubFindAlbumsByOwnerPortWithError(err error) func(t *testing.T) catalog.FindAlbumsByOwnerPort {
	return func(t *testing.T) catalog.FindAlbumsByOwnerPort {
		return catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
			return nil, err
		})
	}
}

func expectCountMediasBySelectorsPortNotCalled() func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
			assert.Failf(t, "unexpected call", "CountMediasBySelectors(%+v)", selectors)
			return 0, nil
		})
	}
}

func expectCountMediasBySelectorsPortCalled(count int, expectedOwner ownermodel.Owner, expectedSelectors ...catalog.MediaSelector) func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
			assert.Equalf(t, expectedOwner, owner, "unexpected call CountMediasBySelectors(%v, %+v)", owner, selectors)
			assert.Equalf(t, expectedSelectors, selectors, "unexpected call CountMediasBySelectors(%v, %+v)", owner, selectors)

			return count, nil
		})
	}
}

func stubCountMediasBySelectorsPortWithError(err error) func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
			return 0, err
		})
	}
}

func stubCountMediasBySelectorsPort(count int) func(t *testing.T) catalog.CountMediasBySelectorsPort {
	return func(t *testing.T) catalog.CountMediasBySelectorsPort {
		return catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
			return count, nil
		})
	}
}

type ExternalTimelineMutationObserver struct {
	Transfers catalog.TransferredMedias
}

func (e *ExternalTimelineMutationObserver) OnTransferredMedias(ctx context.Context, transfers catalog.TransferredMedias) error {
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
		catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
			return []*catalog.Album{&existingAllYearAlbum, &toDeleteAlbum}, nil
		}),
		catalog.CountMediasBySelectorsFunc(func(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
			return 0, nil
		}),
		stubTransferMediaPort(transferredMedias)(t),
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

func stubTransferMediaPort(transferredMedias catalog.TransferredMedias) func(t *testing.T) catalog.TransferMediasRepositoryPort {
	return func(t *testing.T) catalog.TransferMediasRepositoryPort {
		return catalog.TransferMediasFunc(func(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
			return transferredMedias, nil
		})
	}
}
