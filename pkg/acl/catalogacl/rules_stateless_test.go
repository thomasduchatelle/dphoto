package catalogacl_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
	"time"
)

var (
	userEmail  = usermodel.UserId("pepper@stark.com")
	ownerEmail = ownermodel.Owner("tony@stark.com")
	mediaId    = catalog.MediaId("Snap")
)

var (
	folderName = catalog.NewFolderName("InfinityWar")
	albumId    = catalog.AlbumId{Owner: ownerEmail, FolderName: folderName}
)

func Test_rules_CanListMediasFromAlbum(t *testing.T) {
	type args struct {
		owner      ownermodel.Owner
		folderName catalog.FolderName
	}
	tests := []struct {
		name      string
		initMocks func(repository *mocks.ScopeRepository)
		args      args
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should GRANT for the owner of the album",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.MainOwnerScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, folderName: folderName},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT for the visitor of the album",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
						ResourceId:    folderName.String(),
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, folderName: folderName},
			wantErr: assert.NoError,
		},
		{
			name: "it should DENY for others",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
				).Return(nil, nil)
			},
			args: args{owner: ownerEmail, folderName: folderName},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewScopeRepository(t)
			tt.initMocks(repository)

			rules := catalogacl.NewCatalogRules(repository, nil, userEmail)

			err := rules.CanListMediasFromAlbum(catalog.AlbumId{Owner: tt.args.owner, FolderName: tt.args.folderName})
			tt.wantErr(t, err)
		})
	}
}

func Test_rules_CanReadMedia(t *testing.T) {
	type args struct {
		owner   ownermodel.Owner
		mediaId catalog.MediaId
	}
	tests := []struct {
		name      string
		initMocks func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver)
		args      args
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should GRANT for the owner of the media",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(albumId, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId.Value()},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.MainOwnerScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, mediaId: mediaId},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT for the visitor of the album",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(albumId, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId.Value()},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
						ResourceId:    folderName.String(),
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, mediaId: mediaId},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT for the visitor of the media",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(albumId, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId.Value()},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.MediaVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
						ResourceId:    mediaId.Value(),
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, mediaId: mediaId},
			wantErr: assert.NoError,
		},
		{
			name: "it should DENY if media not found",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(catalog.AlbumId{}, catalog.AlbumNotFoundError)
			},
			args: args{owner: ownerEmail, mediaId: mediaId},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
			},
		},
		{
			name: "it should DENY for others",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(albumId, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName.String()},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId.Value()},
				).Return(nil, nil)
			},
			args: args{owner: ownerEmail, mediaId: mediaId},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewScopeRepository(t)
			resolver := mocks.NewMediaAlbumResolver(t)
			tt.initMocks(repository, resolver)

			rules := catalogacl.NewCatalogRules(repository, resolver, userEmail)

			err := rules.CanReadMedia(tt.args.owner, tt.args.mediaId)
			tt.wantErr(t, err)
		})
	}
}

func Test_rules_SharedByUserGrid(t *testing.T) {
	tests := []struct {
		name      string
		initMocks func(repository *mocks.ScopeRepository)
		want      map[string]map[usermodel.UserId]aclcore.ScopeType
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should show an album shared with 2 users, and another with 1",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListScopesByOwner", mock.Anything, ownerEmail, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumContributorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						GrantedTo:     "blackwidow@avengers.com",
						ResourceOwner: ownerEmail,
						ResourceId:    "InfinityWar",
					},
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						GrantedTo:     "blackwidow@avengers.com",
						ResourceOwner: ownerEmail,
						ResourceId:    "Endgame",
					},
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						GrantedTo:     "hulk@avengers.com",
						ResourceOwner: ownerEmail,
						ResourceId:    "InfinityWar",
					},
				}, nil)
			},
			want: map[string]map[usermodel.UserId]aclcore.ScopeType{
				"InfinityWar": {"blackwidow@avengers.com": aclcore.AlbumContributorScope, "hulk@avengers.com": aclcore.AlbumVisitorScope},
				"Endgame":     {"blackwidow@avengers.com": aclcore.AlbumVisitorScope},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return empty map if no album is shared",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListScopesByOwner", mock.Anything, ownerEmail, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope).Return(nil, nil)
			},
			want:    nil,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewScopeRepository(t)
			tt.initMocks(repository)

			rules := catalogacl.NewCatalogRules(repository, nil, userEmail)

			got, err := rules.SharedByUserGrid(ownerEmail)
			if tt.wantErr(t, err) && err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_rules_SharedWithUserAlbum(t *testing.T) {
	tests := []struct {
		name      string
		initMocks func(repository *mocks.ScopeRepository)
		want      []catalog.AlbumId
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should return list of albums shared to current users",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListScopesByUser", mock.Anything, userEmail, aclcore.AlbumVisitorScope).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						GrantedTo:     userEmail,
						ResourceOwner: ownerEmail,
						ResourceId:    "InfinityWar",
					},
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						GrantedTo:     userEmail,
						ResourceOwner: "blackwidow@avengers.com",
						ResourceId:    "Endgame",
					},
				}, nil)
			},
			want: []catalog.AlbumId{
				{Owner: ownerEmail, FolderName: "/InfinityWar"},
				{Owner: "blackwidow@avengers.com", FolderName: "/Endgame"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return nil list if nothing has been shared",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListScopesByUser", mock.Anything, userEmail, aclcore.AlbumVisitorScope).Return(nil, nil)
			},
			want:    nil,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewScopeRepository(t)
			tt.initMocks(repository)

			rules := catalogacl.NewCatalogRules(repository, nil, userEmail)

			got, err := rules.SharedWithUserAlbum()
			if tt.wantErr(t, err) && err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_rules_CanManageAlbum(t *testing.T) {
	type args struct {
		owner      ownermodel.Owner
		folderName catalog.FolderName
	}
	tests := []struct {
		name      string
		initMocks func(repository *mocks.ScopeRepository)
		args      args
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should return no error if the user is a owner of the album",
			args: args{ownerEmail, folderName},
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("FindScopesById", aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail}).Return([]*aclcore.Scope{
					{
						Type:          aclcore.MainOwnerScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						GrantedTo:     userEmail,
						ResourceOwner: ownerEmail,
					},
				}, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return AccessForbiddenError is the user is not a owner",
			args: args{ownerEmail, folderName},
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("FindScopesById", aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail}).Return(nil, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.AccessForbiddenError, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := mocks.NewScopeRepository(t)
			tt.initMocks(repository)

			rules := catalogacl.NewCatalogRules(repository, nil, userEmail)

			err := rules.CanManageAlbum(catalog.AlbumId{Owner: tt.args.owner, FolderName: tt.args.folderName})
			tt.wantErr(t, err)
		})
	}
}
