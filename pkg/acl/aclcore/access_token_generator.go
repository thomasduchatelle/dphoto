package aclcore

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// AccessTokenGenerator generate an access token pre-authorising consumer to perform most operations
type AccessTokenGenerator struct {
	PermissionsReader     ScopesReader
	Config                OAuthConfig
	AccessTokenRepository RefreshTokenRepository
}

func (t *AccessTokenGenerator) GenerateAccessToken(email usermodel.UserId) (*Authentication, error) {
	issuedAt := TimeFunc().UTC()
	expiresAt := issuedAt.Add(t.Config.AccessDuration)

	ctx := context.TODO()
	scopeStrings, _, _, err := LoadUserScopes(ctx, t.PermissionsReader, email)
	if err != nil {
		return nil, err
	}
	if len(scopeStrings) == 0 {
		return nil, NotPreregisteredError
	}

	tokenId, _ := uuid.NewUUID()
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, struct {
		jwt.RegisteredClaims
		customClaims
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.Config.Issuer,
			Subject:   email.Value(),
			Audience:  []string{t.Config.Issuer},
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
			NotBefore: nil,
			IssuedAt:  &jwt.NumericDate{Time: issuedAt},
			ID:        tokenId.String(),
		},
		customClaims: customClaims{
			Scopes: strings.Join(scopeStrings, " "),
		},
	})

	signedJwt, err := generatedToken.SignedString(t.Config.SecretJwtKey)
	return &Authentication{
		AccessToken: signedJwt,
		ExpiryTime:  expiresAt,
		ExpiresIn:   int64(t.Config.AccessDuration.Seconds()),
	}, errors.Wrapf(err, "couldn't sign the generated JWT")
}
