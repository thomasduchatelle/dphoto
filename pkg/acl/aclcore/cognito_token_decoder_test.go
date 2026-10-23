package aclcore

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestCognitoTokenDecoder_Decode(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	expectedIssuer := "https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_TEST123"
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	mockIssuerConfig := OAuth2IssuerConfig{
		ConfigSource: "test",
		PublicKeysLookup: func(method OAuthTokenMethod) (interface{}, error) {
			if method.Algorithm == "RS256" && method.Kid == "test-kid" {
				return &privateKey.PublicKey, nil
			}
			return nil, errors.Errorf("unknown kid: %s", method.Kid)
		},
	}

	decoder := &CognitoTokenDecoder{
		ExpectedIssuer: expectedIssuer,
		IssuerConfig:   mockIssuerConfig,
		Now:            func() time.Time { return fixedTime },
	}

	type tokenArgs struct {
		privateKey *rsa.PrivateKey
		claims     cognitoClaims
		kid        string
		useHS256   bool
	}

	getValidTokenArgs := func() tokenArgs {
		return tokenArgs{
			privateKey: privateKey,
			claims: cognitoClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    expectedIssuer,
					Subject:   "user123",
					ExpiresAt: jwt.NewNumericDate(fixedTime.Add(1 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(fixedTime),
				},
				TokenUse: "access",
				Email:    "test@example.com",
			},
			kid: "test-kid",
		}
	}

	ownerPointer := func(value string) *ownermodel.Owner {
		o := ownermodel.Owner(value)
		return &o
	}

	tests := []struct {
		name            string
		tokenArgs       func() tokenArgs
		wantErr         assert.ErrorAssertionFunc
		wantErrContains string
		wantSubject     usermodel.UserId
		wantScopes      map[string]interface{}
		wantOwner       *ownermodel.Owner
	}{
		{
			name: "it should read the email, owner and scopes from the injected claims",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.Email = "test@example.com"
				args.claims.Owner = "test@example.com"
				args.claims.Scopes = "owner:test@example.com"
				return args
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("test@example.com"),
			wantScopes:  map[string]interface{}{"owner:test@example.com": nil},
			wantOwner:   ownerPointer("test@example.com"),
		},
		{
			name: "it should read an api scope with no owner",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.Email = "admin@example.com"
				args.claims.Scopes = "api:admin"
				return args
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("admin@example.com"),
			wantScopes:  map[string]interface{}{"api:admin": nil},
			wantOwner:   nil,
		},
		{
			name: "it should read both api and owner scopes",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.Email = "owner@example.com"
				args.claims.Owner = "owner@example.com"
				args.claims.Scopes = "api:admin owner:owner@example.com"
				return args
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("owner@example.com"),
			wantScopes:  map[string]interface{}{"api:admin": nil, "owner:owner@example.com": nil},
			wantOwner:   ownerPointer("owner@example.com"),
		},
		{
			name: "it should read a visitor scope with no owner",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.Email = "visitor@example.com"
				args.claims.Scopes = "visitor"
				return args
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("visitor@example.com"),
			wantScopes:  map[string]interface{}{"visitor": nil},
			wantOwner:   nil,
		},
		{
			name: "it should reject a token missing the email claim",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.Email = ""
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "missing the 'dphoto:email' claim",
		},
		{
			name: "it should reject a token with invalid signature",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.privateKey = wrongKey
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "invalid Cognito JWT",
		},
		{
			name: "it should reject a token with wrong issuer",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.Issuer = "https://wrong-issuer.com"
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "doesn't match expected Cognito issuer",
		},
		{
			name: "it should reject a token with wrong token_use",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.TokenUse = "id"
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "token_use 'id' is not supported",
		},
		{
			name: "it should reject an expired token",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.ExpiresAt = jwt.NewNumericDate(fixedTime.Add(-1 * time.Hour))
				args.claims.IssuedAt = jwt.NewNumericDate(fixedTime.Add(-2 * time.Hour))
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "invalid Cognito JWT",
		},
		{
			name: "it should reject a token not yet valid",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.claims.ExpiresAt = jwt.NewNumericDate(fixedTime.Add(2 * time.Hour))
				args.claims.NotBefore = jwt.NewNumericDate(fixedTime.Add(1 * time.Hour))
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "invalid Cognito JWT",
		},
		{
			name: "it should reject a token with missing kid in header",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.kid = ""
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "unknown kid",
		},
		{
			name: "it should reject a token with unsupported algorithm",
			tokenArgs: func() tokenArgs {
				args := getValidTokenArgs()
				args.useHS256 = true
				return args
			},
			wantErr:         assert.Error,
			wantErrContains: "invalid Cognito JWT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.tokenArgs()
			var tokenString string

			if args.useHS256 {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, args.claims)
				tokenString, err = token.SignedString([]byte("secret"))
				assert.NoError(t, err)
			} else {
				tokenString = createCognitoToken(t, args.privateKey, args.claims, args.kid)
			}

			claims, err := decoder.Decode(tokenString)

			if !tt.wantErr(t, err, fmt.Sprintf("Decode(%v)", tt.name)) {
				return
			}

			if err != nil && tt.wantErrContains != "" {
				assert.Contains(t, err.Error(), tt.wantErrContains)
				assert.True(t, errors.Is(err, AccessUnauthorisedError))
			}

			if err == nil {
				assert.Equal(t, tt.wantSubject, claims.Subject)
				assert.Equal(t, tt.wantScopes, claims.Scopes)
				assert.Equal(t, tt.wantOwner, claims.Owner)
			}
		})
	}
}

func createCognitoToken(t *testing.T, privateKey *rsa.PrivateKey, claims cognitoClaims, kid string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if kid != "" {
		token.Header["kid"] = kid
	}

	tokenString, err := token.SignedString(privateKey)
	assert.NoError(t, err)
	return tokenString
}
