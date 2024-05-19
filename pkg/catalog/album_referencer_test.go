package catalog

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"testing"
	"time"
)

func TestNewAlbumAutoPopulateReferencerAcceptance(t *testing.T) {
	owner := ownermodel.Owner("owner-1")
	jan23 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	album23 := Album{
		AlbumId: AlbumId{
			Owner:      owner,
			FolderName: NewFolderName("/2023"),
		},
		Name:  "2023",
		Start: jan23,
		End:   feb24,
	}
	q1Album := Album{
		AlbumId: AlbumId{
			Owner:      owner,
			FolderName: NewFolderName("/2024-Q1"),
		},
		Name:  "Q1 2024",
		Start: jan24,
		End:   apr24,
	}
	q4album := Album{
		AlbumId: AlbumId{
			Owner:      owner,
			FolderName: NewFolderName("/2023-Q4"),
		},
		Name:  "Q4 2023",
		Start: time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	recordsFrom23 := MediaTransferRecords{
		q1Album.AlbumId: {
			{
				FromAlbums: []AlbumId{album23.AlbumId},
				Start:      jan24,
				End:        apr24,
			},
		},
	}
	transferredMediasFrom23 := TransferredMedias{
		q1Album.AlbumId: {MediaId("media-1"), MediaId("media-2")},
	}
	type fields struct {
		owner                     ownermodel.Owner
		findAlbumsByOwner         FindAlbumsByOwnerPort
		transferMediasPort        *TransferMediasRepositoryPortFake
		timelineMutationObservers []TimelineMutationObserver
	}
	type args struct {
		mediaTime time.Time
	}
	type exec struct {
		args    args
		want    AlbumReference
		wantErr assert.ErrorAssertionFunc
	}
	tests := []struct {
		name         string
		fields       fields
		exec         []exec
		wantInserted []*Album
		wantObserved []TransferredMedias
	}{
		{
			name: "it should create a new album (complete journey) but without transferring albums because no overlap with other albums",
			fields: fields{
				owner:              owner,
				findAlbumsByOwner:  make(FindAlbumsByOwnerPortFake),
				transferMediasPort: new(TransferMediasRepositoryPortFake),
			},
			exec: []exec{
				{
					args: args{
						mediaTime: feb24,
					},
					want:    AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr: assert.NoError,
				},
			},
			wantInserted: []*Album{&q1Album},
			wantObserved: nil,
		},
		{
			name: "it should create a new album without transferring medias to it",
			fields: fields{
				owner:              owner,
				findAlbumsByOwner:  FindAlbumsByOwnerPortFake{owner: []*Album{&q4album}},
				transferMediasPort: new(TransferMediasRepositoryPortFake),
			},
			exec: []exec{
				{
					args: args{
						mediaTime: feb24,
					},
					want:    AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr: assert.NoError,
				},
			},
			wantInserted: []*Album{&q1Album},
			wantObserved: nil,
		},
		{
			name: "it should create a new album and transfer medias to it",
			fields: fields{
				owner:             owner,
				findAlbumsByOwner: FindAlbumsByOwnerPortFake{owner: []*Album{&album23}},
				transferMediasPort: &TransferMediasRepositoryPortFake{
					ExpectedRecords:   recordsFrom23,
					TransferredMedias: transferredMediasFrom23,
				},
			},
			exec: []exec{
				{
					args: args{
						mediaTime: feb24,
					},
					want:    AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr: assert.NoError,
				},
			},
			wantInserted: []*Album{&q1Album},
			wantObserved: []TransferredMedias{transferredMediasFrom23},
		},
		{
			name: "it should create the album on the first call, and find it on the second call (without requesting list of albums again)",
			fields: fields{
				owner:              owner,
				findAlbumsByOwner:  make(FindAlbumsByOwnerPortFake),
				transferMediasPort: new(TransferMediasRepositoryPortFake),
			},
			exec: []exec{
				{
					args: args{
						mediaTime: feb24,
					},
					want:    AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: true},
					wantErr: assert.NoError,
				},
				{
					args: args{
						mediaTime: feb24,
					},
					want:    AlbumReference{AlbumId: &q1Album.AlbumId, AlbumJustCreated: false},
					wantErr: assert.NoError,
				},
			},
			wantInserted: []*Album{&q1Album},
			wantObserved: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insertAlbumPortFake := new(InsertAlbumPortFake)
			observer := new(TimelineMutationObserverFake)

			referencer, err := NewAlbumAutoPopulateReferencer(
				tt.fields.owner,
				tt.fields.findAlbumsByOwner,
				insertAlbumPortFake,
				tt.fields.transferMediasPort,
				observer,
			)
			if !assert.NoError(t, err) {
				return
			}

			for _, ex := range tt.exec {
				got, err := referencer.FindReference(context.Background(), ex.args.mediaTime)
				if !ex.wantErr(t, err) {
					return
				}
				assert.Equal(t, ex.want, got, "FindReference(%v, %v)", context.Background(), ex.args.mediaTime)
			}

			assert.Equal(t, tt.wantInserted, insertAlbumPortFake.Albums)
			assert.Equal(t, tt.fields.transferMediasPort.ExpectedRecords, tt.fields.transferMediasPort.GotRecords)
			assert.Equal(t, tt.wantObserved, observer.Requests)
		})
	}
}

func TestNewAlbumDryRunReferencer(t *testing.T) {
	owner := ownermodel.Owner("owner-1")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	album24 := Album{
		AlbumId: AlbumId{
			Owner:      owner,
			FolderName: NewFolderName("/2024"),
		},
		Name:  "2024",
		Start: jan24,
		End:   jan25,
	}

	type fields struct {
		owner             ownermodel.Owner
		findAlbumsByOwner FindAlbumsByOwnerPort
	}
	type args struct {
		mediaTime time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    AlbumReference
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return a reference for an album that has been found",
			fields: fields{
				owner:             owner,
				findAlbumsByOwner: FindAlbumsByOwnerPortFake{owner: []*Album{&album24}},
			},
			args: args{
				mediaTime: feb24,
			},
			want: AlbumReference{
				AlbumId:          &album24.AlbumId,
				AlbumJustCreated: false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should makeup a reference when the album has not been found",
			fields: fields{
				owner:             owner,
				findAlbumsByOwner: make(FindAlbumsByOwnerPortFake),
			},
			args: args{
				mediaTime: jan24,
			},
			want: AlbumReference{
				AlbumId:          &AlbumId{Owner: owner, FolderName: NewFolderName("/new-album")},
				AlbumJustCreated: true,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referencer, err := NewAlbumDryRunReferencer(tt.fields.owner, tt.fields.findAlbumsByOwner)
			if !assert.NoError(t, err) {
				return
			}

			got, err := referencer.FindReference(context.Background(), tt.args.mediaTime)
			if tt.wantErr(t, err) {
				assert.Equalf(t, tt.want, got, "FindReference(%v, %v)", context.Background(), tt.args.mediaTime)
			}
		})
	}
}

func TestTimelineLookupStrategy_LookupAlbum(t1 *testing.T) {
	owner := ownermodel.Owner("owner-1")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	q1Album := Album{
		AlbumId: AlbumId{
			Owner:      owner,
			FolderName: NewFolderName("/2024-Q1"),
		},
		Name:  "Q1 2024",
		Start: jan24,
		End:   apr24,
	}
	febAprAlbum := Album{
		AlbumId: AlbumId{
			Owner:      owner,
			FolderName: NewFolderName("/2024-Feb-Apr"),
		},
		Name:  "Feb-Apr 2024",
		Start: feb24,
		End:   apr24,
	}

	type args struct {
		owner     ownermodel.Owner
		albums    []*Album
		mediaTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    AlbumReference
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should find an album id that exists in a timelines",
			args: args{
				owner:     owner,
				albums:    []*Album{&q1Album},
				mediaTime: feb24,
			},
			want: AlbumReference{
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
			want: AlbumReference{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, NoAlbumLookedUpError, i)
			},
		},
		{
			name: "it should pick the album with highest priority",
			args: args{
				owner: owner,
				albums: []*Album{
					&febAprAlbum,
					&q1Album,
				},
				mediaTime: feb24,
			},
			want: AlbumReference{
				AlbumId:          &febAprAlbum.AlbumId,
				AlbumJustCreated: false,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TimelineLookupStrategy{}
			got, err := t.LookupAlbum(context.Background(), tt.args.owner, NewLazyTimelineAggregate(tt.args.albums), tt.args.mediaTime)
			if !tt.wantErr(t1, err, fmt.Sprintf("LookupAlbum(%v, %v, %v, %v)", context.Background(), tt.args.owner, tt.args.albums, tt.args.mediaTime)) {
				return
			}
			assert.Equalf(t1, tt.want, got, "LookupAlbum(%v, %v, %v, %v)", context.Background(), tt.args.owner, tt.args.albums, tt.args.mediaTime)
		})
	}
}

func TestAlbumAutoCreateLookupStrategy_LookupAlbum(t *testing.T) {
	owner := ownermodel.Owner("owner-1")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	type args struct {
		owner     ownermodel.Owner
		mediaTime time.Time
	}
	tests := []struct {
		name             string
		args             args
		want             AlbumReference
		wantCreateAlbums []CreateAlbumRequest
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "it should initiate an album creation for a quarter",
			args: args{
				owner:     owner,
				mediaTime: feb24,
			},
			want: AlbumReference{
				AlbumId:          &AlbumId{Owner: owner, FolderName: NewFolderName("/2024-Q1")},
				AlbumJustCreated: true,
			},
			wantCreateAlbums: []CreateAlbumRequest{
				{
					Owner:            owner,
					Name:             "Q1 2024",
					Start:            jan24,
					End:              apr24,
					ForcedFolderName: "/2024-Q1",
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := new(CreateAlbumWithTimelineFake)
			a := &AlbumAutoCreateLookupStrategy{
				Delegate: delegate,
			}

			timeline := NewLazyTimelineAggregate(nil)
			got, err := a.LookupAlbum(context.Background(), tt.args.owner, timeline, tt.args.mediaTime)
			if !tt.wantErr(t, err, fmt.Sprintf("LookupAlbum(%v, %v, %v, %v)", context.Background(), tt.args.owner, timeline, tt.args.mediaTime)) {
				return
			}

			assert.Equalf(t, tt.want, got, "LookupAlbum(%v, %v, %v, %v)", context.Background(), tt.args.owner, timeline, tt.args.mediaTime)
		})
	}
}

type CreateAlbumWithTimelineFake struct {
	Requests []CreateAlbumRequest
}

func (a *CreateAlbumWithTimelineFake) Create(ctx context.Context, timeline *TimelineAggregate, request CreateAlbumRequest) (*AlbumId, error) {
	a.Requests = append(a.Requests, request)
	return &AlbumId{Owner: request.Owner, FolderName: NewFolderName(request.ForcedFolderName)}, nil
}

type FindAlbumsByOwnerPortFake map[ownermodel.Owner][]*Album

func (f FindAlbumsByOwnerPortFake) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error) {
	albums, _ := f[owner]
	return albums, nil
}

type InsertAlbumPortFake struct {
	Albums []*Album
}

func (i *InsertAlbumPortFake) InsertAlbum(ctx context.Context, album Album) error {
	i.Albums = append(i.Albums, &album)
	return nil
}

type TransferMediasRepositoryPortFake struct {
	ExpectedRecords   MediaTransferRecords
	GotRecords        MediaTransferRecords
	TransferredMedias TransferredMedias
}

func (t *TransferMediasRepositoryPortFake) TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error) {
	t.GotRecords = records
	return t.TransferredMedias, nil
}

type TimelineMutationObserverFake struct {
	Requests []TransferredMedias
}

func (t *TimelineMutationObserverFake) OnTransferredMedias(ctx context.Context, transfers TransferredMedias) error {
	t.Requests = append(t.Requests, transfers)
	return nil
}
