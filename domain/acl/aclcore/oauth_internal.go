package aclcore

import (
	"github.com/golang-jwt/jwt/v4"
)

type customClaims struct {
	Scopes string
}

type accessTokenClaims struct {
	customClaims
	jwt.RegisteredClaims
}
