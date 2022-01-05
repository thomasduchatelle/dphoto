package oauth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"strings"
	"time"
)

type googleClaims struct {
	jwt.RegisteredClaims
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// Authenticate verifies JWT and create an access token to use with DPhoto APIs
func Authenticate(tokenString string) (oauthmodel.Authentication, oauthmodel.Identity, error) {
	identityClaims := &googleClaims{}
	token, err := jwt.ParseWithClaims(tokenString, identityClaims, keyLookup)

	if err != nil {
		return oauthmodel.Authentication{}, oauthmodel.Identity{}, errors.Wrapf(err, "authenticated failed")
	}
	if !token.Valid {
		return oauthmodel.Authentication{}, oauthmodel.Identity{}, errors.Errorf("authentication failed: token invalid")
	}

	identity := oauthmodel.Identity{
		Email:   identityClaims.Email,
		Name:    identityClaims.Name,
		Picture: identityClaims.Picture,
	}

	duration, err := time.ParseDuration(Config.ValidityDuration)
	if err != nil {
		return oauthmodel.Authentication{}, identity, errors.Wrapf(err, "invalid configuration 'ValidityDuration'")
	}

	tokenId, _ := uuid.NewUUID()
	now := Now().UTC()
	expiresAt := now.Add(duration)

	scopes, err := loadUserScopes(identity.Email)
	if err != nil || len(scopes) == 0 {
		if err != nil {
			log.WithError(err).Errorf("Failed to load user's roles %s", identity.Email)
		}
		return oauthmodel.Authentication{}, identity, errors.Errorf("user must be pre-registered")
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, struct {
		jwt.RegisteredClaims
		oauthmodel.Claims
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Config.Issuer,
			Subject:   identityClaims.Email,
			Audience:  []string{Config.Issuer},
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
			NotBefore: nil,
			IssuedAt:  &jwt.NumericDate{Time: now},
			ID:        tokenId.String(),
		},
		Claims: oauthmodel.Claims{
			Scopes: strings.Join(scopes, " "),
			Owners: map[string]string{
				identityClaims.Email: "ADMIN",
			},
		},
	})

	signedJwt, err := generatedToken.SignedString(Config.SecretJwtKey)
	return oauthmodel.Authentication{
		AccessToken: signedJwt,
		ExpiryTime:  expiresAt,
		ExpiresIn:   int64(duration.Seconds()),
	}, identity, errors.Wrapf(err, "couldn't sign the generated JWT")
}

func loadUserScopes(email string) ([]string, error) {
	var scopes []string

	roles, err := UserRepository.FindUserRoles(email)
	if err != nil {
		return nil, errors.Wrapf(err, "authentication failed while loading user roles")
	}

	if roles.IsUserManager {
		scopes = append(scopes, "user-manager-admin")
	}
	if len(roles.Owners) > 0 {
		scopes = append(scopes, "owner")
	}

	return scopes, nil
}

func keyLookup(token *jwt.Token) (interface{}, error) {
	claims, ok := token.Claims.(*googleClaims)
	if !ok {
		return nil, errors.Errorf("claims are expected to be of googleClaims type.")
	}

	issuerName := claims.Issuer

	if issuerConfig, ok := Config.TrustedIssuers[issuerName]; ok {
		var kid string
		if kidObj, ok := token.Header["kid"]; ok {
			kid, _ = kidObj.(string)
		}

		return issuerConfig.PublicKeysLookup(oauthmodel.TokenMethod{
			Algorithm: token.Method.Alg(),
			Kid:       kid,
		})
	}

	// Issuer not found
	var issuers []string
	for iss, _ := range Config.TrustedIssuers {
		issuers = append(issuers, iss)
	}

	return nil, errors.Errorf("Issuer '%s' is not supported. Trusted issuers are: %s", issuerName, strings.Join(issuers, ", "))
}
