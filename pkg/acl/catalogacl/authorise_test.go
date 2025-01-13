package catalogacl

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

func TestCatalogAuthorizer_IsAuthorisedToListMedias(t *testing.T) {
	owner1 := ownermodel.Owner("owner-1")
	owner2 := ownermodel.Owner("owner-2")
	authenticatedOwner1 := usermodel.CurrentUser{UserId: "user-1", Owner: &owner1}
	visitor2 := usermodel.CurrentUser{UserId: "user-2"}
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/folder-1")}
	albumId2 := catalog.AlbumId{Owner: owner2, FolderName: catalog.NewFolderName("/folder-1")}
	accessDenied := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, ErrAccessDenied, i...)
	}
	noPermissionsStored := &aclcore.ScopeReadRepositoryInMemory{}

	type fields struct {
		HasPermissionPort HasPermissionPort
	}
	type args struct {
		ctx     context.Context
		userId  usermodel.CurrentUser
		albumId catalog.AlbumId
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "it should deny access unless explicitly granted",
			fields: fields{HasPermissionPort: noPermissionsStored},
			args: args{
				ctx:     context.Background(),
				userId:  authenticatedOwner1,
				albumId: albumId2,
			},
			wantErr: accessDenied,
		},
		{
			name:   "it should grant access to a visitor (not the owner himself)",
			fields: fields{HasPermissionPort: noPermissionsStored},
			args: args{
				ctx:     context.Background(),
				userId:  authenticatedOwner1,
				albumId: albumId1,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should deny access to a non-authorised visitor",
			fields: fields{HasPermissionPort: noPermissionsStored},
			args: args{
				ctx:     context.Background(),
				userId:  visitor2,
				albumId: albumId1,
			},
			wantErr: accessDenied,
		},
		{
			name: "it should grant access to a visitor specifically authorised",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.AlbumVisitorScope, GrantedTo: visitor2.UserId, ResourceOwner: albumId1.Owner, ResourceId: albumId1.FolderName.String()},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				userId:  visitor2,
				albumId: albumId1,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CatalogAuthorizer{
				HasPermissionPort: tt.fields.HasPermissionPort,
			}

			err := f.IsAuthorisedToListMedias(tt.args.ctx, tt.args.userId, tt.args.albumId)
			tt.wantErr(t, err, fmt.Sprintf("IsAuthorisedToListMedias(%v, %v, %v)", tt.args.ctx, tt.args.userId, tt.args.albumId))
		})
	}
}

func TestCatalogAuthorizer_IsAuthorisedToViewMedia(t *testing.T) {
	userId1 := usermodel.UserId("user-1")
	owner1 := ownermodel.Owner("owner-1")
	mediaId1 := catalog.MediaId("media-1")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/folder-1")}
	currentUserAsOwner1 := usermodel.CurrentUser{UserId: userId1, Owner: &owner1}
	currentUserAsVisitor := usermodel.CurrentUser{UserId: userId1}

	isAnAccessForbiddenError := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
	}

	type fields struct {
		HasPermissionPort  HasPermissionPort
		CatalogQueriesPort CatalogQueriesPort
	}
	type args struct {
		ctx         context.Context
		currentUser usermodel.CurrentUser
		owner       ownermodel.Owner
		mediaId     catalog.MediaId
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should GRANT access to the media owner",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
				CatalogQueriesPort: &catalog.MediaQueriesInMemory{
					Medias: []catalog.InMemoryMedia{
						catalog.NewInMemoryMedia(mediaId1, albumId1),
					},
				},
			},
			args: args{
				ctx:         context.Background(),
				currentUser: currentUserAsOwner1,
				owner:       owner1,
				mediaId:     mediaId1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT access to a user with OWNER permission",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.MainOwnerScope, GrantedTo: userId1, ResourceOwner: owner1},
					},
				},
				CatalogQueriesPort: &catalog.MediaQueriesInMemory{
					Medias: []catalog.InMemoryMedia{
						catalog.NewInMemoryMedia(mediaId1, albumId1),
					},
				},
			},
			args: args{
				ctx:         context.Background(),
				currentUser: currentUserAsVisitor,
				owner:       owner1,
				mediaId:     mediaId1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT access to a user with VISITOR permission on the album",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.AlbumVisitorScope, GrantedTo: userId1, ResourceOwner: owner1, ResourceId: albumId1.FolderName.String()},
					},
				},
				CatalogQueriesPort: &catalog.MediaQueriesInMemory{
					Medias: []catalog.InMemoryMedia{
						catalog.NewInMemoryMedia(mediaId1, albumId1),
					},
				},
			},
			args: args{
				ctx:         context.Background(),
				currentUser: currentUserAsVisitor,
				owner:       owner1,
				mediaId:     mediaId1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT access to a visitor with a visitor permission on the MEDIA",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.MediaVisitorScope, GrantedTo: userId1, ResourceOwner: owner1, ResourceId: mediaId1.Value()},
					},
				},
				CatalogQueriesPort: &catalog.MediaQueriesInMemory{
					Medias: []catalog.InMemoryMedia{
						catalog.NewInMemoryMedia(mediaId1, albumId1),
					},
				},
			},
			args: args{
				ctx:         context.Background(),
				currentUser: currentUserAsVisitor,
				owner:       owner1,
				mediaId:     mediaId1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should DENY access to a visitor with no permission",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
				CatalogQueriesPort: &catalog.MediaQueriesInMemory{
					Medias: []catalog.InMemoryMedia{
						catalog.NewInMemoryMedia(mediaId1, albumId1),
					},
				},
			},
			args: args{
				ctx:         context.Background(),
				currentUser: currentUserAsVisitor,
				owner:       owner1,
				mediaId:     mediaId1,
			},
			wantErr: isAnAccessForbiddenError,
		},
		{
			name: "it should DENY access if the media is not found",
			fields: fields{
				HasPermissionPort:  &aclcore.ScopeReadRepositoryInMemory{},
				CatalogQueriesPort: &catalog.MediaQueriesInMemory{},
			},
			args: args{
				ctx:         context.Background(),
				currentUser: currentUserAsVisitor,
				owner:       owner1,
				mediaId:     mediaId1,
			},
			wantErr: isAnAccessForbiddenError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CatalogAuthorizer{
				HasPermissionPort:  tt.fields.HasPermissionPort,
				CatalogQueriesPort: tt.fields.CatalogQueriesPort,
			}
			tt.wantErr(t, f.IsAuthorisedToViewMedia(tt.args.ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId), fmt.Sprintf("IsAuthorisedToViewMedia(%v, %v, %v, %v)", tt.args.ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId))
		})
	}
}

func TestCatalogAuthorizer_CanManageAlbum(t *testing.T) {
	owner1 := ownermodel.Owner("owner-1")
	userOfOwner1 := usermodel.CurrentUser{UserId: "user-1", Owner: &owner1}
	user2 := usermodel.CurrentUser{UserId: "user-2"}
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/folder-1")}
	isAccessForbidden := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
	}

	type fields struct {
		HasPermissionPort  HasPermissionPort
		CatalogQueriesPort CatalogQueriesPort
	}
	type args struct {
		ctx     context.Context
		user    usermodel.CurrentUser
		albumId catalog.AlbumId
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should DENY access to a visitor without permissions",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
			},
			args: args{
				ctx:     context.Background(),
				user:    user2,
				albumId: albumId1,
			},
			wantErr: isAccessForbidden,
		},
		{
			name: "it should DENY access to a visitor of this album",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.AlbumVisitorScope, GrantedTo: user2.UserId, ResourceOwner: owner1},
						{Type: aclcore.AlbumVisitorScope, GrantedTo: user2.UserId, ResourceOwner: owner1, ResourceId: albumId1.FolderName.String()},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				user:    user2,
				albumId: albumId1,
			},
			wantErr: isAccessForbidden,
		},
		{
			name: "it should GRANT access to the album owner (from the currentUser)",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
			},
			args: args{
				ctx:     context.Background(),
				user:    userOfOwner1,
				albumId: albumId1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT access to the album owner (from the permissions)",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.MainOwnerScope, GrantedTo: user2.UserId, ResourceOwner: owner1},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				user:    user2,
				albumId: albumId1,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CatalogAuthorizer{
				HasPermissionPort:  tt.fields.HasPermissionPort,
				CatalogQueriesPort: tt.fields.CatalogQueriesPort,
			}
			err := a.CanShareAlbum(tt.args.ctx, tt.args.user, tt.args.albumId)
			tt.wantErr(t, err, fmt.Sprintf("CanShareAlbum(%v, %v, %v)", tt.args.ctx, tt.args.user, tt.args.albumId))
		})
	}
}

func TestCatalogAuthorizer_CanCreateAlbum(t *testing.T) {
	owner1 := ownermodel.Owner("owner-1")
	const user1Id = "user-1"
	isAccessForbidden := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
	}

	type fields struct {
		HasPermissionPort  HasPermissionPort
		CatalogQueriesPort CatalogQueriesPort
	}
	type args struct {
		user usermodel.CurrentUser
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantOwner *ownermodel.Owner
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should DENY access to a user without permission",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
			},
			args: args{
				user: usermodel.CurrentUser{UserId: user1Id},
			},
			wantErr: isAccessForbidden,
		},
		{
			name: "it should GRANT access to user representing the owner (from the currentUser)",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
			},
			args: args{
				user: usermodel.CurrentUser{UserId: user1Id, Owner: &owner1},
			},
			wantOwner: &owner1,
			wantErr:   assert.NoError,
		},
		{
			name: "it should GRANT access to user representing the owner (from the permissions)",
			fields: fields{
				HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
					Scopes: []*aclcore.Scope{
						{Type: aclcore.MainOwnerScope, GrantedTo: user1Id, ResourceOwner: owner1},
					},
				},
			},
			args: args{
				user: usermodel.CurrentUser{UserId: user1Id},
			},
			wantOwner: &owner1,
			wantErr:   assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CatalogAuthorizer{
				HasPermissionPort:  tt.fields.HasPermissionPort,
				CatalogQueriesPort: tt.fields.CatalogQueriesPort,
			}
			gotOwner, err := a.CanCreateAlbum(context.Background(), tt.args.user)
			if tt.wantErr(t, err, fmt.Sprintf("CanCreateAlbum(%v, %v)", context.Background(), tt.args.user)) {
				assert.Equal(t, tt.wantOwner, gotOwner)
			}
		})
	}
}
