import * as ssm from 'aws-cdk-lib/aws-ssm';
import {Construct} from 'constructs';
import {aws_dynamodb} from "aws-cdk-lib";
import {ITable} from "aws-cdk-lib/aws-dynamodb";
import {GoLangLambdaFunction} from "../utils/golang-lambda-function";
import {Bucket, IBucket} from "aws-cdk-lib/aws-s3";
import {ITopic, Topic} from "aws-cdk-lib/aws-sns";
import {IQueue, Queue} from "aws-cdk-lib/aws-sqs";
import {CatalogTableIndexes} from "./catalog-dynamodb";

export interface InfraConfiguration {
    CATALOG_TABLE_NAME: string
    CACHE_BUCKET_NAME: string
    STORAGE_BUCKET_NAME: string
    SNS_ARCHIVE_ARN: string
    SQS_ARCHIVE_RELOCATE_ARN: string
    SQS_ARCHIVE_RELOCATE_URL: string
}

export interface InfrastructureConfigurationProps {
    environmentName: string;
}

export class StoragesConnectorConstruct extends Construct {
    private readonly configuration: InfraConfiguration;
    private catalogTable: ITable;
    private storageBucket: IBucket;
    private cacheBucket: IBucket;
    private snsArchiveTopic: ITopic;
    private sqsRelocateQueue: IQueue;

    constructor(scope: Construct, id: string, props: InfrastructureConfigurationProps) {
        super(scope, id);

        this.configuration = {
            CATALOG_TABLE_NAME: ssm.StringParameter.valueForStringParameter(
                this,
                `/dphoto/${props.environmentName}/dynamodb/catalog/tableName`
            ),
            CACHE_BUCKET_NAME: ssm.StringParameter.valueForStringParameter(
                this,
                `/dphoto/${props.environmentName}/s3/cache/bucketName`
            ),
            STORAGE_BUCKET_NAME: ssm.StringParameter.valueForStringParameter(
                this,
                `/dphoto/${props.environmentName}/s3/storage/bucketName`
            ),
            SNS_ARCHIVE_ARN: ssm.StringParameter.valueForStringParameter(
                this,
                `/dphoto/${props.environmentName}/sns/archive/arn`
            ),
            SQS_ARCHIVE_RELOCATE_ARN: ssm.StringParameter.valueForStringParameter(
                this,
                `/dphoto/${props.environmentName}/sqs/archive_relocate/arn`
            ),
            SQS_ARCHIVE_RELOCATE_URL: ssm.StringParameter.valueForStringParameter(
                this,
                `/dphoto/${props.environmentName}/sqs/archive_relocate/url`
            )
        };

        this.catalogTable = aws_dynamodb.Table.fromTableAttributes(this, 'CatalogTable', {
            tableName: this.configuration.CATALOG_TABLE_NAME,
            globalIndexes: CatalogTableIndexes,
        })

        this.storageBucket = Bucket.fromBucketName(this, 'StorageBucket', this.configuration.STORAGE_BUCKET_NAME);
        this.cacheBucket = Bucket.fromBucketName(this, 'CacheBucket', this.configuration.CACHE_BUCKET_NAME);

        this.snsArchiveTopic = Topic.fromTopicArn(this, "ArchiveTopic", this.configuration.SNS_ARCHIVE_ARN);
        this.sqsRelocateQueue = Queue.fromQueueArn(this, 'RelocateQueue', this.configuration.SQS_ARCHIVE_RELOCATE_ARN)
    }

    public grantRWToCatalogTable(lambda: GoLangLambdaFunction) {
        this.catalogTable.grantReadWriteData(lambda.role);
        lambda.function.addEnvironment("CATALOG_TABLE_NAME", this.catalogTable.tableName)
    }

    public grantReadToCatalogTable(lambda: GoLangLambdaFunction) {
        this.catalogTable.grantReadData(lambda.role);
        lambda.function.addEnvironment("CATALOG_TABLE_NAME", this.catalogTable.tableName)
    }

    public grantTriggerOptimisationToArchive(lambda: GoLangLambdaFunction) {
        this.snsArchiveTopic.grantPublish(lambda.role);
        this.sqsRelocateQueue.grantSendMessages(lambda.role);

        lambda.function.addEnvironment("SNS_ARCHIVE_ARN", this.configuration.SNS_ARCHIVE_ARN);
        lambda.function.addEnvironment("SQS_ARCHIVE_RELOCATE_URL", this.configuration.SQS_ARCHIVE_RELOCATE_URL);
    }

    public grantReadToStorageAndCache(lambda: GoLangLambdaFunction) {
        this.storageBucket.grantRead(lambda.role);
        this.cacheBucket.grantRead(lambda.role);
        lambda.function.addEnvironment("CACHE_BUCKET_NAME", this.configuration.CACHE_BUCKET_NAME);
        lambda.function.addEnvironment("STORAGE_BUCKET_NAME", this.configuration.STORAGE_BUCKET_NAME);
    }

    public grantRWToArchive(lambda: GoLangLambdaFunction) {
        this.storageBucket.grantReadWrite(lambda.role);
        this.cacheBucket.grantReadWrite(lambda.role);
        lambda.function.addEnvironment("CACHE_BUCKET_NAME", this.configuration.CACHE_BUCKET_NAME);
        lambda.function.addEnvironment("STORAGE_BUCKET_NAME", this.configuration.STORAGE_BUCKET_NAME);
    }
}
