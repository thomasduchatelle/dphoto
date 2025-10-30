import * as cdk from 'aws-cdk-lib';
import {Match, Template} from 'aws-cdk-lib/assertions';
import {InfrastructureStack} from './infrastructure-stack';
import {environments} from '../config/environments';

describe('DPhotoInfrastructureStack', () => {
    describe("prod-like", () => {

        let app: cdk.App;
        let stack: InfrastructureStack;
        let template: Template;

        beforeEach(() => {
            app = new cdk.App();
            stack = new InfrastructureStack(app, 'TestStack', {
                environmentName: 'test',
                config: environments.test,
                env: {
                    account: '123456789012',
                    region: 'eu-west-1'
                }
            });
            template = Template.fromStack(stack);
        });

        test('exports all required SSM parameters for application stack', () => {
            // Verify critical SSM parameters exist
            const expectedParameters = [
                '/dphoto/test/dynamodb/catalog/tableName',
                '/dphoto/test/iam/policies/archive_relocate_send/arn',
                '/dphoto/test/iam/policies/archive_sns_publish/arn',
                '/dphoto/test/iam/policies/archive_sqs_send/arn',
                '/dphoto/test/iam/policies/cacheRWArn',
                '/dphoto/test/iam/policies/indexRWArn',
                '/dphoto/test/iam/policies/storageROArn',
                '/dphoto/test/iam/policies/storageRWArn',
                '/dphoto/test/s3/cache/bucketName',
                '/dphoto/test/s3/storage/bucketName',
                '/dphoto/test/sns/archive/arn',
                '/dphoto/test/sqs/archive/arn',
                '/dphoto/test/sqs/archive/url',
                '/dphoto/test/sqs/archive_relocate/arn',
                '/dphoto/test/sqs/archive_relocate/url'
            ];

            expectedParameters.forEach(paramName => {
                template.hasResourceProperties('AWS::SSM::Parameter', {
                    Name: paramName,
                    Type: 'String'
                });
            });
        });

        test('main S3 bucket will not be deleted with the Stack: RETAIN', () => {
            template.hasResource('AWS::S3::Bucket', {
                DeletionPolicy: 'Retain',
                Properties: {
                    BucketName: 'dphoto-test-storage'
                }
            });

            const bucketsWithAutoDelete = Object.entries(template.findResources('Custom::S3AutoDeleteObjects', {})).map(([id, resource]) => {
                return resource.Properties.BucketName;
            })
            expect(bucketsWithAutoDelete).not.toContain(bucketReference(template, 'dphoto-test-storage'))
        });

        test('DynamoDB table has point-in-time recovery enabled', () => {
            template.hasResource('AWS::DynamoDB::GlobalTable', {
                DeletionPolicy: 'Retain',
                Properties: {
                    TableName: 'dphoto-test-index',
                    BillingMode: "PAY_PER_REQUEST",
                    Replicas: [
                        {
                            DeletionProtectionEnabled: true,
                            PointInTimeRecoverySpecification: {
                                PointInTimeRecoveryEnabled: true
                            },
                        }
                    ],
                }
            });
        });

        test('creates FIFO SQS queue for archive jobs', () => {
            template.hasResourceProperties('AWS::SQS::Queue', {
                QueueName: 'dphoto-test-async-archive-caching-jobs.fifo',
                FifoQueue: true,
                ContentBasedDeduplication: true
            });
        });

        test('creates SNS topic for archive notifications', () => {
            template.hasResourceProperties('AWS::SNS::Topic', {
                TopicName: 'dphoto-test-archive-jobs'
            });
        });

        test('creates IAM user with correct policies attached', () => {
            template.hasResourceProperties('AWS::IAM::User', {
                UserName: 'dphoto-test-cli',
                Path: '/dphoto/'
            });

            // Verify managed policies are created
            template.hasResourceProperties('AWS::IAM::ManagedPolicy', {
                ManagedPolicyName: 'dphoto-test-storage-rw',
                Path: '/dphoto/'
            });

            template.hasResourceProperties('AWS::IAM::ManagedPolicy', {
                ManagedPolicyName: 'dphoto-test-index-rw',
                Path: '/dphoto/'
            });
        });

        test('critical resources have pinned logical IDs to prevent recreation', () => {
            const resourceHasPinnedId = (expectedLogicalId: string, resourceType: string, matcher: any) => {
                expect(Object.keys(template.findResources(resourceType, matcher))).toEqual([expectedLogicalId]);
            };

            resourceHasPinnedId('MediaStorageStorageBucket1696C64E', 'AWS::S3::Bucket', {
                Properties: {
                    BucketName: 'dphoto-test-storage'
                }
            });

            resourceHasPinnedId('MediaStorageCacheBucketFB633A3D', 'AWS::S3::Bucket', {
                Properties: {
                    BucketName: 'dphoto-test-cache'
                }
            });

            resourceHasPinnedId('CatalogStoreCatalogTable874E34D1', 'AWS::DynamoDB::GlobalTable', {
                Properties: {
                    TableName: 'dphoto-test-index'
                }
            });

            resourceHasPinnedId('ArchiveMessagingArchiveTopic4F67B9F5', 'AWS::SNS::Topic', {
                Properties: {
                    TopicName: 'dphoto-test-archive-jobs'
                }
            });

            resourceHasPinnedId('ArchiveMessagingArchiveQueue02DEF245', 'AWS::SQS::Queue', {
                Properties: {
                    QueueName: 'dphoto-test-async-archive-caching-jobs.fifo'
                }
            });

            resourceHasPinnedId('ArchiveMessagingArchiveRelocateQueue32B7B729', 'AWS::SQS::Queue', {
                Properties: {
                    QueueName: 'dphoto-test-archive-relocate'
                }
            });

            resourceHasPinnedId('MediaStorageStorageKeyA14447B2', 'AWS::KMS::Key', {});

            resourceHasPinnedId('CliUserA7F35037', 'AWS::IAM::User', {
                Properties: {
                    UserName: 'dphoto-test-cli'
                }
            });
        });

        test('Cognito User Pool has advanced security mode enforced', () => {
            template.hasResourceProperties('AWS::Cognito::UserPool', {
                UserPoolName: 'dphoto-test-users',
                UserPoolAddOns: {
                    AdvancedSecurityMode: 'ENFORCED'
                }
            });
        });

        test('Cognito User Pool creates required user groups', () => {
            template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
                GroupName: 'admins',
                Description: 'Administrators with full system access',
                Precedence: 1
            });

            template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
                GroupName: 'owners',
                Description: 'Content owners with full access to their media',
                Precedence: 2
            });

            template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
                GroupName: 'visitors',
                Description: 'Visitors with limited access to shared albums',
                Precedence: 3
            });
        });

        test('Cognito User Pool configures Google as identity provider', () => {
            template.hasResourceProperties('AWS::Cognito::UserPoolIdentityProvider', {
                ProviderName: 'Google',
                ProviderType: 'Google',
                AttributeMapping: {
                    email: 'email',
                    given_name: 'given_name',
                    family_name: 'family_name',
                    picture: 'picture'
                }
            });
        });

        test('managed policy resources have pinned logical IDs to prevent recreation', () => {
            const resourceHasPinnedId = (expectedLogicalId: string, resourceType: string, matcher: any) => {
                expect(Object.keys(template.findResources(resourceType, matcher))).toEqual([expectedLogicalId]);
            };

            resourceHasPinnedId('MediaStorageStorageRwPolicyC4C10BB9', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-storage-rw'
                }
            });

            resourceHasPinnedId('MediaStorageStorageRoPolicyAE409884', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-storage-ro'
                }
            });

            resourceHasPinnedId('MediaStorageCacheRwPolicyBBDEDD20', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-cache-rw'
                }
            });

            resourceHasPinnedId('CatalogDbIndexRwPolicy18F429CA', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-index-rw'
                }
            });

            resourceHasPinnedId('ArchiveMessagingArchiveSnsPublishPolicy3A80ABCB', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-archive-sns-publish'
                }
            });

            resourceHasPinnedId('ArchiveMessagingArchiveSqsSendPolicy07AA8FDC', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-archive-sqs-send'
                }
            });

            resourceHasPinnedId('ArchiveMessagingArchiveRelocatePolicyFFC7CD89', 'AWS::IAM::ManagedPolicy', {
                Properties: {
                    ManagedPolicyName: 'dphoto-test-archive-relocate-sqs-send'
                }
            });

            resourceHasPinnedId('CliUserAccessKey202507148E1156', 'AWS::IAM::AccessKey', {});
        });

        test('all resources have correct tags', () => {
            const expectedTagsArray = [
                {Key: 'Application', Value: 'dphoto'},
                {Key: 'CreatedBy', Value: 'cdk'},
                {Key: 'Environment', Value: 'test'}
            ];

            // Test resources that use array format for tags
            const arrayTaggedResourceTypes = [
                'AWS::S3::Bucket',
                'AWS::DynamoDB::Table',
                'AWS::SNS::Topic',
                'AWS::SQS::Queue',
                'AWS::IAM::User'
            ];

            arrayTaggedResourceTypes.forEach(resourceType => {
                const resources = template.findResources(resourceType);
                if (Object.keys(resources).length > 0) {
                    try {
                        template.allResourcesProperties(resourceType, {
                            Tags: Match.arrayWith(expectedTagsArray)
                        });
                    } catch (e) {
                        Object.keys(resources).forEach((key, index) => {
                            console.log(`Testing tags for ${key}: ${JSON.stringify(resources[key], null, 2)}`);
                        })
                        throw e
                    }
                }
            });

            // Test SSM parameters that use object format for tags
            template.allResourcesProperties('AWS::SSM::Parameter', {
                Tags: expectedTagsArray.reduce((acc, tag) => {
                    acc[tag.Key] = tag.Value;
                    return acc;
                }, {} as Record<string, string>)
            });
        });

        test('S3 configuration enables KMS encryption for test environment', () => {
            template.resourceCountIs('AWS::KMS::Key', 1);

            template.hasResourceProperties('AWS::S3::Bucket', {
                BucketName: 'dphoto-test-storage',
                BucketEncryption: {
                    ServerSideEncryptionConfiguration: [{
                        ServerSideEncryptionByDefault: {
                            SSEAlgorithm: 'aws:kms'
                        }
                    }]
                }
            });
        });
    })


    describe("non-prod", () => {
        let app: cdk.App;
        let stack: InfrastructureStack;
        let template: Template;

        beforeEach(() => {
            app = new cdk.App();
            stack = new InfrastructureStack(app, 'TestStack', {
                environmentName: 'test',
                config: {
                    ...environments.test,
                    production: false,
                },
                env: {
                    account: '123456789012',
                    region: 'eu-west-1'
                }
            });
            template = Template.fromStack(stack);
        });

        test('S3 bucket for original media is EMPTIED and DELETED with the Stack', () => {
            template.hasResource('AWS::S3::Bucket', {
                DeletionPolicy: 'Delete',
                Properties: {
                    BucketName: 'dphoto-test-storage',
                }
            });
            template.hasResourceProperties('Custom::S3AutoDeleteObjects', {
                BucketName: bucketReference(template, 'dphoto-test-storage'),
            });
        });

        test('S3 bucket for cache is EMPTIED and DELETED with the Stack', () => {
            template.hasResource('AWS::S3::Bucket', {
                DeletionPolicy: 'Delete',
                Properties: {
                    BucketName: 'dphoto-test-cache',
                }
            });
            template.hasResourceProperties('Custom::S3AutoDeleteObjects', {
                BucketName: bucketReference(template, 'dphoto-test-cache'),
            });
        });

        test('DynamoDB table has point-in-time recovery disabled and will be DELETED with the Stack', () => {
            template.hasResource('AWS::DynamoDB::GlobalTable', {
                DeletionPolicy: 'Delete',
                Properties: {
                    TableName: 'dphoto-test-index',
                    BillingMode: "PAY_PER_REQUEST",
                    Replicas: [
                        {
                            DeletionProtectionEnabled: false,
                            PointInTimeRecoverySpecification: {
                                PointInTimeRecoveryEnabled: false
                            },
                        }
                    ],
                }
            });
        });

        test('simple S3 configuration disable KMS encryption for test environment', () => {
            template.resourceCountIs('AWS::KMS::Key', 0);

            template.hasResourceProperties('AWS::S3::Bucket', {
                BucketName: 'dphoto-test-storage',
                BucketEncryption: {
                    ServerSideEncryptionConfiguration: [{
                        ServerSideEncryptionByDefault: {
                            SSEAlgorithm: 'AES256'
                        }
                    }]
                }
            });
        });
    });
});

function bucketReference(template: Template, bucketName: string) {
    let buckets = template.findResources('AWS::S3::Bucket', {
        Properties: {
            BucketName: bucketName,
        }
    });
    return {Ref: Object.keys(buckets)[0]};
}
