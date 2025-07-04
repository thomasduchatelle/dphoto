import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {LetsEncryptLambdaConstruct} from '../constructs/letsencrypt-lambda';
import {ApiGatewayConstruct} from '../constructs/api-gateway';

export interface DPhotoApplicationStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

export class DPhotoApplicationStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props: DPhotoApplicationStackProps) {
        super(scope, id, props);

        // Apply tags to all resources in this stack
        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);

        // Create Let's Encrypt certificate management lambda
        const letsEncryptLambda = new LetsEncryptLambdaConstruct(this, 'LetsEncrypt', {
            environmentName: props.environmentName,
            domainName: this.getDomainName(props.environmentName),
            certificateEmail: 'duchatelle.thomas@gmail.com'
        });

        // Create API Gateway with domain configuration
        const apiGateway = new ApiGatewayConstruct(this, 'ApiGateway', {
            environmentName: props.environmentName,
            domainName: this.getDomainName(props.environmentName),
            rootDomain: this.getRootDomain(props.environmentName)
        });

        // Ensure certificate is created before API Gateway domain
        apiGateway.node.addDependency(letsEncryptLambda);

        // Outputs
        new cdk.CfnOutput(this, 'PublicURL', {
            value: `https://${this.getDomainName(props.environmentName)}`,
            description: 'User friendly HTTPS url where the application has been deployed'
        });
    }

    private getDomainName(environmentName: string): string {
        if (environmentName === 'live') {
            return 'dphoto.duchatelle.net';
        } else if (environmentName === 'next') {
            return 'next.duchatelle.me';
        } else {
            return 'dphoto-dev.duchatelle.net';
        }
    }

    private getRootDomain(environmentName: string): string {
        if (environmentName === 'next') {
            return 'duchatelle.me';
        } else {
            return 'duchatelle.net';
        }
    }
}
