package aclcore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// HTTPClient is an interface for making HTTP requests (for mocking in tests)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CognitoTokenDecoder validates Cognito JWT tokens signed with RSA256 and enriches claims with scopes from the database.
// Unlike legacy DPhoto JWT tokens which have scopes baked in, Cognito tokens require a database lookup to fetch user permissions.
// Implements TokenDecoder interface.
type CognitoTokenDecoder struct {
	ExpectedIssuer string             // e.g., https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_vnHCM2cK8
	IssuerConfig   OAuth2IssuerConfig // JWKS configuration for public key lookup
	HTTPClient     HTTPClient         // HTTPClient for calling UserInfo endpoint to retrieve user email
	ScopesReader   ScopesReader       // ScopesReader for fetching user scopes from database (required)
	Now            func() time.Time   // Now is defaulted to time.Now
}

type cognitoClaims struct {
	jwt.RegisteredClaims
	CognitoGroups []string `json:"cognito:groups"`
	TokenUse      string   `json:"token_use"` // Should be "id" or "access"
}

type userInfoResponse struct {
	Sub                 string `json:"sub"`
	Email               string `json:"email"`
	EmailVerified       string `json:"email_verified"`
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified string `json:"phone_number_verified"`
	Username            string `json:"username"`
}

func (c *CognitoTokenDecoder) Decode(accessToken string) (Claims, error) {
	ctx := context.TODO()

	if c.Now != nil {
		jwt.TimeFunc = c.Now
	} else {
		jwt.TimeFunc = time.Now
	}

	claims := new(cognitoClaims)
	token, err := jwt.ParseWithClaims(accessToken, claims, c.keyLookup)
	if err != nil {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "invalid Cognito JWT: %s", err.Error())
	}

	if !token.Valid {
		return Claims{}, errors.Wrap(AccessUnauthorisedError, "Cognito token is not valid")
	}

	// Validate issuer
	if claims.Issuer != c.ExpectedIssuer {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "issuer '%s' doesn't match expected Cognito issuer '%s'", claims.Issuer, c.ExpectedIssuer)
	}

	// Validate token use (should be "access" token for authorization)
	if claims.TokenUse != "access" {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "token_use '%s' is not supported (expected 'access')", claims.TokenUse)
	}

	// Get email from UserInfo endpoint
	email, err := c.getUserEmail(accessToken)
	if err != nil {
		return Claims{}, errors.Wrap(AccessUnauthorisedError, err.Error())
	}

	// Load scopes from database using the shared canonical implementation
	userId := usermodel.NewUserId(email)
	_, scopeMap, owner, err := LoadUserScopes(ctx, c.ScopesReader, userId)
	if err != nil {
		return Claims{}, errors.Wrap(AccessUnauthorisedError, err.Error())
	}

	return Claims{
		Subject: userId,
		Scopes:  scopeMap,
		Owner:   owner,
	}, nil
}

// getUserEmail calls the Cognito UserInfo endpoint to retrieve the user's email
func (c *CognitoTokenDecoder) getUserEmail(accessToken string) (string, error) {
	// Get UserInfo endpoint URL from issuer config
	userInfoURL := c.IssuerConfig.UserInfoEndpoint
	if userInfoURL == "" {
		return "", errors.New("UserInfo endpoint not configured in issuer config")
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create UserInfo request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Use default HTTP client if not provided (for backward compatibility)
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to call UserInfo endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.Errorf("UserInfo endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var userInfo userInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", errors.Wrapf(err, "failed to parse UserInfo response")
	}

	if userInfo.Email == "" {
		return "", errors.New("email not found in UserInfo response")
	}

	return userInfo.Email, nil
}

func (c *CognitoTokenDecoder) keyLookup(token *jwt.Token) (interface{}, error) {
	// Extract kid from header
	var kid string
	if kidObj, ok := token.Header["kid"]; ok {
		kid, _ = kidObj.(string)
	}

	// Use IssuerConfig to look up public key
	return c.IssuerConfig.PublicKeysLookup(OAuthTokenMethod{
		Algorithm: token.Method.Alg(),
		Kid:       kid,
	})
}
