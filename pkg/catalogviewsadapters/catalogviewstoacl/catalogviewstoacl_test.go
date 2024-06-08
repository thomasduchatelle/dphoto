package catalogviewstoacl

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

func TestCatalogToACLAdapter_ListUsersWhoCanAccessAlbum(t *testing.T) {
	ctx := context.Background()
	owner1 := ownermodel.Owner("owner-1")
	owner2 := ownermodel.Owner("owner-2")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album1")}
	albumId2 := catalog.AlbumId{Owner: owner2, FolderName: catalog.NewFolderName("/album2")}
	albumId3 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album3")}
	user1 := usermodel.UserId("user-1")
	user2 := usermodel.UserId("user-2")
	user3 := usermodel.UserId("user-3")

	type fields struct {
		ScopeRepository ScopeReadRepositoryPort
	}
	type args struct {
		albumIds []catalog.AlbumId
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[catalog.AlbumId][]catalogviews.Availability
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should not return any availability if no permission is present",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1},
			},
			want:    make(map[catalog.AlbumId][]catalogviews.Availability),
			wantErr: assert.NoError,
		},
		{
			name: "it should not return any availability if no permission for the resource is present",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						ownerPermission(user1, owner2),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1},
			},
			want:    make(map[catalog.AlbumId][]catalogviews.Availability),
			wantErr: assert.NoError,
		},
		{
			name: "it should find an availability for the user who's the owner of the album",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						ownerPermission(user1, owner1),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1},
			},
			want: map[catalog.AlbumId][]catalogviews.Availability{
				albumId1: {
					catalogviews.OwnerAvailability(user1),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find an availability for the user who has a visitor permission on the album",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						visitorAlbumPermission(user1, albumId1),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1},
			},
			want: map[catalog.AlbumId][]catalogviews.Availability{
				albumId1: {
					catalogviews.VisitorAvailability(user1),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find an availability for the user who has a contributor permission on the album",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						contributorAlbumPermission(user1, albumId1),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1},
			},
			want: map[catalog.AlbumId][]catalogviews.Availability{
				albumId1: {
					catalogviews.VisitorAvailability(user1),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find an availability for the users who has a visitor permission on the album or is the owner",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						visitorAlbumPermission(user1, albumId1),
						ownerPermission(user2, owner1),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1},
			},
			want: map[catalog.AlbumId][]catalogviews.Availability{
				albumId1: {
					catalogviews.VisitorAvailability(user1),
					catalogviews.OwnerAvailability(user2),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find an availabilities for several albums owned or shared with other users",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						ownerPermission(user1, owner1),
						ownerPermission(user2, owner2),
						visitorAlbumPermission(user1, albumId2),
						visitorAlbumPermission(user2, albumId1),
						visitorAlbumPermission(user3, albumId1),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1, albumId2},
			},
			want: map[catalog.AlbumId][]catalogviews.Availability{
				albumId1: {
					catalogviews.OwnerAvailability(user1),
					catalogviews.VisitorAvailability(user2),
					catalogviews.VisitorAvailability(user3),
				},
				albumId2: {
					catalogviews.OwnerAvailability(user2),
					catalogviews.VisitorAvailability(user1),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find availabilities for 2 albums owned by the same user",
			fields: fields{
				ScopeRepository: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						ownerPermission(user1, owner1),
					},
				},
			},
			args: args{
				albumIds: []catalog.AlbumId{albumId1, albumId3},
			},
			want: map[catalog.AlbumId][]catalogviews.Availability{
				albumId1: {
					catalogviews.OwnerAvailability(user1),
				},
				albumId3: {
					catalogviews.OwnerAvailability(user1),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CatalogToACLAdapter{
				ScopeRepository: tt.fields.ScopeRepository,
			}
			got, err := f.ListUsersWhoCanAccessAlbum(ctx, tt.args.albumIds...)
			if !tt.wantErr(t, err, fmt.Sprintf("ListUsersWhoCanAccessAlbum(%v, %v)", ctx, tt.args.albumIds)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ListUsersWhoCanAccessAlbum(%v, %v)", ctx, tt.args.albumIds)
		})
	}
}

func contributorAlbumPermission(userId usermodel.UserId, albumId catalog.AlbumId) *aclcore.Scope {
	return &aclcore.Scope{
		Type:          aclcore.AlbumContributorScope,
		GrantedTo:     userId,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	}
}

func visitorAlbumPermission(userId usermodel.UserId, albumId catalog.AlbumId) *aclcore.Scope {
	return &aclcore.Scope{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     userId,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	}
}

func ownerPermission(userId usermodel.UserId, owner ownermodel.Owner) *aclcore.Scope {
	return &aclcore.Scope{
		Type:          aclcore.MainOwnerScope,
		GrantedTo:     userId,
		ResourceOwner: owner,
	}
}
