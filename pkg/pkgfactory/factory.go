package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

var (
	AWSNames AWSAdapterNames // Names provides the config required by the adapters

	factory *AWSCloud // factory supports deprecated implementation while migrating out of global OverriddenAWSFactory and AWSNames
)

// Factory is the builder of the application ; any direct variables are deprecated in favour of moving to the appropriate factory implementation.
type Factory interface {
	ArchiveFactory
	CatalogFactory

	// InitArchive shouldn't be used directly but is exposed to support legacy implementation
	InitArchive(ctx context.Context)
}

type ArchiveFactory interface {
	ArchiveAsyncJobAdapter(ctx context.Context) archive.AsyncJobAdapter
}

type AWSAdapterNames interface {
	DynamoDBName() string
	ArchiveMainBucketName() string
	ArchiveCacheBucketName() string
	ArchiveJobsSNSARN() string
	ArchiveJobsSQSURL() string
	ArchiveRelocateJobsSQSURL() string
}

type AWSCloud struct {
	awsfactory.AWSFactory
	ArchiveFactory
	*SimpleCatalogFactory
	Names AWSAdapterNames
}

type AWSCloudBuilder struct {
	advancedAsyncFeatures bool
	names                 AWSAdapterNames
	awsFactory            awsfactory.AWSFactory
	err                   []error
}

// StartAWSCloudBuilder creates a version of the application directly connected to AWS cloud using DynamoDB and S3.
func StartAWSCloudBuilder(names AWSAdapterNames) *AWSCloudBuilder {
	return &AWSCloudBuilder{
		names: names,
	}
}

// OverridesAWSFactory allows to use an alternative AWS configuration where credentials are not contextual (unlike lambdas)
func (a *AWSCloudBuilder) OverridesAWSFactory(factory awsfactory.AWSFactory, err error) *AWSCloudBuilder {
	if err != nil {
		a.err = append(a.err, err)
	} else {
		a.awsFactory = factory
	}
	return a
}

// WithAdvancedAWSAsyncFeatures enable the use of SNS/SQS to process asynchronously the archive jobs. (required lambdas to be listening the messages)
func (a *AWSCloudBuilder) WithAdvancedAWSAsyncFeatures() *AWSCloudBuilder {
	a.advancedAsyncFeatures = true
	return a
}

// Build creates the application factory ; and set legacy global variables
func (a *AWSCloudBuilder) Build(ctx context.Context) (*AWSCloud, error) {
	if len(a.err) > 0 {
		return nil, a.err[0]
	}

	if a.awsFactory == nil {
		var err error
		a.awsFactory, err = awsfactory.ContextualAWSFactory(ctx)
		if err != nil {
			return nil, err
		}
	}

	AWSNames = a.names
	factory = &AWSCloud{
		AWSFactory:     a.awsFactory,
		ArchiveFactory: new(SyncArchiveFactory),
		SimpleCatalogFactory: &SimpleCatalogFactory{
			ArchiveAdapterForCatalog: new(SyncArchiveAdapterForCatalog),
		},
		Names: a.names,
	}

	if a.advancedAsyncFeatures {
		factory.ArchiveFactory = new(AsyncArchiveFactory)
		factory.SimpleCatalogFactory.ArchiveAdapterForCatalog = &ASyncArchiveAdapterForCatalog{
			AWSFactory:      factory.AWSFactory,
			AWSAdapterNames: a.names,
		}
	}

	return factory, nil
}

func AWSFactory(ctx context.Context) awsfactory.AWSFactory {
	return singletons.MustSingletonKey("AWSFactory", func() (awsfactory.AWSFactory, error) {
		return factory.AWSFactory, nil
	})
}

type StaticAWSAdapterNames struct {
	DynamoDBNameValue              string
	ArchiveMainBucketNameValue     string
	ArchiveCacheBucketNameValue    string
	ArchiveJobsSNSARNValue         string
	ArchiveJobsSQSURLValue         string
	ArchiveRelocateJobsSQSURLValue string
}

func (s StaticAWSAdapterNames) DynamoDBName() string {
	return s.DynamoDBNameValue
}

func (s StaticAWSAdapterNames) ArchiveMainBucketName() string {
	return s.ArchiveMainBucketNameValue
}

func (s StaticAWSAdapterNames) ArchiveCacheBucketName() string {
	return s.ArchiveCacheBucketNameValue
}

func (s StaticAWSAdapterNames) ArchiveJobsSNSARN() string {
	return s.ArchiveJobsSNSARNValue
}

func (s StaticAWSAdapterNames) ArchiveJobsSQSURL() string {
	return s.ArchiveJobsSQSURLValue
}

func (s StaticAWSAdapterNames) ArchiveRelocateJobsSQSURL() string {
	return s.ArchiveRelocateJobsSQSURLValue
}
