package aclcore

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// mockHTTPClient is a mock implementation of HTTPClient
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// mockScopesReader is a mock implementation of ScopesReader
type mockScopesReader struct {
	ListScopesByUserFunc func(ctx context.Context, email usermodel.UserId, types ...ScopeType) ([]*Scope, error)
}

func (m *mockScopesReader) ListScopesByUser(ctx context.Context, email usermodel.UserId, types ...ScopeType) ([]*Scope, error) {
	if m.ListScopesByUserFunc != nil {
		return m.ListScopesByUserFunc(ctx, email, types...)
	}
	return nil, nil
}

func (m *mockScopesReader) FindScopesById(ids ...ScopeId) ([]*Scope, error) {
	return nil, errors.New("not implemented")
}

// createMockUserInfoResponse creates a mock UserInfo HTTP response
func createMockUserInfoResponse(email string, statusCode int) *http.Response {
	if statusCode != http.StatusOK {
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error": "invalid_token"}`)),
		}
	}

	userInfo := userInfoResponse{
		Sub:           "user123",
		Email:         email,
		EmailVerified: "true",
		Username:      "testuser",
	}
	body, _ := json.Marshal(userInfo)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(body)),
	}
}

func TestCognitoTokenDecoder_Decode(t *testing.T) {
	// Generate RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	expectedIssuer := "https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_TEST123"
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Mock issuer config
	mockIssuerConfig := OAuth2IssuerConfig{
		ConfigSource:     "test",
		UserInfoEndpoint: "https://test.example.com/oauth2/userInfo",
		PublicKeysLookup: func(method OAuthTokenMethod) (interface{}, error) {
			if method.Algorithm == "RS256" && method.Kid == "test-kid" {
				return &privateKey.PublicKey, nil
			}
			return nil, errors.Errorf("unknown kid: %s", method.Kid)
		},
	}

	// Mock HTTP client that returns successful UserInfo response
	mockHTTP := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return createMockUserInfoResponse("test@example.com", http.StatusOK), nil
		},
	}

	decoder := &CognitoTokenDecoder{
		ExpectedIssuer: expectedIssuer,
		IssuerConfig:   mockIssuerConfig,
		HTTPClient:     mockHTTP,
		ScopesReader:   nil, // Will be set per test
		Now:            func() time.Time { return fixedTime },
	}

	type tokenArgs struct {
		privateKey    *rsa.PrivateKey
		signingMethod jwt.SigningMethod
		claims        cognitoClaims
		kid           string
		useHS256      bool // special flag for HMAC tokens
	}

	type mockHTTPSetup struct {
		email      string
		statusCode int
		err        error
	}

	type mockScopesSetup struct {
		returnScopes []*Scope
		returnError  error
	}

	getValidTokenArgs := func() tokenArgs {
		return tokenArgs{
			privateKey:    privateKey,
			signingMethod: jwt.SigningMethodRS256,
			claims: cognitoClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    expectedIssuer,
					Subject:   "user123",
					ExpiresAt: jwt.NewNumericDate(fixedTime.Add(1 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(fixedTime),
				},
				TokenUse:      "access",
				CognitoGroups: []string{},
			},
			kid:      "test-kid",
			useHS256: false,
		}
	}

	tests := []struct {
		name            string
		tokenArgs       func() tokenArgs
		mockHTTP        *mockHTTPSetup
		mockScopes      *mockScopesSetup
		wantErr         assert.ErrorAssertionFunc
		wantErrContains string
		wantSubject     usermodel.UserId
		wantScopes      map[string]interface{}
		wantOwner       *ownermodel.Owner
	}{
		{
			name: "it should accept a valid Cognito token with owner scope in DB",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				email:      "test@example.com",
				statusCode: http.StatusOK,
			},
			mockScopes: &mockScopesSetup{
				returnScopes: []*Scope{
					{
						Type:          MainOwnerScope,
						GrantedTo:     usermodel.NewUserId("test@example.com"),
						ResourceOwner: ownermodel.Owner("test@example.com"),
					},
				},
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("test@example.com"),
			wantScopes: map[string]interface{}{
				"owner:test@example.com": nil,
			},
			wantOwner: func() *ownermodel.Owner { o := ownermodel.Owner("test@example.com"); return &o }(),
		},
		{
			name: "it should accept a valid Cognito token with api:admin scope in DB",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				email:      "admin@example.com",
				statusCode: http.StatusOK,
			},
			mockScopes: &mockScopesSetup{
				returnScopes: []*Scope{
					{
						Type:       ApiScope,
						GrantedTo:  usermodel.NewUserId("admin@example.com"),
						ResourceId: "admin",
					},
				},
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("admin@example.com"),
			wantScopes: map[string]interface{}{
				"api:admin": nil,
			},
			wantOwner: nil,
		},
		{
			name: "it should accept a valid Cognito token with both owner and api scopes in DB",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				email:      "owner@example.com",
				statusCode: http.StatusOK,
			},
			mockScopes: &mockScopesSetup{
				returnScopes: []*Scope{
					{
						Type:       ApiScope,
						GrantedTo:  usermodel.NewUserId("owner@example.com"),
						ResourceId: "admin",
					},
					{
						Type:          MainOwnerScope,
						GrantedTo:     usermodel.NewUserId("owner@example.com"),
						ResourceOwner: ownermodel.Owner("owner@example.com"),
					},
				},
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("owner@example.com"),
			wantScopes: map[string]interface{}{
				"api:admin":               nil,
				"owner:owner@example.com": nil,
			},
			wantOwner: func() *ownermodel.Owner { o := ownermodel.Owner("owner@example.com"); return &o }(),
		},
		{
			name: "it should accept a visitor with album/media scopes",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				email:      "visitor@example.com",
				statusCode: http.StatusOK,
			},
			mockScopes: &mockScopesSetup{
				returnScopes: []*Scope{
					{
						Type:          AlbumVisitorScope,
						GrantedTo:     usermodel.NewUserId("visitor@example.com"),
						ResourceOwner: ownermodel.Owner("owner@example.com"),
						ResourceId:    "album123",
					},
				},
			},
			wantErr:     assert.NoError,
			wantSubject: usermodel.NewUserId("visitor@example.com"),
			wantScopes: map[string]interface{}{
				"visitor": nil,
			},
			wantOwner: nil,
		},
		{
			name: "it should reject a user not pre-registered (no scopes in DB)",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				email:      "unregistered@example.com",
				statusCode: http.StatusOK,
			},
			mockScopes: &mockScopesSetup{
				returnScopes: []*Scope{}, // No scopes
			},
			wantErr:         assert.Error,
			wantErrContains: "user must be pre-registered",
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
		{
			name: "it should reject when UserInfo endpoint returns error",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				email:      "",
				statusCode: http.StatusUnauthorized,
			},
			wantErr:         assert.Error,
			wantErrContains: "UserInfo endpoint returned status 401",
		},
		{
			name: "it should reject when UserInfo endpoint fails",
			tokenArgs: func() tokenArgs {
				return getValidTokenArgs()
			},
			mockHTTP: &mockHTTPSetup{
				err: errors.New("network error"),
			},
			wantErr:         assert.Error,
			wantErrContains: "failed to call UserInfo endpoint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.tokenArgs()
			var tokenString string

			// Setup mock HTTP client for this test
			if tt.mockHTTP != nil {
				decoder.HTTPClient = &mockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						if tt.mockHTTP.err != nil {
							return nil, tt.mockHTTP.err
						}
						return createMockUserInfoResponse(tt.mockHTTP.email, tt.mockHTTP.statusCode), nil
					},
				}
			}

			// Setup mock scopes reader for this test
			if tt.mockScopes != nil {
				decoder.ScopesReader = &mockScopesReader{
					ListScopesByUserFunc: func(ctx context.Context, email usermodel.UserId, types ...ScopeType) ([]*Scope, error) {
						if tt.mockScopes.returnError != nil {
							return nil, tt.mockScopes.returnError
						}
						// Filter scopes by type if requested
						if len(types) > 0 {
							var filtered []*Scope
							for _, scope := range tt.mockScopes.returnScopes {
								for _, reqType := range types {
									if scope.Type == reqType {
										filtered = append(filtered, scope)
										break
									}
								}
							}
							return filtered, nil
						}
						return tt.mockScopes.returnScopes, nil
					},
				}
			}

			if args.useHS256 {
				// Create HMAC token
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, args.claims)
				tokenString, err = token.SignedString([]byte("secret"))
				assert.NoError(t, err)
			} else {
				// Create RSA token
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

// Helper function to create Cognito JWT tokens for testing
func createCognitoToken(t *testing.T, privateKey *rsa.PrivateKey, claims cognitoClaims, kid string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if kid != "" {
		token.Header["kid"] = kid
	}

	tokenString, err := token.SignedString(privateKey)
	assert.NoError(t, err)
	return tokenString
}
