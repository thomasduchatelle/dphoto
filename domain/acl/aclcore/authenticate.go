package aclcore

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
)

// TokenGenerator generate an access token pre-authorising consumer to perform most operations
type TokenGenerator struct {
	PermissionsReader ScopesReader
	Config            OAuthConfig
}

func (t *TokenGenerator) GenerateAccessToken(email string) (*Authentication, error) {
	issuedAt := TimeFunc().UTC()
	expiresAt := issuedAt.Add(t.Config.ValidityDuration)

	scopes, err := t.loadUserScopes(email)
	if err != nil {
		return nil, err
	}
	if len(scopes) == 0 {
		return nil, NotPreregisteredError
	}

	tokenId, _ := uuid.NewUUID()
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, struct {
		jwt.RegisteredClaims
		customClaims
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.Config.Issuer,
			Subject:   email,
			Audience:  []string{t.Config.Issuer},
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
			NotBefore: nil,
			IssuedAt:  &jwt.NumericDate{Time: issuedAt},
			ID:        tokenId.String(),
		},
		customClaims: customClaims{
			Scopes: strings.Join(scopes, " "),
		},
	})

	signedJwt, err := generatedToken.SignedString(t.Config.SecretJwtKey)
	return &Authentication{
		AccessToken: signedJwt,
		ExpiryTime:  expiresAt,
		ExpiresIn:   int64(t.Config.ValidityDuration.Seconds()),
	}, errors.Wrapf(err, "couldn't sign the generated JWT")
}

func (t *TokenGenerator) loadUserScopes(email string) ([]string, error) {
	grants, err := t.PermissionsReader.ListUserScopes(email, ApiScope, MainOwnerScope)

	var scopes []string
	for _, grant := range grants {
		switch grant.Type {
		case ApiScope:
			scopes = append(scopes, fmt.Sprintf("api:%s", grant.ResourceId))

		case MainOwnerScope:
			scopes = append(scopes, fmt.Sprintf("%s%s", JWTScopeOwnerPrefix, grant.ResourceOwner))
		}
	}

	if len(scopes) > 0 || err != nil {
		return scopes, errors.Wrapf(err, "couldn't list grants of '%s'", email)
	}

	// second change for visitors
	grants, err = t.PermissionsReader.ListUserScopes(email, AlbumVisitorScope, MediaVisitorScope)
	if len(grants) > 0 {
		scopes = []string{"visitor"}
	}
	return scopes, errors.Wrapf(err, "couldn't list [AlbumVisitorScope, MediaVisitorScope] grants of %s", email)
}
