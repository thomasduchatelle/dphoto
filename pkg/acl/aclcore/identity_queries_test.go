package aclcore_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestIdentityQueries_FindOwnerIdentities(t *testing.T) {
	const ironmanOwner = ownermodel.Owner("iroman@avenger.com")
	const avengerOwner = ownermodel.Owner("heroes@avenger.com")
	const tonyUser = usermodel.UserId("tony@stark.com")
	const natashaUser = usermodel.UserId("natasha@banner.com")

	tonyStark := aclcore.Identity{Email: tonyUser, Name: "Tony Stark", Picture: "/tony-stark.jpg"}
	natashaBanner := aclcore.Identity{Email: natashaUser, Name: "Natasha Banner", Picture: "/natasha.png"}

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
			args:    args{[]ownermodel.Owner{ironmanOwner}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should fallback on email/email identity if the user never logged in the application",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(),
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				),
			},
			args: args{[]ownermodel.Owner{ironmanOwner}},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {
					{Email: tonyUser, Name: tonyUser.Value()},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the user identity attached to the owner",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(tonyStark),
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				),
			},
			args: args{[]ownermodel.Owner{ironmanOwner}},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {&tonyStark},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should support the same user to be used by several owners",
			fields: fields{
				IdentityRepository: NewIdentityRepositoryInMemory(tonyStark, natashaBanner),
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: tonyUser, ResourceOwner: avengerOwner},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: natashaUser, ResourceOwner: avengerOwner},
				),
			},
			args: args{[]ownermodel.Owner{ironmanOwner, avengerOwner}},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {&tonyStark},
				avengerOwner: {&tonyStark, &natashaBanner},
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
