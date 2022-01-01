package common

import (
	"encoding/base64"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"github.com/thomasduchatelle/dphoto/domain/oauthadapters/googleoauth"
	"github.com/thomasduchatelle/dphoto/domain/oauthadapters/userrepositorystatic"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"os"
)

func Bootstrap() {
	initOAuthDomain()
}

func initOAuthDomain() {
	secretJwtKeyB64, b := os.LookupEnv("SECRET_JWT_KEY_B64")
	if !b {
		panic("environment variable 'SECRET_JWT_KEY_B64' is mandatory.")
	}

	secretJwtKey, err := base64.StdEncoding.DecodeString(secretJwtKeyB64)
	if err != nil {
		panic("environment variable 'SECRET_JWT_KEY_B64' must be encoded in base 64.")
	}

	oauth.UserRepository = userrepositorystatic.New()
	oauth.Config = oauthmodel.Config{
		TrustedIssuers:   nil,
		Issuer:           "https://dphoto-dev.duchatelle.io",
		ValidityDuration: "8h",
		SecretJwtKey:     secretJwtKey,
	}
	err = googleoauth.NewGoogle().Read(&oauth.Config)
	if err != nil {
		panic(err)
	}
}
