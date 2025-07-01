import * as cdk from 'aws-cdk-lib';
import * as sns from 'aws-cdk-lib/aws-sns';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';

export interface ArchiveMessagingProps {
    environmentName: string;
}

export class ArchiveMessagingConstruct extends Construct {
    public readonly archiveTopic: sns.Topic;
    public readonly archiveQueue: sqs.Queue;
    public readonly archiveRelocateQueue: sqs.Queue;
    public readonly archiveSnsPublishPolicy: iam.ManagedPolicy;
    public readonly archiveSqsSendPolicy: iam.ManagedPolicy;
    public readonly archiveRelocatePolicy: iam.ManagedPolicy;

    constructor(scope: Construct, id: string, props: ArchiveMessagingProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // SNS Topic
        this.archiveTopic = new sns.Topic(this, 'ArchiveTopic', {
            topicName: `${prefix}-archive-jobs`
        });

        // SQS Queue for archive jobs (FIFO)
        this.archiveQueue = new sqs.Queue(this, 'ArchiveQueue', {
            queueName: `${prefix}-async-archive-caching-jobs.fifo`,
            fifo: true,
            contentBasedDeduplication: true,
            visibilityTimeout: cdk.Duration.seconds(900 * 6)
        });

        // SQS Queue for archive relocate
        this.archiveRelocateQueue = new sqs.Queue(this, 'ArchiveRelocateQueue', {
            queueName: `${prefix}-archive-relocate`,
            visibilityTimeout: cdk.Duration.seconds(900 * 6),
            retentionPeriod: cdk.Duration.days(14)
        });

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
    }
}
