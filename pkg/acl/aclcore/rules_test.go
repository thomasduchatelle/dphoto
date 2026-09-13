package aclcore_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestCoreRules_Owner(t *testing.T) {
	ironmanOwner := ownermodel.Owner("ironman")
	tonyEmail := usermodel.UserId("tony@stark.com")

	tests := []struct {
		name    string
		email   usermodel.UserId
		scopes  []aclcore.Scope
		want    *ownermodel.Owner
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:  "it should return resource owner from the ACL",
			email: tonyEmail,
			scopes: []aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Time{},
					GrantedTo:     tonyEmail,
					ResourceOwner: ironmanOwner,
					ResourceId:    "007",
					ResourceName:  "Junior",
				},
			},
			want:    &ironmanOwner,
			wantErr: assert.NoError,
		},
		{
			name:   "it should return an error if no scopes are returned",
			email:  tonyEmail,
			scopes: nil,
			want:   nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) &&
					assert.Contains(t, err.Error(), "is not a main user", i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aclcore.CoreRules{
				ScopeReader: NewScopeRepositoryInMemory(tt.scopes...),
				Email:       tt.email,
			}

			got, err := a.Owner()
			if !tt.wantErr(t, err, fmt.Sprintf("Owner()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Owner()")
		})
	}
}
