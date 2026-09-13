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
	const email = usermodel.UserId("tony@stark.com")

	fullSpec := aclcore.RefreshTokenSpec{
		Email:               email,
		RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
		AbsoluteExpiryTime:  time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		Scopes:              []string{"ironman"},
	}
	storedIdentity := aclcore.Identity{
		Email:   email,
		Name:    "Tony Stark",
		Picture: "/you-know-who-am-i.jpg",
	}
	expectedNewRefreshToken := "rt-" + email.Value() + "-" + string(aclcore.RefreshTokenPurposeWeb)

	newFakesWithTokens := func(tokens map[string]aclcore.RefreshTokenSpec, identities ...aclcore.Identity) (*RefreshTokenRepositoryInMemory, *IdentityRepositoryInMemory, *AccessTokenGeneratorInMemory, *RefreshTokenGeneratorInMemory) {
		refreshTokenRepository := NewRefreshTokenRepositoryInMemory()
		for token, spec := range tokens {
			_ = refreshTokenRepository.StoreRefreshToken(token, spec)
		}
		identityRepository := NewIdentityRepositoryInMemory(identities...)
		accessTokenGenerator := NewAccessTokenGeneratorInMemory()
		refreshTokenGenerator := NewRefreshTokenGeneratorInMemory()
		refreshTokenGenerator.Repository = refreshTokenRepository
		return refreshTokenRepository, identityRepository, accessTokenGenerator, refreshTokenGenerator
	}

	type fields struct {
		AccessTokenGenerator   *AccessTokenGeneratorInMemory
		RefreshTokenGenerator  *RefreshTokenGeneratorInMemory
		RefreshTokenRepository *RefreshTokenRepositoryInMemory
		IdentityDetailsStore   *IdentityRepositoryInMemory
	}
	type args struct {
		refreshToken string
	}

	newFields := func(tokens map[string]aclcore.RefreshTokenSpec, identities ...aclcore.Identity) fields {
		repo, identities1, at, rt := newFakesWithTokens(tokens, identities...)
		return fields{
			AccessTokenGenerator:   at,
			RefreshTokenGenerator:  rt,
			RefreshTokenRepository: repo,
			IdentityDetailsStore:   identities1,
		}
	}

	tests := []struct {
		name               string
		fields             fields
		args               args
		wantAuthentication *aclcore.Authentication
		wantIdentity       *aclcore.Identity
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name:   "it should generate a new access token and refresh token with same spec",
			fields: newFields(map[string]aclcore.RefreshTokenSpec{refreshToken: fullSpec}, storedIdentity),
			args:   args{refreshToken: refreshToken},
			wantAuthentication: &aclcore.Authentication{
				AccessToken:  "at-" + email.Value(),
				RefreshToken: expectedNewRefreshToken,
				ExpiryTime:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresIn:    42,
			},
			wantIdentity: &storedIdentity,
			wantErr:      assert.NoError,
		},
		{
			name:   "it should generate tokens with fallback identity when identity is not stored",
			fields: newFields(map[string]aclcore.RefreshTokenSpec{refreshToken: fullSpec}),
			args:   args{refreshToken: refreshToken},
			wantAuthentication: &aclcore.Authentication{
				AccessToken:  "at-" + email.Value(),
				RefreshToken: expectedNewRefreshToken,
				ExpiryTime:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresIn:    42,
			},
			wantIdentity: &aclcore.Identity{
				Email:   email,
				Name:    email.Value(),
				Picture: "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not issue tokens if refresh token expired",
			fields: newFields(map[string]aclcore.RefreshTokenSpec{
				refreshToken: {AbsoluteExpiryTime: time.Date(2020, 12, 31, 23, 59, 59, 999, time.UTC)},
			}),
			args:               args{refreshToken: refreshToken},
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.ExpiredRefreshTokenError)
			},
		},
		{
			name:               "it should not issue tokens if refresh token not found",
			fields:             newFields(nil),
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

			if err == nil {
				_, findErr := tt.fields.RefreshTokenRepository.FindRefreshToken(tt.args.refreshToken)
				assert.ErrorIs(t, findErr, aclcore.InvalidRefreshTokenError, "old refresh token should be deleted")

				newSpec, findErr := tt.fields.RefreshTokenRepository.FindRefreshToken(gotToken.RefreshToken)
				if assert.NoError(t, findErr, "new refresh token should be stored") {
					assert.Equal(t, fullSpec, *newSpec, "new refresh token spec")
				}

				assert.Equal(t, []aclcore.RefreshTokenSpec{fullSpec}, tt.fields.RefreshTokenGenerator.GeneratedFor, "RefreshTokenGenerator.GeneratedFor")
				assert.Equal(t, []usermodel.UserId{email}, tt.fields.AccessTokenGenerator.GeneratedFor, "AccessTokenGenerator.GeneratedFor")
			}
		})
	}
}
