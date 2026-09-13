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

	trustedIssuers := map[string]aclcore.OAuth2IssuerConfig{
		"accounts.google.com": {
			ConfigSource: "unitTest",
			PublicKeysLookup: func(method aclcore.OAuthTokenMethod) (interface{}, error) {
				if method.Algorithm == "HS512" {
					return []byte("ExternalJWTSecret"), nil
				}
				return nil, errors.Errorf("key for %s not found", method)
			},
		},
	}

	tonyIdentity := aclcore.Identity{
		Email:   email,
		Name:    "Tony Stark aka Ironman",
		Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
	}

	type fields struct {
		ScopeRepository       *ScopeRepositoryInMemory
		IdentityRepository    *IdentityRepositoryInMemory
		RefreshTokenGenerator *RefreshTokenGeneratorInMemory
	}
	type args struct {
		identityJWT         string
		refreshTokenPurpose aclcore.RefreshTokenPurpose
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		assertAuth        func(*testing.T, string, aclcore.Authentication)
		wantIdentity      aclcore.Identity
		wantRefreshedSpec *aclcore.RefreshTokenSpec
		wantErrContains   string
	}{
		{
			name: "it should exchange a valid identity JWT into an access token",
			fields: fields{
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.ApiScope, GrantedTo: email, ResourceId: "admin"},
					aclcore.Scope{Type: aclcore.MainOwnerScope, GrantedTo: email, ResourceOwner: owner},
				),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorInMemory(),
			},
			args: args{identityJWT: okJwtString, refreshTokenPurpose: aclcore.RefreshTokenPurposeWeb},
			assertAuth: func(t *testing.T, name string, auth aclcore.Authentication) {
				assert.Equal(t, time.Date(1980, 1, 1, 0, 0, 12, 0, time.UTC), auth.ExpiryTime, name)
				assert.Equal(t, int64(12), auth.ExpiresIn, name)
				assert.Equal(t, "rt-"+email.Value()+"-"+string(aclcore.RefreshTokenPurposeWeb), auth.RefreshToken, name)

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
					assert.Equal(t, email.Value(), claims.Subject, name)

					scopes := strings.Split(claims.Scopes, " ")
					sort.Slice(scopes, func(i, j int) bool {
						return scopes[i] < scopes[j]
					})
					assert.Equal(t, []string{"api:admin", "owner:tony@stark.com"}, scopes)
				}
			},
			wantIdentity: tonyIdentity,
			wantRefreshedSpec: &aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
			},
		},
		{
			name: "it should let a pure visitor authenticate",
			fields: fields{
				ScopeRepository: NewScopeRepositoryInMemory(
					aclcore.Scope{Type: aclcore.AlbumVisitorScope, GrantedTo: email},
				),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorInMemory(),
			},
			args: args{identityJWT: okJwtString, refreshTokenPurpose: aclcore.RefreshTokenPurposeWeb},
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
					assert.Equal(t, email.Value(), claims.Subject, name)

					scopes := strings.Split(claims.Scopes, " ")
					sort.Slice(scopes, func(i, j int) bool {
						return scopes[i] < scopes[j]
					})
					assert.Equal(t, []string{"visitor"}, scopes)
				}
			},
			wantIdentity: tonyIdentity,
			wantRefreshedSpec: &aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
			},
		},
		{
			name: "it should not let unregistered user log in",
			fields: fields{
				ScopeRepository:       NewScopeRepositoryInMemory(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorInMemory(),
			},
			args:            args{identityJWT: unregisteredJwtString, refreshTokenPurpose: aclcore.RefreshTokenPurposeWeb},
			wantErrContains: "must be pre-registered",
		},
		{
			name: "it should not accept JWT from non-approved issuers",
			fields: fields{
				ScopeRepository:       NewScopeRepositoryInMemory(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorInMemory(),
			},
			args:            args{identityJWT: wrongISSJwtString, refreshTokenPurpose: aclcore.RefreshTokenPurposeWeb},
			wantErrContains: "Issuer 'wrongISS' is not supported",
		},
		{
			name: "it should not accept expired JWT",
			fields: fields{
				ScopeRepository:       NewScopeRepositoryInMemory(),
				IdentityRepository:    NewIdentityRepositoryInMemory(),
				RefreshTokenGenerator: NewRefreshTokenGeneratorInMemory(),
			},
			args:            args{identityJWT: expiredJwtString, refreshTokenPurpose: aclcore.RefreshTokenPurposeWeb},
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
				RefreshTokenGenerator:  tt.fields.RefreshTokenGenerator,
				IdentityDetailsStore:   tt.fields.IdentityRepository,
				TrustedIdentityIssuers: trustedIssuers,
			}

			gotAuth, gotIdentity, err := authenticator.AuthenticateFromExternalIDProvider(tt.args.identityJWT, tt.args.refreshTokenPurpose)
			if tt.wantErrContains != "" && a.Error(err, tt.name) {
				a.Contains(err.Error(), tt.wantErrContains, tt.name)
				a.Empty(tt.fields.RefreshTokenGenerator.GeneratedFor, "no refresh token should be generated on error")
				a.Empty(tt.fields.IdentityRepository.Identities, "no identity should be stored on error")

			} else if tt.wantErrContains == "" && a.NoError(err, tt.name) {
				a.Equal(tt.wantIdentity, *gotIdentity, tt.name)
				tt.assertAuth(t, tt.name, *gotAuth)

				storedIdentity, findErr := tt.fields.IdentityRepository.FindIdentity(email)
				if a.NoError(findErr, "identity should be stored") {
					a.Equal(tt.wantIdentity, *storedIdentity, "stored identity")
				}
				a.Equal([]aclcore.RefreshTokenSpec{*tt.wantRefreshedSpec}, tt.fields.RefreshTokenGenerator.GeneratedFor, "RefreshTokenGenerator.GeneratedFor")
			}
		})
	}
}
