import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import {Construct} from 'constructs';
import {Workload} from '../utils/workload';
import {pinLogicalId} from '../utils/override-logical-ids';

export const CatalogTableIndexes = ["AlbumIndex", "ReverseLocationIndex", "ReverseGrantIndex", "RefreshTokenExpiration"];

export interface CatalogStoreConstructProps {
    environmentName: string;
    production: boolean;
}

export class CatalogStoreConstruct extends Construct {
    public readonly table: dynamodb.TableV2;

    constructor(scope: Construct, id: string, props: CatalogStoreConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;
        const tableName = `${prefix}-index`;

        this.table = new dynamodb.TableV2(this, 'CatalogTable', {
            tableName: tableName,
            partitionKey: {
                name: 'PK',
                type: dynamodb.AttributeType.STRING
            },
            sortKey: {
                name: 'SK',
                type: dynamodb.AttributeType.STRING
            },
            billing: dynamodb.Billing.onDemand(),
            removalPolicy: props.production ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY,
            pointInTimeRecoverySpecification: props.production ? {
                pointInTimeRecoveryEnabled: true
            } : {
                pointInTimeRecoveryEnabled: false
            },
            deletionProtection: props.production,
            globalSecondaryIndexes: [
                {
                    indexName: 'AlbumIndex',
                    partitionKey: {
                        name: 'AlbumIndexPK',
                        type: dynamodb.AttributeType.STRING
                    },
                    sortKey: {
                        name: 'AlbumIndexSK',
                        type: dynamodb.AttributeType.STRING
                    },
                    projectionType: dynamodb.ProjectionType.ALL
                },
                {
                    indexName: 'ReverseLocationIndex',
                    partitionKey: {
                        name: 'LocationKeyPrefix',
                        type: dynamodb.AttributeType.STRING
                    },
                    sortKey: {
                        name: 'LocationId',
                        type: dynamodb.AttributeType.STRING
                    },
                    projectionType: dynamodb.ProjectionType.ALL
                },
                {
                    indexName: 'ReverseGrantIndex',
                    partitionKey: {
                        name: 'ResourceOwner',
                        type: dynamodb.AttributeType.STRING
                    },
                    sortKey: {
                        name: 'SK',
                        type: dynamodb.AttributeType.STRING
                    },
                    projectionType: dynamodb.ProjectionType.ALL
                },
                {
                    indexName: 'RefreshTokenExpiration',
                    partitionKey: {
                        name: 'SK',
                        type: dynamodb.AttributeType.STRING
                    },
                    sortKey: {
                        name: 'AbsoluteExpiryTime',
                        type: dynamodb.AttributeType.STRING
                    },
                    projectionType: dynamodb.ProjectionType.INCLUDE,
                    nonKeyAttributes: ['PK']
                }
            ],
        });
        pinLogicalId(this.table, "CatalogStoreCatalogTable874E34D1");
    }

    public grantReadAccess(workload: Workload): void {
        this.table.grantReadData(workload.role);
        workload.function?.addEnvironment("CATALOG_TABLE_NAME", this.table.tableName);
    }

    public grantReadWriteAccess(workload: Workload): void {
        this.table.grantReadWriteData(workload.role);
        workload.function?.addEnvironment("CATALOG_TABLE_NAME", this.table.tableName);
    }
}
