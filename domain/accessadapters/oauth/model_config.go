package oauth

import "fmt"

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

func (t *TokenMethod) String() string {
	return fmt.Sprintf("TokenMethod(alg=%s, kid=%s)", t.Algorithm, t.Kid)
}
func (i *IssuerOAuth2Config) String() string {
	return fmt.Sprintf("%s", i.ConfigSource)
}
