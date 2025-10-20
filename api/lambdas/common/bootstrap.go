package common

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclidentitydynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclrefreshdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclscopedynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/jwks"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var (
	jwtDecoder      *aclcore.AccessTokenDecoder
	grantRepository aclscopedynamodb.GrantRepository
	Factory         pkgfactory.Factory
)

func init() {
	initViper()

	var err error
	Factory, err = pkgfactory.StartAWSCloudBuilder(new(LambdaViperNames)).WithAdvancedAWSAsyncFeatures().Build(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to start AWS cloud factory: %v", err))
	}
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

func CognitoTokenDecoder() (*aclcore.CognitoTokenDecoder, error) {
	userPoolId := viper.GetString(CognitoUserPoolId)
	region := viper.GetString(CognitoRegion)
	
	if userPoolId == "" || region == "" {
		return nil, fmt.Errorf("COGNITO_USER_POOL_ID and COGNITO_REGION must be set")
	}
	
	// Cognito JWKS URL format: https://cognito-idp.{region}.amazonaws.com/{userPoolId}/.well-known/jwks.json
	cognitoIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolId)
	openIdConfigUrl := fmt.Sprintf("%s/.well-known/openid-configuration", cognitoIssuer)
	
	config, err := jwks.LoadIssuerConfig(openIdConfigUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to load Cognito JWKS config: %w", err)
	}
	
	return &aclcore.CognitoTokenDecoder{
		CognitoIssuers: config,
	}, nil
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
	ctx := context.TODO()
	Factory.InitArchive(ctx)
	return Factory.ArchiveAsyncJobAdapter(ctx)
}

func newV2Config() aws.Config {
	ctx := context.TODO()
	return MustAWSFactory(ctx).GetCfg()
}

func MustAWSFactory(ctx context.Context) awsfactory.AWSFactory {
	return pkgfactory.AWSFactory(ctx)
}
