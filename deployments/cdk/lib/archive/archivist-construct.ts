import * as cdk from 'aws-cdk-lib';
import * as sns from 'aws-cdk-lib/aws-sns';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';
import {Workload} from '../utils/workload';
import {pinLogicalId} from '../utils/override-logical-ids';

export interface ArchiveMessagingConstructProps {
    environmentName: string;
}

export class ArchivistConstruct extends Construct {
    public readonly archiveTopic: sns.Topic;
    public readonly archiveQueue: sqs.Queue;
    public readonly archiveRelocateQueue: sqs.Queue;
    public readonly archiveSnsPublishPolicy: iam.ManagedPolicy;
    public readonly archiveSqsSendPolicy: iam.ManagedPolicy;
    public readonly archiveRelocatePolicy: iam.ManagedPolicy;

    constructor(scope: Construct, id: string, props: ArchiveMessagingConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // SNS Topic
        this.archiveTopic = new sns.Topic(this, 'ArchiveTopic', {
            topicName: `${prefix}-archive-jobs`
        });
        pinLogicalId(this.archiveTopic, "ArchiveMessagingArchiveTopic4F67B9F5");

        // SQS Queue for archive jobs (FIFO)
        this.archiveQueue = new sqs.Queue(this, 'ArchiveQueue', {
            queueName: `${prefix}-async-archive-caching-jobs.fifo`,
            fifo: true,
            contentBasedDeduplication: true,
            visibilityTimeout: cdk.Duration.seconds(900 * 6)
        });
        pinLogicalId(this.archiveQueue, "ArchiveMessagingArchiveQueue02DEF245");

        // SQS Queue for archive relocate
        this.archiveRelocateQueue = new sqs.Queue(this, 'ArchiveRelocateQueue', {
            queueName: `${prefix}-archive-relocate`,
            visibilityTimeout: cdk.Duration.seconds(900 * 6),
            retentionPeriod: cdk.Duration.days(14)
        });
        pinLogicalId(this.archiveRelocateQueue, "ArchiveMessagingArchiveRelocateQueue32B7B729");

        // Allow SNS to publish to SQS
        this.archiveQueue.addToResourcePolicy(
            new iam.PolicyStatement({
                sid: 'Allow SNS to publish messages',
                effect: iam.Effect.ALLOW,
                principals: [new iam.ServicePrincipal('sns.amazonaws.com')],
                actions: ['sqs:SendMessage'],
                resources: [this.archiveQueue.queueArn],
                conditions: {
                    ArnEquals: {
                        'aws:SourceArn': this.archiveTopic.topicArn
                    }
                }
            })
        );

        // IAM Policies
        this.archiveSnsPublishPolicy = new iam.ManagedPolicy(this, 'ArchiveSnsPublishPolicy', {
            managedPolicyName: `${prefix}-archive-sns-publish`,
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['sns:Publish'],
                    resources: [this.archiveTopic.topicArn]
                })
            ]
        });
        pinLogicalId(this.archiveSnsPublishPolicy, "ArchiveMessagingArchiveSnsPublishPolicy3A80ABCB");

        this.archiveSqsSendPolicy = new iam.ManagedPolicy(this, 'ArchiveSqsSendPolicy', {
            managedPolicyName: `${prefix}-archive-sqs-send`,
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['sqs:SendMessage'],
                    resources: [this.archiveQueue.queueArn]
                })
            ]
        });
        pinLogicalId(this.archiveSqsSendPolicy, "ArchiveMessagingArchiveSqsSendPolicy07AA8FDC");

        this.archiveRelocatePolicy = new iam.ManagedPolicy(this, 'ArchiveRelocatePolicy', {
            managedPolicyName: `${prefix}-archive-relocate-sqs-send`,
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: ['sqs:SendMessage'],
                    resources: [this.archiveRelocateQueue.queueArn]
                })
            ]
        });
        pinLogicalId(this.archiveRelocatePolicy, "ArchiveMessagingArchiveRelocatePolicyFFC7CD89");
    }

    public grantAccessToAsyncArchivist(workload: Workload): void {
        this.archiveTopic.grantPublish(workload.role);
        this.archiveQueue.grantSendMessages(workload.role);
        this.archiveRelocateQueue.grantSendMessages(workload.role);
        workload.function?.addEnvironment("SNS_ARCHIVE_ARN", this.archiveTopic.topicArn);
        workload.function?.addEnvironment("SQS_ARCHIVE_URL", this.archiveQueue.queueUrl);
        workload.function?.addEnvironment("SQS_ARCHIVE_RELOCATE_URL", this.archiveRelocateQueue.queueUrl);
    }
}
