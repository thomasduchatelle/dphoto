package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type albumDates struct {
	start time.Time
	end   time.Time
}

func TestAmendAlbumDates_AmendAlbumDates(t *testing.T) {
	const owner = "ironman"

	jan26 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	may01 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	may03 := time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC)
	may04 := time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)
	may05 := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	jan27 := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	may26Id := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/may-26")}
	may26Album := &catalog.Album{AlbumId: may26Id, Name: "may 26", Start: may01, End: may05}
	allYearId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/2026")}
	allYearAlbum := &catalog.Album{AlbumId: allYearId, Name: "2026", Start: jan26, End: jan27}

	photo2may26 := photoAt("photo2may26", may01.AddDate(0, 0, 1))
	photo3may26 := photoAt("photo3may26", may03)
	photo4may26 := photoAt("photo4may26", may04)

	type fields struct {
		Catalog *CatalogInMemory
		Covers  *CoverRepositoryInMemory
	}
	type args struct {
		albumId catalog.AlbumId
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectAlbumDates    map[catalog.AlbumId]albumDates
		expectMediasByAlbum map[catalog.AlbumId][]*catalog.MediaMeta
		expectAmendedEvents []catalog.AlbumDatesAmended
		expectStoredCovers  map[catalog.AlbumId][]catalog.Cover
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should shrink the album, transfer the medias that fall outside the new range to the surrounding album, and reconcile covers across both albums",
			fields: fields{
				Catalog: NewCatalogInMemory(
					withAlbum(allYearAlbum),
					withAlbum(may26Album, photo2may26, photo3may26, photo4may26),
				),
				Covers: NewCoverRepositoryInMemory(
					coversFor(may26Id, pickedCover(photo4may26), randomCover(photo2may26)),
				),
			},
			args: args{albumId: may26Id, start: may01, end: may04},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				may26Id:   {start: may01, end: may04},
				allYearId: {start: jan26, end: jan27},
			},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				may26Id:   {photo2may26, photo3may26},
				allYearId: {photo4may26},
			},
			expectAmendedEvents: []catalog.AlbumDatesAmended{{
				DatesUpdate: catalog.DatesUpdate{
					UpdatedAlbum:  catalog.Album{AlbumId: may26Id, Name: "may 26", Start: may01, End: may04},
					PreviousStart: may01,
					PreviousEnd:   may05,
				},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{allYearId: {photo4may26.Id}},
					FromAlbums: []catalog.AlbumId{may26Id},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					may26Id:   {randomCover(photo2may26), randomCover(photo3may26)},
					allYearId: {pickedCover(photo4may26)},
				},
			}},
			expectStoredCovers: map[catalog.AlbumId][]catalog.Cover{
				may26Id:   {randomCover(photo2may26), randomCover(photo3may26)},
				allYearId: {pickedCover(photo4may26)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return OrphanedMediasErr and not persist nor transfer when medias would be orphaned",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(may26Album, photo4may26)),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args: args{albumId: may26Id, start: may01, end: may03},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				may26Id: {start: may01, end: may05},
			},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{may26Id: {photo4may26}},
			expectStoredCovers:  map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.OrphanedMediasErr, i...)
			},
		},
		{
			name: "it should return without firing any event when the dates have not changed",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(may26Album)),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args: args{albumId: may26Id, start: may01, end: may05},
			expectAlbumDates: map[catalog.AlbumId]albumDates{
				may26Id: {start: may01, end: may05},
			},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{may26Id: nil},
			expectStoredCovers:  map[catalog.AlbumId][]catalog.Cover{},
			wantErr:             assert.NoError,
		},
		{
			name: "it should return AlbumNotFoundErr when the album does not exist",
			fields: fields{
				Catalog: NewCatalogInMemory(),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args:                args{albumId: may26Id, start: may01, end: may05},
			expectAlbumDates:    map[catalog.AlbumId]albumDates{},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{},
			expectStoredCovers:  map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &AlbumDatesAmendedObserverInMemory{}
			coverService := &catalog.CoverService{
				CoverRepository:     tt.fields.Covers,
				MediaReadRepository: tt.fields.Catalog,
				Randomiser:          deterministicRandomiser,
			}

			amendAlbumDates := catalog.NewAmendAlbumDates(
				tt.fields.Catalog,
				tt.fields.Catalog,
				&catalog.TransferMediasFromRepository{TransferMediasRepository: tt.fields.Catalog},
				coverService,
				observer,
			)

			err := amendAlbumDates.AmendAlbumDates(context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
			if !tt.wantErr(t, err, fmt.Sprintf("AmendAlbumDates(%v, %v, %v)", tt.args.albumId, tt.args.start, tt.args.end)) {
				return
			}

			dates := make(map[catalog.AlbumId]albumDates, len(tt.fields.Catalog.AlbumsByIds()))
			for id, album := range tt.fields.Catalog.AlbumsByIds() {
				dates[id] = albumDates{start: album.Start, end: album.End}
			}
			assert.Equal(t, tt.expectAlbumDates, dates, "album dates in the catalog")
			assert.Equal(t, tt.expectMediasByAlbum, tt.fields.Catalog.MediasByAlbum(), "medias remaining per album")
			assert.Equal(t, tt.expectAmendedEvents, observer.Events, "AlbumDatesAmended events fired")
			assert.Equal(t, tt.expectStoredCovers, tt.fields.Covers.Covers, "covers stored in the repository")
		})
	}
}
