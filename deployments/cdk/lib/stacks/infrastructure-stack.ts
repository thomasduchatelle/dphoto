import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {MediaStorageConstruct} from '../constructs-storages/media-storage';
import {CatalogDynamoDbConstruct} from '../constructs-storages/catalog-dynamodb';
import {ArchiveMessagingConstruct} from '../constructs-storages/archive-messaging';
import {CliUserConstruct} from '../constructs-cli/cli-user';
import {exportSsmParameters} from '../utils/ssm-parameters';

export interface DPhotoInfrastructureStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

export class InfrastructureStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props: DPhotoInfrastructureStackProps) {
        super(scope, id, props);

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoInfrastructureStack");

        this.createFullInfrastructure(props);
    }

    private createFullInfrastructure(props: DPhotoInfrastructureStackProps): void {
        const mediaStorage = new MediaStorageConstruct(this, 'MediaStorage', {
            environmentName: props.environmentName,
            simpleS3: !props.config.production
        });

        const catalogDb = new CatalogDynamoDbConstruct(this, 'CatalogDb', {
            environmentName: props.environmentName,
            production: props.config.production,
        });

        const archiveMessaging = new ArchiveMessagingConstruct(this, 'ArchiveMessaging', {
            environmentName: props.environmentName
        });

        const cliUser = new CliUserConstruct(this, 'CliUser', {
            environmentName: props.environmentName,
            cliAccessKeys: props.config.cliAccessKeys || ['2024-04'],
            storageRwPolicyArn: mediaStorage.storageRwPolicy.managedPolicyArn,
            cacheRwPolicyArn: mediaStorage.cacheRwPolicy.managedPolicyArn,
            indexRwPolicyArn: catalogDb.indexRwPolicy.managedPolicyArn,
            archiveSnsPublishPolicyArn: archiveMessaging.archiveSnsPublishPolicy.managedPolicyArn,
            archiveSqsSendPolicyArn: archiveMessaging.archiveSqsSendPolicy.managedPolicyArn,
            archiveRelocatePolicyArn: archiveMessaging.archiveRelocatePolicy.managedPolicyArn
        });

        exportSsmParameters(
            this,
            props.environmentName,
            mediaStorage,
            catalogDb,
            archiveMessaging
        );

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