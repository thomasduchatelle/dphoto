package aclcore

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
	"time"
)

func TestAccessTokenDecoder_Decode(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	cognitoIssuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TEST"
	userEmail := "user@example.com"
	now := time.Now()

	createToken := func(email string, groups []string, tokenUse string) string {
		claims := cognitoClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    cognitoIssuer,
				Subject:   "user-sub-123",
				ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			Email:         email,
			CognitoGroups: groups,
			TokenUse:      tokenUse,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = "test-kid"
		signedToken, _ := token.SignedString(privateKey)
		return signedToken
	}

	cognitoIssuers := map[string]OAuth2IssuerConfig{
		cognitoIssuer: {
			ConfigSource: "test",
			PublicKeysLookup: func(method OAuthTokenMethod) (interface{}, error) {
				return &privateKey.PublicKey, nil
			},
		},
	}

	type args struct {
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    Claims
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should decode a valid token with owners group",
			args: args{
				accessToken: createToken(userEmail, []string{"owners"}, "access"),
			},
			want: Claims{
				Subject: usermodel.NewUserId(userEmail),
				Scopes: map[string]interface{}{
					"owner:user@example.com": nil,
				},
				Owner: func() *ownermodel.Owner {
					o := ownermodel.Owner(userEmail)
					return &o
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should decode a valid token with admins group",
			args: args{
				accessToken: createToken(userEmail, []string{"admins"}, "access"),
			},
			want: Claims{
				Subject: usermodel.NewUserId(userEmail),
				Scopes: map[string]interface{}{
					"api:admin": nil,
				},
				Owner: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should decode a valid token with visitors group",
			args: args{
				accessToken: createToken(userEmail, []string{"visitors"}, "access"),
			},
			want: Claims{
				Subject: usermodel.NewUserId(userEmail),
				Scopes: map[string]interface{}{
					"visitor": nil,
				},
				Owner: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should decode a token with multiple groups",
			args: args{
				accessToken: createToken(userEmail, []string{"admins", "owners"}, "access"),
			},
			want: Claims{
				Subject: usermodel.NewUserId(userEmail),
				Scopes: map[string]interface{}{
					"api:admin":              nil,
					"owner:user@example.com": nil,
				},
				Owner: func() *ownermodel.Owner {
					o := ownermodel.Owner(userEmail)
					return &o
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should reject a token with no groups",
			args: args{
				accessToken: createToken(userEmail, []string{}, "access"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, AccessUnauthorisedError)
			},
		},
		{
			name: "it should reject a token with wrong token_use",
			args: args{
				accessToken: createToken(userEmail, []string{"owners"}, "id"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, AccessUnauthorisedError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := &AccessTokenDecoder{
				CognitoIssuers: cognitoIssuers,
				Now:            func() time.Time { return now },
			}

			got, err := decoder.Decode(tt.args.accessToken)

			if !tt.wantErr(t, err) {
				return
			}

			if err == nil {
				assert.Equal(t, tt.want.Subject, got.Subject)
				assert.Equal(t, tt.want.Scopes, got.Scopes)
				assert.Equal(t, tt.want.Owner, got.Owner)
			}
		})
	}
}
