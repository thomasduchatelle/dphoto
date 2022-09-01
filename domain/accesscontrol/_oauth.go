// Package accesscontrol provides ACL features for the rest of the application, and OAUTH support
package accesscontrol

import (
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

type oauth struct {
	config     OAuthConfig
	now        func() time.Time
	listGrants func(email string, resourceType ...access.PermissionType) ([]*access.Permission, error)
}

func NewOAuth(config OAuthConfig) Oauth {
	return NewOAuthOverride(config, time.Now, ScopesReader)
}

func NewOAuthOverride(config OAuthConfig, now func() time.Time, listGrants func(email string, resourceType ...access.PermissionType) ([]*access.Permission, error)) Oauth {
	return &oauth{
		config:     config,
		now:        now,
		listGrants: listGrants,
	}
}

// TODO add convenient method to read scopes (HasApiAccess(api string) error ; IsOwnerOf(owner string) error)
