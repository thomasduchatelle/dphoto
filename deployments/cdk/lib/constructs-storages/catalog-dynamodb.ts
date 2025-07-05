import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';

export const CatalogTableIndexes = ["AlbumIndex", "ReverseLocationIndex", "ReverseGrantIndex", "RefreshTokenExpiration"];

export interface CatalogDynamoDbProps {
    environmentName: string;
}

export class CatalogDynamoDbConstruct extends Construct {
    public readonly table: dynamodb.Table;
    public readonly indexRwPolicy: iam.ManagedPolicy;

    constructor(scope: Construct, id: string, props: { environmentName: string; production: boolean }) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;
        const tableName = `${prefix}-index`;

        // DynamoDB Table
        this.table = new dynamodb.Table(this, 'CatalogTable', {
            tableName: tableName,
            partitionKey: {
                name: 'PK',
                type: dynamodb.AttributeType.STRING
            },
            sortKey: {
                name: 'SK',
                type: dynamodb.AttributeType.STRING
            },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: props.production ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY,
            pointInTimeRecovery: props.production,
        });

        // IAM Policy for DynamoDB access
        this.indexRwPolicy = new iam.ManagedPolicy(this, 'IndexRwPolicy', {
            managedPolicyName: `${prefix}-index-rw`,
            path: '/dphoto/',
            statements: [
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: [
                        'dynamodb:List*',
                        'dynamodb:DescribeReservedCapacity*',
                        'dynamodb:DescribeLimits',
                        'dynamodb:DescribeTimeToLive'
                    ],
                    resources: ['*']
                }),
                new iam.PolicyStatement({
                    effect: iam.Effect.ALLOW,
                    actions: [
                        'dynamodb:BatchGet*',
                        'dynamodb:DescribeStream',
                        'dynamodb:DescribeTable',
                        'dynamodb:Get*',
                        'dynamodb:Query',
                        'dynamodb:Scan',
                        'dynamodb:BatchWrite*',
                        'dynamodb:CreateTable',
                        'dynamodb:Delete*',
                        'dynamodb:Update*',
                        'dynamodb:PutItem',
                        'dynamodb:TagResource'
                    ],
                    resources: [
                        this.table.tableArn,
                        `${this.table.tableArn}/*`
                    ]
                })
            ]
        });
    }
}

