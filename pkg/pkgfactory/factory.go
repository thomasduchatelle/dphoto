package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

var (
	AWSConfigFactory = awsfactory.NewContextualConfigFactory() // AWSConfigFactory can be overridden to use other AWS authentication means (default on AWS Default config)
	AWSNames         AWSAdapterNames                           // Names provides the config required by the adapters

)

type AWSAdapterNames interface {
	DynamoDBName() string
	ArchiveMainBucketName() string
	ArchiveCacheBucketName() string
	ArchiveJobsSNSARN() string
	ArchiveJobsSQSURL() string
}

func AWSFactory(ctx context.Context) *awsfactory.AWSFactory {
	return singletons.MustSingleton(func() (*awsfactory.AWSFactory, error) {
		return awsfactory.NewAWSFactory(ctx, AWSConfigFactory)
	})
}
