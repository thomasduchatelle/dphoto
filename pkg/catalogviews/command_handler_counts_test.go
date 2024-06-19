package catalogviews

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

func TestCommandHandlerAlbumSize_OnTransferredMedias(t *testing.T) {
	owner1 := ownermodel.Owner("owner-1")
	owner1User := usermodel.UserId("user-of-owner-1")
	user2 := usermodel.UserId("user-2")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album1")}
	albumId2 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album2")}

	var r MediaCounterPort = MediaCounterPortFake(nil)
	var r2 MediaCounterPort = MediaCounterPortFake(map[catalog.AlbumId]int{
		albumId1: 1,
	})
	var r3 MediaCounterPort = MediaCounterPortFake(map[catalog.AlbumId]int{
		albumId1: 1,
	})
	var r4 MediaCounterPort = MediaCounterPortFake(map[catalog.AlbumId]int{
		albumId1: 1,
	})

	type fields struct {
		MediaCounterPort              MediaCounterPort
		ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	}
	type args struct {
		transfers catalog.TransferredMedias
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRepo []UserAlbumSize
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "it should return nil and do nothing if the transfer is empty",
			fields: fields{
				MediaCounterPort:              r,
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(nil),
			},
			args: args{
				transfers: catalog.TransferredMedias{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should update the album size for the owner view only if not shared",
			fields: fields{
				MediaCounterPort: r2,
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User)},
				}),
			},
			args: args{
				transfers: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{
						albumId1: {"media1"},
					},
				},
			},
			wantRepo: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: OwnerAvailability(owner1User)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should insert the album size for one album and one user",
			fields: fields{
				MediaCounterPort: r3,
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User), VisitorAvailability(user2)},
				}),
			},
			args: args{
				transfers: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{
						albumId1: {"media1"},
					},
				},
			},
			wantRepo: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: OwnerAvailability(owner1User)},
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: VisitorAvailability(user2)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should also recount the albums in the origins",
			fields: fields{
				MediaCounterPort: r4,
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User), VisitorAvailability(user2)},
					albumId2: {OwnerAvailability(owner1User)},
				}),
			},
			args: args{
				transfers: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{
						albumId1: {"media1"},
					},
					FromAlbums: []catalog.AlbumId{albumId2},
				},
			},
			wantRepo: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: OwnerAvailability(owner1User)},
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: VisitorAvailability(user2)},
				{AlbumSize: AlbumSize{AlbumId: albumId2, MediaCount: 0}, Availability: OwnerAvailability(owner1User)},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(AlbumSizeInMemoryRepository)
			c := &CommandHandlerAlbumSize{
				MediaCounterPort:              tt.fields.MediaCounterPort,
				ListUserWhoCanAccessAlbumPort: tt.fields.ListUserWhoCanAccessAlbumPort,
				ViewWriteRepository:           repository,
			}
			err := c.OnTransferredMedias(context.Background(), tt.args.transfers)
			if !tt.wantErr(t, err, fmt.Sprintf("OnTransferredMedias(%v, %v)", context.Background(), tt.args.transfers)) {
				return
			}

			assert.ElementsMatchf(t, repository.Sizes, tt.wantRepo, "AlbumSizes should be %v", tt.wantRepo)
		})
	}
}

func TestCommandHandlerAlbumSize_OnMediasInserted(t *testing.T) {
	owner1 := ownermodel.Owner("owner-1")
	owner1User := usermodel.UserId("user-of-owner-1")
	user2 := usermodel.UserId("user-2")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album1")}

	type fields struct {
		ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	}
	type args struct {
		medias map[catalog.AlbumId][]catalog.MediaId
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRepo []UserAlbumSize
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "it shouldn't call anything if there is no media inserted",
			fields: fields{
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(nil),
			},
			args: args{
				medias: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should update the count for the owner",
			fields: fields{
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User)},
				}),
			},
			args: args{
				medias: map[catalog.AlbumId][]catalog.MediaId{
					albumId1: {"media1"},
				},
			},
			wantRepo: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: OwnerAvailability(owner1User)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should update the count for the owner and the user(s) that have access to the album",
			fields: fields{
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User), VisitorAvailability(user2)},
				}),
			},
			args: args{
				medias: map[catalog.AlbumId][]catalog.MediaId{
					albumId1: {"media1", "media2", "media3"},
				},
			},
			wantRepo: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 3}, Availability: OwnerAvailability(owner1User)},
				{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 3}, Availability: VisitorAvailability(user2)},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(AlbumSizeInMemoryRepository)

			var r MediaCounterPort = MediaCounterPortFake(nil)
			c := &CommandHandlerAlbumSize{
				MediaCounterPort:              r,
				ListUserWhoCanAccessAlbumPort: tt.fields.ListUserWhoCanAccessAlbumPort,
				ViewWriteRepository:           repository,
			}

			err := c.OnMediasInserted(context.Background(), tt.args.medias)
			if !tt.wantErr(t, err, fmt.Sprintf("OnMediasInserted(%v, %v)", context.Background(), tt.args.medias)) {
				return
			}

			assert.ElementsMatchf(t, repository.Sizes, tt.wantRepo, "AlbumSizeDiffs should be %v", tt.wantRepo)
		})
	}
}

type ListUserWhoCanAccessAlbumPortFake struct {
	Values map[catalog.AlbumId][]Availability
}

func (l *ListUserWhoCanAccessAlbumPortFake) ListUsersWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]Availability, error) {
	if albumId == nil {
		return nil, errors.Errorf("ListUsersWhoCanAccessAlbum(nil): albumId should not be nil")
	}

	result := make(map[catalog.AlbumId][]Availability)
	for _, id := range albumId {
		result[id] = l.Values[id]
	}
	return result, nil
}

func stubListUserWhoCanAccessAlbumPort(values map[catalog.AlbumId][]Availability) ListUserWhoCanAccessAlbumPort {
	return &ListUserWhoCanAccessAlbumPortFake{Values: values}
}
