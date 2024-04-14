package common

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclidentitydynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclrefreshdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclscopedynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/jwks"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/asyncjobadapter"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/s3store"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
	"os"
	"time"
)

const (
	JWTIssuer            = "DPHOTO_JWT_ISSUER"
	JWTKeyB64            = "DPHOTO_JWT_KEY_B64"
	JWTValidity          = "DPHOTO_JWT_VALIDITY"
	RefreshTokenValidity = "DPHOTO_REFRESH_TOKEN_VALIDITY"
	DynamoDBTableName    = "CATALOG_TABLE_NAME"
)

var (
	jwtDecoder      *aclcore.AccessTokenDecoder
	grantRepository aclscopedynamodb.GrantRepository
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault(JWTValidity, "15m")
	viper.SetDefault(RefreshTokenValidity, "")
}

func appAuthConfig() aclcore.OAuthConfig {
	jwtKey, err := base64.StdEncoding.DecodeString(viper.GetString(JWTKeyB64))
	if err != nil {
		panic(fmt.Sprintf("environment variable '%s' must be encoded in base 64 [value was %s]", JWTKeyB64, viper.GetString(JWTKeyB64)))
	}

	refreshDurations := map[aclcore.RefreshTokenPurpose]time.Duration{
		aclcore.RefreshTokenPurposeWeb: 90 * 24 * time.Hour,
	}
	for token, validity := range viper.GetStringMapString(RefreshTokenValidity) {
		refreshDurations[aclcore.RefreshTokenPurpose(token)], err = time.ParseDuration(validity)
	}

	return aclcore.OAuthConfig{
		AccessDuration:  viper.GetDuration(JWTValidity),
		RefreshDuration: refreshDurations,
		Issuer:          viper.GetString(JWTIssuer),
		SecretJwtKey:    jwtKey,
	}
}

func ssoAuthenticatorPermissionReader() aclscopedynamodb.GrantRepository {
	return aclscopedynamodb.Must(aclscopedynamodb.New(newV1Session(), viper.GetString(DynamoDBTableName)))
}

func newRefreshTokenRepository() aclcore.RefreshTokenRepository {
	return aclrefreshdynamodb.Must(aclrefreshdynamodb.New(newV1Session(), viper.GetString(DynamoDBTableName)))
}

func NewAuthenticators() (*aclcore.SSOAuthenticator, *aclcore.RefreshTokenAuthenticator) {
	config, err := jwks.LoadIssuerConfig(aclcore.TrustedIdentityProvider...)
	if err != nil {
		panic(err)
	}

	identityDetailsStore := getIdentityDetailsStore()
	refreshTokenRepository := newRefreshTokenRepository()

	refreshTokenGenerator := aclcore.RefreshTokenGenerator{
		RefreshTokenRepository: refreshTokenRepository,
		RefreshDuration:        appAuthConfig().RefreshDuration,
	}
	accessTokenGenerator := aclcore.AccessTokenGenerator{
		PermissionsReader: ssoAuthenticatorPermissionReader(),
		Config:            appAuthConfig(),
	}

	return &aclcore.SSOAuthenticator{
			AccessTokenGenerator:   accessTokenGenerator,
			RefreshTokenGenerator:  &refreshTokenGenerator,
			IdentityDetailsStore:   identityDetailsStore,
			TrustedIdentityIssuers: config,
		},
		&aclcore.RefreshTokenAuthenticator{
			AccessTokenGenerator:   &accessTokenGenerator,
			RefreshTokenGenerator:  &refreshTokenGenerator,
			RefreshTokenRepository: refreshTokenRepository,
			IdentityDetailsStore:   identityDetailsStore,
		}
}

func getIdentityDetailsStore() aclidentitydynamodb.IdentityRepository {
	return aclidentitydynamodb.Must(aclidentitydynamodb.New(newV1Session(), viper.GetString(DynamoDBTableName)))
}

func NewLogout() *aclcore.Logout {
	return &aclcore.Logout{RevokeAccessTokenAdapter: newRefreshTokenRepository()}
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
	dynamoAdapter := catalogdynamo.Must(catalogdynamo.NewRepository(newV1Session(), tableName))
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

	sess := newV1Session()
	cfg := newV2Config()
	archiveAsyncAdapter := asyncjobadapter.New(cfg, archiveJobsSnsARN, archiveJobsSqsURL, asyncjobadapter.DefaultImagesPerMessage)
	archive.Init(
		archivedynamo.Must(archivedynamo.New(sess, tableName)),
		s3store.Must(s3store.New(cfg, storeBucketName)),
		s3store.Must(s3store.New(cfg, cacheBucketName)),
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

func newV1Session() *session.Session {
	return session.Must(session.NewSession())
}
func newV2Config() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return cfg
}
