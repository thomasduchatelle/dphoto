package config

import (
	"github.com/spf13/viper"
)

type ViperAWSName struct{}

func (v *ViperAWSName) DynamoDBName() string {
	return viper.GetString(CatalogDynamodbTable)
}

func (v *ViperAWSName) ArchiveMainBucketName() string {
	return viper.GetString(ArchiveMainBucketName)

}

func (v *ViperAWSName) ArchiveCacheBucketName() string {
	return viper.GetString(ArchiveCacheBucketName)
}

func (v *ViperAWSName) ArchiveJobsSNSARN() string {
	return viper.GetString(ArchiveJobsSNSARN)
}

func (v *ViperAWSName) ArchiveJobsSQSURL() string {
	return viper.GetString(ArchiveJobsSQSURL)
}

func (v *ViperAWSName) ArchiveRelocateJobsSQSURL() string {
	panic("ArchiveRelocateJobsSQSURL is not defined outside AWS LAMBDA environment ; advanced async configuration cannot be used.")
}

func (v *ViperAWSName) CognitoUserPoolId() string {
	return viper.GetString(CognitoUserPoolId)
}
