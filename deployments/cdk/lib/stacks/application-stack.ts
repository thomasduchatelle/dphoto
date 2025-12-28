import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ApiGatewayConstruct} from '../utils/api-gateway-construct';
import {CatalogEndpointsConstruct} from '../catalog/catalog-endpoints-construct';
import {WakuWebUiConstruct} from '../utils/waku-web-ui-construct';
import {ArchiveEndpointsConstruct} from "../archive/archive-endpoints-construct";
import {VersionEndpointConstruct} from "../utils/version-endpoint-construct";
import {UserEndpointsConstruct} from "../access/user-endpoints-construct";
import {LambdaAuthoriserConstruct} from "../access/lambda-authoriser-construct";
import {CognitoStackExports} from "./cognito-stack";
import {ArchiveAccessManager} from "../archive/archive-access-manager";
import {CatalogAccessManager} from "../catalog/catalog-access-manager";
import {ArchivistAccessManager} from "../archive/archivist-access-manager";
import {AuthenticationEndpointsConstruct} from "../access/authentication-endpoints";
import * as apigatewayv2_integrations from "aws-cdk-lib/aws-apigatewayv2-integrations";
import * as apigatewayv2 from "aws-cdk-lib/aws-apigatewayv2";
import {HttpRoute, HttpRouteKey, MappingValue, ParameterMapping} from "aws-cdk-lib/aws-apigatewayv2";
import {SSTIntegrationConstruct} from "../utils/sst-integration-construct";

export interface DPhotoApplicationStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
    archiveAccessManager: ArchiveAccessManager;
    catalogAccessManager: CatalogAccessManager;
    archivistAccessManager: ArchivistAccessManager;
    oauth2ClientConfig: CognitoStackExports;
}

export class ApplicationStack extends cdk.Stack {
    constructor(scope: Construct, id: string, {
        environmentName,
        config,
        archiveAccessManager,
        catalogAccessManager,
        archivistAccessManager,
        oauth2ClientConfig,
        ...props
    }: DPhotoApplicationStackProps) {
        super(scope, id, props);

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoApplicationStack");

        const apiGateway = new ApiGatewayConstruct(this, 'ApiGateway', {
            environmentName,
            ...config,
        });

        if (config.featureFlags?.useNextJS) {
            new HttpRoute(this, 'NextJsRoute', {
                httpApi: apiGateway.httpApi,
                routeKey: HttpRouteKey.DEFAULT,
                integration: new apigatewayv2_integrations.HttpUrlIntegration(
                    `NextJsIntegration`,
                    `https://${config.nextjsDomainName}`,
                )
            });

        } else {
            new WakuWebUiConstruct(this, 'WakuWebUi', {
                environmentName,
                httpApi: apiGateway.httpApi,
                oauth2ClientConfig: oauth2ClientConfig,
            });

            const nextJsHttpIntegration = new apigatewayv2_integrations.HttpUrlIntegration(
                `NextJsWithBasePathIntegration`,
                `https://${config.nextjsDomainName}`,
                {
                    parameterMapping: new ParameterMapping().overwritePath(MappingValue.requestPath()),
                },
            );
            new apigatewayv2.HttpRoute(this, 'NextJsEagerRoute', {
                httpApi: apiGateway.httpApi,
                routeKey: apigatewayv2.HttpRouteKey.with('/nextjs/{proxy+}', apigatewayv2.HttpMethod.ANY),
                integration: nextJsHttpIntegration
            });
            new apigatewayv2.HttpRoute(this, 'NextJsBaseRoute', {
                httpApi: apiGateway.httpApi,
                routeKey: apigatewayv2.HttpRouteKey.with('/nextjs', apigatewayv2.HttpMethod.ANY),
                integration: nextJsHttpIntegration
            });
        }

        new VersionEndpointConstruct(this, 'VersionEndpoint', {
            environmentName,
            httpApi: apiGateway.httpApi,
        })

        // Create Lambda Authoriser
        const lambdaAuthorizer = new LambdaAuthoriserConstruct(this, 'LambdaAuthoriser', {
            environmentName,
            catalogStore: catalogAccessManager,
            issuerUrl: oauth2ClientConfig.cognitoIssuer,
        });

        new UserEndpointsConstruct(this, 'UserEndpoints', {
            environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore: catalogAccessManager,
            archiveStore: archiveAccessManager,
            authorizer: lambdaAuthorizer.authorizer,
        });

        new CatalogEndpointsConstruct(this, 'CatalogEndpoints', {
            environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore: catalogAccessManager,
            archiveStore: archiveAccessManager,
            archiveMessaging: archivistAccessManager,
            authorizer: lambdaAuthorizer.authorizer,
        });

        new ArchiveEndpointsConstruct(this, 'ArchiveEndpoints', {
            environmentName,
            httpApi: apiGateway.httpApi,
            archiveStore: archiveAccessManager,
            catalogStore: catalogAccessManager,
            archivist: archivistAccessManager,
            authorizer: lambdaAuthorizer.authorizer,
            queryParamAuthorizer: lambdaAuthorizer.queryParamAuthorizer,
        });

        new SSTIntegrationConstruct(this, 'SSTIntegration', {
            environmentName,
            oauth2ClientConfig,
            config,
        });

        // TODO AGENTS - Remove the construct (class definition and this instantiation) after Cognito switch over (it won't be used).
        new AuthenticationEndpointsConstruct(this, 'AuthenticationEndpoints', {
            environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore: catalogAccessManager,
            archiveStore: archiveAccessManager,
            googleLoginClientId: config.googleLoginClientId,
        });

        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${config.domainName}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });
    }
}
