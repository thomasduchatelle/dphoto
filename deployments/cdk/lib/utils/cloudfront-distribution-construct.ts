import * as cdk from 'aws-cdk-lib';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53_targets from 'aws-cdk-lib/aws-route53-targets';
import {Construct} from 'constructs';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';

export interface CloudFrontDistributionConstructProps {
    environmentName: string;
    nextjsDomainName: string;
    rootDomain: string;
    httpApi: apigatewayv2.HttpApi;
    certificate: ICertificate;
}

export class CloudFrontDistributionConstruct extends Construct {
    public readonly distribution: cloudfront.Distribution;
    public readonly distributionId: string;

    constructor(scope: Construct, id: string, props: CloudFrontDistributionConstructProps) {
        super(scope, id);

        // Create a cache policy that never caches
        const noCachePolicy = new cloudfront.CachePolicy(this, 'ApiNoCachePolicy', {
            cachePolicyName: `dphoto-${props.environmentName}-api-no-cache`,
            comment: 'Policy to never cache API responses',
            defaultTtl: cdk.Duration.seconds(0),
            minTtl: cdk.Duration.seconds(0),
            maxTtl: cdk.Duration.seconds(0),
            cookieBehavior: cloudfront.CacheCookieBehavior.all(),
            headerBehavior: cloudfront.CacheHeaderBehavior.allowList(
                'Authorization',
                'Content-Type',
                'Accept',
                'Origin',
                'Referer',
                'User-Agent'
            ),
            queryStringBehavior: cloudfront.CacheQueryStringBehavior.all(),
            enableAcceptEncodingGzip: true,
            enableAcceptEncodingBrotli: true,
        });

        // Create origin request policy to forward all necessary headers, cookies, and query strings
        const apiOriginRequestPolicy = new cloudfront.OriginRequestPolicy(this, 'ApiOriginRequestPolicy', {
            originRequestPolicyName: `dphoto-${props.environmentName}-api-forward-all`,
            comment: 'Policy to forward all headers, cookies, and query strings to API',
            cookieBehavior: cloudfront.OriginRequestCookieBehavior.all(),
            headerBehavior: cloudfront.OriginRequestHeaderBehavior.all(),
            queryStringBehavior: cloudfront.OriginRequestQueryStringBehavior.all(),
        });

        // Extract the API Gateway domain from the httpApi
        // The API Gateway V2 HTTP API URL format is: https://{api-id}.execute-api.{region}.amazonaws.com
        const apiGatewayDomainName = `${props.httpApi.apiId}.execute-api.${cdk.Stack.of(this).region}.amazonaws.com`;

        // Create CloudFront distribution
        this.distribution = new cloudfront.Distribution(this, 'Distribution', {
            comment: `DPhoto ${props.environmentName} - CloudFront distribution for API and NextJS`,
            domainNames: [props.nextjsDomainName],
            certificate: props.certificate,
            defaultBehavior: {
                origin: new origins.HttpOrigin(apiGatewayDomainName, {
                    protocolPolicy: cloudfront.OriginProtocolPolicy.HTTPS_ONLY,
                }),
                viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
                allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
                cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
                originRequestPolicy: cloudfront.OriginRequestPolicy.ALL_VIEWER,
            },
            additionalBehaviors: {
                '/api/*': {
                    origin: new origins.HttpOrigin(apiGatewayDomainName, {
                        protocolPolicy: cloudfront.OriginProtocolPolicy.HTTPS_ONLY,
                    }),
                    viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
                    allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
                    cachePolicy: noCachePolicy,
                    originRequestPolicy: apiOriginRequestPolicy,
                },
            },
            enableLogging: false,
            priceClass: cloudfront.PriceClass.PRICE_CLASS_100,
        });

        this.distributionId = this.distribution.distributionId;

        // Create Route53 A record for the custom domain
        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.rootDomain
        });

        new route53.ARecord(this, 'DnsRecord', {
            zone: hostedZone,
            recordName: props.nextjsDomainName,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.CloudFrontTarget(this.distribution)
            )
        });

        // Tag the distribution
        cdk.Tags.of(this.distribution).add('Name', `dphoto-${props.environmentName}-cdn`);
        cdk.Tags.of(this.distribution).add('Environment', props.environmentName);

        // Output the distribution ID
        new cdk.CfnOutput(this, 'DistributionId', {
            value: this.distributionId,
            description: 'CloudFront Distribution ID',
            exportName: `dphoto-${props.environmentName}-distribution-id`,
        });

        // Output the distribution domain name
        new cdk.CfnOutput(this, 'DistributionDomainName', {
            value: this.distribution.distributionDomainName,
            description: 'CloudFront Distribution Domain Name',
        });

        // Output the custom domain name
        new cdk.CfnOutput(this, 'CustomDomainName', {
            value: props.nextjsDomainName,
            description: 'CloudFront Custom Domain Name',
        });
    }
}
