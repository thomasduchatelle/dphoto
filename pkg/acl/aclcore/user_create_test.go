package aclcore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestCreateUser_CreateUser(t *testing.T) {
	mockedDate := time.Date(2022, 12, 11, 12, 42, 0, 0, time.UTC)
	aclcore.TimeFunc = func() time.Time {
		return mockedDate
	}
	const tonyEmail = "tony@stark.com"
	const ironmanOwner = "ironman"
	const tonyUserId = usermodel.UserId(tonyEmail)

	type args struct {
		email string
		owner string
	}
	tests := []struct {
		name           string
		initialScopes  []aclcore.Scope
		args           args
		wantErr        assert.ErrorAssertionFunc
		wantUserScopes []*aclcore.Scope
	}{
		{
			name:          "it should create the scope when no scope already exists",
			initialScopes: nil,
			args:          args{email: tonyEmail, owner: ironmanOwner},
			wantErr:       assert.NoError,
			wantUserScopes: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     mockedDate,
					GrantedTo:     tonyUserId,
					ResourceOwner: ironmanOwner,
				},
			},
		},
		{
			name: "it should override a scope for a different owner (and remove noise)",
			initialScopes: []aclcore.Scope{
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: tonyEmail},
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: "someoneelse"},
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ironmanOwner, ResourceId: "the suit"},
				{Type: aclcore.AlbumVisitorScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ironmanOwner},
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: "pepper@stark.com", ResourceOwner: ironmanOwner},
			},
			args:    args{email: tonyEmail, owner: ironmanOwner},
			wantErr: assert.NoError,
			wantUserScopes: []*aclcore.Scope{
				{Type: aclcore.AlbumVisitorScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ironmanOwner},
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ironmanOwner},
			},
		},
		{
			name: "it should skip if the scope already exists",
			initialScopes: []aclcore.Scope{
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ironmanOwner},
			},
			args:    args{email: tonyEmail, owner: ironmanOwner},
			wantErr: assert.NoError,
			wantUserScopes: []*aclcore.Scope{
				{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ironmanOwner},
			},
		},
		{
			name:          "it should default the owner to the email",
			initialScopes: nil,
			args:          args{email: tonyEmail},
			wantErr:       assert.NoError,
			wantUserScopes: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     mockedDate,
					GrantedTo:     tonyUserId,
					ResourceOwner: ownermodel.Owner(tonyEmail),
				},
			},
		},
		{
			name:          "it should return an error if the email is empty / invalid",
			initialScopes: nil,
			args:          args{email: "   "},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, usermodel.InvalidUserEmailError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewScopeRepositoryInMemory(tt.initialScopes...)
			c := &aclcore.CreateUser{
				ScopesReader: repository,
				ScopeWriter:  repository,
			}
			err := c.CreateUser(tt.args.email, tt.args.owner)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateUser(%v, %v)", tt.args.email, tt.args.owner)) {
				return
			}
			if tt.wantUserScopes != nil {
				got, listErr := repository.ListScopesByUser(context.Background(), tonyUserId)
				assert.NoError(t, listErr)
				assert.ElementsMatch(t, tt.wantUserScopes, got, "user scopes after CreateUser")
			}
		})
	}
}
