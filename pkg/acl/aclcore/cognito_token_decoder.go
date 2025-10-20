package aclcore

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"time"
)

type CognitoTokenDecoder struct {
	CognitoIssuers map[string]OAuth2IssuerConfig
	Now            func() time.Time
}

type cognitoClaims struct {
	jwt.RegisteredClaims
	Email          string   `json:"email"`
	CognitoGroups  []string `json:"cognito:groups"`
	TokenUse       string   `json:"token_use"`
	CognitoUsername string  `json:"cognito:username"`
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
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "invalid JWT, %s", err.Error())
	}

	if !token.Valid {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "token is not valid")
	}

	if claims.TokenUse != "access" {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "token_use must be 'access', got '%s'", claims.TokenUse)
	}

	// Check if user belongs to at least one required group
	if len(claims.CognitoGroups) == 0 {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "user does not belong to any group")
	}

	// Build scopes from Cognito groups
	scopes := make(map[string]interface{})
	var owner *ownermodel.Owner

	for _, group := range claims.CognitoGroups {
		switch group {
		case "admins":
			scopes["api:admin"] = nil
		case "owners":
			// For owners, create a scope with their owner ID
			// The owner ID is derived from the username (email)
			ownerValue := ownermodel.Owner(claims.Email)
			owner = &ownerValue
			scopes[JWTScopeOwnerPrefix+string(ownerValue)] = nil
		case "visitors":
			scopes["api:visitor"] = nil
		}
	}

	return Claims{
		Subject: usermodel.NewUserId(claims.Email),
		Scopes:  scopes,
		Owner:   owner,
	}, nil
}

func (c *CognitoTokenDecoder) keyLookup(token *jwt.Token) (interface{}, error) {
	claims, ok := token.Claims.(*cognitoClaims)
	if !ok {
		return nil, errors.Errorf("claims are expected to be of cognitoClaims type")
	}

	issuerName := claims.Issuer

	if issuerConfig, ok := c.CognitoIssuers[issuerName]; ok {
		var kid string
		if kidObj, ok := token.Header["kid"]; ok {
			kid, _ = kidObj.(string)
		}

		return issuerConfig.PublicKeysLookup(OAuthTokenMethod{
			Algorithm: token.Method.Alg(),
			Kid:       kid,
		})
	}

	return nil, errors.Errorf("issuer '%s' is not a trusted Cognito issuer", issuerName)
}
