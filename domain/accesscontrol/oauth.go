// Package accesscontrol provides ACL features for the rest of the application, and OAUTH support
package accesscontrol

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/access"
	"time"
)

var (
	AccessForbiddenError = errors.Errorf("access forbidden") // AccessForbiddenError should be used directly, or wrapped, when a function is failing because of too limited access.
)

type Oauth interface {
	// AuthenticateFromExternalIDProvider create an access token for DPhoto API from an identity token of an external provider
	AuthenticateFromExternalIDProvider(identityJWT string) (Authentication, Identity, error)

	// DecodeAndValidate tests the validity of the JWT token (signature and issuer), and the presence of the scopes.
	DecodeAndValidate(accessJWT string) (Claims, error)
}

type OAuthConfig struct {
	TrustedIssuers   map[string]OAuth2IssuerConfig // TrustedIssuers is the list of accepted 'iss', and their public key
	Issuer           string                        // Issuer is user in the generated JWT for both 'iss' and 'aud'
	ValidityDuration string                        // ValidityDuration uses time.ParseDuration format (ex: '8h')
	SecretJwtKey     []byte                        // SecretJwtKey is the key used to sign and validate DPhoto JWT
}

type OAuthTokenMethod struct {
	Algorithm string
	Kid       string
}

func (t *OAuthTokenMethod) String() string {
	return fmt.Sprintf("OAuthTokenMethod(alg=%s, kid=%s)", t.Algorithm, t.Kid)
}

type OAuth2IssuerConfig struct {
	ConfigSource     string
	PublicKeysLookup func(method OAuthTokenMethod) (interface{}, error)
}

func (i *OAuth2IssuerConfig) String() string {
	return fmt.Sprintf("%s", i.ConfigSource)
}

type oauth struct {
	config     OAuthConfig
	now        func() time.Time
	listGrants func(email string, resourceType ...access.PermissionType) ([]*access.Permission, error)
}

func NewOAuth(config OAuthConfig) Oauth {
	return NewOAuthOverride(config, time.Now, ListUserPermissions)
}

func NewOAuthOverride(config OAuthConfig, now func() time.Time, listGrants func(email string, resourceType ...access.PermissionType) ([]*access.Permission, error)) Oauth {
	return &oauth{
		config:     config,
		now:        now,
		listGrants: listGrants,
	}
}

// Authentication is generated upon successful authentication
type Authentication struct {
	AccessToken string
	ExpiryTime  time.Time
	ExpiresIn   int64 // ExpiresIn is the number of seconds before access token expires
}

// Identity is read from token created by the Identity Provider (google, ...)
type Identity struct {
	Email   string
	Name    string
	Picture string
}

type Claims struct {
	Subject string                 // Subject is the user id (its email)
	Scopes  map[string]interface{} // Scopes is the list of permissions stored eagerly in access token
}

// TODO add convenient method to read scopes (HasApiAccess(api string) error ; IsOwnerOf(owner string) error)

// IsInvalidTokenError returns true if authenticated failed because of an invalid token.
func IsInvalidTokenError(err error) bool {
	return err.Error() == invalidTokenError || err.Error() == invalidTokenExplicitError
}

// IsNotPreregisteredError returns true when authentication failed because user is not pre-registered
func IsNotPreregisteredError(err error) bool {
	return err.Error() == notPreregisteredError
}