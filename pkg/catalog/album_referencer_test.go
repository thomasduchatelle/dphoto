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

func TestNewAlbumAutoPopulateReferencer(t *testing.T) {
	const owner = ownermodel.Owner("owner-1")
	jan26 := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb26 := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr26 := time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)
	oct25 := time.Date(2025, time.October, 1, 0, 0, 0, 0, time.UTC)
	jan23 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	q1_26Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/2026-Q1")},
		Name:    "Q1 2026",
		Start:   jan26,
		End:     apr26,
	}
	q4_25Album := &catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/2025-Q4")},
		Name:    "Q4 2025",
		Start:   oct25,
		End:     jan26,
	}
	wideAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/wide")},
		Name:    "wide",
		Start:   jan23,
		End:     feb26,
	}

	photo5feb26 := photoAt("photo5feb26", feb26.AddDate(0, 0, 5))

	type fields struct {
		Catalog *CatalogInMemory
	}
	type exec struct {
		mediaTime time.Time
		want      catalog.AlbumReference
		wantErr   assert.ErrorAssertionFunc
	}
	tests := []struct {
		name                string
		fields              fields
		exec                []exec
		expectAlbumIds      []catalog.AlbumId
		expectMediasByAlbum map[catalog.AlbumId][]*catalog.MediaMeta
		expectCreatedEvents []catalog.AlbumCreated
	}{
		{
			name:   "it should look up an existing album by media time without creating anything",
			fields: fields{Catalog: NewCatalogInMemory(withAlbum(&q1_26Album))},
			exec: []exec{
				{
					mediaTime: feb26,
					want:      catalog.AlbumReference{AlbumId: &q1_26Album.AlbumId, AlbumJustCreated: false},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds:      []catalog.AlbumId{q1_26Album.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{q1_26Album.AlbumId: nil},
		},
		{
			name:   "it should create a new quarterly album with no transfer when an adjacent album does not overlap",
			fields: fields{Catalog: NewCatalogInMemory(withAlbum(q4_25Album))},
			exec: []exec{
				{
					mediaTime: feb26,
					want:      catalog.AlbumReference{AlbumId: &q1_26Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds: []catalog.AlbumId{q4_25Album.AlbumId, q1_26Album.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				q4_25Album.AlbumId: nil,
				q1_26Album.AlbumId: nil,
			},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      q1_26Album,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
		},
		{
			name:   "it should create a new quarterly album AND transfer overlapping medias when the new album overlaps an existing one",
			fields: fields{Catalog: NewCatalogInMemory(withAlbum(wideAlbum, photo5feb26))},
			exec: []exec{
				{
					mediaTime: feb26,
					want:      catalog.AlbumReference{AlbumId: &q1_26Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds: []catalog.AlbumId{wideAlbum.AlbumId, q1_26Album.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{
				wideAlbum.AlbumId:  nil,
				q1_26Album.AlbumId: {photo5feb26},
			},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum: q1_26Album,
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{q1_26Album.AlbumId: {photo5feb26.Id}},
					FromAlbums: []catalog.AlbumId{wideAlbum.AlbumId},
				},
				Covers: map[catalog.AlbumId][]catalog.Cover{
					q1_26Album.AlbumId: {randomCover(photo5feb26)},
				},
			}},
		},
		{
			name:   "it should reuse the auto-created album for a second media in the same quarter via the cached timeline",
			fields: fields{Catalog: NewCatalogInMemory()},
			exec: []exec{
				{
					mediaTime: feb26,
					want:      catalog.AlbumReference{AlbumId: &q1_26Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
				{
					mediaTime: feb26,
					want:      catalog.AlbumReference{AlbumId: &q1_26Album.AlbumId, AlbumJustCreated: false},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds:      []catalog.AlbumId{q1_26Album.AlbumId},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{q1_26Album.AlbumId: nil},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      q1_26Album,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &AlbumCreatedObserverInMemory{}
			coverService := &catalog.CoverService{
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: tt.fields.Catalog,
				Randomiser:          deterministicRandomiser,
			}

			referencer, err := catalog.NewAlbumAutoPopulateReferencer(
				owner,
				tt.fields.Catalog,
				&catalog.TransferMediasFromRepository{TransferMediasRepository: tt.fields.Catalog},
				coverService,
				observer,
			)
			if !assert.NoError(t, err) {
				return
			}

			for _, ex := range tt.exec {
				got, err := referencer.FindReference(context.Background(), ex.mediaTime)
				if !ex.wantErr(t, err, "FindReference(%v)", ex.mediaTime) {
					return
				}
				assert.Equal(t, ex.want, got, "FindReference(%v)", ex.mediaTime)
			}

			assert.ElementsMatch(t, tt.expectAlbumIds, tt.fields.Catalog.AlbumIds(), "albums in the catalog")
			assert.Equal(t, tt.expectMediasByAlbum, tt.fields.Catalog.MediasByAlbum(), "medias per album")
			assert.Equal(t, tt.expectCreatedEvents, observer.Events, "AlbumCreated events fired")
		})
	}
}

func TestNewAlbumDryRunReferencer(t *testing.T) {
	const owner = ownermodel.Owner("owner-1")
	jan26 := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan27 := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb26 := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	album26 := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2026"),
		},
		Name:  "2026",
		Start: jan26,
		End:   jan27,
	}

	type fields struct {
		Catalog *CatalogInMemory
	}
	type args struct {
		mediaTime time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    catalog.AlbumReference
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "it should return a reference for an album that has been found",
			fields: fields{Catalog: NewCatalogInMemory(withAlbum(album26))},
			args:   args{mediaTime: feb26},
			want: catalog.AlbumReference{
				AlbumId:          &album26.AlbumId,
				AlbumJustCreated: false,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should makeup a reference when the album has not been found",
			fields: fields{Catalog: NewCatalogInMemory()},
			args:   args{mediaTime: jan26},
			want: catalog.AlbumReference{
				AlbumId:          &catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/new-album")},
				AlbumJustCreated: true,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referencer, err := catalog.NewAlbumDryRunReferencer(owner, tt.fields.Catalog)
			if !assert.NoError(t, err) {
				return
			}

			got, err := referencer.FindReference(context.Background(), tt.args.mediaTime)
			if tt.wantErr(t, err) {
				assert.Equalf(t, tt.want, got, "FindReference(%v)", tt.args.mediaTime)
			}
		})
	}
}

func TestTimelineLookupStrategy_LookupAlbum(t1 *testing.T) {
	const owner = ownermodel.Owner("owner-1")
	jan26 := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb26 := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr26 := time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)
	q1Album := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2026-Q1"),
		},
		Name:  "Q1 2026",
		Start: jan26,
		End:   apr26,
	}
	febAprAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2026-Feb-Apr"),
		},
		Name:  "Feb-Apr 2026",
		Start: feb26,
		End:   apr26,
	}

	type args struct {
		owner     ownermodel.Owner
		albums    []*catalog.Album
		mediaTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    catalog.AlbumReference
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should find an album id that exists in a timelines",
			args: args{
				owner:     owner,
				albums:    []*catalog.Album{&q1Album},
				mediaTime: feb26,
			},
			want: catalog.AlbumReference{
				AlbumId:          &q1Album.AlbumId,
				AlbumJustCreated: false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should pass to album creation when no album fits the time",
			args: args{
				owner:     owner,
				albums:    nil,
				mediaTime: feb26,
			},
			want: catalog.AlbumReference{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.NoAlbumLookedUpError, i)
			},
		},
		{
			name: "it should pick the album with highest priority",
			args: args{
				owner: owner,
				albums: []*catalog.Album{
					&febAprAlbum,
					&q1Album,
				},
				mediaTime: feb26,
			},
			want: catalog.AlbumReference{
				AlbumId:          &febAprAlbum.AlbumId,
				AlbumJustCreated: false,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			timeline, err := catalog.NewTimelineAggregate(tt.args.albums)
			if !assert.NoError(t1, err) {
				return
			}
			strategy := catalog.TimelineLookupStrategy{}
			got, err := strategy.LookupAlbum(context.Background(), tt.args.owner, timeline, tt.args.mediaTime)
			if !tt.wantErr(t1, err, fmt.Sprintf("LookupAlbum(%v, %v, %v)", tt.args.owner, tt.args.albums, tt.args.mediaTime)) {
				return
			}
			assert.Equalf(t1, tt.want, got, "LookupAlbum(%v, %v, %v)", tt.args.owner, tt.args.albums, tt.args.mediaTime)
		})
	}
}
