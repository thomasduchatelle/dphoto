import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as kms from 'aws-cdk-lib/aws-kms';
import {Construct} from 'constructs';
import {Workload} from '../utils/workload';
import {pinLogicalId} from "../utils/override-logical-ids";

export interface ArchiveStoreConstructProps {
    environmentName: string;
    simpleS3?: boolean;
}

export class ArchiveStoreConstruct extends Construct {
    public readonly storageBucket: s3.Bucket;
    public readonly cacheBucket: s3.Bucket;
    public readonly storageKey?: kms.Key;

    constructor(scope: Construct, id: string, props: ArchiveStoreConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        if (!props.simpleS3) {
            this.storageKey = new kms.Key(this, 'StorageKey', {
                description: `${prefix} encryption key`,
                removalPolicy: cdk.RemovalPolicy.RETAIN,
                pendingWindow: cdk.Duration.days(30)
            });
            this.storageKey.addAlias(`${prefix}-archive`)
            pinLogicalId(this.storageKey, "MediaStorageStorageKeyA14447B2");

            cdk.Tags.of(this.storageKey).add('Name', `${prefix}-encryption-key`);
        }

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
        pinLogicalId(this.storageBucket, "MediaStorageStorageBucket1696C64E");

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
        pinLogicalId(this.cacheBucket, "MediaStorageCacheBucketFB633A3D");
    }

    public grantReadAccessToRawAndCacheMedias(workload: Workload): void {
        this.storageBucket.grantRead(workload.role);
        this.cacheBucket.grantRead(workload.role);
        if (this.storageKey) {
            this.storageKey.grantDecrypt(workload.role)
        }

        workload.function?.addEnvironment("CACHE_BUCKET_NAME", this.cacheBucket.bucketName);
        workload.function?.addEnvironment("STORAGE_BUCKET_NAME", this.storageBucket.bucketName);
    }

    public grantWriteAccessToRawAndCachedMedias(workload: Workload): void {
        this.storageBucket.grantReadWrite(workload.role);
        this.cacheBucket.grantReadWrite(workload.role);
        if (this.storageKey) {
            this.storageKey.grantEncryptDecrypt(workload.role)
        }

        workload.function?.addEnvironment("CACHE_BUCKET_NAME", this.cacheBucket.bucketName);
        workload.function?.addEnvironment("STORAGE_BUCKET_NAME", this.storageBucket.bucketName);
    }
}

