import * as ssm from 'aws-cdk-lib/aws-ssm';
import {Construct} from 'constructs';
import {CatalogStoreConstruct} from './catalog-store-construct';
import * as iam from "aws-cdk-lib/aws-iam";
import {pinLogicalId} from "../utils/override-logical-ids";

export interface CatalogServerlessIntegrationConstructProps {
    environmentName: string;
    catalogStore: CatalogStoreConstruct;
}

export class CatalogServerlessIntegrationConstruct extends Construct {
    constructor(scope: Construct, id: string, props: CatalogServerlessIntegrationConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        const indexRwPolicy = new iam.ManagedPolicy(this, 'IndexRwPolicy', {
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
                        props.catalogStore.table.tableArn,
                        `${props.catalogStore.table.tableArn}/*`
                    ]
                })
            ]
        });
        pinLogicalId(indexRwPolicy, "CatalogDbIndexRwPolicy18F429CA");

        new ssm.StringParameter(scope, 'IamPolicyIndexRwArnSSM', {
            parameterName: `/dphoto/${props.environmentName}/iam/policies/indexRWArn`,
            stringValue: indexRwPolicy.managedPolicyArn,
            description: 'ARN of the index read-write policy'
        });

        new ssm.StringParameter(scope, 'CatalogTableNameSSM', {
            parameterName: `/dphoto/${props.environmentName}/dynamodb/catalog/tableName`,
            stringValue: props.catalogStore.table.tableName,
            description: 'Name of the catalog table'
        });
    }
}
