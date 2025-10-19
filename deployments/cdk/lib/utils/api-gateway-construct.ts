import * as cdk from 'aws-cdk-lib';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {CfnStage, HttpApi} from 'aws-cdk-lib/aws-apigatewayv2';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53_targets from 'aws-cdk-lib/aws-route53-targets';
import * as logs from 'aws-cdk-lib/aws-logs';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from "./simple-go-endpoint";
import {ICertificate} from "aws-cdk-lib/aws-certificatemanager";
import {LetsEncryptCertificateConstruct} from "./letsencrypt-certificate-construct";

export interface ApiGatewayConstructProps {
    environmentName: string;
    domainName: string;
    rootDomain: string;
    certificateEmail: string;
    accessLogEnabled?: boolean;
}

export class ApiGatewayConstruct extends Construct {
    public readonly httpApi: apigatewayv2.HttpApi;
    public readonly domainName: apigatewayv2.DomainName;
    public readonly certificate: ICertificate;
    private readonly accessLogEnabled: boolean;

    constructor(scope: Construct, id: string, {accessLogEnabled = false, ...props}: ApiGatewayConstructProps) {
        super(scope, id);
        this.accessLogEnabled = accessLogEnabled;

        const letsEncryptCertificate = new LetsEncryptCertificateConstruct(this, 'LetsEncryptCertificate', {
            environmentName: props.environmentName,
            domainName: props.domainName,
            certificateEmail: props.certificateEmail,
        });
        this.certificate = letsEncryptCertificate.certificate;

        const {httpApi, domainName} = this.createAPIGateway(props, this.certificate)
        this.httpApi = httpApi;
        this.domainName = domainName;

        this.addsDefaultRoutes(httpApi, props.environmentName);
    }



    private createAPIGateway(props: ApiGatewayConstructProps, certificate: ICertificate) {

        const httpApi = new apigatewayv2.HttpApi(this, 'HttpApi', {
            apiName: `dphoto-${props.environmentName}-api`,
            description: `DPhoto API for ${props.environmentName} environment`,
            disableExecuteApiEndpoint: true,
            corsPreflight: {
                allowOrigins: ['*'],
                allowMethods: [apigatewayv2.CorsHttpMethod.ANY],
                allowHeaders: ['*']
            }
        });

        if (this.accessLogEnabled) {
            const apiAccessLogGroup = new logs.LogGroup(this, 'ApiAccessLogs', {
                logGroupName: `/dphoto/${props.environmentName}/api-gateway-access-logs`,
                retention: logs.RetentionDays.ONE_WEEK,
                removalPolicy: cdk.RemovalPolicy.DESTROY,
            });
            // Enable access logging on the default stage
            if (httpApi.defaultStage?.node?.defaultChild) {
                const child = httpApi.defaultStage.node.defaultChild as CfnStage
                child.accessLogSettings = {
                    destinationArn: apiAccessLogGroup.logGroupArn,
                    format: JSON.stringify({
                        requestId: "$context.requestId",
                        ip: "$context.identity.sourceIp",
                        requestTime: "$context.requestTime",
                        httpMethod: "$context.httpMethod",
                        routeKey: "$context.routeKey",
                        status: "$context.status",
                        protocol: "$context.protocol",
                        responseLength: "$context.responseLength",
                        authorizerError: "$context.authorizer.error",
                        authorizerUserId: "$context.authorizer.claims.userId",
                        errorMessage: "$context.error.message",
                        integrationErrorMessage: "$context.integration.error",
                        // authorizer: "$context.authorizer"
                    })
                };
            }
        }
        const domainName = new apigatewayv2.DomainName(this, 'DomainName', {
            domainName: props.domainName,
            certificate: certificate
        });

        new apigatewayv2.ApiMapping(this, 'ApiMapping', {
            api: httpApi,
            domainName: domainName,
            stage: httpApi.defaultStage
        });

        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.rootDomain
        });

        new route53.ARecord(this, 'DnsRecord', {
            zone: hostedZone,
            recordName: props.domainName,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.ApiGatewayv2DomainProperties(
                    domainName.regionalDomainName,
                    domainName.regionalHostedZoneId
                )
            )
        });

        return {httpApi, domainName};
    }

    private addsDefaultRoutes(httpApi: HttpApi, environmentName: string) {
        createSingleRouteEndpoint(this, 'NotFound', {
            functionName: 'not-found',
            path: '/api/{path+}',
            method: apigatewayv2.HttpMethod.ANY,
            httpApi,
            environmentName,
        });
    }
}