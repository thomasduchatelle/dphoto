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

func TestAlbumAutoPopulateReferencer_FindReference(t *testing.T) {
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
	newAlbumId := AlbumId{
		Owner:      owner,
		FolderName: NewFolderName("/new-album-q1"),
	}

	type fields struct {
		owner               ownermodel.Owner
		timelineAggregate   *TimelineAggregate
		AlbumCreateHandover *AlbumCreateHandoverFake
	}
	type args struct {
		mediaTime time.Time
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		want             AlbumReference
		wantCreateAlbums []CreateAlbumRequest
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "it should find an album id that exists in a timelines",
			fields: fields{
				owner:               owner,
				timelineAggregate:   NewLazyTimelineAggregate([]*Album{&q1Album}),
				AlbumCreateHandover: new(AlbumCreateHandoverFake),
			},
			args: args{
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
			fields: fields{
				owner:             owner,
				timelineAggregate: NewLazyTimelineAggregate(nil),
				AlbumCreateHandover: &AlbumCreateHandoverFake{
					AlbumId: &newAlbumId,
				},
			},
			args: args{
				mediaTime: feb24,
			},
			want: AlbumReference{
				AlbumId:          &newAlbumId,
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
			a := &AutoCreateAlbumReferencer{
				owner:             tt.fields.owner,
				timelineAggregate: tt.fields.timelineAggregate,
				Observer:          tt.fields.AlbumCreateHandover,
			}
			got, err := a.FindReference(context.Background(), tt.args.mediaTime)
			if !tt.wantErr(t, err, fmt.Sprintf("FindReference(%v, %v)", context.Background(), tt.args.mediaTime)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindReference(%v, %v)", context.Background(), tt.args.mediaTime)
			assert.Equalf(t, tt.wantCreateAlbums, tt.fields.AlbumCreateHandover.Requests, "FindReference(%v, %v)", context.Background(), tt.args.mediaTime)
		})
	}
}

type AlbumCreateHandoverFake struct {
	AlbumId  *AlbumId
	Requests []CreateAlbumRequest
}

func (a *AlbumCreateHandoverFake) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	a.Requests = append(a.Requests, request)
	return a.AlbumId, nil
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

func (t *TimelineMutationObserverFake) Observe(ctx context.Context, transfers TransferredMedias) error {
	t.Requests = append(t.Requests, transfers)
	return nil
}
