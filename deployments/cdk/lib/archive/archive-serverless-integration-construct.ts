import * as ssm from 'aws-cdk-lib/aws-ssm';
import {Construct} from 'constructs';
import {ArchiveStoreConstruct} from './archive-store-construct';
import {ArchivistConstruct} from './archivist-construct';
import * as iam from "aws-cdk-lib/aws-iam";
import {ManagedPolicy} from "aws-cdk-lib/aws-iam";
import {pinLogicalId} from "../utils/override-logical-ids";

export interface ArchiveServerlessIntegrationConstructProps {
    environmentName: string;
    archiveStore: ArchiveStoreConstruct;
    archivist: ArchivistConstruct;
}

export class ArchiveServerlessIntegrationConstruct extends Construct {
    constructor(scope: Construct, id: string, props: ArchiveServerlessIntegrationConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        const policies = this.createManagedPolicies(prefix, props);
        this.exportsReferencesIntoSSM(scope, props, policies);
    }

    private createManagedPolicies(prefix: string, props: ArchiveServerlessIntegrationConstructProps) {
        const storageRwPolicy = new iam.ManagedPolicy(this, 'StorageRwPolicy', {
            managedPolicyName: `${prefix}-storage-rw`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:ListBucket'],
                    resources: [props.archiveStore.storageBucket.bucketArn]
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:*Object'],
                    resources: [`${props.archiveStore.storageBucket.bucketArn}/*`]
                }),
                ...(props.archiveStore.storageKey ? [
                    new iam.PolicyStatement({
                        effect: iam.Effect.ALLOW,
                        actions: ['kms:Decrypt', 'kms:GenerateDataKey'],
                        resources: [props.archiveStore.storageKey.keyArn]
                    })
                ] : [])
            ]
        });
        pinLogicalId(storageRwPolicy, "MediaStorageStorageRwPolicyC4C10BB9");

        const storageRoPolicy = new iam.ManagedPolicy(this, 'StorageRoPolicy', {
            managedPolicyName: `${prefix}-storage-ro`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:ListBucket'],
                    resources: [props.archiveStore.storageBucket.bucketArn]
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:GetObject'],
                    resources: [`${props.archiveStore.storageBucket.bucketArn}/*`]
                }),
                ...(props.archiveStore.storageKey ? [
                    new iam.PolicyStatement({
                        effect: iam.Effect.ALLOW,
                        actions: ['kms:Decrypt', 'kms:GenerateDataKey'],
                        resources: [props.archiveStore.storageKey.keyArn]
                    })
                ] : [])
            ]
        });
        pinLogicalId(storageRoPolicy, "MediaStorageStorageRoPolicyAE409884");

        const cacheRwPolicy = new iam.ManagedPolicy(this, 'CacheRwPolicy', {
            managedPolicyName: `${prefix}-cache-rw`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:ListBucket'],
                    resources: [props.archiveStore.cacheBucket.bucketArn]
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:*Object'],
                    resources: [`${props.archiveStore.cacheBucket.bucketArn}/*`]
                })
            ]
        });
        pinLogicalId(cacheRwPolicy, "MediaStorageCacheRwPolicyBBDEDD20");

        return {storageRwPolicy, storageRoPolicy, cacheRwPolicy};
    }

    private exportsReferencesIntoSSM(scope: Construct, props: ArchiveServerlessIntegrationConstructProps, policies: {
        storageRwPolicy: ManagedPolicy;
        storageRoPolicy: ManagedPolicy;
        cacheRwPolicy: ManagedPolicy
    }) {
        new ssm.StringParameter(scope, 'IamPolicyArchiveSnsPublishArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/archive_sns_publish/arn`,
            stringValue: props.archivist.archiveSnsPublishPolicy.managedPolicyArn,
            description: 'ARN of the archive SNS publish policy'
        });

        new ssm.StringParameter(scope, 'IamPolicyArchiveSqsSendArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/archive_sqs_send/arn`,
            stringValue: props.archivist.archiveSqsSendPolicy.managedPolicyArn,
            description: 'ARN of the archive SQS send policy'
        });

        new ssm.StringParameter(scope, 'IamPolicyArchiveRelocateArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/archive_relocate_send/arn`,
            stringValue: props.archivist.archiveRelocatePolicy.managedPolicyArn,
            description: 'ARN of the archive relocate policy'
        });

        // IAM Policy ARNs for archive storage
        new ssm.StringParameter(scope, 'IamPolicyStorageRoArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/storageROArn`,
            stringValue: policies.storageRoPolicy.managedPolicyArn,
            description: 'ARN of the storage read-only policy'
        });

        new ssm.StringParameter(scope, 'IamPolicyStorageRwArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/storageRWArn`,
            stringValue: policies.storageRwPolicy.managedPolicyArn,
            description: 'ARN of the storage read-write policy'
        });

        new ssm.StringParameter(scope, 'IamPolicyCacheRwArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/cacheRWArn`,
            stringValue: policies.cacheRwPolicy.managedPolicyArn,
            description: 'ARN of the cache read-write policy'
        });

        // S3 Bucket Names
        new ssm.StringParameter(scope, 'StorageBucketNameSSM', {
            parameterName: `/dphoto/${props.environmentName}/s3/storage/bucketName`,
            stringValue: props.archiveStore.storageBucket.bucketName,
            description: 'Name of the storage bucket'
        });

        new ssm.StringParameter(scope, 'CacheBucketNameSSM', {
            parameterName: `/dphoto/${props.environmentName}/s3/cache/bucketName`,
            stringValue: props.archiveStore.cacheBucket.bucketName,
            description: 'Name of the cache bucket'
        });

        // SNS Topic ARN
        new ssm.StringParameter(scope, 'SnsArchiveArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/sns/archive/arn`,
            stringValue: props.archivist.archiveTopic.topicArn,
            description: 'ARN of the archive SNS topic'
        });

        // SQS Queue ARNs and URLs
        new ssm.StringParameter(scope, 'SqsArchiveArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/sqs/archive/arn`,
            stringValue: props.archivist.archiveQueue.queueArn,
            description: 'ARN of the archive SQS queue'
        });

        new ssm.StringParameter(scope, 'SqsArchiveUrlSSM', {
            parameterName: `/dphoto/${props.environmentName}/sqs/archive/url`,
            stringValue: props.archivist.archiveQueue.queueUrl,
            description: 'URL of the archive SQS queue'
        });

        new ssm.StringParameter(scope, 'SqsArchiveRelocateArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/sqs/archive_relocate/arn`,
            stringValue: props.archivist.archiveRelocateQueue.queueArn,
            description: 'ARN of the archive relocate SQS queue'
        });

        new ssm.StringParameter(scope, 'SqsArchiveRelocateUrlSSM', {
            parameterName: `/dphoto/${props.environmentName}/sqs/archive_relocate/url`,
            stringValue: props.archivist.archiveRelocateQueue.queueUrl,
            description: 'URL of the archive relocate SQS queue'
        });
    }
}
