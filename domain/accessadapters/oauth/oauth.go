// Package oauth is a subdomain wrapping 'access' domain into a OAuth compatible JWT implementation
package oauth

import (
	"github.com/thomasduchatelle/dphoto/domain/access"
	"time"
)

type Oauth interface {
	// AuthenticateFromExternalIDProvider create an access token for DPhoto API from an identity token of an external provider
	AuthenticateFromExternalIDProvider(identityJWT string) (Authentication, Identity, error)

	// Authorise tests the validity of the JWT token (signature and issuer), and the presence of the scopes.
	Authorise(accessJWT string, validator func(claims Claims) error) error
}

type oauth struct {
	config     Config
	now        func() time.Time
	listGrants func(email string, resourceType ...access.ResourceType) ([]*access.Resource, error)
}

func New(config Config) Oauth {
	return NewOverride(config, time.Now, access.ListGrants)
}

func NewOverride(config Config, now func() time.Time, listGrants func(email string, resourceType ...access.ResourceType) ([]*access.Resource, error)) Oauth {
	return &oauth{
		config:     config,
		now:        now,
		listGrants: listGrants,
	}
}

// Authentication is generated open successful authentication
type Authentication struct {
	AccessToken string
	ExpiryTime  time.Time
	ExpiresIn   int64 // ExpiresIn is the number of seconds before access token expires
}

type Identity struct {
	Email   string
	Name    string
	Picture string
}

type Claims interface {
	HasApiAccess(api string) error
	IsOwnerOf(owner string) error
}
