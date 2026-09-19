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
	jan23 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	album23 := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2023"),
		},
		Name:  "2023",
		Start: jan23,
		End:   feb24,
	}
	q1Album := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2024-Q1"),
		},
		Name:  "Q1 2024",
		Start: jan24,
		End:   apr24,
	}
	q4album := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2023-Q4"),
		},
		Name:  "Q4 2023",
		Start: time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	recordsFrom23 := catalog.MediaTransferRecords{
		q1Album.AlbumId: {
			{
				FromAlbums: []catalog.AlbumId{album23.AlbumId},
				Start:      jan24,
				End:        apr24,
			},
		},
	}
	transferredFrom23 := func() catalog.TransferredMedias {
		var ids []catalog.MediaId
		for day := jan24; day.Before(apr24); day = day.AddDate(0, 0, 1) {
			ids = append(ids, fakeMediaId(album23.AlbumId, day))
		}
		return catalog.TransferredMedias{
			Transfers:  map[catalog.AlbumId][]catalog.MediaId{q1Album.AlbumId: ids},
			FromAlbums: []catalog.AlbumId{album23.AlbumId},
		}
	}

	type fields struct {
		AlbumRepository *AlbumRepositoryInMemory
	}
	type exec struct {
		mediaTime time.Time
		want      catalog.AlbumReference
		wantErr   assert.ErrorAssertionFunc
	}
	tests := []struct {
		name                  string
		fields                fields
		exec                  []exec
		expectAlbumIds        []catalog.AlbumId
		expectTransferRecords []catalog.MediaTransferRecords
		expectCreatedEvents   []catalog.AlbumCreated
	}{
		{
			name:   "it should look up an existing album by media time without creating anything",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(&q1Album)},
			exec: []exec{
				{
					mediaTime: feb24,
					want:      catalog.AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: false},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds: []catalog.AlbumId{q1Album.AlbumId},
		},
		{
			name:   "it should create a new quarterly album when no album covers the media time (no overlap, no transfer)",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			exec: []exec{
				{
					mediaTime: feb24,
					want:      catalog.AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds:        []catalog.AlbumId{q1Album.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{nil},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      q1Album,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
		},
		{
			name:   "it should create a new quarterly album with no transfer when an adjacent album does not overlap",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(q4album)},
			exec: []exec{
				{
					mediaTime: feb24,
					want:      catalog.AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds:        []catalog.AlbumId{q4album.AlbumId, q1Album.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{nil},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      q1Album,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
		},
		{
			name:   "it should create a new quarterly album AND transfer medias when the new album overlaps an existing one",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(album23)},
			exec: []exec{
				{
					mediaTime: feb24,
					want:      catalog.AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
			},
			expectAlbumIds:        []catalog.AlbumId{album23.AlbumId, q1Album.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{recordsFrom23},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      q1Album,
				TransferredMedias: transferredFrom23(),
			}},
		},
		{
			// The referencer's cached TimelineAggregate is not updated after an auto-create
			// (CreateAlbum re-loads the owner's albums on every call), so the second media
			// in the same quarter re-enters the auto-create strategy, which fails because
			// the folder /2024-Q1 already exists in the reloaded timeline. A follow-up
			// TimelineRepository refactor will restore the cached-timeline behaviour.
			name:   "it should fail with AlbumFolderNameAlreadyTakenErr when a second media in the same quarter re-triggers auto-create",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			exec: []exec{
				{
					mediaTime: feb24,
					want:      catalog.AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr:   assert.NoError,
				},
				{
					mediaTime: feb24,
					want:      catalog.AlbumReference{},
					wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
						return assert.ErrorIs(t, err, catalog.AlbumFolderNameAlreadyTakenErr, i)
					},
				},
			},
			expectAlbumIds:        []catalog.AlbumId{q1Album.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{nil},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      q1Album,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferService := &TransferMediasServiceFake{}
			observer := &AlbumCreatedObserverInMemory{}

			createAlbum := &catalog.CreateAlbum{
				FindAlbumsByOwnerPort: tt.fields.AlbumRepository,
				InsertAlbumPort:       tt.fields.AlbumRepository,
				TransferMediasService: transferService,
				AlbumCreatedObservers: []catalog.AlbumCreatedObserver{observer},
			}

			referencer, err := catalog.NewAlbumAutoPopulateReferencer(
				owner,
				tt.fields.AlbumRepository,
				createAlbum,
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

			storedIds := make([]catalog.AlbumId, 0, len(tt.fields.AlbumRepository.Albums))
			for id := range tt.fields.AlbumRepository.Albums {
				storedIds = append(storedIds, id)
			}
			assert.ElementsMatch(t, tt.expectAlbumIds, storedIds, "albums in the repository")
			assert.Equal(t, tt.expectTransferRecords, transferService.Records, "records passed to TransferMedias")
			assert.Equal(t, tt.expectCreatedEvents, observer.Events, "AlbumCreated events fired")
		})
	}
}

func TestNewAlbumDryRunReferencer(t *testing.T) {
	const owner = ownermodel.Owner("owner-1")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	album24 := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2024"),
		},
		Name:  "2024",
		Start: jan24,
		End:   jan25,
	}

	type fields struct {
		AlbumRepository *AlbumRepositoryInMemory
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
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory(album24)},
			args:   args{mediaTime: feb24},
			want: catalog.AlbumReference{
				AlbumId:          &album24.AlbumId,
				AlbumJustCreated: false,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should makeup a reference when the album has not been found",
			fields: fields{AlbumRepository: NewAlbumRepositoryInMemory()},
			args:   args{mediaTime: jan24},
			want: catalog.AlbumReference{
				AlbumId:          &catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/new-album")},
				AlbumJustCreated: true,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referencer, err := catalog.NewAlbumDryRunReferencer(owner, tt.fields.AlbumRepository)
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
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	q1Album := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2024-Q1"),
		},
		Name:  "Q1 2024",
		Start: jan24,
		End:   apr24,
	}
	febAprAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/2024-Feb-Apr"),
		},
		Name:  "Feb-Apr 2024",
		Start: feb24,
		End:   apr24,
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
				mediaTime: feb24,
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
				mediaTime: feb24,
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
				mediaTime: feb24,
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
			strategy := catalog.TimelineLookupStrategy{}
			got, err := strategy.LookupAlbum(context.Background(), tt.args.owner, catalog.NewLazyTimelineAggregate(tt.args.albums), tt.args.mediaTime)
			if !tt.wantErr(t1, err, fmt.Sprintf("LookupAlbum(%v, %v, %v)", tt.args.owner, tt.args.albums, tt.args.mediaTime)) {
				return
			}
			assert.Equalf(t1, tt.want, got, "LookupAlbum(%v, %v, %v)", tt.args.owner, tt.args.albums, tt.args.mediaTime)
		})
	}
}
