package aclcore_test

import (
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

	type fields struct {
		ScopeReader *ScopeRepositoryInMemory
		Email       usermodel.UserId
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ownermodel.Owner
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return resource owner from the ACL",
			fields: fields{
				ScopeReader: NewScopeRepositoryInMemory(
					aclcore.Scope{
						Type:          aclcore.MainOwnerScope,
						GrantedAt:     time.Time{},
						GrantedTo:     tonyEmail,
						ResourceOwner: ironmanOwner,
						ResourceId:    "007",
						ResourceName:  "Junior",
					},
				),
				Email: tonyEmail,
			},
			want:    &ironmanOwner,
			wantErr: assert.NoError,
		},
		{
			name: "it should return an error if no scopes are returned",
			fields: fields{
				ScopeReader: NewScopeRepositoryInMemory(),
				Email:       tonyEmail,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) &&
					assert.Contains(t, err.Error(), "is not a main user", i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aclcore.CoreRules{
				ScopeReader: tt.fields.ScopeReader,
				Email:       tt.fields.Email,
			}

			got, err := a.Owner()
			if !tt.wantErr(t, err, "Owner()") {
				return
			}
			assert.Equalf(t, tt.want, got, "Owner()")
		})
	}
}
