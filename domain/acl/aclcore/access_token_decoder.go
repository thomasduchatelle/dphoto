package aclcore

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type AccessTokenDecoder struct {
	Config OAuthConfig
	Now    func() time.Time // Now is defaulted to time.Now
}

func (a *AccessTokenDecoder) Decode(accessToken string) (Claims, error) {
	if a.Now != nil {
		jwt.TimeFunc = a.Now
	} else {
		jwt.TimeFunc = time.Now
	}

	claims := new(accessTokenClaims)
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return a.Config.SecretJwtKey, nil
	})
	if err != nil {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "invalid JWT, %s", err.Error())
	}

	if claims.Issuer != a.Config.Issuer {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "'%s' issuer not accepted", claims.Issuer)
	}
	if !containsAudience(claims.Audience, a.Config.Issuer) {
		return Claims{}, errors.Wrapf(AccessUnauthorisedError, "%s is not in the audience list %s", a.Config.Issuer, strings.Join(claims.Audience, ", "))
	}

	scopes := make(map[string]interface{})
	owner := ""
	for _, scope := range strings.Split(claims.Scopes, " ") {
		scopes[scope] = nil
		if strings.HasPrefix(scope, JWTScopeOwnerPrefix) {
			owner = scope[len(JWTScopeOwnerPrefix):]
		}
	}

	return Claims{
		Subject: claims.Subject,
		Scopes:  scopes,
		Owner:   owner,
	}, nil
}

func containsAudience(audiences jwt.ClaimStrings, issuer string) bool {
	for _, audience := range audiences {
		if audience == issuer {
			return true
		}
	}

	return false
}
