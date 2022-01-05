package googleoauth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"math/big"
	"net/http"
	"strings"
)

type OAuth2Config struct {
	Issuers     map[string]interface{}
	KeySupplier func() ([]byte, error)
}

type OAuth2ConfigReader struct {
	OpenIdConfigUrl string
}

type openIdConfiguration struct {
	Issuer  string `json:"issuer"`
	JwksUri string `json:"jwks_uri"`
}

type jwksKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	E   string `json:"e"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
}

type jwksResponse struct {
	Keys []jwksKey `json:"keys"`
}

func NewGoogle() *OAuth2ConfigReader {
	return &OAuth2ConfigReader{
		OpenIdConfigUrl: "https://accounts.google.com/.well-known/openid-configuration",
	}
}

func (o *OAuth2ConfigReader) Read(config *oauthmodel.Config) error {
	if config == nil {
		return errors.Errorf("config argument must not be nil")
	}

	index, err := o.readConfigIndex(o.OpenIdConfigUrl)
	if err != nil {
		return errors.Wrapf(err, "failed to read JWKS config from %s", o.OpenIdConfigUrl)
	}

	jwks, err := o.readJWKS(index.JwksUri)
	if err != nil {
		return errors.Wrapf(err, "invalid JWKS URL %s", index.JwksUri)
	}

	if len(config.TrustedIssuers) == 0 {
		config.TrustedIssuers = make(map[string]oauthmodel.IssuerOAuth2Config)
	}

	config.TrustedIssuers[strings.TrimLeft(index.Issuer, "https://")] = oauthmodel.IssuerOAuth2Config{
		ConfigSource: o.OpenIdConfigUrl,
		PublicKeysLookup: func(method oauthmodel.TokenMethod) (interface{}, error) {
			if method.Algorithm != jwt.SigningMethodRS256.Alg() {
				return nil, errors.Errorf("[OAuth2JwksConfigReader] %s algorithm is not supported.", method.Algorithm)
			}

			var kids []string
			for _, key := range jwks.Keys {
				kids = append(kids, key.Kid)

				if key.Kid == method.Kid {
					return o.parseJwks(key)
				}
			}

			return nil, errors.Errorf("kid '%s' is not defined in %s [%s] config. Available kids are: %s.", method.Kid, index.Issuer, o.OpenIdConfigUrl, strings.Join(kids, ", "))
		},
	}

	return nil
}

func (o *OAuth2ConfigReader) readConfigIndex(url string) (*openIdConfiguration, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	config := &openIdConfiguration{}
	err = json.NewDecoder(response.Body).Decode(config)
	return config, err
}

func (o *OAuth2ConfigReader) readJWKS(uri string) (*jwksResponse, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jwks := &jwksResponse{}
	err = json.NewDecoder(response.Body).Decode(jwks)
	return jwks, err
}

func (o *OAuth2ConfigReader) parseJwks(key jwksKey) (*rsa.PublicKey, error) {
	n, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}

	e, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	ei := big.NewInt(0).SetBytes(e).Int64()
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(n),
		E: int(ei),
	}, nil
}
