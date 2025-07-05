import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ApiGatewayConstruct} from '../utils/api-gateway-construct';
import {MetadataEndpointsConstruct} from '../constructs-web/metadata-endpoints-construct';
import {AuthenticationEndpointsConstruct} from '../constructs-users/authentication-endpoints-construct';
import {CatalogEndpointsConstruct} from '../constructs-catalog/catalog-endpoints-construct';
import {StaticWebsiteEndpointConstruct} from '../constructs-web/static-website-endpoint';
import {StoragesConnectorConstruct} from '../constructs-storages/storages-connector-construct';
import {ArchiveEndpointsConstruct} from "../constructs-archive/archive-endpoints-construct";
import {UserEndpointsConstruct} from "../constructs-users/users-endpoints-construct";

export interface DPhotoApplicationStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

export class ApplicationStack extends cdk.Stack {
    constructor(scope: Construct, id: string, {config, ...props}: DPhotoApplicationStackProps) {
        super(scope, id, props);

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoApplicationStack");


        const connector = new StoragesConnectorConstruct(this, 'InfrastructureConfiguration', {
            environmentName: props.environmentName
        });

        const apiGateway = new ApiGatewayConstruct(this, 'ApiGateway', {
            environmentName: props.environmentName,
            ...config,
        });

        new StaticWebsiteEndpointConstruct(this, 'StaticWebsite', {
            environmentName: props.environmentName,
            domainName: config.domainName,
            httpApi: apiGateway.httpApi
        });

        new MetadataEndpointsConstruct(this, 'MetadataEndpoints', {
            environmentName: props.environmentName,
            apiGateway: apiGateway,
        });

        new AuthenticationEndpointsConstruct(this, 'AuthenticationEndpoints', {
            environmentName: props.environmentName,
            apiGateway: apiGateway,
            googleLoginClientId: config.googleLoginClientId,
            context: connector,
        });

        new CatalogEndpointsConstruct(this, 'CatalogEndpoints', {
            environmentName: props.environmentName,
            apiGateway: apiGateway,
            context: connector,
        });

        new ArchiveEndpointsConstruct(this, 'ArchiveEndpoints', {
            environmentName: props.environmentName,
            apiGateway: apiGateway,
            context: connector,
        });

        new UserEndpointsConstruct(this, 'UserEndpoints', {
            environmentName: props.environmentName,
            apiGateway: apiGateway,
            context: connector,
        });

        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${config.domainName}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });
    }
}
