package common

import (
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

const (
	JWTIssuer            = "DPHOTO_JWT_ISSUER"
	JWTKeyB64            = "DPHOTO_JWT_KEY_B64"
	JWTValidity          = "DPHOTO_JWT_VALIDITY"
	RefreshTokenValidity = "DPHOTO_REFRESH_TOKEN_VALIDITY"
	DynamoDBTableName    = "CATALOG_TABLE_NAME"
	StorageBucketName    = "STORAGE_BUCKET_NAME"
	CacheBucketName      = "CACHE_BUCKET_NAME"
	SNSArchiveARN        = "SNS_ARCHIVE_ARN"
	SQSArchiveURL        = "SQS_ARCHIVE_URL"
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault(JWTValidity, "15m")
	viper.SetDefault(RefreshTokenValidity, "")

	pkgfactory.AWSNames = new(LambdaViperNames)
}

type LambdaViperNames struct{}

func (l *LambdaViperNames) DynamoDBName() string {
	return viper.GetString(DynamoDBTableName)
}

func (l *LambdaViperNames) ArchiveMainBucketName() string {
	return viper.GetString(StorageBucketName)
}

func (l *LambdaViperNames) ArchiveCacheBucketName() string {
	return viper.GetString(CacheBucketName)
}

func (l *LambdaViperNames) ArchiveJobsSNSARN() string {
	return viper.GetString(SNSArchiveARN)
}

func (l *LambdaViperNames) ArchiveJobsSQSURL() string {
	return viper.GetString(SQSArchiveURL)
}
