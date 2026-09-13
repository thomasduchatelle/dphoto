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

	tests := []struct {
		name       string
		scopes     []aclcore.Scope
		identities []aclcore.Identity
		owners     []ownermodel.Owner
		want       map[ownermodel.Owner][]*aclcore.Identity
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "it should return empty if no identity is attached to the owner",
			owners:  []ownermodel.Owner{ironmanOwner},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should fallback on email/email identity if the user never logged in the application",
			scopes: []aclcore.Scope{
				{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
			},
			owners: []ownermodel.Owner{ironmanOwner},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {
					{Email: tonyUser, Name: tonyUser},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the user identity attached to the owner",
			scopes: []aclcore.Scope{
				{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
			},
			identities: []aclcore.Identity{tonyIdentity},
			owners:     []ownermodel.Owner{ironmanOwner},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {&tonyIdentity},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should support the same user to be used by several owners",
			scopes: []aclcore.Scope{
				{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: ironmanOwner},
				{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: tonyUser, ResourceOwner: avengerOwner},
				{Type: aclcore.MainOwnerScope, GrantedAt: time.Time{}, GrantedTo: natashaUser, ResourceOwner: avengerOwner},
			},
			identities: []aclcore.Identity{tonyIdentity, natashaIdentity},
			owners:     []ownermodel.Owner{ironmanOwner, avengerOwner},
			want: map[ownermodel.Owner][]*aclcore.Identity{
				ironmanOwner: {&tonyIdentity},
				avengerOwner: {&tonyIdentity, &natashaIdentity},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopeRepository := NewScopeRepositoryInMemory(tt.scopes...)
			identityRepository := NewIdentityRepositoryInMemory(tt.identities...)
			i := &aclcore.IdentityQueries{
				IdentityRepository: identityRepository,
				ScopeRepository:    scopeRepository,
			}
			got, err := i.FindOwnerIdentities(tt.owners)
			if !tt.wantErr(t, err, fmt.Sprintf("FindOwnerIdentities(%v)", tt.owners)) {
				return
			}
			if err == nil {
				assert.Equalf(t, tt.want, got, "FindOwnerIdentities(%v)", tt.owners)
			}
		})
	}
}
