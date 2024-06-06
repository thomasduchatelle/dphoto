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
		wantRepo []AlbumSize
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "it should return nil and do nothing if the transfer is empty",
			fields: fields{
				MediaCounterPort:              stubMediaCounterPort(nil),
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(nil),
			},
			args: args{
				transfers: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should update the album size for the owner view onluy if not shared",
			fields: fields{
				MediaCounterPort: stubMediaCounterPort(map[catalog.AlbumId]int{
					albumId1: 1,
				}),
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User)},
				}),
			},
			args: args{
				transfers: catalog.TransferredMedias{
					albumId1: []catalog.MediaId{"media1"},
				},
			},
			wantRepo: []AlbumSize{
				{
					AlbumId:    albumId1,
					Users:      []Availability{OwnerAvailability(owner1User)},
					MediaCount: 1,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should insert the album size for one album and one user",
			fields: fields{
				MediaCounterPort: stubMediaCounterPort(map[catalog.AlbumId]int{
					albumId1: 1,
				}),
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					albumId1: {OwnerAvailability(owner1User), VisitorAvailability(user2)},
				}),
			},
			args: args{
				transfers: catalog.TransferredMedias{
					albumId1: []catalog.MediaId{"media1"},
				},
			},
			wantRepo: []AlbumSize{
				{
					AlbumId:    albumId1,
					Users:      []Availability{OwnerAvailability(owner1User), VisitorAvailability(user2)},
					MediaCount: 1,
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(ViewWriteRepositoryFake)
			c := &CommandHandlerAlbumSize{
				MediaCounterPort:              tt.fields.MediaCounterPort,
				ListUserWhoCanAccessAlbumPort: tt.fields.ListUserWhoCanAccessAlbumPort,
				ViewWriteRepository:           repository,
			}
			err := c.OnTransferredMedias(context.Background(), tt.args.transfers)
			if !tt.wantErr(t, err, fmt.Sprintf("OnTransferredMedias(%v, %v)", context.Background(), tt.args.transfers)) {
				return
			}

			assert.ElementsMatchf(t, repository.AlbumSizes, tt.wantRepo, "AlbumSizes should be %v", tt.wantRepo)
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
		wantRepo []AlbumSizeDiff
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
			wantRepo: []AlbumSizeDiff{
				{
					AlbumId:        albumId1,
					Users:          []Availability{OwnerAvailability(owner1User)},
					MediaCountDiff: 1,
				},
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
			wantRepo: []AlbumSizeDiff{
				{
					AlbumId:        albumId1,
					Users:          []Availability{OwnerAvailability(owner1User), VisitorAvailability(user2)},
					MediaCountDiff: 3,
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(ViewWriteRepositoryFake)

			c := &CommandHandlerAlbumSize{
				MediaCounterPort:              stubMediaCounterPort(nil),
				ListUserWhoCanAccessAlbumPort: tt.fields.ListUserWhoCanAccessAlbumPort,
				ViewWriteRepository:           repository,
			}

			err := c.OnMediasInserted(context.Background(), tt.args.medias)
			if !tt.wantErr(t, err, fmt.Sprintf("OnMediasInserted(%v, %v)", context.Background(), tt.args.medias)) {
				return
			}

			assert.ElementsMatchf(t, repository.Updated, tt.wantRepo, "AlbumSizeDiffs should be %v", tt.wantRepo)
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

type CountKey struct {
	Availability Availability
	AlbumId      catalog.AlbumId
}

type ViewWriteRepositoryFake struct {
	AlbumSizes []AlbumSize
	Updated    []AlbumSizeDiff
	Deleted    []CountKey
}

func (v *ViewWriteRepositoryFake) InsertAlbumSize(ctx context.Context, albumSize []AlbumSize) error {
	if albumSize == nil {
		return errors.Errorf("InsertAlbumSize(nil): albumSize should not be nil")
	}
	v.AlbumSizes = append(v.AlbumSizes, albumSize...)
	return nil
}

func (v *ViewWriteRepositoryFake) UpdateAlbumSize(ctx context.Context, albumCountUpdates []AlbumSizeDiff) error {
	if albumCountUpdates == nil {
		return errors.Errorf("UpdateAlbumSize(nil): albumCountUpdates should not be nil")
	}

	v.Updated = append(v.Updated, albumCountUpdates...)
	return nil
}

func (v *ViewWriteRepositoryFake) DeleteAlbumSize(ctx context.Context, availability Availability, albumId catalog.AlbumId) error {
	v.Deleted = append(v.Deleted, CountKey{Availability: availability, AlbumId: albumId})
	return nil
}
