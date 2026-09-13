package aclcore_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestAccessTokenAuthenticator_AuthenticateFromAccessToken(t *testing.T) {
	aclcore.TimeFunc = func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	const refreshToken = "1234567890qwertyuiop"
	const expiredRefreshToken = "expired-refresh-token"
	const email = usermodel.UserId("tony@stark.com")
	newRefreshToken := "rt-" + string(email) + "-" + string(aclcore.RefreshTokenPurposeWeb)

	fullSpec := aclcore.RefreshTokenSpec{
		Email:               email,
		RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
		AbsoluteExpiryTime:  time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		Scopes:              []string{"ironman"},
	}
	tonyIdentity := aclcore.Identity{
		Email:   email,
		Name:    "Tony Stark",
		Picture: "/you-know-who-am-i.jpg",
	}

	refreshTokenRepositoryWithFullSpec := func() *RefreshTokenRepositoryInMemory {
		repo := NewRefreshTokenRepositoryInMemory()
		_ = repo.StoreRefreshToken(refreshToken, fullSpec)
		return repo
	}

	type fields struct {
		AccessTokenGenerator   *AccessTokenGeneratorFake
		RefreshTokenGenerator  *RefreshTokenGeneratorFake
		RefreshTokenRepository *RefreshTokenRepositoryInMemory
		IdentityDetailsStore   *IdentityRepositoryInMemory
	}
	type args struct {
		refreshToken string
	}
	tests := []struct {
		name                       string
		fields                     fields
		args                       args
		wantAuthentication         *aclcore.Authentication
		wantIdentity               *aclcore.Identity
		expectAccessGeneratedFor   []usermodel.UserId
		expectRefreshGeneratedFor  []aclcore.RefreshTokenSpec
		expectOriginalTokenDeleted bool
		wantErr                    assert.ErrorAssertionFunc
	}{
		{
			name: "it should generate a new access token and refresh token with same spec",
			fields: fields{
				AccessTokenGenerator:   NewAccessTokenGeneratorFake(),
				RefreshTokenGenerator:  NewRefreshTokenGeneratorFake(),
				RefreshTokenRepository: refreshTokenRepositoryWithFullSpec(),
				IdentityDetailsStore:   NewIdentityRepositoryInMemory(tonyIdentity),
			},
			args: args{refreshToken: refreshToken},
			wantAuthentication: &aclcore.Authentication{
				AccessToken:  "at-" + string(email),
				RefreshToken: newRefreshToken,
				ExpiryTime:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresIn:    42,
			},
			wantIdentity:               &tonyIdentity,
			expectAccessGeneratedFor:   []usermodel.UserId{email},
			expectRefreshGeneratedFor:  []aclcore.RefreshTokenSpec{fullSpec},
			expectOriginalTokenDeleted: true,
			wantErr:                    assert.NoError,
		},
		{
			name: "it should generates tokens with fallback identity when the identity is unknown",
			fields: fields{
				AccessTokenGenerator:   NewAccessTokenGeneratorFake(),
				RefreshTokenGenerator:  NewRefreshTokenGeneratorFake(),
				RefreshTokenRepository: refreshTokenRepositoryWithFullSpec(),
				IdentityDetailsStore:   NewIdentityRepositoryInMemory(),
			},
			args: args{refreshToken: refreshToken},
			wantAuthentication: &aclcore.Authentication{
				AccessToken:  "at-" + string(email),
				RefreshToken: newRefreshToken,
				ExpiryTime:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresIn:    42,
			},
			wantIdentity: &aclcore.Identity{
				Email:   email,
				Name:    email.Value(),
				Picture: "",
			},
			expectAccessGeneratedFor:   []usermodel.UserId{email},
			expectRefreshGeneratedFor:  []aclcore.RefreshTokenSpec{fullSpec},
			expectOriginalTokenDeleted: true,
			wantErr:                    assert.NoError,
		},
		{
			name: "it should not issue tokens if refresh token expired",
			fields: fields{
				AccessTokenGenerator:  NewAccessTokenGeneratorFake(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorFake(),
				RefreshTokenRepository: func() *RefreshTokenRepositoryInMemory {
					repo := NewRefreshTokenRepositoryInMemory()
					_ = repo.StoreRefreshToken(expiredRefreshToken, aclcore.RefreshTokenSpec{
						AbsoluteExpiryTime: time.Date(2020, 12, 31, 23, 59, 59, 999, time.UTC),
					})
					return repo
				}(),
				IdentityDetailsStore: NewIdentityRepositoryInMemory(),
			},
			args:               args{refreshToken: expiredRefreshToken},
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.ExpiredRefreshTokenError)
			},
		},
		{
			name: "it should not issue tokens if refresh token not found",
			fields: fields{
				AccessTokenGenerator:   NewAccessTokenGeneratorFake(),
				RefreshTokenGenerator:  NewRefreshTokenGeneratorFake(),
				RefreshTokenRepository: NewRefreshTokenRepositoryInMemory(),
				IdentityDetailsStore:   NewIdentityRepositoryInMemory(),
			},
			args:               args{refreshToken: refreshToken},
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.InvalidRefreshTokenError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &aclcore.RefreshTokenAuthenticator{
				AccessTokenGenerator:   tt.fields.AccessTokenGenerator,
				RefreshTokenGenerator:  tt.fields.RefreshTokenGenerator,
				RefreshTokenRepository: tt.fields.RefreshTokenRepository,
				IdentityDetailsStore:   tt.fields.IdentityDetailsStore,
			}
			gotToken, gotIdentity, err := s.AuthenticateFromRefreshToken(tt.args.refreshToken)
			if !tt.wantErr(t, err, fmt.Sprintf("AuthenticateFromRefreshToken(%v)", tt.args.refreshToken)) {
				return
			}
			assert.Equalf(t, tt.wantAuthentication, gotToken, "AuthenticateFromRefreshToken(%v)", tt.args.refreshToken)
			assert.Equalf(t, tt.wantIdentity, gotIdentity, "AuthenticateFromRefreshToken(%v)", tt.args.refreshToken)
			assert.Equal(t, tt.expectAccessGeneratedFor, tt.fields.AccessTokenGenerator.GeneratedFor, "AccessTokenGenerator.GeneratedFor")
			assert.Equal(t, tt.expectRefreshGeneratedFor, tt.fields.RefreshTokenGenerator.GeneratedFor, "RefreshTokenGenerator.GeneratedFor")

			if tt.expectOriginalTokenDeleted {
				_, err := tt.fields.RefreshTokenRepository.FindRefreshToken(tt.args.refreshToken)
				assert.ErrorIs(t, err, aclcore.InvalidRefreshTokenError, "original refresh token should be deleted")
			}
		})
	}
}
