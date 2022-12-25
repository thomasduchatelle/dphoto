package config

const (
	AwsRegion                   = "aws.region"
	AwsKey                      = "aws.key"
	AwsSecret                   = "aws.secret"
	ArchiveDynamodbTable        = "archive.dynamodb.table"
	ArchiveMainBucketName       = "archive.main.bucketName"
	ArchiveCacheBucketName      = "archive.cache.bucketName"
	ArchiveJobsSNSARN           = "archive.sns.arn"
	ArchiveJobsSQSURL           = "archive.sqs.url"
	BackupConcurrencyAnalyser   = "backup.concurrency.analyser"
	BackupConcurrencyCataloguer = "backup.concurrency.cataloguer"
	BackupConcurrencyUploader   = "backup.concurrency.uploader"
	CatalogDynamodbTable        = "catalog.dynamodb.table"
	LocalHome                   = "home.dir"
	Owner                       = "owner"
)
