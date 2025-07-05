import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as apigatewayv2_integrations from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import {Construct} from 'constructs';

export interface StaticWebsiteEndpointProps {
    environmentName: string;
    domainName: string;
    httpApi: apigatewayv2.HttpApi;
}

export class StaticWebsiteEndpointConstruct extends Construct {
    public readonly uiBucket: s3.Bucket;

    constructor(scope: Construct, id: string, props: StaticWebsiteEndpointProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // S3 Bucket for static UI content
        this.uiBucket = new s3.Bucket(this, 'UiBucket', {
            bucketName: `${prefix}-ui-static-public`,
            publicReadAccess: true,
            blockPublicAccess: new s3.BlockPublicAccess({
                blockPublicAcls: false,
                blockPublicPolicy: false,
                ignorePublicAcls: false,
                restrictPublicBuckets: false
            }),
            websiteIndexDocument: 'index.html',
            websiteErrorDocument: 'index.html',
            websiteRoutingRules: [
                {
                    condition: {
                        httpErrorCodeReturnedEquals: '404'
                    },
                    redirectRule: {
                        hostName: props.domainName,
                        protocol: s3.RedirectProtocol.HTTPS,
                        httpRedirectCode: '307',
                        replaceKey: s3.ReplaceKey.prefixWith('?path=')
                    }
                }
            ],
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            autoDeleteObjects: true
        });

        // Create HTTP integration for static content
        const staticIntegration = new apigatewayv2_integrations.HttpUrlIntegration(
            'StaticIntegration',
            this.uiBucket.bucketWebsiteUrl,
            {
                method: apigatewayv2.HttpMethod.GET
            }
        );

        // Add default route to API Gateway that redirects to S3 website
        new apigatewayv2.HttpRoute(this, 'StaticDefaultRoute', {
            httpApi: props.httpApi,
            routeKey: apigatewayv2.HttpRouteKey.DEFAULT,
            integration: staticIntegration
        });
    }
}
