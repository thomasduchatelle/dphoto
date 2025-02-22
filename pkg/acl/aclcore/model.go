// Package aclcore manage what an API consumer is granted to do, and issue access tokens to represent it.
//
// This package supports the following use-cases:
//
// 1. an admin can create a new owner and administrate the site
// 2. an owner can back up and retrieve his medias through the API (WEB, CLI and MOBILE)
// 3. a guest can see albums and medias that have been shared with him
// 4. an owner can share an album to the rest of the family, and contribute to family albums
package aclcore

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"time"
)

var TimeFunc = time.Now

var (
	InvalidTokenError         = errors.New("authenticated failed")
	InvalidTokenExplicitError = errors.New("authentication failed: token invalid")
	NotPreregisteredError     = errors.New("user must be pre-registered")
	AccessUnauthorisedError   = errors.New("access unauthorised") // AccessUnauthorisedError is used when the request doesn't have valid credentials (no bearer token, or invalid token)
	AccessForbiddenError      = errors.New("access forbidden")    // AccessForbiddenError is used when the request has valid credentials, but the access to the resource has been denied

	ExpiredRefreshTokenError     = errors.New("refresh token has expired")
	InvalidRefreshTokenError     = errors.New("refresh token is not valid")
	IdentityDetailsNotFoundError = errors.New("no identity details stored for this identity") // IdentityDetailsNotFoundError is an internal error between the domain and the repository

	// TrustedIdentityProvider is the default list of trusted identity providers
	TrustedIdentityProvider = []string{
		"https://accounts.google.com/.well-known/openid-configuration", // Google Identity Provider
	}
)

const (
	ApiScope              ScopeType = "api"               // ApiScope represents a set of API endpoints, like 'admin'
	MainOwnerScope        ScopeType = "owner:main"        // MainOwnerScope is limited to 1 per user, it's the tenant all backups of the user will be stored against
	AlbumVisitorScope     ScopeType = "album:visitor"     // AlbumVisitorScope gives read access to an album and the media it contains FIXME The role is NOT YET MAPPED
	AlbumContributorScope ScopeType = "album:contributor" // AlbumContributorScope gives read access and ability to contribute (add medias) to an album
	MediaVisitorScope     ScopeType = "media:visitor"     // MediaVisitorScope gives read access to medias directly

	JWTScopeOwnerPrefix = "owner:"

	RefreshTokenPurposeWeb RefreshTokenPurpose = "web" // RefreshTokenPurposeWeb is used for WEB sessions
)

// ScopeType is a type of API (admin) or a catalog resource (owner, album, ...)
type ScopeType string
type RefreshTokenPurpose string

// Scope is attached to a user (a consumer of the API) and define the role it has on resource basis
type Scope struct {
	Type          ScopeType        // Type is mandatory, it defines what fields on this structure is used and allow to filter the results
	GrantedAt     time.Time        // GrantedAt is the date the scope has been granted to the user for the first time
	GrantedTo     usermodel.UserId // GrantedTo is the consumer, usually an email address
	ResourceOwner ownermodel.Owner // ResourceOwner (optional) is used has part of the ID of the catalog resources
	ResourceId    string           // ResourceId if a unique identifier of the resource (in conjunction of the ResourceOwner for most catalog resources) ; ex: 'admin' (for 'api' type)
	ResourceName  string           // ResourceName (optional) used for user-friendly display of the shared albums // TODO is it necessary ?
}

func (s Scope) Id() ScopeId {
	return ScopeId{
		Type:          s.Type,
		GrantedTo:     s.GrantedTo,
		ResourceOwner: s.ResourceOwner,
		ResourceId:    s.ResourceId,
	}
}

// ScopeId are the properties of a Scope that identity it
type ScopeId struct {
	Type          ScopeType        // Type is mandatory, it defines what fields on this structure is used and allow to filter the results
	GrantedTo     usermodel.UserId // GrantedTo is the consumer, usually an email address
	ResourceOwner ownermodel.Owner // ResourceOwner (optional) is used has part of the ID of the catalog resources
	ResourceId    string           // ResourceId if a unique identifier of the resource (in conjunction of the ResourceOwner for most catalog resources) ; ex: 'admin' (for 'api' type)
}

type OAuthConfig struct {
	AccessDuration  time.Duration                         // AccessDuration for generated access tokens
	RefreshDuration map[RefreshTokenPurpose]time.Duration // RefreshDuration for generated refresh token (based on the purpose)
	Issuer          string                                // Issuer is the application instance ID, used in both 'iss' and 'aud'
	SecretJwtKey    []byte                                // SecretJwtKey is the key used to sign and validate DPhoto JWT
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

// Authentication is generated upon successful authentication
type Authentication struct {
	AccessToken  string
	RefreshToken string // RefreshToken is optional
	ExpiryTime   time.Time
	ExpiresIn    int64 // ExpiresIn is the number of seconds before access token expires
}

// Identity is read from token created by the Identity Provider (google, ...)
type Identity struct {
	Email   usermodel.UserId
	Name    string
	Picture string
}

type Claims struct {
	Subject usermodel.UserId       // Subject is the user id (its email)
	Scopes  map[string]interface{} // Scopes is the list of permissions stored eagerly in access token
	Owner   *ownermodel.Owner      // Owner is deviated from Scopes (extract of the MainOwnerScope)
}

func (c *Claims) AsCurrentUser() usermodel.CurrentUser {
	return usermodel.CurrentUser{
		UserId: c.Subject,
		Owner:  c.Owner,
	}
}

type RefreshTokenSpec struct {
	Email               usermodel.UserId
	RefreshTokenPurpose RefreshTokenPurpose // RefreshTokenPurpose is mandatory
	AbsoluteExpiryTime  time.Time           // AbsoluteExpiryTime will be generated from RefreshTokenPurpose if not defined
	Scopes              []string            // Scopes is the list of scopes for which an access token can be generated
}
