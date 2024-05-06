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

func TestCreateAlbum_Create(t *testing.T) {
	const owner = "tonystark"
	album := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.FolderName("/2024-04_Ironman_1"),
		},
		Name:  "Ironman 1",
		Start: time.Date(2024, 04, 28, 8, 33, 42, 0, time.UTC),
		End:   time.Date(2024, 05, 1, 0, 0, 0, 0, time.UTC),
	}
	lifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/lifetime"),
		},
		Name:  "lifetime",
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	remainingLifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/remaining-lifetime"),
		},
		Name:  "remaining-lifetime",
		Start: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	highPriorityAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.NewFolderName("/high-priority"),
		},
		Name:  "remaining-lifetime",
		Start: time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
	}

	standardRequest := catalog.CreateAlbumRequest{
		Owner: owner,
		Name:  album.Name,
		Start: album.Start,
		End:   album.End,
	}

	type fields struct {
		FindAlbumsByOwnerPort func(t *testing.T) catalog.FindAlbumsByOwnerPort
		InsertAlbumPort       func(t *testing.T) catalog.InsertAlbumPort
	}
	type args struct {
		request catalog.CreateAlbumRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    catalog.MediaTransferRecords
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the album with a generated name",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				InsertAlbumPort:       expectAlbumInserted(album),
			},
			args: args{
				request: standardRequest,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should create the album with a forced name",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner),
				InsertAlbumPort: expectAlbumInserted(catalog.Album{
					AlbumId: catalog.AlbumId{
						Owner:      owner,
						FolderName: "/Phase_1_Avenger",
					},
					Name:  "Avenger 1",
					Start: album.Start,
					End:   album.End,
				}),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             "Avenger 1",
					Start:            album.Start,
					End:              album.End,
					ForcedFolderName: "Phase_1_Avenger",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from a lower priority album",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner, lifetimeAlbum),
				InsertAlbumPort:       expectAlbumInserted(album),
			},
			args: args{
				request: standardRequest,
			},
			want: catalog.MediaTransferRecords{
				album.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
						Start:      album.Start,
						End:        album.End,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 2 lower priority albums ; selector still in one single block",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner, lifetimeAlbum, remainingLifetimeAlbum),
				InsertAlbumPort:       expectAlbumInserted(album),
			},
			args: args{
				request: standardRequest,
			},
			want: catalog.MediaTransferRecords{
				album.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{remainingLifetimeAlbum.AlbumId, lifetimeAlbum.AlbumId},
						Start:      album.Start,
						End:        album.End,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-allocate medias from 1 lower priority albums, avoiding 1 high priority (selectors in two blocks)",
			fields: fields{
				FindAlbumsByOwnerPort: returnListOfAlbums(owner, lifetimeAlbum, highPriorityAlbum),
				InsertAlbumPort:       expectAlbumInserted(album),
			},
			args: args{
				request: standardRequest,
			},
			want: catalog.MediaTransferRecords{
				album.AlbumId: {
					{
						FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
						Start:      album.Start,
						End:        highPriorityAlbum.Start,
					},
					{
						FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
						Start:      highPriorityAlbum.End,
						End:        album.End,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := new(TimelineMutationCollector)
			c := &catalog.CreateAlbum{
				FindAlbumsByOwnerPort:     tt.fields.FindAlbumsByOwnerPort(t),
				InsertAlbumPort:           tt.fields.InsertAlbumPort(t),
				TransferMediasPort:        collector,
				TimelineMutationObservers: nil,
			}

			err := c.Create(context.TODO(), tt.args.request)
			if tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.args.request)) && err == nil {
				assert.Equal(t, tt.want, collector.Records, fmt.Sprintf("Create(%v)", tt.args.request))
			}
		})
	}
}

type TimelineMutationCollector struct {
	T       *testing.T
	Records catalog.MediaTransferRecords
}

func (t *TimelineMutationCollector) TransferMediasFromRecords(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
	if records == nil {
		assert.Fail(t.T, "TransferMediasFromRecords(nil): TransferMediasFromRecords should not be called with nil value.")
	}
	t.Records = records
	return nil, nil
}

func returnListOfAlbums(expectedOwner catalog.Owner, albums ...*catalog.Album) func(t *testing.T) catalog.FindAlbumsByOwnerPort {
	return func(t *testing.T) catalog.FindAlbumsByOwnerPort {
		return catalog.FindAlbumsByOwnerFunc(func(ctx context.Context, owner catalog.Owner) ([]*catalog.Album, error) {
			if owner == expectedOwner && len(albums) > 0 {
				return albums, nil
			}

			return nil, nil
		})
	}
}

func expectAlbumInserted(album catalog.Album) func(t *testing.T) catalog.InsertAlbumPort {
	return func(t *testing.T) catalog.InsertAlbumPort {
		adapter := mocks.NewInsertAlbumPort(t)
		adapter.On("InsertAlbum", mock.Anything, album).Return(nil)
		return adapter
	}
}
