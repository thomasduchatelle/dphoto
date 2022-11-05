package accesscontrol

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/access"
	"strings"
	"time"
)

func (o *oauth) AuthenticateFromExternalIDProvider(identityJWT string) (Authentication, Identity, error) {
	identity, err := o.parseGoogleIdentity(identityJWT)
	if err != nil {
		return Authentication{}, identity, err
	}

	duration, err := time.ParseDuration(o.config.ValidityDuration)
	if err != nil {
		return Authentication{}, identity, errors.Wrapf(err, "invalid configuration 'ValidityDuration'")
	}

	tokenId, _ := uuid.NewUUID()
	now := o.now().UTC()
	expiresAt := now.Add(duration)

	scopes, err := o.loadUserScopes(identity.Email)
	if err != nil || len(scopes) == 0 {
		if err != nil {
			log.WithError(err).Errorf("couldn't load user's roles %s", identity.Email)
		}
		return Authentication{}, identity, errors.Errorf(notPreregisteredError)
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, struct {
		jwt.RegisteredClaims
		customClaims
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    o.config.Issuer,
			Subject:   identity.Email,
			Audience:  []string{o.config.Issuer},
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
			NotBefore: nil,
			IssuedAt:  &jwt.NumericDate{Time: now},
			ID:        tokenId.String(),
		},
		customClaims: customClaims{
			Scopes: strings.Join(scopes, " "),
		},
	})

	signedJwt, err := generatedToken.SignedString(o.config.SecretJwtKey)
	return Authentication{
		AccessToken: signedJwt,
		ExpiryTime:  expiresAt,
		ExpiresIn:   int64(duration.Seconds()),
	}, identity, errors.Wrapf(err, "couldn't sign the generated JWT")
}

func (o *oauth) loadUserScopes(email string) ([]string, error) {
	var scopes []string

	grants, err := o.listGrants(email, access.ApiRole, access.OwnerRole)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't list grants for %s", email)
	}
	for _, grant := range grants {
		var scope string
		switch grant.Type {
		case access.ApiRole:
			scope = fmt.Sprintf("api:%s", grant.ResourceId)

		case access.OwnerRole:
			scope = fmt.Sprintf("owner:self:%s", grant.ResourceOwner)

		}
		if scope != "" {
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

type googleClaims struct {
	jwt.RegisteredClaims
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (o *oauth) parseGoogleIdentity(identityJWT string) (Identity, error) {
	identityClaims := &googleClaims{}
	token, err := jwt.ParseWithClaims(identityJWT, identityClaims, o.keyLookup)

	identity := Identity{
		Email:   identityClaims.Email,
		Name:    identityClaims.Name,
		Picture: identityClaims.Picture,
	}
	if err != nil {
		return identity, errors.Wrapf(err, invalidTokenError)
	}
	if !token.Valid {
		return identity, errors.Errorf(invalidTokenExplicitError)
	}

	return identity, nil
}

func (o *oauth) keyLookup(token *jwt.Token) (interface{}, error) {
	claims, ok := token.Claims.(*googleClaims)
	if !ok {
		return nil, errors.Errorf("claims are expected to be of googleClaims type.")
	}

	issuerName := claims.Issuer

	if issuerConfig, ok := o.config.TrustedIssuers[issuerName]; ok {
		var kid string
		if kidObj, ok := token.Header["kid"]; ok {
			kid, _ = kidObj.(string)
		}

		return issuerConfig.PublicKeysLookup(OAuthTokenMethod{
			Algorithm: token.Method.Alg(),
			Kid:       kid,
		})
	}

	// Issuer not found
	var issuers []string
	for iss, _ := range o.config.TrustedIssuers {
		issuers = append(issuers, iss)
	}

	return nil, errors.Errorf("Issuer '%s' is not supported. Trusted issuers are: %s", issuerName, strings.Join(issuers, ", "))
}
