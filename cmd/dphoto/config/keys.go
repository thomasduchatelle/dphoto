package config

const (
	Localstack                  = "aws.localstack" // Localstack set to true will ignore AwsRegion, AwsKey, and AwsSecret
	AwsRegion                   = "aws.region"
	AwsKey                      = "aws.key"
	AwsSecret                   = "aws.secret"
	ArchiveDynamodbTable        = "archive.dynamodb.table"
	ArchiveMainBucketName       = "archive.main.bucketName"
	ArchiveCacheBucketName      = "archive.cache.bucketName"
	ArchiveJobsSNSARN           = "archive.sns.arn"
	ArchiveJobsSQSURL           = "archive.sqs.url"
	BackupCacheDirectory        = "backup.cache.dir"
	BackupConcurrencyAnalyser   = "backup.concurrency.analyser"
	BackupConcurrencyCataloguer = "backup.concurrency.cataloguer"
	BackupConcurrencyUploader   = "backup.concurrency.uploader"
	CatalogDynamodbTable        = "catalog.dynamodb.table"
	LocalHome                   = "home.dir"
	Owner                       = "owner"
)
