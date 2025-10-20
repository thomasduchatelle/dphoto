package common

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	JWTIssuer             = "DPHOTO_JWT_ISSUER"
	JWTKeyB64             = "DPHOTO_JWT_KEY_B64"
	JWTValidity           = "DPHOTO_JWT_VALIDITY"
	RefreshTokenValidity  = "DPHOTO_REFRESH_TOKEN_VALIDITY"
	DynamoDBTableName     = "CATALOG_TABLE_NAME"
	StorageBucketName     = "STORAGE_BUCKET_NAME"
	CacheBucketName       = "CACHE_BUCKET_NAME"
	SNSArchiveARN         = "SNS_ARCHIVE_ARN"
	SQSArchiveURL         = "SQS_ARCHIVE_URL"
	SQSArchiveRelocateURL = "SQS_ARCHIVE_RELOCATE_URL"
	CognitoUserPoolId     = "COGNITO_USER_POOL_ID"
)

func initViper() {
	viper.AutomaticEnv()

	viper.SetDefault(JWTValidity, "15m")
	viper.SetDefault(RefreshTokenValidity, "")
}

type LambdaViperNames struct{}

func (l *LambdaViperNames) DynamoDBName() string {
	tableName := viper.GetString(DynamoDBTableName)
	if tableName == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", DynamoDBTableName))
	}
	return tableName
}

func (l *LambdaViperNames) ArchiveMainBucketName() string {
	storeBucketName := viper.GetString(StorageBucketName)
	if storeBucketName == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", StorageBucketName))
	}
	return storeBucketName
}

func (l *LambdaViperNames) ArchiveCacheBucketName() string {
	cacheBucketName := viper.GetString(CacheBucketName)
	if cacheBucketName == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", CacheBucketName))
	}

	return cacheBucketName
}

func (l *LambdaViperNames) ArchiveJobsSNSARN() string {
	archiveJobsSnsARN := viper.GetString(SNSArchiveARN)
	if archiveJobsSnsARN == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", SNSArchiveARN))
	}

	return archiveJobsSnsARN
}

func (l *LambdaViperNames) ArchiveJobsSQSURL() string {
	archiveJobsSqsURL := viper.GetString(SQSArchiveURL)
	if archiveJobsSqsURL == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", SQSArchiveURL))
	}
	return archiveJobsSqsURL
}

func (l *LambdaViperNames) ArchiveRelocateJobsSQSURL() string {
	archiveJobsSqsURL := viper.GetString(SQSArchiveRelocateURL)
	if archiveJobsSqsURL == "" {
		panic(fmt.Sprintf("%s must be set and non-empty", SQSArchiveRelocateURL))
	}
	return archiveJobsSqsURL
}

func (l *LambdaViperNames) CognitoUserPoolId() string {
	// Optional - returns empty string if not configured
	return viper.GetString(CognitoUserPoolId)
}
