package aclcore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestCognitoClaimsEnricher_Enrich(t *testing.T) {
	const tonyEmail = "tony@stark.com"
	tonyOwner := ownermodel.Owner("tony@stark.com")

	isNotPreregisteredError := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, aclcore.NotPreregisteredError)
	}

	type fields struct {
		ScopesReader aclcore.ScopesReader
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    aclcore.EnrichedClaims
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should resolve the owner and scopes of a main owner",
			fields: fields{ScopesReader: NewScopeRepositoryInMemory(aclcore.Scope{
				Type:          aclcore.MainOwnerScope,
				GrantedTo:     tonyEmail,
				ResourceOwner: tonyOwner,
			})},
			args: args{email: tonyEmail},
			want: aclcore.EnrichedClaims{
				Email:  tonyEmail,
				Owner:  "tony@stark.com",
				Scopes: "owner:tony@stark.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should resolve api and owner scopes together",
			fields: fields{ScopesReader: NewScopeRepositoryInMemory(
				aclcore.Scope{Type: aclcore.ApiScope, GrantedTo: tonyEmail, ResourceId: "admin"},
				aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: tonyEmail, ResourceOwner: tonyOwner},
			)},
			args: args{email: tonyEmail},
			want: aclcore.EnrichedClaims{
				Email:  tonyEmail,
				Owner:  "tony@stark.com",
				Scopes: "api:admin owner:tony@stark.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should resolve a visitor with no owner",
			fields: fields{ScopesReader: NewScopeRepositoryInMemory(aclcore.Scope{
				Type:          aclcore.AlbumVisitorScope,
				GrantedTo:     tonyEmail,
				ResourceOwner: "someoneelse",
				ResourceId:    "wedding",
			})},
			args: args{email: tonyEmail},
			want: aclcore.EnrichedClaims{
				Email:  tonyEmail,
				Owner:  "",
				Scopes: "visitor",
			},
			wantErr: assert.NoError,
		},
		{
			name:    "it should normalise the email to lowercase and trimmed",
			fields:  fields{ScopesReader: NewScopeRepositoryInMemory(aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: tonyEmail, ResourceOwner: tonyOwner})},
			args:    args{email: "  Tony@Stark.com "},
			want:    aclcore.EnrichedClaims{Email: tonyEmail, Owner: "tony@stark.com", Scopes: "owner:tony@stark.com"},
			wantErr: assert.NoError,
		},
		{
			name:    "it should reject a user who is not pre-registered",
			fields:  fields{ScopesReader: NewScopeRepositoryInMemory()},
			args:    args{email: tonyEmail},
			wantErr: isNotPreregisteredError,
		},
		{
			name:    "it should reject an empty email",
			fields:  fields{ScopesReader: NewScopeRepositoryInMemory()},
			args:    args{email: ""},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enricher := &aclcore.CognitoClaimsEnricher{ScopesReader: tt.fields.ScopesReader}

			got, err := enricher.Enrich(context.Background(), tt.args.email)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
