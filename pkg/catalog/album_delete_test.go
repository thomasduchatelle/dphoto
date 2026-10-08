package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// catalogTransferMediasFailing delegates every TransferMediasRepositoryPort call to the
// embedded CatalogInMemory except TransferMediasFromRecords, which returns Err. It lets
// tests prove that an in-flight transfer failure leaves the catalog in a recoverable
// state (source medias not moved, destination album not touched, downstream steps like
// DeleteAlbum / event-fire never execute).
type catalogTransferMediasFailing struct {
	*CatalogInMemory
	Err error
}

func (c *catalogTransferMediasFailing) TransferMediasFromRecords(_ context.Context, _ catalog.MediaTransferRecords) (map[catalog.AlbumId][]catalog.MediaId, error) {
	return nil, c.Err
}

func TestDeleteAlbum_DeleteAlbum(t *testing.T) {
	const owner = "ironman"

	jan26 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	mar26 := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	apr26 := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	jan27 := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	marAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/mar-26")}
	marAlbum := &catalog.Album{AlbumId: marAlbumId, Name: "mar 26", Start: mar26, End: apr26}
	lifetimeId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")}
	lifetimeAlbum := &catalog.Album{AlbumId: lifetimeId, Name: "lifetime", Start: jan26, End: jan27}
	jan_feb_26_Id := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/jan-feb-26")}
	jan_feb_26_Album := &catalog.Album{AlbumId: jan_feb_26_Id, Name: "jan-feb 26", Start: jan26, End: mar26}

	photo10mar26 := photoAt("photo10mar26", mar26.AddDate(0, 0, 10))
	photo20mar26 := photoAt("photo20mar26", mar26.AddDate(0, 0, 20))

	transferErr := errors.New("TEST transfer medias failure")

	// transferPort is an optional per-case override: when nil the test uses
	// tt.fields.Catalog (the real TransferMediasFromRepository path); when set it is
	// wrapped by TransferMediasFromRepository so a failing wrapper can short-circuit
	// the transfer step.
	type fields struct {
		Catalog      *CatalogInMemory
		Covers       *CoverRepositoryInMemory
		TransferPort catalog.TransferMediasRepositoryPort
	}
	type args struct {
		albumId catalog.AlbumId
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectAlbumIds      []catalog.AlbumId
		expectMediasByAlbum map[catalog.AlbumId][]*catalog.MediaMeta
		expectCoversByAlbum map[catalog.AlbumId][]catalog.Cover
		expectDeletedEvents []catalog.AlbumDeleted
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should transfer medias to the surrounding album and inherit the deleted album's CHERRY_PICKED cover",
			fields: fields{
				Catalog: NewCatalogInMemory(
					withAlbum(lifetimeAlbum),
					withAlbum(marAlbum, photo10mar26, photo20mar26),
				),
				Covers: NewCoverRepositoryInMemory(coversFor(marAlbumId, pickedCover(photo10mar26), randomCover(photo20mar26))),
			},
			args:           args{albumId: marAlbumId},
			expectAlbumIds: []catalog.AlbumId{lifetimeId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				lifetimeId: {photo10mar26, photo20mar26},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{
				lifetimeId: {pickedCover(photo10mar26), randomCover(photo20mar26)},
			},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId: marAlbumId,
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{lifetimeId: {photo10mar26.Id, photo20mar26.Id}},
					FromAlbums: []catalog.AlbumId{marAlbumId},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					lifetimeId:  {pickedCover(photo10mar26), randomCover(photo20mar26)},
					marAlbumId: nil,
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete an empty album, and fire an empty-transfer event when no surrounding album covers it",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(marAlbum)),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args:                args{albumId: marAlbumId},
			expectAlbumIds:      []catalog.AlbumId{},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			expectDeletedEvents: []catalog.AlbumDeleted{{
				DeletedAlbumId:    marAlbumId,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
				Covers:            nil,
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should return OrphanedMediasErr without touching any album when the deletion would orphan medias (jan-feb album ends before mar)",
			fields: fields{
				Catalog: NewCatalogInMemory(
					withAlbum(jan_feb_26_Album),
					withAlbum(marAlbum, photo10mar26),
				),
				Covers: NewCoverRepositoryInMemory(),
			},
			args:           args{albumId: marAlbumId},
			expectAlbumIds: []catalog.AlbumId{jan_feb_26_Id, marAlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				jan_feb_26_Id: nil,
				marAlbumId:    {photo10mar26},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
		{
			name: "it should return AlbumNotFoundErr when the album does not exist",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(lifetimeAlbum)),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args:                args{albumId: marAlbumId},
			expectAlbumIds:      []catalog.AlbumId{lifetimeId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{lifetimeId: nil},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
		{
			name: "it should not delete the album when the media transfer fails: medias stay in the source so a retry can re-run the transfer",
			fields: func() fields {
				backing := NewCatalogInMemory(
					withAlbum(lifetimeAlbum),
					withAlbum(marAlbum, photo10mar26, photo20mar26),
				)
				return fields{
					Catalog:      backing,
					Covers:       NewCoverRepositoryInMemory(coversFor(marAlbumId, pickedCover(photo10mar26), randomCover(photo20mar26))),
					TransferPort: &catalogTransferMediasFailing{CatalogInMemory: backing, Err: transferErr},
				}
			}(),
			args:           args{albumId: marAlbumId},
			expectAlbumIds: []catalog.AlbumId{lifetimeId, marAlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				lifetimeId: nil,
				marAlbumId: {photo10mar26, photo20mar26},
			},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{
				marAlbumId: {pickedCover(photo10mar26), randomCover(photo20mar26)},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, transferErr, i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferPort := tt.fields.TransferPort
			if transferPort == nil {
				transferPort = tt.fields.Catalog
			}
			observer := &AlbumDeletedObserverInMemory{}
			deleteAlbum := catalog.NewDeleteAlbum(
				tt.fields.Catalog,
				tt.fields.Catalog,
				&catalog.TransferMediasFromRepository{TransferMediasRepository: transferPort},
				&catalog.CoverService{
					CoverRepository:     tt.fields.Covers,
					MediaReadRepository: tt.fields.Catalog,
					Randomiser:          deterministicRandomiser,
				},
				observer,
			)

			err := deleteAlbum.DeleteAlbum(context.Background(), tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteAlbum(%v)", tt.args.albumId)) {
				return
			}

			assert.ElementsMatch(t, tt.expectAlbumIds, tt.fields.Catalog.AlbumIds(), "albums remaining in the catalog")
			assert.Equal(t, tt.expectMediasByAlbum, tt.fields.Catalog.MediasByAlbum(), "medias remaining per album")
			assert.Equal(t, tt.expectCoversByAlbum, tt.fields.Covers.Covers, "canonical covers stored")
			assert.Equal(t, tt.expectDeletedEvents, observer.Events, "AlbumDeleted events fired")
		})
	}
}
