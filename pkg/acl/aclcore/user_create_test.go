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
	const tonyOwner = ownermodel.Owner(ironmanOwner)

	type fields struct {
		ScopeRepository *ScopeRepositoryInMemory
	}
	type args struct {
		email string
		owner string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectOwner ownermodel.Owner
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "it should create the scope when no scope already exists",
			fields:      fields{ScopeRepository: NewScopeRepositoryInMemory()},
			args:        args{email: tonyEmail, owner: ironmanOwner},
			expectOwner: tonyOwner,
			wantErr:     assert.NoError,
		},
		{
			name: "it should override a scope for a different owner (and remove noise)",
			fields: fields{
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: ownermodel.Owner(tonyEmail)},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: "someoneelse"},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: tonyOwner, ResourceId: "the suit"},
					aclcore.Scope{Type: aclcore.AlbumVisitorScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: tonyOwner},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: "pepper@stark.com", ResourceOwner: tonyOwner},
				),
			},
			args:        args{email: tonyEmail, owner: ironmanOwner},
			expectOwner: tonyOwner,
			wantErr:     assert.NoError,
		},
		{
			name: "it should skip if the scope already exists",
			fields: fields{
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: mockedDate, GrantedTo: tonyUserId, ResourceOwner: tonyOwner},
				),
			},
			args:        args{email: tonyEmail, owner: ironmanOwner},
			expectOwner: tonyOwner,
			wantErr:     assert.NoError,
		},
		{
			name:        "it should default the owner to the email",
			fields:      fields{ScopeRepository: NewScopeRepositoryInMemory()},
			args:        args{email: tonyEmail},
			expectOwner: ownermodel.Owner(tonyEmail),
			wantErr:     assert.NoError,
		},
		{
			name:   "it should return an error if the email is empty / invalid",
			fields: fields{ScopeRepository: NewScopeRepositoryInMemory()},
			args:   args{email: "   "},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, usermodel.InvalidUserEmailError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &aclcore.CreateUser{
				ScopesReader: tt.fields.ScopeRepository,
				ScopeWriter:  tt.fields.ScopeRepository,
			}
			if !tt.wantErr(t, c.CreateUser(tt.args.email, tt.args.owner), fmt.Sprintf("CreateUser(%v, %v)", tt.args.email, tt.args.owner)) {
				return
			}
			if tt.expectOwner == "" {
				return
			}
			scopes, err := tt.fields.ScopeRepository.ListScopesByUser(context.Background(), tonyUserId, aclcore.MainOwnerScope)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     mockedDate,
					GrantedTo:     tonyUserId,
					ResourceOwner: tt.expectOwner,
				},
			}, scopes, "MainOwnerScope for %s", tonyUserId)
		})
	}
}
