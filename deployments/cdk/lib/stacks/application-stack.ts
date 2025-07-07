import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ApiGatewayConstruct} from '../utils/api-gateway-construct';
import {AuthenticationEndpointsConstruct} from '../access/authentication-endpoints-construct';
import {CatalogEndpointsConstruct} from '../catalog/catalog-endpoints-construct';
import {StaticWebUiConstruct} from '../utils/static-web-ui-construct';
import {ArchiveEndpointsConstruct} from "../archive/archive-endpoints-construct";
import {VersionEndpointConstruct} from "../utils/version-endpoint-construct";
import {UserEndpointsConstruct} from "../access/user-endpoints-construct";
import {ArchiveStoreConstruct} from "../archive/archive-store-construct";
import {CatalogStoreConstruct} from "../catalog/catalog-store-construct";
import {ArchivistConstruct} from "../archive/archivist-construct";

export interface DPhotoApplicationStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
    archiveStore: ArchiveStoreConstruct;
    catalogStore: CatalogStoreConstruct;
    archivist: ArchivistConstruct;
}

export class ApplicationStack extends cdk.Stack {
    constructor(scope: Construct, id: string, {config, archiveStore, catalogStore, archivist, ...props}: DPhotoApplicationStackProps) {
        super(scope, id, props);

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoApplicationStack");

        const apiGateway = new ApiGatewayConstruct(this, 'ApiGateway', {
            environmentName: props.environmentName,
            ...config,
        });

        new StaticWebUiConstruct(this, 'StaticWebUi', {
            environmentName: props.environmentName,
            domainName: config.domainName,
            apiGateway: apiGateway
        });

        new VersionEndpointConstruct(this, 'VersionEndpoint', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
        })

        new AuthenticationEndpointsConstruct(this, 'AuthenticationEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore,
            archiveStore,
            googleLoginClientId: config.googleLoginClientId,
        });

        new UserEndpointsConstruct(this, 'UserEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore,
            archiveStore,
            googleLoginClientId: config.googleLoginClientId,
        });

        new CatalogEndpointsConstruct(this, 'CatalogEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            catalogStore,
            archiveMessaging: archivist,
        });

        new ArchiveEndpointsConstruct(this, 'ArchiveEndpoints', {
            environmentName: props.environmentName,
            httpApi: apiGateway.httpApi,
            archiveStore,
            catalogStore,
            archivist: archivist,
        });

        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${config.domainName}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });
    }
}
