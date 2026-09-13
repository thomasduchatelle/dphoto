package aclcore_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestAuthenticate(t *testing.T) {
	aclcore.TimeFunc = func() time.Time {
		return time.Unix(315532800, 0)
	}
	jwt.TimeFunc = aclcore.TimeFunc

	const email = usermodel.UserId("tony@stark.com")
	const owner = ownermodel.Owner("tony@stark.com")

	okJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyODk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.m4fmV7k63JhFwT_ZNtAg6O5xvJZQvGt3yx_Xrr5Yjln4PeXF70jcp31A3qDwIA5ah2X9ZmjZWRbU3_Xbm3LTlg"
	unregisteredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoicGV0ZXJAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.0-6HL6oW7MyCyXq-yXtYTXThvk90AIAQaJ9MkARiE4I6ixXF-UQnCQtl29jBA-xrwFet6D9NCFmBR95KUNOI4w"
	wrongISSJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ3cm9uZ0lTUyIsImF6cCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsImF1ZCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJlbWFpbCI6InRvbnlAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.Olo5ok8FOMk3jP1aJlQG2l4rrVtiPfVxKXky_tcRWM3BCZmdQvJ-m5sgmztc-hwy6Dm-SQpSeWvc7Jgf_neC-w"
	expiredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyNzk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.8DWKVd3Xh2WXgwTaRONhpxLZs_G2dRLJkz6Qtw3VC-SJYzVivHGyWbUt1TG8GrKq5-a_CZC_UvbpSxzV68skng"

	config := aclcore.OAuthConfig{
		Issuer:         "https://dphoto.unit.test",
		AccessDuration: 12 * time.Second,
		SecretJwtKey:   []byte("DPhotoJwtSecret"),
	}

	tonyGoogleIdentity := aclcore.Identity{
		Email:   email,
		Name:    "Tony Stark aka Ironman",
		Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
	}
	expectedRefreshSpec := aclcore.RefreshTokenSpec{
		Email:               email,
		RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
	}

	scopeRepositoryWithOwnerAndAdmin := func() *ScopeRepositoryInMemory {
		return NewScopeRepositoryInMemory(
			aclcore.Scope{Type: aclcore.ApiScope, GrantedTo: email, ResourceId: "admin"},
			aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: email, ResourceOwner: owner},
		)
	}
	scopeRepositoryWithVisitor := func() *ScopeRepositoryInMemory {
		return NewScopeRepositoryInMemory(
			aclcore.Scope{Type: aclcore.AlbumVisitorScope, GrantedTo: email},
		)
	}

	type fields struct {
		ScopeRepository       *ScopeRepositoryInMemory
		RefreshTokenGenerator *RefreshTokenGeneratorFake
		IdentityRepository    *IdentityRepositoryInMemory
	}

	tests := []struct {
		name                      string
		fields                    fields
		argToken                  string
		assertAuth                func(*testing.T, string, aclcore.Authentication)
		wantIdentity              aclcore.Identity
		expectStoredIdentity      *aclcore.Identity
		expectRefreshGeneratedFor []aclcore.RefreshTokenSpec
		wantErrContains           string
	}{
		{
			name: "it should exchange a valid identity JWT into an access token",
			fields: fields{
				ScopeRepository:       scopeRepositoryWithOwnerAndAdmin(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorFake(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
			},
			argToken: okJwtString,
			assertAuth: func(t *testing.T, name string, auth aclcore.Authentication) {
				assert.Equal(t, time.Date(1980, 1, 1, 0, 0, 12, 0, time.UTC), auth.ExpiryTime, name)
				assert.Equal(t, int64(12), auth.ExpiresIn, name)
				assert.Equal(t, "rt-"+string(email)+"-"+string(aclcore.RefreshTokenPurposeWeb), auth.RefreshToken, name)

				assertAccessTokenClaims(t, auth.AccessToken, config, name, []string{"api:admin", "owner:tony@stark.com"})
			},
			wantIdentity:              tonyGoogleIdentity,
			expectStoredIdentity:      &tonyGoogleIdentity,
			expectRefreshGeneratedFor: []aclcore.RefreshTokenSpec{expectedRefreshSpec},
		},
		{
			name: "it should let a pure visitor authenticate",
			fields: fields{
				ScopeRepository:       scopeRepositoryWithVisitor(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorFake(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
			},
			argToken: okJwtString,
			assertAuth: func(t *testing.T, name string, auth aclcore.Authentication) {
				assert.Equal(t, time.Date(1980, 1, 1, 0, 0, 12, 0, time.UTC), auth.ExpiryTime, name)
				assert.Equal(t, int64(12), auth.ExpiresIn, name)

				assertAccessTokenClaims(t, auth.AccessToken, config, name, []string{"visitor"})
			},
			wantIdentity:              tonyGoogleIdentity,
			expectStoredIdentity:      &tonyGoogleIdentity,
			expectRefreshGeneratedFor: []aclcore.RefreshTokenSpec{expectedRefreshSpec},
		},
		{
			name: "it should not let unregistered user log in",
			fields: fields{
				ScopeRepository:       NewScopeRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorFake(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
			},
			argToken:        unregisteredJwtString,
			wantErrContains: "must be pre-registered",
		},
		{
			name: "it should not accept JWT from non-approved issuers",
			fields: fields{
				ScopeRepository:       NewScopeRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorFake(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
			},
			argToken:        wrongISSJwtString,
			wantErrContains: "Issuer 'wrongISS' is not supported",
		},
		{
			name: "it should not accept expired JWT",
			fields: fields{
				ScopeRepository:       NewScopeRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorFake(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
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
					PermissionsReader: tt.fields.ScopeRepository,
					Config:            config,
				},
				RefreshTokenGenerator: tt.fields.RefreshTokenGenerator,
				IdentityDetailsStore:  tt.fields.IdentityRepository,
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
			if tt.wantErrContains != "" {
				if a.Error(err, tt.name) {
					a.Contains(err.Error(), tt.wantErrContains, tt.name)
				}
				return
			}

			if !a.NoError(err, tt.name) {
				return
			}
			a.Equal(tt.wantIdentity, *gotIdentity, tt.name)
			tt.assertAuth(t, tt.name, *gotAuth)

			a.Equal(tt.expectRefreshGeneratedFor, tt.fields.RefreshTokenGenerator.GeneratedFor, "RefreshTokenGenerator.GeneratedFor")

			if tt.expectStoredIdentity != nil {
				storedIdentity, findErr := tt.fields.IdentityRepository.FindIdentity(tt.expectStoredIdentity.Email)
				if a.NoError(findErr, "FindIdentity(%v)", tt.expectStoredIdentity.Email) {
					a.Equal(*tt.expectStoredIdentity, *storedIdentity, "stored identity")
				}
			}
		})
	}
}

func assertAccessTokenClaims(t *testing.T, accessToken string, config aclcore.OAuthConfig, name string, expectedScopes []string) {
	type decodedClaims struct {
		Scopes string
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(accessToken, &decodedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.SecretJwtKey, nil
	})
	if !assert.NoError(t, err, name) {
		return
	}
	claims := token.Claims.(*decodedClaims)
	assert.Equal(t, config.Issuer, claims.Issuer, name)
	assert.Equal(t, jwt.ClaimStrings{config.Issuer}, claims.Audience, name)
	assert.Equal(t, "tony@stark.com", claims.Subject, name)

	scopes := strings.Split(claims.Scopes, " ")
	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i] < scopes[j]
	})
	assert.Equal(t, expectedScopes, scopes, name)
}
