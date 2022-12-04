package catalogacl_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/domain/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
)

func Test_rulesWithAccessToken_CanListMediasFromAlbum(t *testing.T) {
	nopeError := errors.Errorf("an error")

	type args struct {
		owner      string
		folderName string
	}
	tests := []struct {
		name    string
		mocks   func(t *testing.T) catalogacl.CatalogRules
		claims  aclcore.Claims
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should authorise access if albums belongs to pre-authorised owner",
			mocks: func(t *testing.T) catalogacl.CatalogRules {
				return mocks.NewCatalogRules(t)
			},
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes:  nil,
				Owner:   ownerEmail,
			},
			args:    args{owner: ownerEmail, folderName: folderName},
			wantErr: assert.NoError,
		},
		{
			name: "it should delegate access if albums belongs a different owner",
			mocks: func(t *testing.T) catalogacl.CatalogRules {
				rules := mocks.NewCatalogRules(t)
				rules.On("CanListMediasFromAlbum", ownerEmail, folderName).Return(nopeError)
				return rules
			},
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes:  nil,
				Owner:   "some@one.else",
			},
			args: args{owner: ownerEmail, folderName: folderName},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && assert.Equal(t, nopeError, err, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := catalogacl.OptimiseRulesWithAccessToken(tt.mocks(t), tt.claims)
			tt.wantErr(t, rules.CanListMediasFromAlbum(tt.args.owner, tt.args.folderName), fmt.Sprintf("CanListMediasFromAlbum(%v, %v)", tt.args.owner, tt.args.folderName))
		})
	}
}

func Test_rulesWithAccessToken_CanReadMedia(t *testing.T) {
	nopeError := errors.Errorf("an error")

	type fields struct {
		catalogRules func(t *testing.T) catalogacl.CatalogRules
	}
	type args struct {
		owner   string
		mediaId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		claims  aclcore.Claims
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should directly authorise if the owner is pre-authorised in the token",
			fields: fields{
				catalogRules: func(t *testing.T) catalogacl.CatalogRules {
					return mocks.NewCatalogRules(t)
				},
			},
			args: args{owner: ownerEmail, mediaId: mediaId},
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes:  nil,
				Owner:   ownerEmail,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delegate authorisation if not in the token [APPROVED]",
			fields: fields{
				catalogRules: func(t *testing.T) catalogacl.CatalogRules {
					rules := mocks.NewCatalogRules(t)
					rules.On("CanReadMedia", ownerEmail, mediaId).Return(nil)
					return rules
				},
			},
			args: args{owner: ownerEmail, mediaId: mediaId},
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes:  nil,
				Owner:   "some@else.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delegate authorisation if not in the token [DENIED]",
			fields: fields{
				catalogRules: func(t *testing.T) catalogacl.CatalogRules {
					rules := mocks.NewCatalogRules(t)
					rules.On("CanReadMedia", ownerEmail, mediaId).Return(nopeError)
					return rules
				},
			},
			args: args{owner: ownerEmail, mediaId: mediaId},
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes:  nil,
				Owner:   "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && assert.Equal(t, nopeError, nopeError, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := catalogacl.OptimiseRulesWithAccessToken(tt.fields.catalogRules(t), tt.claims)
			err := rules.CanReadMedia(tt.args.owner, tt.args.mediaId)
			tt.wantErr(t, err, fmt.Sprintf("Owner()"))
		})
	}
}

func Test_rulesWithAccessToken_Owner(t *testing.T) {
	tests := []struct {
		name    string
		claims  aclcore.Claims
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return the Owner from the token",
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes: map[string]interface{}{
					"owner:main:foobar": nil, // should be ignored
				},
				Owner: ownerEmail,
			},
			want:    ownerEmail,
			wantErr: assert.NoError,
		},
		{
			name: "it should return the empty if there is no owner",
			claims: aclcore.Claims{
				Subject: userEmail,
				Scopes: map[string]interface{}{
					"owner:main:foobar": nil, // should be ignored
				},
				Owner: "",
			},
			want:    "",
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := catalogacl.OptimiseRulesWithAccessToken(mocks.NewCatalogRules(t), tt.claims)
			got, err := rules.Owner()
			if !tt.wantErr(t, err, fmt.Sprintf("Owner()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Owner()")
		})
	}
}

type mockConstructorTestingTNewCatalogRules interface {
	mock.TestingT
	Cleanup(func())
}
