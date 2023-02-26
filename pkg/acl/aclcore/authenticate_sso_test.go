package aclcore_test

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	aclcore.TimeFunc = func() time.Time {
		return time.Unix(315532800, 0)
	}
	jwt.TimeFunc = aclcore.TimeFunc

	const email = "tony@stark.com"
	const refreshToken = "1234567890qwertyuiop"

	okJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyODk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.m4fmV7k63JhFwT_ZNtAg6O5xvJZQvGt3yx_Xrr5Yjln4PeXF70jcp31A3qDwIA5ah2X9ZmjZWRbU3_Xbm3LTlg"
	unregisteredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoicGV0ZXJAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.0-6HL6oW7MyCyXq-yXtYTXThvk90AIAQaJ9MkARiE4I6ixXF-UQnCQtl29jBA-xrwFet6D9NCFmBR95KUNOI4w"
	wrongISSJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ3cm9uZ0lTUyIsImF6cCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsImF1ZCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJlbWFpbCI6InRvbnlAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.Olo5ok8FOMk3jP1aJlQG2l4rrVtiPfVxKXky_tcRWM3BCZmdQvJ-m5sgmztc-hwy6Dm-SQpSeWvc7Jgf_neC-w"
	expiredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyNzk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.8DWKVd3Xh2WXgwTaRONhpxLZs_G2dRLJkz6Qtw3VC-SJYzVivHGyWbUt1TG8GrKq5-a_CZC_UvbpSxzV68skng"

	config := aclcore.OAuthConfig{
		Issuer:         "https://dphoto.unit.test",
		AccessDuration: 12 * time.Second,
		SecretJwtKey:   []byte("DPhotoJwtSecret"),
	}

	type fields struct {
		ScopesReader          func(t *testing.T) aclcore.ScopesReader
		RefreshTokenGenerator func(t *testing.T) aclcore.IRefreshTokenGenerator
		IdentityDetailsStore  func(t *testing.T) aclcore.IdentityDetailsStore
	}

	tests := []struct {
		name            string
		fields          fields
		argToken        string
		assertAuth      func(*testing.T, string, aclcore.Authentication)
		wantIdentity    aclcore.Identity
		wantErrContains string
	}{
		{
			name: "it should exchange a valid identity JWT into an access token",
			fields: fields{
				ScopesReader: func(t *testing.T) aclcore.ScopesReader {
					reader := mocks.NewScopesReader(t)
					reader.On("ListUserScopes", email, aclcore.ApiScope, aclcore.MainOwnerScope).Return([]*aclcore.Scope{
						{
							Type:       aclcore.ApiScope,
							ResourceId: "admin",
						},
						{
							Type:          aclcore.MainOwnerScope,
							ResourceOwner: email,
						},
					}, nil)
					return reader
				},
				RefreshTokenGenerator: refreshTokenGenerator(email, refreshToken),
				IdentityDetailsStore: identityDetailsStore(aclcore.Identity{
					Email:   email,
					Name:    "Tony Stark aka Ironman",
					Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
				}),
			},
			argToken: okJwtString,
			assertAuth: func(t *testing.T, name string, auth aclcore.Authentication) {
				assert.Equal(t, time.Date(1980, 1, 1, 0, 0, 12, 0, time.UTC), auth.ExpiryTime, name)
				assert.Equal(t, int64(12), auth.ExpiresIn, name)
				assert.Equal(t, refreshToken, auth.RefreshToken, name)

				type decodedClaims struct {
					Scopes string
					jwt.RegisteredClaims
				}

				token, err := jwt.ParseWithClaims(auth.AccessToken, &decodedClaims{}, func(token *jwt.Token) (interface{}, error) {
					return config.SecretJwtKey, nil
				})
				if assert.NoError(t, err, name) {
					claims := token.Claims.(*decodedClaims)
					assert.Equal(t, config.Issuer, claims.Issuer, name)
					assert.Equal(t, jwt.ClaimStrings{config.Issuer}, claims.Audience, name)
					assert.Equal(t, email, claims.Subject, name)

					scopes := strings.Split(claims.Scopes, " ")
					sort.Slice(scopes, func(i, j int) bool {
						return scopes[i] < scopes[j]
					})
					assert.Equal(t, []string{"api:admin", "owner:tony@stark.com"}, scopes)
				}
			},
			wantIdentity: aclcore.Identity{
				Email:   email,
				Name:    "Tony Stark aka Ironman",
				Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
			},
		},
		{
			name: "it should let a pure visitor authenticate",
			fields: fields{
				ScopesReader: func(t *testing.T) aclcore.ScopesReader {
					reader := mocks.NewScopesReader(t)
					reader.On("ListUserScopes", mock.Anything, aclcore.ApiScope, aclcore.MainOwnerScope).Return(nil, nil)
					reader.On("ListUserScopes", mock.Anything, aclcore.AlbumVisitorScope, aclcore.MediaVisitorScope).Return([]*aclcore.Scope{
						{},
					}, nil)
					return reader
				},
				RefreshTokenGenerator: refreshTokenGenerator(email, refreshToken),
				IdentityDetailsStore: identityDetailsStore(aclcore.Identity{
					Email:   email,
					Name:    "Tony Stark aka Ironman",
					Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
				}),
			},
			argToken: okJwtString,
			assertAuth: func(t *testing.T, name string, auth aclcore.Authentication) {
				assert.Equal(t, time.Date(1980, 1, 1, 0, 0, 12, 0, time.UTC), auth.ExpiryTime, name)
				assert.Equal(t, int64(12), auth.ExpiresIn, name)

				type decodedClaims struct {
					Scopes string
					jwt.RegisteredClaims
				}

				token, err := jwt.ParseWithClaims(auth.AccessToken, &decodedClaims{}, func(token *jwt.Token) (interface{}, error) {
					return config.SecretJwtKey, nil
				})
				if assert.NoError(t, err, name) {
					claims := token.Claims.(*decodedClaims)
					assert.Equal(t, config.Issuer, claims.Issuer, name)
					assert.Equal(t, jwt.ClaimStrings{config.Issuer}, claims.Audience, name)
					assert.Equal(t, email, claims.Subject, name)

					scopes := strings.Split(claims.Scopes, " ")
					sort.Slice(scopes, func(i, j int) bool {
						return scopes[i] < scopes[j]
					})
					assert.Equal(t, []string{"visitor"}, scopes)
				}
			},
			wantIdentity: aclcore.Identity{
				Email:   email,
				Name:    "Tony Stark aka Ironman",
				Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
			},
		},
		{
			name: "it should not let unregistered user log in",
			fields: fields{
				ScopesReader: func(t *testing.T) aclcore.ScopesReader {
					reader := mocks.NewScopesReader(t)
					reader.On("ListUserScopes", mock.Anything, aclcore.ApiScope, aclcore.MainOwnerScope).Return(nil, nil)
					reader.On("ListUserScopes", mock.Anything, aclcore.AlbumVisitorScope, aclcore.MediaVisitorScope).Return(nil, nil)
					return reader
				},
				RefreshTokenGenerator: refreshTokenGeneratorNotCalled(),
				IdentityDetailsStore:  identityDetailsStoreNotCalled(),
			},
			argToken:        unregisteredJwtString,
			wantErrContains: "must be pre-registered",
		},
		{
			name: "it should not accept JWT from non-approved issuers",
			fields: fields{
				ScopesReader:          scopeReaderNotCalled(),
				RefreshTokenGenerator: refreshTokenGeneratorNotCalled(),
				IdentityDetailsStore:  identityDetailsStoreNotCalled(),
			},
			argToken:        wrongISSJwtString,
			wantErrContains: "Issuer 'wrongISS' is not supported",
		},
		{
			name: "it should not accept expired JWT",
			fields: fields{
				ScopesReader:          scopeReaderNotCalled(),
				RefreshTokenGenerator: refreshTokenGeneratorNotCalled(),
				IdentityDetailsStore:  identityDetailsStoreNotCalled(),
			},
			argToken:        expiredJwtString,
			wantErrContains: "token is expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			authenticator := aclcore.SSOAuthenticator{
				AccessTokenGenerator: aclcore.AccessTokenGenerator{
					PermissionsReader: tt.fields.ScopesReader(t),
					Config:            config,
				},
				RefreshTokenGenerator: tt.fields.RefreshTokenGenerator(t),
				IdentityDetailsStore:  tt.fields.IdentityDetailsStore(t),
				TrustedIdentityIssuers: map[string]aclcore.OAuth2IssuerConfig{
					"accounts.google.com": {
						ConfigSource: "unitTest",
						PublicKeysLookup: func(method aclcore.OAuthTokenMethod) (interface{}, error) {
							if method.Algorithm == "HS512" {
								return []byte("ExternalJWTSecret"), nil
							}
							return nil, errors.Errorf("key for %s not found", method)
						},
					},
				},
			}

			gotAuth, gotIdentity, err := authenticator.AuthenticateFromExternalIDProvider(tt.argToken, aclcore.RefreshTokenPurposeWeb)
			if tt.wantErrContains != "" && a.Error(err, tt.name) {
				a.Contains(err.Error(), tt.wantErrContains, tt.name)

			} else if tt.wantErrContains == "" && a.NoError(err, tt.name) {
				a.Equal(tt.wantIdentity, *gotIdentity, tt.name)
				tt.assertAuth(t, tt.name, *gotAuth)
			}
		})
	}
}

func scopeReaderNotCalled() func(t *testing.T) aclcore.ScopesReader {
	return func(t *testing.T) aclcore.ScopesReader {
		return mocks.NewScopesReader(t)
	}
}

func refreshTokenGenerator(email string, refreshToken string) func(t *testing.T) aclcore.IRefreshTokenGenerator {
	return func(t *testing.T) aclcore.IRefreshTokenGenerator {
		generator := mocks.NewIRefreshTokenGenerator(t)
		generator.On("GenerateRefreshToken", aclcore.RefreshTokenSpec{
			Email:               email,
			RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
		}).Return(refreshToken, nil)
		return generator
	}
}

func identityDetailsStore(identity aclcore.Identity) func(t *testing.T) aclcore.IdentityDetailsStore {
	return func(t *testing.T) aclcore.IdentityDetailsStore {
		store := mocks.NewIdentityDetailsStore(t)
		store.On("StoreIdentity", identity).Return(nil)
		return store
	}
}

func refreshTokenGeneratorNotCalled() func(t *testing.T) aclcore.IRefreshTokenGenerator {
	return func(t *testing.T) aclcore.IRefreshTokenGenerator {
		return mocks.NewIRefreshTokenGenerator(t)
	}
}

func identityDetailsStoreNotCalled() func(t *testing.T) aclcore.IdentityDetailsStore {
	return func(t *testing.T) aclcore.IdentityDetailsStore {
		return mocks.NewIdentityDetailsStore(t)
	}
}
