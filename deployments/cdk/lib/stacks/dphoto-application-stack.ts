import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ApiGatewayConstruct} from '../constructs/api-gateway';
import {MetadataEndpoints} from '../constructs/metadata-endpoints';
import {StaticWebsiteEndpointConstruct} from '../constructs/static-website-endpoint';

export interface DPhotoApplicationStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

export class DPhotoApplicationStack extends cdk.Stack {
    constructor(scope: Construct, id: string, {config, ...props}: DPhotoApplicationStackProps) {
        super(scope, id, props);

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "DPhotoApplicationStack");


        // Create API Gateway with domain configuration
        const apiGateway = new ApiGatewayConstruct(this, 'ApiGateway', {
            environmentName: props.environmentName,
            ...config,
        });

        // Add metadata endpoints (version, not-found)
        new MetadataEndpoints(this, 'MetadataEndpoints', {
            environmentName: props.environmentName,
            apiGateway: apiGateway
        });

        // Add static website endpoint
        const staticWebsite = new StaticWebsiteEndpointConstruct(this, 'StaticWebsite', {
            environmentName: props.environmentName,
            domainName: config.domainName,
            httpApi: apiGateway.httpApi
        });

        // Outputs
        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${config.domainName}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });

        new cdk.CfnOutput(this, 'ViewerUiBucketName', {
            value: staticWebsite.uiBucket.bucketName,
            description: 'Bucket name where static resources of DPhoto are stored'
        });
    }
}
