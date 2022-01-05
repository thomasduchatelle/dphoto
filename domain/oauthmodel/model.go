package oauthmodel

import (
	"fmt"
	"time"
)

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

// Claims custom to DPhoto used for authorisation
type Claims struct {
	Owners map[string]string `json:"owners,omitempty"`
	Scopes string            `json:"scopes"`
}

type TokenMethod struct {
	Algorithm string
	Kid       string
}

type IssuerOAuth2Config struct {
	ConfigSource     string
	PublicKeysLookup func(method TokenMethod) (interface{}, error)
}

type Config struct {
	TrustedIssuers   map[string]IssuerOAuth2Config // TrustedIssuers is the list of accepted 'iss', and their public key
	Issuer           string                        // Issuer is user in the generated JWT for both 'iss' and 'aud'
	ValidityDuration string                        // ValidityDuration uses time.ParseDuration format (ex: '8h')
	SecretJwtKey     []byte                        // SecretJwtKey is the key used to sign and validate DPhoto JWT
}

type UserRoles struct {
	IsUserManager bool
	Owners        map[string]string
}

type UserRepository interface {
	FindUserRoles(email string) (*UserRoles, error)
}

func (t *TokenMethod) String() string {
	return fmt.Sprintf("TokenMethod(alg=%s, kid=%s)", t.Algorithm, t.Kid)
}
func (i *IssuerOAuth2Config) String() string {
	return fmt.Sprintf("%s", i.ConfigSource)
}
