package aclcore_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
	"time"
)

func TestAccessTokenAuthenticator_AuthenticateFromAccessToken(t *testing.T) {
	aclcore.TimeFunc = func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	const refreshToken = "1234567890qwertyuiop"
	const newRefreshToken = "a new refresh token"
	const email = usermodel.UserId("tony@stark.com")

	type fields struct {
		AccessTokenGenerator  func(t *testing.T) aclcore.IAccessTokenGenerator
		RefreshTokenGenerator func(t *testing.T) aclcore.IRefreshTokenGenerator
		AccessTokenRepository func(t *testing.T) aclcore.RefreshTokenRepository
		IdentityDetailsStore  func(t *testing.T) aclcore.IdentityDetailsStore
	}
	type args struct {
		refreshToken string
	}
	fullSpec := aclcore.RefreshTokenSpec{
		Email:               email,
		RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
		AbsoluteExpiryTime:  time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		Scopes:              []string{"ironman"},
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
			name: "it should generate a new access token and refresh token with same spec",
			fields: fields{
				AccessTokenGenerator: func(t *testing.T) aclcore.IAccessTokenGenerator {
					generator := mocks.NewIAccessTokenGenerator(t)
					generator.On("GenerateAccessToken", email).Return(&aclcore.Authentication{
						AccessToken:  "test-access-token",
						RefreshToken: "should be empty",
						ExpiryTime:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						ExpiresIn:    42,
					}, nil)
					return generator
				},
				RefreshTokenGenerator: func(t *testing.T) aclcore.IRefreshTokenGenerator {
					generator := mocks.NewIRefreshTokenGenerator(t)
					generator.On("GenerateRefreshToken", fullSpec).Return(newRefreshToken, nil)
					return generator
				},
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewRefreshTokenRepository(t)
					repository.On("FindRefreshToken", refreshToken).Return(&fullSpec, nil)
					repository.On("DeleteRefreshToken", refreshToken).Return(nil)
					return repository
				},
				IdentityDetailsStore: func(t *testing.T) aclcore.IdentityDetailsStore {
					store := mocks.NewIdentityDetailsStore(t)
					store.On("FindIdentity", email).Return(&aclcore.Identity{
						Email:   email,
						Name:    "Tony Stark",
						Picture: "/you-know-who-am-i.jpg",
					}, nil)
					return store
				},
			},
			args: args{
				refreshToken: refreshToken,
			},
			wantAuthentication: &aclcore.Authentication{
				AccessToken:  "test-access-token",
				RefreshToken: newRefreshToken,
				ExpiryTime:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresIn:    42,
			},
			wantIdentity: &aclcore.Identity{
				Email:   email,
				Name:    "Tony Stark",
				Picture: "/you-know-who-am-i.jpg",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should generates tokens with fallback identity",
			fields: fields{
				AccessTokenGenerator: func(t *testing.T) aclcore.IAccessTokenGenerator {
					generator := mocks.NewIAccessTokenGenerator(t)
					generator.On("GenerateAccessToken", email).Return(&aclcore.Authentication{}, nil)
					return generator
				},
				RefreshTokenGenerator: func(t *testing.T) aclcore.IRefreshTokenGenerator {
					generator := mocks.NewIRefreshTokenGenerator(t)
					generator.On("GenerateRefreshToken", fullSpec).Return(newRefreshToken, nil)
					return generator
				},
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewRefreshTokenRepository(t)
					repository.On("FindRefreshToken", refreshToken).Return(&fullSpec, nil)
					repository.On("DeleteRefreshToken", refreshToken).Return(nil)
					return repository
				},
				IdentityDetailsStore: func(t *testing.T) aclcore.IdentityDetailsStore {
					store := mocks.NewIdentityDetailsStore(t)
					store.On("FindIdentity", email).Return(nil, aclcore.IdentityDetailsNotFoundError)
					return store
				},
			},
			args: args{
				refreshToken: refreshToken,
			},
			wantAuthentication: &aclcore.Authentication{RefreshToken: newRefreshToken},
			wantIdentity: &aclcore.Identity{
				Email:   email,
				Name:    email.Value(),
				Picture: "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should pass through error with identity repository",
			fields: fields{
				AccessTokenGenerator: func(t *testing.T) aclcore.IAccessTokenGenerator {
					return mocks.NewIAccessTokenGenerator(t)
				},
				RefreshTokenGenerator: func(t *testing.T) aclcore.IRefreshTokenGenerator {
					return mocks.NewIRefreshTokenGenerator(t)
				},
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewRefreshTokenRepository(t)
					repository.On("FindRefreshToken", refreshToken).Return(&fullSpec, nil)
					return repository
				},
				IdentityDetailsStore: func(t *testing.T) aclcore.IdentityDetailsStore {
					store := mocks.NewIdentityDetailsStore(t)
					store.On("FindIdentity", email).Return(nil, errors.Errorf("TEST - identity error"))
					return store
				},
			},
			args: args{
				refreshToken: refreshToken,
			},
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "TEST - identity error")
			},
		},
		{
			name: "it should not issue tokens if refresh token expired",
			fields: fields{
				AccessTokenGenerator: func(t *testing.T) aclcore.IAccessTokenGenerator {
					return mocks.NewIAccessTokenGenerator(t)
				},
				RefreshTokenGenerator: func(t *testing.T) aclcore.IRefreshTokenGenerator {
					return mocks.NewIRefreshTokenGenerator(t)
				},
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewRefreshTokenRepository(t)
					repository.On("FindRefreshToken", refreshToken).Return(&aclcore.RefreshTokenSpec{
						AbsoluteExpiryTime: time.Date(2020, 12, 31, 23, 59, 59, 999, time.UTC),
					}, nil)
					repository.On("HouseKeepRefreshToken").Return(1, nil)
					return repository
				},
				IdentityDetailsStore: func(t *testing.T) aclcore.IdentityDetailsStore {
					return mocks.NewIdentityDetailsStore(t)
				},
			},
			args: args{
				refreshToken: refreshToken,
			},
			wantAuthentication: nil,
			wantIdentity:       nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.ExpiredRefreshTokenError)
			},
		},
		{
			name: "it should not issue tokens if refresh token not found",
			fields: fields{
				AccessTokenGenerator: func(t *testing.T) aclcore.IAccessTokenGenerator {
					return mocks.NewIAccessTokenGenerator(t)
				},
				RefreshTokenGenerator: func(t *testing.T) aclcore.IRefreshTokenGenerator {
					return mocks.NewIRefreshTokenGenerator(t)
				},
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewRefreshTokenRepository(t)
					repository.On("FindRefreshToken", refreshToken).Return(nil, aclcore.InvalidRefreshTokenError)
					return repository
				},
				IdentityDetailsStore: func(t *testing.T) aclcore.IdentityDetailsStore {
					return mocks.NewIdentityDetailsStore(t)
				},
			},
			args: args{
				refreshToken: refreshToken,
			},
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
				AccessTokenGenerator:   tt.fields.AccessTokenGenerator(t),
				RefreshTokenGenerator:  tt.fields.RefreshTokenGenerator(t),
				RefreshTokenRepository: tt.fields.AccessTokenRepository(t),
				IdentityDetailsStore:   tt.fields.IdentityDetailsStore(t),
			}
			gotToken, gotIdentity, err := s.AuthenticateFromRefreshToken(tt.args.refreshToken)
			if !tt.wantErr(t, err, fmt.Sprintf("AuthenticateFromRefreshToken(%v)", tt.args.refreshToken)) {
				return
			}
			assert.Equalf(t, tt.wantAuthentication, gotToken, "AuthenticateFromRefreshToken(%v)", tt.args.refreshToken)
			assert.Equalf(t, tt.wantIdentity, gotIdentity, "AuthenticateFromRefreshToken(%v)", tt.args.refreshToken)
		})
	}
}
