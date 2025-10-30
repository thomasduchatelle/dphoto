import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {CognitoUserPoolConstruct} from '../access/cognito-user-pool-construct';
import {CognitoClientConstruct} from '../access/cognito-client-construct';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';

export interface CognitoStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
    cognitoCertificate: ICertificate;
}

export interface CognitoStackExports {
    userPoolId: string;
    userPoolClientId: string;
    userPoolClientSecret: cdk.SecretValue;
    cognitoDomain: string;
    cognitoIssuer: string;
}

export class CognitoStack extends cdk.Stack {
    public readonly userPoolConstruct: CognitoUserPoolConstruct;
    public readonly clientConstruct: CognitoClientConstruct;

    constructor(scope: Construct, id: string, props: CognitoStackProps) {
        super(scope, id, {
            ...props,
            crossRegionReferences: true,
        });

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "CognitoStack");

        // Create User Pool
        this.userPoolConstruct = new CognitoUserPoolConstruct(this, 'CognitoUserPool', {
            environmentName: props.environmentName,
            googleClientId: props.config.googleLoginClientId,
        });

        // Create Client with custom domain
        this.clientConstruct = new CognitoClientConstruct(this, 'CognitoClient', {
            environmentName: props.environmentName,
            userPool: this.userPoolConstruct.userPool,
            cognitoDomainName: props.config.cognitoDomainName,
            rootDomain: props.config.rootDomain,
            domainName: props.config.domainName,
            cognitoExtraRedirectURLs: props.config.cognitoExtraRedirectURLs,
            cognitoCertificate: props.cognitoCertificate,
        });

        // Outputs for easy reference
        new cdk.CfnOutput(this, 'UserPoolId', {
            value: this.userPoolConstruct.userPool.userPoolId,
            description: 'Cognito User Pool ID',
            exportName: `${props.environmentName}-cognito-user-pool-id`,
        });

        new cdk.CfnOutput(this, 'UserPoolClientId', {
            value: this.clientConstruct.userPoolClient.userPoolClientId,
            description: 'Cognito User Pool Client ID',
            exportName: `${props.environmentName}-cognito-user-pool-client-id`,
        });

        new cdk.CfnOutput(this, 'CognitoDomain', {
            value: `https://${this.clientConstruct.cognitoDomainName}`,
            description: 'Cognito Custom Domain URL',
            exportName: `${props.environmentName}-cognito-domain`,
        });

        new cdk.CfnOutput(this, 'CognitoIssuer', {
            value: `https://cognito-idp.${this.region}.amazonaws.com/${this.userPoolConstruct.userPool.userPoolId}`,
            description: 'Cognito Issuer URL',
            exportName: `${props.environmentName}-cognito-issuer`,
        });
    }

    /**
     * Returns environment variables needed by the web application
     */
    public getWebEnvironmentVariables(): Record<string, string> {
        return {
            COGNITO_USER_POOL_ID: this.userPoolConstruct.userPool.userPoolId,
            COGNITO_CLIENT_ID: this.clientConstruct.userPoolClient.userPoolClientId,
            COGNITO_CLIENT_SECRET: this.clientConstruct.userPoolClient.userPoolClientSecret.unsafeUnwrap(),
            COGNITO_DOMAIN: `https://${this.clientConstruct.cognitoDomainName}`,
            COGNITO_ISSUER: `https://cognito-idp.${this.region}.amazonaws.com/${this.userPoolConstruct.userPool.userPoolId}`,
        };
    }
}
