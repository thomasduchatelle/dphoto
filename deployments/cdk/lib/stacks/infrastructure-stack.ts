import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ArchiveStoreConstruct} from '../archive/archive-store-construct';
import {ArchivistConstruct} from '../archive/archivist-construct';
import {CatalogStoreConstruct} from '../catalog/catalog-store-construct';
import {CliUserAccessConstruct} from '../access/cli-user-access-construct';
import {ArchiveServerlessIntegrationConstruct} from '../archive/archive-serverless-integration-construct';
import {CatalogServerlessIntegrationConstruct} from '../catalog/catalog-serverless-integration-construct';

export interface DPhotoInfrastructureStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

export class InfrastructureStack extends cdk.Stack {
    public readonly archiveStore: ArchiveStoreConstruct;
    public readonly catalogStore: CatalogStoreConstruct;
    public readonly archivist: ArchivistConstruct;

    constructor(scope: Construct, id: string, props: DPhotoInfrastructureStackProps) {
        super(scope, id, props);

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoInfrastructureStack");

        const {archiveStore, catalogStore, archivist} = this.createFullInfrastructure(props);
        this.archiveStore = archiveStore;
        this.catalogStore = catalogStore;
        this.archivist = archivist
    }

    private createFullInfrastructure(props: DPhotoInfrastructureStackProps): {
        archiveStore: ArchiveStoreConstruct;
        catalogStore: CatalogStoreConstruct;
        archivist: ArchivistConstruct
    } {
        const archiveStore = new ArchiveStoreConstruct(this, 'ArchiveStore', {
            environmentName: props.environmentName,
            simpleS3: !props.config.production
        });

        const catalogStore = new CatalogStoreConstruct(this, 'CatalogStore', {
            environmentName: props.environmentName,
            production: props.config.production,
        });

        const archivist = new ArchivistConstruct(this, 'Archivist', {
            environmentName: props.environmentName
        });

        const accessStore = new CliUserAccessConstruct(this, 'AccessStore', {
            environmentName: props.environmentName,
            cliAccessKeys: props.config.cliAccessKeys || ['2024-04'],
            archiveStore,
            catalogStore,
            archivist,
        });

        new ArchiveServerlessIntegrationConstruct(this, 'ArchiveServerlessIntegration', {
            environmentName: props.environmentName,
            archiveStore: archiveStore,
            archivist: archivist
        });

        new CatalogServerlessIntegrationConstruct(this, 'CatalogServerlessIntegration', {
            environmentName: props.environmentName,
            catalogStore: catalogStore
        });

        new cdk.CfnOutput(this, 'ArchiveBucketName', {
            value: archiveStore.storageBucket.bucketName,
            description: 'Name of the bucket where medias can be uploaded'
        });

        new cdk.CfnOutput(this, 'CacheBucketName', {
            value: archiveStore.cacheBucket.bucketName,
            description: 'Name of the bucket where miniatures are cached'
        });

        new cdk.CfnOutput(this, 'DynamodbName', {
            value: catalogStore.table.tableName,
            description: 'Name of the table that need to be created'
        });

        new cdk.CfnOutput(this, 'Region', {
            value: this.region,
            description: 'AWS Region'
        });

        new cdk.CfnOutput(this, 'SqsAsyncArchiveJobsArn', {
            value: archivist.archiveQueue.queueArn,
            description: 'SQS topic ARN where are dispatched asynchronous jobs'
        });

        new cdk.CfnOutput(this, 'SnsArchiveArn', {
            value: archivist.archiveTopic.topicArn,
            description: 'SNS topic ARN where are dispatched asynchronous jobs'
        });

        new cdk.CfnOutput(this, 'SqsArchiveUrl', {
            value: archivist.archiveQueue.queueUrl,
            description: 'SQS topic URL where are de-duplicated messages'
        });

        Object.entries(accessStore.accessKeys).forEach(([keyDate, accessKey]) => {
            new cdk.CfnOutput(this, `DelegateAccessKeyId${keyDate.replace('-', '')}`, {
                value: accessKey.accessKeyId,
                description: `AWS access Key ID for ${keyDate}`
            });

            new cdk.CfnOutput(this, `DelegateSecretAccessKey${keyDate.replace('-', '')}`, {
                value: accessKey.secretAccessKey.unsafeUnwrap(),
                description: `AWS secret access Key for ${keyDate}`
            });
        });

        return {archiveStore, catalogStore, archivist: archivist}
    }
}
