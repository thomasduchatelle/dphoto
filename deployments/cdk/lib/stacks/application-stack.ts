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
import {ArchiveStoreConstruct} from "../archive/archive-store-construct";
import {CatalogStoreConstruct} from "../catalog/catalog-store-construct";
import {ArchivistConstruct} from "../archive/archivist-construct";
import {LambdaAuthoriserConstruct} from "../access/lambda-authoriser-construct";
import {CognitoUserPoolConstruct} from "../access/cognito-user-pool-construct";
import {CognitoClientConstruct} from "../access/cognito-client-construct";

export interface DPhotoApplicationStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
    archiveStore: ArchiveStoreConstruct;
    catalogStore: CatalogStoreConstruct;
    archivist: ArchivistConstruct;
    cognitoUserPool: CognitoUserPoolConstruct;
    cognitoCertificate: cdk.aws_certificatemanager.ICertificate;
}

export class ApplicationStack extends cdk.Stack {
    constructor(scope: Construct, id: string, {
        config,
        archiveStore,
        catalogStore,
        archivist,
        cognitoUserPool,
        cognitoCertificate,
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

        const cognitoClient = new CognitoClientConstruct(this, 'CognitoClient', {
            environmentName: props.environmentName,
            userPool: cognitoUserPool.userPool,
            cognitoDomainName: config.cognitoDomainName,
            rootDomain: config.rootDomain,
            domainName: config.domainName,
            cognitoExtraRedirectDomains: config.cognitoExtraRedirectDomains,
            cognitoCertificate: cognitoCertificate,
        });

        new WakuWebUiConstruct(this, 'WakuWebUi', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            userPool: cognitoUserPool.userPool,
            userPoolClient: cognitoClient.userPoolClient,
            cognitoDomainName: cognitoClient.cognitoDomainName,
            googleLoginClientId: config.googleLoginClientId,
        });

        new VersionEndpointConstruct(this, 'VersionEndpoint', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
        })

        // Create Lambda Authoriser
        const lambdaAuthoriser = new LambdaAuthoriserConstruct(this, 'LambdaAuthoriser', {
            environmentName: props.environmentName,
            catalogStore,
            cognitoUserPool,
        });

        new AuthenticationEndpointsConstruct(this, 'AuthenticationEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore,
            archiveStore,
        });

        new UserEndpointsConstruct(this, 'UserEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore,
            archiveStore,
            authorizer: lambdaAuthoriser.authorizer,
        });

        new CatalogEndpointsConstruct(this, 'CatalogEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore,
            archiveStore,
            archiveMessaging: archivist,
            authorizer: lambdaAuthoriser.authorizer,
        });

        new ArchiveEndpointsConstruct(this, 'ArchiveEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            archiveStore,
            catalogStore,
            archivist: archivist,
            authorizer: lambdaAuthoriser.authorizer,
            queryParamAuthorizer: lambdaAuthoriser.queryParamAuthorizer,
        });

        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${config.domainName}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });
    }
}
