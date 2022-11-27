package accesscontrol_test

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/mocks"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	accesscontrol.TimeFunc = func() time.Time {
		return time.Unix(315532800, 0)
	}
	jwt.TimeFunc = accesscontrol.TimeFunc

	okJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyODk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.m4fmV7k63JhFwT_ZNtAg6O5xvJZQvGt3yx_Xrr5Yjln4PeXF70jcp31A3qDwIA5ah2X9ZmjZWRbU3_Xbm3LTlg"
	unregisteredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoicGV0ZXJAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.0-6HL6oW7MyCyXq-yXtYTXThvk90AIAQaJ9MkARiE4I6ixXF-UQnCQtl29jBA-xrwFet6D9NCFmBR95KUNOI4w"
	wrongISSJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ3cm9uZ0lTUyIsImF6cCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsImF1ZCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJlbWFpbCI6InRvbnlAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.Olo5ok8FOMk3jP1aJlQG2l4rrVtiPfVxKXky_tcRWM3BCZmdQvJ-m5sgmztc-hwy6Dm-SQpSeWvc7Jgf_neC-w"
	expiredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyNzk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.8DWKVd3Xh2WXgwTaRONhpxLZs_G2dRLJkz6Qtw3VC-SJYzVivHGyWbUt1TG8GrKq5-a_CZC_UvbpSxzV68skng"

	config := accesscontrol.OAuthConfig{
		Issuer:           "https://dphoto.unit.test",
		ValidityDuration: 12 * time.Second,
		SecretJwtKey:     []byte("DPhotoJwtSecret"),
	}

	tests := []struct {
		name            string
		initMocks       func(permissionReaderMock *mocks.ScopesReader)
		argToken        string
		assertAuth      func(*testing.T, string, accesscontrol.Authentication)
		wantIdentity    accesscontrol.Identity
		wantErrContains string
	}{
		{
			name: "it should exchange a valid identity JWT into an access token",
			initMocks: func(permissionReaderMock *mocks.ScopesReader) {
				permissionReaderMock.On("ListUserScopes", "tony@stark.com", accesscontrol.ApiScope, accesscontrol.MainOwnerScope).Return([]*accesscontrol.Scope{
					{
						Type:       accesscontrol.ApiScope,
						ResourceId: "admin",
					},
					{
						Type:          accesscontrol.MainOwnerScope,
						ResourceOwner: "tony@stark.com",
					},
				}, nil)
			},
			argToken: okJwtString,
			assertAuth: func(t *testing.T, name string, auth accesscontrol.Authentication) {
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
					assert.Equal(t, "tony@stark.com", claims.Subject, name)

					scopes := strings.Split(claims.Scopes, " ")
					sort.Slice(scopes, func(i, j int) bool {
						return scopes[i] < scopes[j]
					})
					assert.Equal(t, []string{"api:admin", "owner:tony@stark.com"}, scopes)
				}
			},
			wantIdentity: accesscontrol.Identity{
				Email:   "tony@stark.com",
				Name:    "Tony Stark aka Ironman",
				Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
			},
		},
		{
			name: "it should not let unregistered user log in",
			initMocks: func(permissionReaderMock *mocks.ScopesReader) {
				permissionReaderMock.On("ListUserScopes", mock.Anything, accesscontrol.ApiScope, accesscontrol.MainOwnerScope).Return(nil, nil)
			},
			argToken:        unregisteredJwtString,
			wantErrContains: "must be pre-registered",
		},
		{
			name:            "it should not accept JWT from non-approved issuers",
			initMocks:       func(permissionReaderMock *mocks.ScopesReader) {},
			argToken:        wrongISSJwtString,
			wantErrContains: "Issuer 'wrongISS' is not supported",
		},
		{
			name:            "it should not accept expired JWT",
			initMocks:       func(permissionReaderMock *mocks.ScopesReader) {},
			argToken:        expiredJwtString,
			wantErrContains: "token is expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			reader := mocks.NewScopesReader(t)
			authenticator := accesscontrol.SSOAuthenticator{
				TokenGenerator: accesscontrol.TokenGenerator{
					PermissionsReader: reader,
					Config:            config,
				},
				TrustedIdentityIssuers: map[string]accesscontrol.OAuth2IssuerConfig{
					"accounts.google.com": {
						ConfigSource: "unitTest",
						PublicKeysLookup: func(method accesscontrol.OAuthTokenMethod) (interface{}, error) {
							if method.Algorithm == "HS512" {
								return []byte("ExternalJWTSecret"), nil
							}
							return nil, errors.Errorf("key for %s not found", method)
						},
					},
				},
			}

			tt.initMocks(authenticator.TokenGenerator.PermissionsReader.(*mocks.ScopesReader))

			gotAuth, gotIdentity, err := authenticator.AuthenticateFromExternalIDProvider(tt.argToken)
			if tt.wantErrContains != "" && a.Error(err, tt.name) {
				a.Contains(err.Error(), tt.wantErrContains, tt.name)

			} else if tt.wantErrContains == "" && a.NoError(err, tt.name) {
				a.Equal(tt.wantIdentity, *gotIdentity, tt.name)
				tt.assertAuth(t, tt.name, *gotAuth)
			}
		})
	}
}
