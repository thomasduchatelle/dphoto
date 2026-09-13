package aclcore_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestIdentityQueries_FindOwnerIdentities(t *testing.T) {
	const ironmanOwner = ownermodel.Owner("iroman@avenger.com")
	const avengerOwner = ownermodel.Owner("heroes@avenger.com")
	const tonyUser = "tony@stark.com"
	const natashaUser = "natasha@banner.com"

	tonyIdentity := aclcore.Identity{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"}
	natashaIdentity := aclcore.Identity{Email: natashaUser, Name: "Natasha Banner", Picture: "/natasha.png"}

	type fields struct {
		IdentityRepository *IdentityRepositoryInMemory
		ScopeRepository    *ScopeRepositoryInMemory
	}
	type args struct {
		owners []ownermodel.Owner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[ownermodel.Owner][]*aclcore.Identity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return empty if no identity is attached to the owner",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(),
				ScopeRepository:    NewScopeRepositoryInMemory(),
			},
			args:    args{owners: []ownermodel.Owner{ironmanOwner}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should fallback on email/email identity if the user never logged in the application",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(),
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				),
			},
			args: args{owners: []ownermodel.Owner{ironmanOwner}},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {
					{Email: tonyUser, Name: tonyUser},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the user identity attached to the owner",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(tonyIdentity),
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				),
			},
			args: args{owners: []ownermodel.Owner{ironmanOwner}},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {&tonyIdentity},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should support the same user to be used by several owners",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(tonyIdentity, natashaIdentity),
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: avengerOwner},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: natashaUser, ResourceOwner: avengerOwner},
				),
			},
			args: args{owners: []ownermodel.Owner{ironmanOwner, avengerOwner}},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {&tonyIdentity},
				avengerOwner: {&tonyIdentity, &natashaIdentity},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &aclcore.IdentityQueries{
				IdentityRepository: tt.fields.IdentityRepository,
				ScopeRepository:    tt.fields.ScopeRepository,
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
