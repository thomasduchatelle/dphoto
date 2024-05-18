package aclcore_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
	"time"
)

func TestCoreRules_Owner(t *testing.T) {
	ironmanOwner := ownermodel.Owner("ironman")

	tests := []struct {
		name      string
		email     usermodel.UserId
		initMocks func(scopesReader *mocks.ScopesReader)
		want      *ownermodel.Owner
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:  "it should return resource owner from the ACL",
			email: "tony@stark.com",
			initMocks: func(scopesReader *mocks.ScopesReader) {
				scopesReader.On("ListScopesByUser", mock.Anything, usermodel.UserId("tony@stark.com"), aclcore.MainOwnerScope).Return([]*aclcore.Scope{
					{
						Type:          aclcore.MainOwnerScope,
						GrantedAt:     time.Time{},
						GrantedTo:     "tony@stark.com",
						ResourceOwner: ironmanOwner,
						ResourceId:    "007",
						ResourceName:  "Junior",
					},
				}, nil)
			},
			want:    &ironmanOwner,
			wantErr: assert.NoError,
		},
		{
			name:  "it should return an error if no scopes are returned",
			email: "tony@stark.com",
			initMocks: func(scopesReader *mocks.ScopesReader) {
				scopesReader.On("ListScopesByUser", mock.Anything, usermodel.UserId("tony@stark.com"), aclcore.MainOwnerScope).Return(nil, nil)
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
			scopesReader := mocks.NewScopesReader(t)
			tt.initMocks(scopesReader)

			a := &aclcore.CoreRules{
				ScopeReader: scopesReader,
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
