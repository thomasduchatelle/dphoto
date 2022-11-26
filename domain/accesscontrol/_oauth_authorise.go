package accesscontrol

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"strings"
)

func (o *oauth) DecodeAndValidate(accessJWT string, validator func(claims Claims) error) error {
	claims, err := o.parseAccessToken(accessJWT)
	if err != nil {
		return err
	}

	return validator(claims)
}

func (o *oauth) parseAccessToken(accessJWT string) (*accessTokenClaims, error) {
	claims := new(accessTokenClaims)
	_, err := jwt.ParseWithClaims(accessJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return o.config.SecretJwtKey, nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "invalid token")
	}

	if claims.Issuer != o.config.Issuer {
		return nil, errors.Errorf("issuer %s not accepted.", claims.Issuer)
	}
	if !containsAudience(claims.Audience, o.config.Issuer) {
		return nil, errors.Errorf("%s is not in the audience list %s", o.config.Issuer, strings.Join(claims.Audience, ", "))
	}

	return claims, nil
}

func containsAudience(audiences jwt.ClaimStrings, issuer string) bool {
	for _, audience := range audiences {
		if audience == issuer {
			return true
		}
	}

	return false
}
