package aclcore_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"testing"
	"time"
)

func TestIdentityQueries_FindOwnerIdentities(t *testing.T) {
	const ironmanOwner = "iroman@avenger.com"
	const avengerOwner = "heroes@avenger.com"
	const tonyUser = "tony@stark.com"
	const natashaUser = "natasha@banner.com"

	type fields struct {
		IdentityRepository aclcore.IdentityQueriesIdentityRepository
		ScopeRepository    aclcore.IdentityQueriesScopeRepository
	}
	type args struct {
		owners []string
	}
	tests := []struct {
		name    string
		fields  func(t *testing.T) fields
		args    args
		want    map[string][]*aclcore.Identity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return empty if no identity is attached to the owner",
			fields: func(t *testing.T) fields {
				scopeRepository := mocks.NewIdentityQueriesScopeRepository(t)
				scopeRepository.On("ListScopesByOwners", []string{ironmanOwner}, aclcore.MainOwnerScope).Return(nil, nil)

				return fields{
					IdentityRepository: mocks.NewIdentityQueriesIdentityRepository(t),
					ScopeRepository:    scopeRepository,
				}
			},
			args:    args{[]string{ironmanOwner}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should fallback on email/email identity if the user never logged in the application",
			fields: func(t *testing.T) fields {
				scopeRepository := mocks.NewIdentityQueriesScopeRepository(t)
				scopeRepository.On("ListScopesByOwners", []string{ironmanOwner}, aclcore.MainOwnerScope).Return([]*aclcore.Scope{
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				}, nil)

				identityRepository := mocks.NewIdentityQueriesIdentityRepository(t)
				identityRepository.On("FindIdentities", []string{tonyUser}).Return(nil, nil)

				return fields{
					IdentityRepository: identityRepository,
					ScopeRepository:    scopeRepository,
				}
			},
			args: args{[]string{ironmanOwner}},
			want: map[string][]*aclcore.Identity{
				ironmanOwner: {
					{Email: tonyUser, Name: tonyUser},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the user identity attached to the owner",
			fields: func(t *testing.T) fields {
				scopeRepository := mocks.NewIdentityQueriesScopeRepository(t)
				scopeRepository.On("ListScopesByOwners", []string{ironmanOwner}, aclcore.MainOwnerScope).Return([]*aclcore.Scope{
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				}, nil)

				identityRepository := mocks.NewIdentityQueriesIdentityRepository(t)
				identityRepository.On("FindIdentities", []string{tonyUser}).Return([]*aclcore.Identity{
					{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"},
				}, nil)

				return fields{
					IdentityRepository: identityRepository,
					ScopeRepository:    scopeRepository,
				}
			},
			args: args{[]string{ironmanOwner}},
			want: map[string][]*aclcore.Identity{
				ironmanOwner: {
					{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should support the same user to be used by several owners",
			fields: func(t *testing.T) fields {
				scopeRepository := mocks.NewIdentityQueriesScopeRepository(t)
				scopeRepository.On("ListScopesByOwners", []string{ironmanOwner}, aclcore.MainOwnerScope).Return([]*aclcore.Scope{
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: avengerOwner},
					{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: natashaUser, ResourceOwner: avengerOwner},
				}, nil)

				identityRepository := mocks.NewIdentityQueriesIdentityRepository(t)
				identityRepository.On("FindIdentities", []string{tonyUser, tonyUser, natashaUser}).Return([]*aclcore.Identity{
					{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"},
					{Email: natashaUser, Name: "Natasha Banner", Picture: "/natasha.png"},
				}, nil)

				return fields{
					IdentityRepository: identityRepository,
					ScopeRepository:    scopeRepository,
				}
			},
			args: args{[]string{ironmanOwner}},
			want: map[string][]*aclcore.Identity{
				ironmanOwner: {
					{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"},
				},
				avengerOwner: {
					{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"},
					{Email: natashaUser, Name: "Natasha Banner", Picture: "/natasha.png"},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields(t)
			i := &aclcore.IdentityQueries{
				IdentityRepository: fields.IdentityRepository,
				ScopeRepository:    fields.ScopeRepository,
			}
			got, err := i.FindOwnerIdentities(tt.args.owners)
			if !tt.wantErr(t, err, fmt.Sprintf("FindOwnerIdentities(%v)", tt.args.owners)) {
				return
			}
			if err == nil {
				assert.Equalf(t, tt.want, got, "FindOwnerIdentities(%v)", tt.args.owners)
			}
		})
	}
}
