package oauth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"strings"
)

const (
	invalidTokenError         = "authenticated failed"
	invalidTokenExplicitError = "authentication failed: token invalid"
	notPreregisteredError     = "user must be pre-registered"
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

// IsInvalidTokenError returns true if authenticated failed because of an invalid token.
func IsInvalidTokenError(err error) bool {
	return err.Error() == invalidTokenError || err.Error() == invalidTokenExplicitError
}

// IsNotPreregisteredError returns true when authentication failed because user is not pre-registered
func IsNotPreregisteredError(err error) bool {
	return err.Error() == notPreregisteredError
}

// MergeErrors is a utility to use in the Oauth.Authorise validator to merge several errors
func MergeErrors(errorsList ...error) error {
	var first error
	var messages []string

	for _, err := range errorsList {
		if err != nil {
			if first == nil {
				first = err
			}
			messages = append(messages, err.Error())
		}
	}

	return errors.Wrapf(first, strings.Join(messages, ", "))
}
