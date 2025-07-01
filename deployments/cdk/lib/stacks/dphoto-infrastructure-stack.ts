import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {MediaStorageConstruct} from '../constructs/media-storage';
import {CatalogDynamoDbConstruct} from '../constructs/catalog-dynamodb';
import {ArchiveMessagingConstruct} from '../constructs/archive-messaging';
import {CliUserConstruct} from '../constructs/cli-user';
import {exportSsmParameters} from '../utils/ssm-parameters';

export interface DPhotoInfrastructureStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

export class DPhotoInfrastructureStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props: DPhotoInfrastructureStackProps) {
        super(scope, id, props);

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);

        const importOnly = this.node.tryGetContext('importOnly') === 'true';

        if (importOnly) {
            // skipped - resource will be imported first
        } else {
            // Full stack creation
            this.createFullInfrastructure(props);
        }
    }

    private createFullInfrastructure(props: DPhotoInfrastructureStackProps): void {
        // Create media storage (S3 buckets and policies)
        const mediaStorage = new MediaStorageConstruct(this, 'MediaStorage', {
            environmentName: props.environmentName,
            simpleS3: !props.config.production
        });

        // Create catalog database (DynamoDB table and policies)
        const catalogDb = new CatalogDynamoDbConstruct(this, 'CatalogDb', {
            environmentName: props.environmentName,
            production: props.config.production,
        });

        // Create archive messaging (SNS/SQS and policies)
        const archiveMessaging = new ArchiveMessagingConstruct(this, 'ArchiveMessaging', {
            environmentName: props.environmentName
        });

        // Create CLI user with access keys
        const cliUser = new CliUserConstruct(this, 'CliUser', {
            environmentName: props.environmentName,
            cliAccessKeys: props.config.cliAccessKeys || ['2024-04'],
            keybaseUser: props.config.keybaseUser || 'keybase:thomasduchatelle',
            storageRwPolicyArn: mediaStorage.storageRwPolicy.managedPolicyArn,
            cacheRwPolicyArn: mediaStorage.cacheRwPolicy.managedPolicyArn,
            indexRwPolicyArn: catalogDb.indexRwPolicy.managedPolicyArn,
            archiveSnsPublishPolicyArn: archiveMessaging.archiveSnsPublishPolicy.managedPolicyArn,
            archiveSqsSendPolicyArn: archiveMessaging.archiveSqsSendPolicy.managedPolicyArn,
            archiveRelocatePolicyArn: archiveMessaging.archiveRelocatePolicy.managedPolicyArn
        });

        // Export SSM parameters for application stack
        exportSsmParameters(
            this,
            props.environmentName,
            mediaStorage,
            catalogDb,
            archiveMessaging
        );

        // Outputs (matching Terraform outputs)
        new cdk.CfnOutput(this, 'ArchiveBucketName', {
            value: mediaStorage.storageBucket.bucketName,
            description: 'Name of the bucket where medias can be uploaded'
        });

        new cdk.CfnOutput(this, 'CacheBucketName', {
            value: mediaStorage.cacheBucket.bucketName,
            description: 'Name of the bucket where miniatures are cached'
        });

        new cdk.CfnOutput(this, 'DynamodbName', {
            value: catalogDb.table.tableName,
            description: 'Name of the table that need to be created'
        });

        new cdk.CfnOutput(this, 'Region', {
            value: this.region,
            description: 'AWS Region'
        });

        new cdk.CfnOutput(this, 'SqsAsyncArchiveJobsArn', {
            value: archiveMessaging.archiveQueue.queueArn,
            description: 'SQS topic ARN where are dispatched asynchronous jobs'
        });

        new cdk.CfnOutput(this, 'SnsArchiveArn', {
            value: archiveMessaging.archiveTopic.topicArn,
            description: 'SNS topic ARN where are dispatched asynchronous jobs'
        });

        new cdk.CfnOutput(this, 'SqsArchiveUrl', {
            value: archiveMessaging.archiveQueue.queueUrl,
            description: 'SQS topic URL where are de-duplicated messages'
        });

        // Access key outputs (encrypted)
        Object.entries(cliUser.accessKeys).forEach(([keyDate, accessKey]) => {
            new cdk.CfnOutput(this, `DelegateAccessKeyId${keyDate.replace('-', '')}`, {
                value: accessKey.accessKeyId,
                description: `AWS access Key ID for ${keyDate}`
            });

            new cdk.CfnOutput(this, `DelegateSecretAccessKey${keyDate.replace('-', '')}`, {
                value: accessKey.secretAccessKey.unsafeUnwrap(),
                description: `AWS secret access Key for ${keyDate}`
            });
        });
    }
}