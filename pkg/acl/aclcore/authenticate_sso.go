package aclcore

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"strings"
)

// SSOAuthenticator use a known identity token issued by a known and trusted identity provider (google, facebook, ...) to create an access token
type SSOAuthenticator struct {
	AccessTokenGenerator
	RefreshTokenGenerator  IRefreshTokenGenerator
	IdentityDetailsStore   IdentityDetailsStore
	TrustedIdentityIssuers map[string]OAuth2IssuerConfig // TrustedIdentityIssuers is the list of accepted 'iss', and their public key
}

type googleClaims struct {
	jwt.RegisteredClaims
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (s *SSOAuthenticator) AuthenticateFromExternalIDProvider(identityJWT string, refreshTokenPurpose RefreshTokenPurpose) (*Authentication, *Identity, error) {
	identity, err := s.parseGoogleIdentity(identityJWT)
	if err != nil {
		return nil, nil, err
	}

	authentication, err := s.AccessTokenGenerator.GenerateAccessToken(identity.Email)
	if err != nil {
		return nil, nil, err
	}

	if err = s.IdentityDetailsStore.StoreIdentity(identity); err != nil {
		log.WithError(err).Warnf("failed to store identity of %s. Skipping.", identity.Email)
	}

	authentication.RefreshToken, err = s.RefreshTokenGenerator.GenerateRefreshToken(RefreshTokenSpec{
		Email:               identity.Email,
		RefreshTokenPurpose: refreshTokenPurpose,
	})

	return authentication, &identity, err
}

func (s *SSOAuthenticator) parseGoogleIdentity(identityJWT string) (Identity, error) {
	identityClaims := &googleClaims{}
	token, err := jwt.ParseWithClaims(identityJWT, identityClaims, s.keyLookup)

	identity := Identity{
		Email:   usermodel.UserId(identityClaims.Email),
		Name:    identityClaims.Name,
		Picture: identityClaims.Picture,
	}
	if err != nil {
		return identity, errors.Wrapf(InvalidTokenError, "(caused by %s)", err.Error())
	}
	if !token.Valid {
		return identity, InvalidTokenExplicitError
	}

	return identity, nil
}

func (s *SSOAuthenticator) keyLookup(token *jwt.Token) (interface{}, error) {
	claims, ok := token.Claims.(*googleClaims)
	if !ok {
		return nil, errors.Errorf("claims are expected to be of googleClaims type.")
	}

	issuerName := claims.Issuer

	if issuerConfig, ok := s.TrustedIdentityIssuers[issuerName]; ok {
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
	for iss := range s.TrustedIdentityIssuers {
		issuers = append(issuers, iss)
	}

	return nil, errors.Errorf("Issuer '%s' is not supported. Trusted issuers are: %s", issuerName, strings.Join(issuers, ", "))
}
