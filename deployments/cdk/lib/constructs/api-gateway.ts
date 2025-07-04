import * as cdk from 'aws-cdk-lib';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as apigatewayv2_integrations from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53_targets from 'aws-cdk-lib/aws-route53-targets';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as ssm from 'aws-cdk-lib/aws-ssm';
import * as logs from 'aws-cdk-lib/aws-logs';
import {Construct} from 'constructs';

export interface ApiGatewayConstructProps {
    environmentName: string;
    domainName: string;
    rootDomain: string;
}

export class ApiGatewayConstruct extends Construct {
    public readonly httpApi: apigatewayv2.HttpApi;
    public readonly domainName: apigatewayv2.DomainName;
    public readonly uiBucket: s3.Bucket;

    constructor(scope: Construct, id: string, props: ApiGatewayConstructProps) {
        super(scope, id);

        // Create S3 bucket for UI static files
        this.uiBucket = new s3.Bucket(this, 'UiBucket', {
            bucketName: `dphoto-${props.environmentName}-ui-static-public`,
            publicReadAccess: true,
            blockPublicAccess: new s3.BlockPublicAccess({
                blockPublicAcls: false,
                blockPublicPolicy: false,
                ignorePublicAcls: false,
                restrictPublicBuckets: false
            }),
            websiteIndexDocument: 'index.html',
            websiteErrorDocument: 'index.html',
            websiteRoutingRules: [{
                condition: {
                    httpErrorCodeReturnedEquals: '404'
                },
                redirect: {
                    hostName: props.domainName,
                    protocol: s3.RedirectProtocol.HTTPS,
                    httpRedirectCode: '307',
                    replaceKey: s3.ReplaceKey.prefixWith('?path=')
                }
            }]
        });

        // Create version lambda
        const versionLambda = new lambda.Function(this, 'VersionLambda', {
            functionName: `dphoto-${props.environmentName}-version`,
            runtime: lambda.Runtime.PROVIDED_AL2,
            architecture: lambda.Architecture.ARM_64,
            handler: 'bootstrap',
            code: lambda.Code.fromAsset('../../bin/version.zip'),
            timeout: cdk.Duration.seconds(30),
            memorySize: 256,
            logRetention: logs.RetentionDays.ONE_WEEK
        });

        // Create not-found lambda
        const notFoundLambda = new lambda.Function(this, 'NotFoundLambda', {
            functionName: `dphoto-${props.environmentName}-not-found`,
            runtime: lambda.Runtime.PROVIDED_AL2,
            architecture: lambda.Architecture.ARM_64,
            handler: 'bootstrap',
            code: lambda.Code.fromAsset('../../bin/not-found.zip'),
            timeout: cdk.Duration.seconds(30),
            memorySize: 256,
            logRetention: logs.RetentionDays.ONE_WEEK
        });

        // Create HTTP API
        this.httpApi = new apigatewayv2.HttpApi(this, 'HttpApi', {
            apiName: `dphoto-${props.environmentName}-api`,
            description: `DPhoto API for ${props.environmentName} environment`,
            corsPreflight: {
                allowOrigins: ['*'],
                allowMethods: [apigatewayv2.CorsHttpMethod.ANY],
                allowHeaders: ['*']
            }
        });

        // Add version route
        this.httpApi.addRoutes({
            path: '/api/v1/version',
            methods: [apigatewayv2.HttpMethod.GET],
            integration: new apigatewayv2_integrations.HttpLambdaIntegration('VersionIntegration', versionLambda)
        });

        // Add not-found route for API paths
        this.httpApi.addRoutes({
            path: '/api/{path+}',
            methods: [apigatewayv2.HttpMethod.ANY],
            integration: new apigatewayv2_integrations.HttpLambdaIntegration('NotFoundIntegration', notFoundLambda)
        });

        // Add default route to redirect to S3 website
        const defaultIntegration = new apigatewayv2.HttpIntegration(this, 'DefaultIntegration', {
            httpApi: this.httpApi,
            integrationType: apigatewayv2.HttpIntegrationType.HTTP_PROXY,
            integrationUri: this.uiBucket.bucketWebsiteUrl,
            method: apigatewayv2.HttpMethod.GET,
            payloadFormatVersion: apigatewayv2.PayloadFormatVersion.VERSION_1_0
        });

        new apigatewayv2.HttpRoute(this, 'DefaultRoute', {
            httpApi: this.httpApi,
            routeKey: '$default',
            integration: defaultIntegration
        });

        // Get certificate ARN from SSM
        const certificateArn = ssm.StringParameter.valueForStringParameter(
            this,
            `/dphoto/${props.environmentName}/acm/domainCertARN`
        );

        // Create custom domain
        this.domainName = new apigatewayv2.DomainName(this, 'DomainName', {
            domainName: props.domainName,
            certificate: cdk.aws_certificatemanager.Certificate.fromCertificateArn(
                this,
                'Certificate',
                certificateArn
            )
        });

        // Create API mapping
        new apigatewayv2.ApiMapping(this, 'ApiMapping', {
            api: this.httpApi,
            domainName: this.domainName,
            stage: this.httpApi.defaultStage
        });

        // Get hosted zone and create DNS record
        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.rootDomain
        });

        new route53.ARecord(this, 'DnsRecord', {
            zone: hostedZone,
            recordName: props.domainName,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.ApiGatewayv2DomainProperties(
                    this.domainName.regionalDomainName,
                    this.domainName.regionalHostedZoneId
                )
            )
        });

        // Output bucket name for S3 sync
        new cdk.CfnOutput(this, 'ViewerUiBucketName', {
            value: this.uiBucket.bucketName,
            description: 'Bucket name where static resources of DPhoto are stored'
        });
    }
}
