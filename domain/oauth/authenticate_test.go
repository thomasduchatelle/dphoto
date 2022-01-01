package oauth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	a := assert.New(t)

	Config = oauthmodel.Config{
		TrustedIssuers: map[string]oauthmodel.IssuerOAuth2Config{
			"accounts.google.com": {
				Name:         "mocked-google",
				ConfigSource: "unitTest",
				PublicKeysLookup: func(method oauthmodel.TokenMethod) (interface{}, error) {
					if method.Algorithm == "HS512" {
						return []byte("ExternalJWTSecret"), nil
					}
					return nil, errors.Errorf("key for %s not found", method)
				},
			},
		},
		Issuer:           "https://dphoto.unit.test",
		ValidityDuration: "12s",
		SecretJwtKey:     []byte("DPhotoJwtSecret"),
	}
	Now = func() time.Time {
		return time.Unix(315532800, 0)
	}
	jwt.TimeFunc = Now
	userRepositoryMock := new(mocks.UserRepository)
	UserRepository = userRepositoryMock

	userRepositoryMock.On("FindUserRoles", "tony@stark.com").Return(&oauthmodel.UserRoles{
		IsUserManager: true,
		Owners: map[string]string{
			"tony@stark.com": "ADMIN",
		},
	}, nil)
	userRepositoryMock.On("FindUserRoles", "peter@stark.com").Return(&oauthmodel.UserRoles{
		IsUserManager: false,
		Owners:        nil,
	}, nil)

	okJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyODk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.m4fmV7k63JhFwT_ZNtAg6O5xvJZQvGt3yx_Xrr5Yjln4PeXF70jcp31A3qDwIA5ah2X9ZmjZWRbU3_Xbm3LTlg"
	unregisteredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoicGV0ZXJAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.0-6HL6oW7MyCyXq-yXtYTXThvk90AIAQaJ9MkARiE4I6ixXF-UQnCQtl29jBA-xrwFet6D9NCFmBR95KUNOI4w"
	wrongISSJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJ3cm9uZ0lTUyIsImF6cCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsImF1ZCI6InF3ZXJ0eS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJlbWFpbCI6InRvbnlAc3RhcmsuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJRQVpXU1hFRENSRlYiLCJuYW1lIjoiVG9ueSBTdGFyayBha2EgSXJvbm1hbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vdG9ueXN0YXJrLXBpY3R1cmUiLCJnaXZlbl9uYW1lIjoiVG9ueSIsImZhbWlseV9uYW1lIjoiU3RhcmsiLCJsb2NhbGUiOiJlbi1HQiIsImlhdCI6MzE1NTMyNzAwLCJleHAiOjMxNTUzMjg5OSwianRpIjoiM2RlNzk5ODYxM2ExYTRmYjhkYTNkZTc5OTg2MTNhMWE0ZmI4ZGEifQ.Olo5ok8FOMk3jP1aJlQG2l4rrVtiPfVxKXky_tcRWM3BCZmdQvJ-m5sgmztc-hwy6Dm-SQpSeWvc7Jgf_neC-w"
	expiredJwtString := "eyJhbGciOiJIUzUxMiIsImtpZCI6IjAzZTg0YWVkNGVmNDQzMTAxNGU4NjE3NTY3ODY0YzRlZmFhYWVkZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoicXdlcnR5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTIzNDU2Nzg5MCIsImVtYWlsIjoidG9ueUBzdGFyay5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFBWldTWEVEQ1JGViIsIm5hbWUiOiJUb255IFN0YXJrIGFrYSBJcm9ubWFuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS90b255c3RhcmstcGljdHVyZSIsImdpdmVuX25hbWUiOiJUb255IiwiZmFtaWx5X25hbWUiOiJTdGFyayIsImxvY2FsZSI6ImVuLUdCIiwiaWF0IjozMTU1MzI3MDAsImV4cCI6MzE1NTMyNzk5LCJqdGkiOiIzZGU3OTk4NjEzYTFhNGZiOGRhM2RlNzk5ODYxM2ExYTRmYjhkYSJ9.8DWKVd3Xh2WXgwTaRONhpxLZs_G2dRLJkz6Qtw3VC-SJYzVivHGyWbUt1TG8GrKq5-a_CZC_UvbpSxzV68skng"

	tests := []struct {
		name            string
		tokenString     string
		assertAuth      func(string, oauthmodel.Authentication)
		wantIden        oauthmodel.Identity
		wantErrContains string
	}{
		{
			"it should exchange a valid identity JWT into an access token",
			okJwtString,
			func(name string, auth oauthmodel.Authentication) {
				a.Equal(time.Date(1980, 1, 1, 0, 0, 12, 0, time.UTC), auth.ExpiryTime, name)
				a.Equal(int64(12), auth.ExpiresIn, name)

				type internalClaims struct {
					oauthmodel.Claims
					jwt.RegisteredClaims
				}

				token, err := jwt.ParseWithClaims(auth.AccessToken, &internalClaims{}, func(token *jwt.Token) (interface{}, error) {
					return Config.SecretJwtKey, nil
				})
				if a.NoError(err, name) {
					claims := token.Claims.(*internalClaims)
					a.Equal(map[string]string{
						"tony@stark.com": "ADMIN",
					}, claims.Owners, name)
					a.Equal(Config.Issuer, claims.Issuer, name)
					a.Equal(jwt.ClaimStrings{Config.Issuer}, claims.Audience, name)
					a.Equal("tony@stark.com", claims.Subject, name)
					a.Equal("user-manager-admin owner", claims.Scopes, name)
				}
			},
			oauthmodel.Identity{
				Email:   "tony@stark.com",
				Name:    "Tony Stark aka Ironman",
				Picture: "https://lh3.googleusercontent.com/a-/tonystark-picture",
			},
			"",
		},
		{
			"it should not let unregistered user log in",
			unregisteredJwtString,
			nil,
			oauthmodel.Identity{},
			"must be pre-registered",
		},
		{
			"it should not accept JWT from non-approved issuers",
			wrongISSJwtString,
			nil,
			oauthmodel.Identity{},
			"Issuer 'wrongISS' is not supported",
		},
		{
			"it should not accept expired JWT",
			expiredJwtString,
			nil,
			oauthmodel.Identity{},
			"token is expired",
		},
	}

	for _, tt := range tests {
		gotAuth, gotIden, err := Authenticate(tt.tokenString)
		if tt.wantErrContains != "" && a.Error(err, tt.name) {
			a.Contains(err.Error(), tt.wantErrContains, tt.name)

		} else if tt.wantErrContains == "" && a.NoError(err, tt.name) {
			a.Equal(tt.wantIden, gotIden, tt.name)
			tt.assertAuth(tt.name, gotAuth)
		}
	}
}
