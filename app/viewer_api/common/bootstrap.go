package common

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/stores3"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/carchivesync"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/repositorydynamo"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"github.com/thomasduchatelle/dphoto/domain/oauthadapters/googleoauth"
	"github.com/thomasduchatelle/dphoto/domain/oauthadapters/userrepositorystatic"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"os"
)

// BootstrapCatalogAndArchiveDomains bootstraps all domains
func BootstrapCatalogAndArchiveDomains() {
	BootstrapOAuthDomain()
	bootstrapCatalogDomain()
	BootstrapArchiveDomain()
}

// BootstrapCatalogDomain bootstraps both oauth and catalog
func BootstrapCatalogDomain() {
	BootstrapOAuthDomain()
	bootstrapCatalogDomain()
}

// BootstrapOAuthDomain only bootstraps oauth
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

func bootstrapCatalogDomain() {
	tableName, ok := os.LookupEnv("CATALOG_TABLE_NAME")
	if !ok || tableName == "" {
		panic("CATALOG_TABLE_NAME environment variable must be set.")
	}
	dynamoAdapter := repositorydynamo.Must(repositorydynamo.NewRepository(newSession(), tableName))
	catalog.Init(dynamoAdapter, carchivesync.New())
}

func BootstrapArchiveDomain() archive.AsyncJobAdapter {
	tableName, ok := os.LookupEnv("CATALOG_TABLE_NAME")
	if !ok || tableName == "" {
		panic("CATALOG_TABLE_NAME environment variable must be set.")
	}
	storeBucketName, ok := os.LookupEnv("STORAGE_BUCKET_NAME")
	if !ok || storeBucketName == "" {
		panic("STORAGE_BUCKET_NAME must be set and non-empty")
	}
	cacheBucketName, ok := os.LookupEnv("CACHE_BUCKET_NAME")
	if !ok || cacheBucketName == "" {
		panic("CACHE_BUCKET_NAME must be set and non-empty")
	}
	archiveJobsSnsARN, ok := os.LookupEnv("SNS_ARCHIVE_ARN")
	if !ok || archiveJobsSnsARN == "" {
		panic("SNS_ARCHIVE_ARN must be set and non-empty")
	}
	archiveJobsSqsURL, ok := os.LookupEnv("SQS_ARCHIVE_URL")
	if !ok || archiveJobsSnsARN == "" {
		panic("SQS_ARCHIVE_URL must be set and non-empty")
	}

	sess := newSession()
	archiveAsyncAdapter := asyncjobsns.New(sess, archiveJobsSnsARN, archiveJobsSqsURL, asyncjobsns.DefaultImagesPerMessage)
	archive.Init(
		arepositorydynamo.Must(arepositorydynamo.New(sess, tableName, false)),
		stores3.Must(stores3.New(sess, storeBucketName)),
		stores3.Must(stores3.New(sess, cacheBucketName)),
		archiveAsyncAdapter,
	)

	return archiveAsyncAdapter
}

func newSession() *session.Session {
	return session.Must(session.NewSession())
}
