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

        // Extract the API Gateway domain from the httpApi
        // The API Gateway V2 HTTP API URL format is: https://{api-id}.execute-api.{region}.amazonaws.com
        // TODO AGENT lock the API Gateway down to only be accessible by the CloudFront distribution ; use a dedicated domain (`api.{domain}`) only if necessary.
        // const apiGatewayDomainName = `${props.httpApi.apiId}.execute-api.${cdk.Stack.of(this).region}.amazonaws.com`;
        const apiGatewayDomainName = 'dev.duchatelle.me';
        // const apiGatewayDomainName = props.httpApi.apiEndpoint!.replace("https://", "");
        // console.log("apiGatewayDomainName =", apiGatewayDomainName);

        // TODO Make sure the distribution is FLAT PRICING -> FREE.
        this.distribution = new cloudfront.Distribution(this, 'Distribution', {
            comment: `DPhoto ${props.environmentName} - CloudFront distribution for API and NextJS`,
            domainNames: [props.domainName],
            certificate: props.certificate,
            defaultBehavior: {
                origin: new origins.HttpOrigin('placeholder.dphoto.duchatelle.me', {
                    originId: "default"
                })
            },
            enableLogging: false,
            priceClass: cloudfront.PriceClass.PRICE_CLASS_100,
        });

        cdk.Tags.of(this.distribution).add("sst:app", "foo")
        cdk.Tags.of(this.distribution).add("sst:ref:kv", "foo")
        cdk.Tags.of(this.distribution).add("sst:ref:kv-namespace", "foo")
        cdk.Tags.of(this.distribution).add("sst:ref:version", "foo")
        cdk.Tags.of(this.distribution).add("sst:stage", "stage")

        this.distribution.addBehavior('/api/*',
            new origins.HttpOrigin(apiGatewayDomainName, {
                protocolPolicy: cloudfront.OriginProtocolPolicy.HTTPS_ONLY,
            }),
            {
                viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.HTTPS_ONLY,
                allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
                cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
                originRequestPolicy: cloudfront.OriginRequestPolicy.ALL_VIEWER,
            }
        )

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
        new route53.AaaaRecord(this, 'QuadDnsRecord', {
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
