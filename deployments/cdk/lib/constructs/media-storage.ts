import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as kms from 'aws-cdk-lib/aws-kms';
import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';

export interface MediaStorageProps {
    environmentName: string;
    simpleS3?: boolean;
}

export class MediaStorageConstruct extends Construct {
    public readonly storageBucket: s3.Bucket;
    public readonly cacheBucket: s3.Bucket;
    public readonly storageKey?: kms.Key;
    public readonly storageRwPolicy: iam.ManagedPolicy;
    public readonly storageRoPolicy: iam.ManagedPolicy;
    public readonly cacheRwPolicy: iam.ManagedPolicy;

    constructor(scope: Construct, id: string, props: MediaStorageProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // KMS Key for encryption (only if not simple S3)
        if (!props.simpleS3) {
            this.storageKey = new kms.Key(this, 'StorageKey', {
                description: `${prefix} encryption key`,
                removalPolicy: cdk.RemovalPolicy.RETAIN,
                pendingWindow: cdk.Duration.days(30)
            });

            cdk.Tags.of(this.storageKey).add('Name', `${prefix}-encryption-key`);
        }

        // Storage Bucket
        this.storageBucket = new s3.Bucket(this, 'StorageBucket', {
            bucketName: `${prefix}-storage`,
            blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
            objectOwnership: s3.ObjectOwnership.BUCKET_OWNER_ENFORCED,
            versioned: !props.simpleS3,
            encryption: this.storageKey ? s3.BucketEncryption.KMS : s3.BucketEncryption.S3_MANAGED,
            encryptionKey: this.storageKey,
            removalPolicy: props.simpleS3 ? cdk.RemovalPolicy.DESTROY : cdk.RemovalPolicy.RETAIN,
            autoDeleteObjects: props.simpleS3,
            lifecycleRules: [
                {
                    id: 'deleted-eviction',
                    enabled: true,
                    noncurrentVersionTransitions: [
                        {
                            storageClass: s3.StorageClass.GLACIER,
                            transitionAfter: cdk.Duration.days(0)
                        }
                    ],
                    noncurrentVersionExpiration: cdk.Duration.days(30)
                },
                ...(props.simpleS3 ? [] : [{
                    id: 'current-cost-saving',
                    enabled: true,
                    transitions: [
                        {
                            storageClass: s3.StorageClass.GLACIER_INSTANT_RETRIEVAL,
                            transitionAfter: cdk.Duration.days(7)
                        }
                    ]
                }])
            ]
        });

        // Cache Bucket
        this.cacheBucket = new s3.Bucket(this, 'CacheBucket', {
            bucketName: `${prefix}-cache`,
            blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
            objectOwnership: s3.ObjectOwnership.BUCKET_OWNER_ENFORCED,
            versioned: false,
            encryption: this.storageKey ? s3.BucketEncryption.KMS : s3.BucketEncryption.S3_MANAGED,
            encryptionKey: this.storageKey,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            autoDeleteObjects: true,
            lifecycleRules: [
                {
                    id: '1-month-eviction',
                    enabled: true,
                    prefix: 'w=',
                    expiration: cdk.Duration.days(120)
                }
            ]
        });

        // IAM Policies
        this.storageRwPolicy = new iam.ManagedPolicy(this, 'StorageRwPolicy', {
            managedPolicyName: `${prefix}-storage-rw`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:ListBucket'],
                    resources: [this.storageBucket.bucketArn]
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:*Object'],
                    resources: [`${this.storageBucket.bucketArn}/*`]
                }),
                ...(this.storageKey ? [
                    new iam.PolicyStatement({
                        effect: iam.Effect.ALLOW,
                        actions: ['kms:Decrypt', 'kms:GenerateDataKey'],
                        resources: [this.storageKey.keyArn]
                    })
                ] : [])
            ]
        });

        this.storageRoPolicy = new iam.ManagedPolicy(this, 'StorageRoPolicy', {
            managedPolicyName: `${prefix}-storage-ro`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:ListBucket'],
                    resources: [this.storageBucket.bucketArn]
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:GetObject'],
                    resources: [`${this.storageBucket.bucketArn}/*`]
                }),
                ...(this.storageKey ? [
                    new iam.PolicyStatement({
                        effect: iam.Effect.ALLOW,
                        actions: ['kms:Decrypt', 'kms:GenerateDataKey'],
                        resources: [this.storageKey.keyArn]
                    })
                ] : [])
            ]
        });

        this.cacheRwPolicy = new iam.ManagedPolicy(this, 'CacheRwPolicy', {
            managedPolicyName: `${prefix}-cache-rw`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:ListBucket'],
                    resources: [this.cacheBucket.bucketArn]
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['s3:*Object'],
                    resources: [`${this.cacheBucket.bucketArn}/*`]
                })
            ]
        });
    }
}
