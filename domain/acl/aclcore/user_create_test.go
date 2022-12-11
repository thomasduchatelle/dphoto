package aclcore_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

func TestCreateUser_CreateUser(t *testing.T) {
	mockedDate := time.Date(2022, 12, 11, 12, 42, 0, 0, time.UTC)
	aclcore.TimeFunc = func() time.Time {
		return mockedDate
	}
	const tonyEmail = "tony@stark.com"
	const ironmanOwner = "ironman"

	type fields struct {
		ScopesReader aclcore.ScopesReader
		ScopeWriter  aclcore.ScopeWriter
	}
	type args struct {
		email string
		owner string
	}
	tests := []struct {
		name      string
		initMocks func(t *testing.T) fields
		args      args
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the scope when no scope already exists",
			initMocks: func(t *testing.T) fields {
				reader := mocks.NewScopesReader(t)
				reader.On("ListUserScopes", tonyEmail, aclcore.MainOwnerScope).Return(nil, nil)

				writer := mocks.NewScopeWriter(t)
				writer.On("SaveIfNewScope", aclcore.Scope{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     mockedDate,
					GrantedTo:     tonyEmail,
					ResourceOwner: ironmanOwner,
				}).Return(nil)

				return fields{reader, writer}
			},
			args:    args{email: tonyEmail, owner: ironmanOwner},
			wantErr: assert.NoError,
		},
		{
			name: "it should override a scope for a different owner (and remove noise)",
			initMocks: func(t *testing.T) fields {
				reader := mocks.NewScopesReader(t)
				reader.On("ListUserScopes", tonyEmail, aclcore.MainOwnerScope).Return([]*aclcore.Scope{
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Now(), GrantedTo: tonyEmail, ResourceOwner: tonyEmail},
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Now(), GrantedTo: tonyEmail, ResourceOwner: "someoneelse"},
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Now(), GrantedTo: tonyEmail, ResourceOwner: ironmanOwner, ResourceId: "the suit"},
					// blast is contained if repository returns something unexpected
					{Type: aclcore.AlbumVisitorScope, GrantedAt: time.Now(), GrantedTo: tonyEmail, ResourceOwner: ironmanOwner},
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Now(), GrantedTo: "pepper@stark.com", ResourceOwner: ironmanOwner},
				}, nil)

				writer := mocks.NewScopeWriter(t)
				writer.On(
					"DeleteScopes",
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: tonyEmail, ResourceOwner: tonyEmail},
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: tonyEmail, ResourceOwner: "someoneelse"},
					aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: tonyEmail, ResourceOwner: ironmanOwner, ResourceId: "the suit"},
				).Return(nil)
				writer.On("SaveIfNewScope", aclcore.Scope{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     mockedDate,
					GrantedTo:     tonyEmail,
					ResourceOwner: ironmanOwner,
				}).Return(nil)

				return fields{reader, writer}
			},
			args:    args{email: tonyEmail, owner: ironmanOwner},
			wantErr: assert.NoError,
		},
		{
			name: "it should skip if the scope already exists",
			initMocks: func(t *testing.T) fields {
				reader := mocks.NewScopesReader(t)
				reader.On("ListUserScopes", tonyEmail, aclcore.MainOwnerScope).Return([]*aclcore.Scope{
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Now(), GrantedTo: tonyEmail, ResourceOwner: ironmanOwner},
				}, nil)

				return fields{reader, mocks.NewScopeWriter(t)}
			},
			args:    args{email: tonyEmail, owner: ironmanOwner},
			wantErr: assert.NoError,
		},
		{
			name: "it should default the owner to the email",
			initMocks: func(t *testing.T) fields {
				reader := mocks.NewScopesReader(t)
				reader.On("ListUserScopes", tonyEmail, aclcore.MainOwnerScope).Return(nil, nil)

				writer := mocks.NewScopeWriter(t)
				writer.On("SaveIfNewScope", aclcore.Scope{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     mockedDate,
					GrantedTo:     tonyEmail,
					ResourceOwner: tonyEmail,
				}).Return(nil)

				return fields{reader, writer}
			},
			args:    args{email: tonyEmail},
			wantErr: assert.NoError,
		},
		{
			name: "it should return an error if the email is empty / invalid",
			initMocks: func(t *testing.T) fields {
				return fields{mocks.NewScopesReader(t), mocks.NewScopeWriter(t)}
			},
			args: args{email: "   "},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.InvalidUserEmailError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedFields := tt.initMocks(t)
			c := &aclcore.CreateUser{
				ScopesReader: mockedFields.ScopesReader,
				ScopeWriter:  mockedFields.ScopeWriter,
			}
			tt.wantErr(t, c.CreateUser(tt.args.email, tt.args.owner), fmt.Sprintf("CreateUser(%v, %v)", tt.args.email, tt.args.owner))
		})
	}
}
