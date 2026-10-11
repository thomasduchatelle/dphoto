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

func TestCreateAlbum_Create(t *testing.T) {
	const owner = "tonystark"

	apr28_26 := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC)
	may01_26 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	ironmanAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/2026-04_Ironman_1"),
		},
		Name:  "Ironman 1",
		Start: apr28_26,
		End:   may01_26,
	}
	standardRequest := catalog.CreateAlbumRequest{
		Owner: ownermodel.Owner(owner),
		Name:  ironmanAlbum.Name,
		Start: ironmanAlbum.Start,
		End:   ironmanAlbum.End,
	}

	lifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/lifetime"),
		},
		Name:  "lifetime",
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	photo28apr26 := photoAt("photo28apr26", apr28_26)
	photo29apr26 := photoAt("photo29apr26", apr28_26.AddDate(0, 0, 1))
	photo30apr26 := photoAt("photo30apr26", apr28_26.AddDate(0, 0, 2))
	photo10jun26 := photoAt("photo10jun26", may01_26.AddDate(0, 0, 40))

	type fields struct {
		Catalog *CatalogInMemory
		Covers  *CoverRepositoryInMemory
	}
	type args struct {
		request catalog.CreateAlbumRequest
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectAlbumIds      []catalog.AlbumId
		expectMediasByAlbum map[catalog.AlbumId][]*catalog.MediaMeta
		expectCreatedEvents []catalog.AlbumCreated
		expectSavedCovers   map[catalog.AlbumId][]catalog.Cover
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should reject an invalid request (empty name) without touching the catalog",
			fields: fields{
				Catalog: NewCatalogInMemory(),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: ownermodel.Owner(owner),
					Name:  "",
					Start: apr28_26,
					End:   may01_26,
				},
			},
			expectAlbumIds:      []catalog.AlbumId{},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{},
			expectSavedCovers:   map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr, i...)
			},
		},
		{
			name: "it should reject a request whose forced folder name is already taken",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(lifetimeAlbum)),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            ownermodel.Owner(owner),
					Name:             "A different name",
					Start:            apr28_26,
					End:              may01_26,
					ForcedFolderName: lifetimeAlbum.AlbumId.FolderName.String(),
				},
			},
			expectAlbumIds:      []catalog.AlbumId{lifetimeAlbum.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{lifetimeAlbum.AlbumId: nil},
			expectSavedCovers:   map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumFolderNameAlreadyTakenErr, i...)
			},
		},
		{
			name: "it should insert the album without transferring any media when no other album overlaps, and not touch covers",
			fields: fields{
				Catalog: NewCatalogInMemory(),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args:                args{request: standardRequest},
			expectAlbumIds:      []catalog.AlbumId{ironmanAlbum.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{ironmanAlbum.AlbumId: nil},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      ironmanAlbum,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:           assert.NoError,
		},
		{
			name: "it should move overlapping medias into the new album, carry the lifetime's CHERRY_PICKED cover across, and backfill both albums' covers from their remaining medias",
			fields: fields{
				Catalog: NewCatalogInMemory(
					withAlbum(lifetimeAlbum, photo28apr26, photo29apr26, photo30apr26, photo10jun26),
				),
				Covers: NewCoverRepositoryInMemory(coversFor(lifetimeAlbum.AlbumId,
					pickedCover(photo28apr26),
					randomCover(photo10jun26),
				)),
			},
			args:           args{request: standardRequest},
			expectAlbumIds: []catalog.AlbumId{lifetimeAlbum.AlbumId, ironmanAlbum.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				lifetimeAlbum.AlbumId: {photo10jun26},
				ironmanAlbum.AlbumId:  {photo28apr26, photo29apr26, photo30apr26},
			},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum: ironmanAlbum,
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{ironmanAlbum.AlbumId: {photo28apr26.Id, photo29apr26.Id, photo30apr26.Id}},
					FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					ironmanAlbum.AlbumId: {
						pickedCover(photo28apr26),
						randomCover(photo29apr26),
						randomCover(photo30apr26),
					},
					lifetimeAlbum.AlbumId: {randomCover(photo10jun26)},
				},
			}},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				ironmanAlbum.AlbumId: {
					pickedCover(photo28apr26),
					randomCover(photo29apr26),
					randomCover(photo30apr26),
				},
				lifetimeAlbum.AlbumId: {randomCover(photo10jun26)},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &AlbumCreatedObserverInMemory{}
			coverService := &catalog.CoverService{
				CoverRepository:     tt.fields.Covers,
				MediaReadRepository: tt.fields.Catalog,
				Randomiser:          deterministicRandomiser,
			}

			createAlbum := catalog.NewAlbumCreate(
				tt.fields.Catalog,
				&catalog.TransferMediasFromRepository{TransferMediasRepository: tt.fields.Catalog},
				coverService,
				observer,
			)

			_, err := createAlbum.Create(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.args.request)) {
				return
			}

			assert.ElementsMatch(t, tt.expectAlbumIds, tt.fields.Catalog.AlbumIds(), "stored albums")
			assert.Equal(t, tt.expectMediasByAlbum, tt.fields.Catalog.MediasByAlbum(), "medias per album")
			assert.Equal(t, tt.expectCreatedEvents, observer.Events, "AlbumCreated events fired")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.Covers.Covers, "covers persisted")
		})
	}
}
