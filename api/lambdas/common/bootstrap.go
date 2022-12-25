package common

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/ephoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/ephoto/pkg/acl/acldynamodb"
	"github.com/thomasduchatelle/ephoto/pkg/acl/jwks"
	"github.com/thomasduchatelle/ephoto/pkg/archive"
	"github.com/thomasduchatelle/ephoto/pkg/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/ephoto/pkg/archiveadapters/asyncjobadapter"
	"github.com/thomasduchatelle/ephoto/pkg/archiveadapters/s3store"
	"github.com/thomasduchatelle/ephoto/pkg/catalog"
	"github.com/thomasduchatelle/ephoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/ephoto/pkg/catalogadapters/catalogdynamo"
	"os"
)

const (
	JWTIssuer         = "DPHOTO_JWT_ISSUER"
	JWTKeyB64         = "DPHOTO_JWT_KEY_B64"
	JWTValidity       = "DPHOTO_JWT_VALIDITY"
	DynamoDBTableName = "CATALOG_TABLE_NAME"
)

var (
	jwtDecoder      *aclcore.AccessTokenDecoder
	grantRepository acldynamodb.GrantRepository
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault(JWTValidity, "8h")
}

func appAuthConfig() aclcore.OAuthConfig {
	jwtKey, err := base64.StdEncoding.DecodeString(viper.GetString(JWTKeyB64))
	if err != nil {
		panic(fmt.Sprintf("environment variable '%s' must be encoded in base 64 [value was %s]", JWTKeyB64, viper.GetString(JWTKeyB64)))
	}

	return aclcore.OAuthConfig{
		ValidityDuration: viper.GetDuration(JWTValidity),
		Issuer:           viper.GetString(JWTIssuer),
		SecretJwtKey:     jwtKey,
	}
}

func ssoAuthenticatorPermissionReader() acldynamodb.GrantRepository {
	return acldynamodb.Must(acldynamodb.New(newSession(), viper.GetString(DynamoDBTableName), true))
}

func NewSSOAuthenticator() *aclcore.SSOAuthenticator {
	config, err := jwks.LoadIssuerConfig(aclcore.TrustedIdentityProvider...)
	if err != nil {
		panic(err)
	}

	return &aclcore.SSOAuthenticator{
		TokenGenerator: aclcore.TokenGenerator{
			PermissionsReader: ssoAuthenticatorPermissionReader(),
			Config:            appAuthConfig(),
		},
		TrustedIdentityIssuers: config,
	}
}

func AccessTokenDecoder() *aclcore.AccessTokenDecoder {
	return &aclcore.AccessTokenDecoder{
		Config: appAuthConfig(),
	}
}

// BootstrapCatalogAndArchiveDomains bootstraps all domains
func BootstrapCatalogAndArchiveDomains() archive.AsyncJobAdapter {
	BootstrapOAuthDomain()
	bootstrapCatalogDomain()
	archiveJobAdapter := BootstrapArchiveDomain()

	return archiveJobAdapter
}

// BootstrapCatalogDomain bootstraps both oauth and catalog
func BootstrapCatalogDomain() {
	BootstrapOAuthDomain()
	bootstrapCatalogDomain()
}

// BootstrapOAuthDomain only bootstraps oauth
func BootstrapOAuthDomain() {
	jwtDecoder = AccessTokenDecoder()
	grantRepository = ssoAuthenticatorPermissionReader()
}

func bootstrapCatalogDomain() {
	tableName, ok := os.LookupEnv(DynamoDBTableName)
	if !ok || tableName == "" {
		panic("CATALOG_TABLE_NAME environment variable must be set.")
	}
	dynamoAdapter := catalogdynamo.Must(catalogdynamo.NewRepository(newSession(), tableName))
	catalog.Init(dynamoAdapter, catalogarchivesync.New())
}

func BootstrapArchiveDomain() archive.AsyncJobAdapter {
	tableName, ok := os.LookupEnv(DynamoDBTableName)
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
	archiveAsyncAdapter := asyncjobadapter.New(sess, archiveJobsSnsARN, archiveJobsSqsURL, asyncjobadapter.DefaultImagesPerMessage)
	archive.Init(
		archivedynamo.Must(archivedynamo.New(sess, tableName, false)),
		s3store.Must(s3store.New(sess, storeBucketName)),
		s3store.Must(s3store.New(sess, cacheBucketName)),
		archiveAsyncAdapter,
	)

	return archiveAsyncAdapter
}

type mediaAlbumResolver struct{}

func (m mediaAlbumResolver) FindAlbumOfMedia(owner, mediaId string) (string, error) {
	return catalog.FindMediaOwnership(owner, mediaId)
}

type catalogAdapter struct{}

func (c catalogAdapter) FindAllAlbums(owner string) ([]*catalog.Album, error) {
	return catalog.FindAllAlbums(owner)
}

func (c catalogAdapter) FindAlbums(keys []catalog.AlbumId) ([]*catalog.Album, error) {
	return catalog.FindAlbums(keys)
}

func (c catalogAdapter) ListMedias(owner string, folderName string, request catalog.PageRequest) (*catalog.MediaPage, error) {
	return catalog.ListMedias(owner, folderName, request)
}

func newSession() *session.Session {
	return session.Must(session.NewSession())
}
