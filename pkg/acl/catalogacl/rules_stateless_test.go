package catalogacl_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

const (
	userEmail  = "pepper@stark.com"
	ownerEmail = "tony@stark.com"
	folderName = "InfinityWar"
	mediaId    = "Snap"
)

func Test_rules_CanListMediasFromAlbum(t *testing.T) {
	type args struct {
		owner      string
		folderName string
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
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
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
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
						ResourceId:    folderName,
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
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
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

			err := rules.CanListMediasFromAlbum(tt.args.owner, tt.args.folderName)
			tt.wantErr(t, err)
		})
	}
}

func Test_rules_CanReadMedia(t *testing.T) {
	type args struct {
		owner   string
		mediaId string
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
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(folderName, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId},
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
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(folderName, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
						ResourceId:    folderName,
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, mediaId: mediaId},
			wantErr: assert.NoError,
		},
		{
			name: "it should GRANT for the visitor of the media",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(folderName, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId},
				).Return([]*aclcore.Scope{
					{
						Type:          aclcore.MediaVisitorScope,
						GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
						ResourceOwner: ownerEmail,
						ResourceId:    mediaId,
					},
				}, nil)
			},
			args:    args{owner: ownerEmail, mediaId: mediaId},
			wantErr: assert.NoError,
		},
		{
			name: "it should DENY if media not found",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return("", catalog.NotFoundError)
			},
			args: args{owner: ownerEmail, mediaId: mediaId},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
			},
		},
		{
			name: "it should DENY for others",
			initMocks: func(repository *mocks.ScopeRepository, resolver *mocks.MediaAlbumResolver) {
				resolver.On("FindAlbumOfMedia", ownerEmail, mediaId).Return(folderName, nil)

				repository.On("FindScopesById",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: userEmail, ResourceOwner: ownerEmail},
					aclcore.ScopeId{Type: aclcore.AlbumVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: folderName},
					aclcore.ScopeId{Type: aclcore.MediaVisitorScope, GrantedTo: userEmail, ResourceOwner: ownerEmail, ResourceId: mediaId},
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
		want      map[string][]string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should show an album shared with 2 users, and another with 1",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListOwnerScopes", ownerEmail, aclcore.AlbumVisitorScope).Return([]*aclcore.Scope{
					{
						Type:          aclcore.AlbumVisitorScope,
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
			want: map[string][]string{
				"InfinityWar": {"blackwidow@avengers.com", "hulk@avengers.com"},
				"Endgame":     {"blackwidow@avengers.com"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return empty map if no album is shared",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListOwnerScopes", ownerEmail, aclcore.AlbumVisitorScope).Return(nil, nil)
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
				repository.On("ListUserScopes", userEmail, aclcore.AlbumVisitorScope).Return([]*aclcore.Scope{
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
				{Owner: ownerEmail, FolderName: "InfinityWar"},
				{Owner: "blackwidow@avengers.com", FolderName: "Endgame"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return nil list if nothing has been shared",
			initMocks: func(repository *mocks.ScopeRepository) {
				repository.On("ListUserScopes", userEmail, aclcore.AlbumVisitorScope).Return(nil, nil)
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
