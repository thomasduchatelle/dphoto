package accesscontrol

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"strings"
)

type customClaims struct {
	Scopes string
}

type accessTokenClaims struct {
	customClaims
	jwt.RegisteredClaims
}

func (a *accessTokenClaims) HasApiAccess(api string) error {
	return a.containsScope("api:" + api)
}

func (a *accessTokenClaims) IsOwnerOf(owner string) error {
	return a.containsScope("owner:" + owner)
}

func (a *accessTokenClaims) containsScope(expectedScope string) error {
	for _, scope := range strings.Split(a.customClaims.Scopes, " ") {
		if scope == expectedScope {
			return nil
		}
	}

	return errors.Errorf("'%s' scope not found in [%s]", expectedScope, a.customClaims.Scopes)
}
