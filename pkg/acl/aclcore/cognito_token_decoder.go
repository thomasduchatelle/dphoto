package aclcore

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// CognitoTokenDecoder validates Cognito JWT access tokens signed with RSA256 and reads the DPhoto claims that the
// pre-token-generation trigger (see CognitoClaimsEnricher) stamped into the token at issuance: the user email, the
// owner, and the scopes. It therefore needs no UserInfo call nor database lookup on the authorization hot path.
// Implements TokenDecoder interface.
type CognitoTokenDecoder struct {
	ExpectedIssuer string             // e.g., https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_vnHCM2cK8
	IssuerConfig   OAuth2IssuerConfig // JWKS configuration for public key lookup
	Now            func() time.Time   // Now is defaulted to time.Now
}

type cognitoClaims struct {
	jwt.RegisteredClaims
	TokenUse string `json:"token_use"`
	Email    string `json:"dphoto:email"`
	Owner    string `json:"dphoto:owner"`
	Scopes   string `json:"dphoto:scopes"`
}

func (c *CognitoTokenDecoder) Decode(accessToken string) (Claims, error) {
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

	if claims.Issuer != c.ExpectedIssuer {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "issuer '%s' doesn't match expected Cognito issuer '%s'", claims.Issuer, c.ExpectedIssuer)
	}

	if claims.TokenUse != "access" {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "token_use '%s' is not supported (expected 'access')", claims.TokenUse)
	}

	if claims.Email == "" {
		return Claims{}, errors.Wrap(AccessUnauthorisedError, "Cognito token is missing the 'dphoto:email' claim")
	}

	scopeMap := make(map[string]interface{})
	var owner *ownermodel.Owner
	for _, scope := range strings.Fields(claims.Scopes) {
		scopeMap[scope] = nil
	}
	if claims.Owner != "" {
		value := ownermodel.Owner(claims.Owner)
		owner = &value
	}

	return Claims{
		Subject: usermodel.NewUserId(claims.Email),
		Scopes:  scopeMap,
		Owner:   owner,
	}, nil
}

func (c *CognitoTokenDecoder) keyLookup(token *jwt.Token) (interface{}, error) {
	var kid string
	if kidObj, ok := token.Header["kid"]; ok {
		kid, _ = kidObj.(string)
	}

	return c.IssuerConfig.PublicKeysLookup(OAuthTokenMethod{
		Algorithm: token.Method.Alg(),
		Kid:       kid,
	})
}
