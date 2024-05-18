package common

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
	"time"
)

var (
	jwtDecoder      *aclcore.AccessTokenDecoder
	grantRepository aclscopedynamodb.GrantRepository
)

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
	ctx := context.TODO()
	return pkgfactory.AclRepository(ctx)
}

func newRefreshTokenRepository() aclcore.RefreshTokenRepository {
	return aclrefreshdynamodb.Must(aclrefreshdynamodb.New(newV2Config(), viper.GetString(DynamoDBTableName)))
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
	return aclidentitydynamodb.Must(aclidentitydynamodb.New(newV2Config(), viper.GetString(DynamoDBTableName)))
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
	ctx := context.TODO()

	catalog.Init(pkgfactory.CatalogRepository(ctx))
}

func BootstrapArchiveDomain() archive.AsyncJobAdapter {
	tableName := viper.GetString(DynamoDBTableName)
	if tableName == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", DynamoDBTableName))
	}
	storeBucketName := viper.GetString(StorageBucketName)
	if storeBucketName == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", StorageBucketName))
	}
	cacheBucketName := viper.GetString(CacheBucketName)
	if cacheBucketName == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", CacheBucketName))
	}
	archiveJobsSnsARN := viper.GetString(SNSArchiveARN)
	if archiveJobsSnsARN == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", SNSArchiveARN))
	}
	archiveJobsSqsURL := viper.GetString(SQSArchiveURL)
	if archiveJobsSqsURL == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", SQSArchiveURL))
	}

	ctx := context.TODO()
	cfg := newV2Config()
	archiveAsyncAdapter := asyncjobadapter.New(cfg, archiveJobsSnsARN, archiveJobsSqsURL, asyncjobadapter.DefaultImagesPerMessage)
	archive.Init(
		must(archivedynamo.New(MustAWSFactory(ctx).GetDynamoDBClient(), tableName)),
		s3store.Must(s3store.New(cfg, storeBucketName)),
		s3store.Must(s3store.New(cfg, cacheBucketName)),
		archiveAsyncAdapter,
	)

	return archiveAsyncAdapter
}

type mediaAlbumResolver struct{}

func (m mediaAlbumResolver) FindAlbumOfMedia(owner ownermodel.Owner, mediaId catalog.MediaId) (catalog.AlbumId, error) {
	ownership, err := catalog.FindMediaOwnership(owner, mediaId)
	if err != nil {
		return catalog.AlbumId{}, err
	}
	return *ownership, nil
}

type catalogAdapter struct{}

func (c catalogAdapter) FindAllAlbums(owner ownermodel.Owner) ([]*catalog.Album, error) {
	return catalog.FindAllAlbums(ownermodel.Owner(owner))
}

func (c catalogAdapter) FindAlbums(keys []catalog.AlbumId) ([]*catalog.Album, error) {
	return catalog.FindAlbums(keys)
}

func (c catalogAdapter) ListMedias(albumId catalog.AlbumId, request catalog.PageRequest) (*catalog.MediaPage, error) {
	return catalog.ListMedias(albumId, request)
}

func newV2Config() aws.Config {
	ctx := context.TODO()
	return MustAWSFactory(ctx).Cfg
}

func MustAWSFactory(ctx context.Context) *awsfactory.AWSFactory {
	return must(singletons.Singleton(func() (*awsfactory.AWSFactory, error) {
		return awsfactory.NewAWSFactory(ctx, awsfactory.NewContextualConfigFactory())
	}))
}

func must[M any](value M, err error) M {
	if err != nil {
		panic(fmt.Sprintf("PANIC - %T couldn't be built: %s", *new(M), err))
	}

	return value
}
