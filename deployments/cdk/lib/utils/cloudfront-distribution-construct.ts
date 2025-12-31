import * as cdk from 'aws-cdk-lib';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53_targets from 'aws-cdk-lib/aws-route53-targets';
import {Construct} from 'constructs';
import {ICertificate} from "aws-cdk-lib/aws-certificatemanager";

export interface CloudFrontDistributionConstructProps {
    environmentName: string;
    rootDomain: string;
    domainName: string;
    httpApi: apigatewayv2.HttpApi;
    certificate: ICertificate;
}

export class CloudFrontDistributionConstruct extends Construct {
    public readonly distribution: cloudfront.Distribution;
    public readonly distributionId: string;
    public readonly url: string;

    constructor(scope: Construct, id: string, props: CloudFrontDistributionConstructProps) {
        super(scope, id);

        const apiOriginRequestPolicy = new cloudfront.OriginRequestPolicy(this, 'ApiOriginRequestPolicy', {
            originRequestPolicyName: `dphoto-${props.environmentName}-api-forward-all`,
            comment: 'Policy to forward all headers, cookies, and query strings to API',
            cookieBehavior: cloudfront.OriginRequestCookieBehavior.all(),
            headerBehavior: cloudfront.OriginRequestHeaderBehavior.all(),
            queryStringBehavior: cloudfront.OriginRequestQueryStringBehavior.all(),
        });

        // Extract the API Gateway domain from the httpApi
        // The API Gateway V2 HTTP API URL format is: https://{api-id}.execute-api.{region}.amazonaws.com
        // TODO AGENT lock the API Gateway down to only be accessible by the CloudFront distribution ; use a dedicated domain (`api.{domain}`) only if necessary.
        const apiGatewayDomainName = `${props.httpApi.apiId}.execute-api.${cdk.Stack.of(this).region}.amazonaws.com`;

        // TODO Make sure the distribution is FLAT PRICING -> FREE.
        this.distribution = new cloudfront.Distribution(this, 'Distribution', {
            comment: `DPhoto ${props.environmentName} - CloudFront distribution for API and NextJS`,
            domainNames: [props.domainName],
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
                    cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
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
            recordName: props.domainName,
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

        this.url = `https://${this.distribution.domainName}`;
    }
}
