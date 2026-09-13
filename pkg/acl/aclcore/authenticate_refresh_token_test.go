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

	tests := []struct {
		name               string
		seedRefreshTokens  map[string]aclcore.RefreshTokenSpec
		seedIdentities     []aclcore.Identity
		refreshToken       string
		wantAuthentication *aclcore.Authentication
		wantIdentity       *aclcore.Identity
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name: "it should generate a new access token and refresh token with same spec",
			seedRefreshTokens: map[string]aclcore.RefreshTokenSpec{
				refreshToken: fullSpec,
			},
			seedIdentities: []aclcore.Identity{storedIdentity},
			refreshToken:   refreshToken,
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
			name: "it should generate tokens with fallback identity when identity is not stored",
			seedRefreshTokens: map[string]aclcore.RefreshTokenSpec{
				refreshToken: fullSpec,
			},
			seedIdentities: nil,
			refreshToken:   refreshToken,
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
			seedRefreshTokens: map[string]aclcore.RefreshTokenSpec{
				refreshToken: {
					AbsoluteExpiryTime: time.Date(2020, 12, 31, 23, 59, 59, 999, time.UTC),
				},
			},
			refreshToken:       refreshToken,
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.ExpiredRefreshTokenError)
			},
		},
		{
			name:               "it should not issue tokens if refresh token not found",
			seedRefreshTokens:  nil,
			refreshToken:       refreshToken,
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.InvalidRefreshTokenError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshTokenRepository := NewRefreshTokenRepositoryInMemory()
			for token, spec := range tt.seedRefreshTokens {
				assert.NoError(t, refreshTokenRepository.StoreRefreshToken(token, spec))
			}
			identityRepository := NewIdentityRepositoryInMemory(tt.seedIdentities...)
			accessTokenGenerator := NewAccessTokenGeneratorInMemory()
			refreshTokenGenerator := NewRefreshTokenGeneratorInMemory()
			refreshTokenGenerator.Repository = refreshTokenRepository

			s := &aclcore.RefreshTokenAuthenticator{
				AccessTokenGenerator:   accessTokenGenerator,
				RefreshTokenGenerator:  refreshTokenGenerator,
				RefreshTokenRepository: refreshTokenRepository,
				IdentityDetailsStore:   identityRepository,
			}
			gotToken, gotIdentity, err := s.AuthenticateFromRefreshToken(tt.refreshToken)
			if !tt.wantErr(t, err, fmt.Sprintf("AuthenticateFromRefreshToken(%v)", tt.refreshToken)) {
				return
			}
			assert.Equalf(t, tt.wantAuthentication, gotToken, "AuthenticateFromRefreshToken(%v)", tt.refreshToken)
			assert.Equalf(t, tt.wantIdentity, gotIdentity, "AuthenticateFromRefreshToken(%v)", tt.refreshToken)

			if err == nil {
				_, findErr := refreshTokenRepository.FindRefreshToken(tt.refreshToken)
				assert.ErrorIs(t, findErr, aclcore.InvalidRefreshTokenError, "old refresh token should be deleted")

				newSpec, findErr := refreshTokenRepository.FindRefreshToken(gotToken.RefreshToken)
				if assert.NoError(t, findErr, "new refresh token should be stored") {
					assert.Equal(t, fullSpec, *newSpec, "new refresh token spec")
				}

				assert.Equal(t, []aclcore.RefreshTokenSpec{fullSpec}, refreshTokenGenerator.GeneratedFor, "RefreshTokenGenerator.GeneratedFor")
				assert.Equal(t, []usermodel.UserId{email}, accessTokenGenerator.GeneratedFor, "AccessTokenGenerator.GeneratedFor")
			}
		})
	}
}
