import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ApiGatewayConstruct} from '../utils/api-gateway-construct';
import {AuthenticationEndpointsConstruct} from '../access/authentication-endpoints-construct';
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
        config,
        archiveAccessManager,
        catalogAccessManager,
        archivistAccessManager,
        oauth2ClientConfig,
        ...props
    }: DPhotoApplicationStackProps) {
        super(scope, id, {
            ...props,
            crossRegionReferences: true,
        });

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoApplicationStack");

        const apiGateway = new ApiGatewayConstruct(this, 'ApiGateway', {
            environmentName: props.environmentName,
            ...config,
        });

        new WakuWebUiConstruct(this, 'WakuWebUi', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            oauth2ClientConfig: oauth2ClientConfig,
        });

        new VersionEndpointConstruct(this, 'VersionEndpoint', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
        })

        // Create Lambda Authoriser
        const lambdaAuthorizer = new LambdaAuthoriserConstruct(this, 'LambdaAuthoriser', {
            environmentName: props.environmentName,
            catalogStore: catalogAccessManager,
            issuerUrl: oauth2ClientConfig.cognitoIssuer,
        });

        new AuthenticationEndpointsConstruct(this, 'AuthenticationEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore: catalogAccessManager,
        });

        new UserEndpointsConstruct(this, 'UserEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore: catalogAccessManager,
            archiveStore: archiveAccessManager,
            authorizer: lambdaAuthorizer.authorizer,
        });

        new CatalogEndpointsConstruct(this, 'CatalogEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore: catalogAccessManager,
            archiveStore: archiveAccessManager,
            archiveMessaging: archivistAccessManager,
            authorizer: lambdaAuthorizer.authorizer,
        });

        new ArchiveEndpointsConstruct(this, 'ArchiveEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            archiveStore: archiveAccessManager,
            catalogStore: catalogAccessManager,
            archivist: archivistAccessManager,
            authorizer: lambdaAuthorizer.authorizer,
            queryParamAuthorizer: lambdaAuthorizer.queryParamAuthorizer,
        });

        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${config.domainName}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });
    }
}
