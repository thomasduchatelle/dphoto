import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3_deployment from 'aws-cdk-lib/aws-s3-deployment';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as apigatewayv2_integrations from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import {Construct} from 'constructs';
import {ApiGatewayConstruct} from './api-gateway-construct';

export interface WakuWebUiConstructProps {
    environmentName: string;
    domainName: string;
    apiGateway: ApiGatewayConstruct;
}

export class WakuWebUiConstruct extends Construct {
    public readonly wakuBucket: s3.Bucket;

    constructor(scope: Construct, id: string, props: WakuWebUiConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        this.wakuBucket = new s3.Bucket(this, 'WakuBucket', {
            bucketName: `${prefix}-waku-static-public`,
            publicReadAccess: true,
            blockPublicAccess: new s3.BlockPublicAccess({
                blockPublicAcls: false,
                blockPublicPolicy: false,
                ignorePublicAcls: false,
                restrictPublicBuckets: false
            }),
            websiteIndexDocument: 'index.html',
            websiteErrorDocument: 'index.html',
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            autoDeleteObjects: true
        });

        const wakuIntegration = new apigatewayv2_integrations.HttpUrlIntegration(
            'WakuIntegration',
            this.wakuBucket.bucketWebsiteUrl,
            {
                method: apigatewayv2.HttpMethod.GET
            }
        );

        // Route for /waku (exact match) - serves index.html
        new apigatewayv2.HttpRoute(this, 'WakuIndexRoute', {
            httpApi: props.apiGateway.httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with('/waku', apigatewayv2.HttpMethod.GET),
            integration: wakuIntegration
        });

        // Route for /waku/* (all subroutes) - serves static assets and SPA routes
        new apigatewayv2.HttpRoute(this, 'WakuProxyRoute', {
            httpApi: props.apiGateway.httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with('/waku/{proxy+}', apigatewayv2.HttpMethod.GET),
            integration: wakuIntegration
        });

        new s3_deployment.BucketDeployment(this, 'WakuBucketDeployment', {
            sources: [s3_deployment.Source.asset('../../web-waku/dist/public')],
            destinationBucket: this.wakuBucket,
            prune: true,
            retainOnDelete: false,
        });
    }
}
