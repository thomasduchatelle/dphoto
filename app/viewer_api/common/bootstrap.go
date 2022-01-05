package common

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamo"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"github.com/thomasduchatelle/dphoto/domain/oauthadapters/googleoauth"
	"github.com/thomasduchatelle/dphoto/domain/oauthadapters/userrepositorystatic"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"os"
)

func Bootstrap() {
	BootstrapOAuthDomain()
}

func BootstrapCatalogDomain(owner string) {
	//bucketName, _ := os.LookupEnv("STORAGE_BUCKET_NAME")
	tableName, ok := os.LookupEnv("CATALOG_TABLE_NAME")
	if !ok || tableName == "" {
		panic("CATALOG_TABLE_NAME environment variable must be set.")
	}
	catalog.Repository = dynamo.Must(dynamo.NewRepository(session.Must(session.NewSession()), owner, tableName))
}

func BootstrapOAuthDomain() {
	secretJwtKeyB64, b := os.LookupEnv("DPHOTO_JWT_KEY_B64")
	if !b {
		panic("environment variable 'DPHOTO_JWT_KEY_B64' is mandatory.")
	}
	secretJwtKey, err := base64.StdEncoding.DecodeString(secretJwtKeyB64)
	if err != nil {
		panic("environment variable 'DPHOTO_JWT_KEY_B64' must be encoded in base 64.")
	}

	jwtIssuer, b := os.LookupEnv("DPHOTO_JWT_ISSUER")
	if !b {
		panic("environment variable 'DPHOTO_JWT_ISSUER' is mandatory.")
	}
	jwtValidity, b := os.LookupEnv("DPHOTO_JWT_VALIDITY")
	if !b {
		jwtValidity = "8h"
	}

	oauth.UserRepository = userrepositorystatic.New()
	oauth.Config = oauthmodel.Config{
		Issuer:           jwtIssuer,
		ValidityDuration: jwtValidity,
		SecretJwtKey:     secretJwtKey,
	}

	err = googleoauth.NewGoogle().Read(&oauth.Config)
	if err != nil {
		panic(err)
	}
}
