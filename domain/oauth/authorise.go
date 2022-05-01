// Package oauth exposes function to forge and validate JWT tokens
package oauth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"strings"
)

type AuthoriseQuery struct {
	Owners map[string]string
	Scopes []string
}

type accessTokenClaims struct {
	oauthmodel.Claims
	jwt.RegisteredClaims
}

// Authorise tests the validity of the JWT token (signature and issuer), and the presence of the scopes.
func Authorise(tokenString string, queryPt *AuthoriseQuery) (*oauthmodel.Claims, error) {
	query := toObj(queryPt)

	claims := new(accessTokenClaims)
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Config.SecretJwtKey, nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "invalid token")
	}

	if claims.Issuer != Config.Issuer {
		return nil, errors.Errorf("issuer %s not accepted.", claims.Issuer)
	}
	if !contains(claims.Audience, Config.Issuer) {
		return nil, errors.Errorf("%s is not in the audience list %s", Config.Issuer, strings.Join(claims.Audience, ", "))
	}

	scopes := strings.Split(claims.Scopes, " ")
	for _, scope := range query.Scopes {
		if !contains(scopes, scope) {
			return nil, errors.Errorf("scopes too restrictive (%s is not included in %s)", strings.Join(query.Scopes, ", "), strings.Join(scopes, ", "))
		}
	}

	for owner, permission := range query.Owners {
		if granted, ok := claims.Owners[owner]; !ok || ownerPermissionLevel(granted) < ownerPermissionLevel(permission) {
			return nil, errors.Errorf("access to owners %+v hasn't been granted (granted: %+v)", query.Owners, claims.Owners)
		}
	}

	return &claims.Claims, nil
}

func toObj(queryPt *AuthoriseQuery) AuthoriseQuery {
	query := AuthoriseQuery{}
	if queryPt != nil {
		query = *queryPt
	}
	if query.Owners == nil {
		query.Owners = make(map[string]string)
	}
	return query
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}

	return false
}

func ownerPermissionLevel(permission string) int {
	switch permission {
	case "ADMIN":
		return 0100
	case "WRITE":
		return 0010
	case "READ":
		return 0001
	default:
		return 0000
	}
}

func NewAuthoriseQuery(scopes ...string) *AuthoriseQuery {
	return &AuthoriseQuery{
		Scopes: scopes,
		Owners: make(map[string]string),
	}
}

func (q *AuthoriseQuery) WithOwner(owner string, permission string) *AuthoriseQuery {
	q.Owners[owner] = permission
	return q
}
