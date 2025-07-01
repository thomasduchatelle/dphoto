import * as ssm from 'aws-cdk-lib/aws-ssm';
import {Construct} from 'constructs';
import {MediaStorageConstruct} from '../constructs/media-storage';
import {CatalogDynamoDbConstruct} from '../constructs/catalog-dynamodb';
import {ArchiveMessagingConstruct} from '../constructs/archive-messaging';

export function exportSsmParameters(
    scope: Construct,
    environmentName: string,
    mediaStorage: MediaStorageConstruct,
    catalogDb: CatalogDynamoDbConstruct,
    archiveMessaging: ArchiveMessagingConstruct,
): void {
    // IAM Policy ARNs
    new ssm.StringParameter(scope, 'IamPolicyArchiveSnsPublishArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/archive_sns_publish/arn`,
        stringValue: archiveMessaging.archiveSnsPublishPolicy.managedPolicyArn,
        description: 'ARN of the archive SNS publish policy'
    });

    new ssm.StringParameter(scope, 'IamPolicyArchiveSqsSendArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/archive_sqs_send/arn`,
        stringValue: archiveMessaging.archiveSqsSendPolicy.managedPolicyArn,
        description: 'ARN of the archive SQS send policy'
    });

    new ssm.StringParameter(scope, 'IamPolicyStorageRoArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/storageROArn`,
        stringValue: mediaStorage.storageRoPolicy.managedPolicyArn,
        description: 'ARN of the storage read-only policy'
    });

    new ssm.StringParameter(scope, 'IamPolicyStorageRwArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/storageRWArn`,
        stringValue: mediaStorage.storageRwPolicy.managedPolicyArn,
        description: 'ARN of the storage read-write policy'
    });

    new ssm.StringParameter(scope, 'IamPolicyIndexRwArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/indexRWArn`,
        stringValue: catalogDb.indexRwPolicy.managedPolicyArn,
        description: 'ARN of the index read-write policy'
    });

    new ssm.StringParameter(scope, 'IamPolicyCacheRwArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/cacheRWArn`,
        stringValue: mediaStorage.cacheRwPolicy.managedPolicyArn,
        description: 'ARN of the cache read-write policy'
    });

    new ssm.StringParameter(scope, 'IamPolicyArchiveRelocateArnSSM', {
        parameterName: `/dphoto/${environmentName}/iam/policies/archive_relocate_send/arn`,
        stringValue: archiveMessaging.archiveRelocatePolicy.managedPolicyArn,
        description: 'ARN of the archive relocate policy'
    });

    // S3 Bucket Names
    new ssm.StringParameter(scope, 'StorageBucketNameSSM', {
        parameterName: `/dphoto/${environmentName}/s3/storage/bucketName`,
        stringValue: mediaStorage.storageBucket.bucketName,
        description: 'Name of the storage bucket'
    });

    new ssm.StringParameter(scope, 'CacheBucketNameSSM', {
        parameterName: `/dphoto/${environmentName}/s3/cache/bucketName`,
        stringValue: mediaStorage.cacheBucket.bucketName,
        description: 'Name of the cache bucket'
    });

    // DynamoDB Table Name
    new ssm.StringParameter(scope, 'CatalogTableNameSSM', {
        parameterName: `/dphoto/${environmentName}/dynamodb/catalog/tableName`,
        stringValue: catalogDb.table.tableName,
        description: 'Name of the catalog table'
    });

    // SNS Topic ARN
    new ssm.StringParameter(scope, 'SnsArchiveArnSSM', {
        parameterName: `/dphoto/${environmentName}/sns/archive/arn`,
        stringValue: archiveMessaging.archiveTopic.topicArn,
        description: 'ARN of the archive SNS topic'
    });

    // SQS Queue ARNs and URLs
    new ssm.StringParameter(scope, 'SqsArchiveArnSSM', {
        parameterName: `/dphoto/${environmentName}/sqs/archive/arn`,
        stringValue: archiveMessaging.archiveQueue.queueArn,
        description: 'ARN of the archive SQS queue'
    });

    new ssm.StringParameter(scope, 'SqsArchiveUrlSSM', {
        parameterName: `/dphoto/${environmentName}/sqs/archive/url`,
        stringValue: archiveMessaging.archiveQueue.queueUrl,
        description: 'URL of the archive SQS queue'
    });

    new ssm.StringParameter(scope, 'SqsArchiveRelocateArnSSM', {
        parameterName: `/dphoto/${environmentName}/sqs/archive_relocate/arn`,
        stringValue: archiveMessaging.archiveRelocateQueue.queueArn,
        description: 'ARN of the archive relocate SQS queue'
    });

    new ssm.StringParameter(scope, 'SqsArchiveRelocateUrlSSM', {
        parameterName: `/dphoto/${environmentName}/sqs/archive_relocate/url`,
        stringValue: archiveMessaging.archiveRelocateQueue.queueUrl,
        description: 'URL of the archive relocate SQS queue'
    });
}