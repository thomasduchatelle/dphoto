import * as cdk from "aws-cdk-lib";
import {Construct} from "constructs";
import {EnvironmentConfig} from "../config/environments";
import {CognitoStackExports} from "../stacks/cognito-stack";

export class SSTIntegrationConstruct extends Construct {
    constructor(scope: Construct, id: string, {environmentName, oauth2ClientConfig, config}: {
        environmentName: string,
        config: EnvironmentConfig;
        oauth2ClientConfig: CognitoStackExports;
    }) {
        super(scope, id);

        // Outputs for SST deployment
        new cdk.CfnOutput(this, 'SSTCloudFrontDomain', {
            value: config.nextjsDomainName,
            description: 'Domain under which the NextJs application must be exposed.',
            exportName: `dphoto-${environmentName}-sst-cloudfront-domain`,
        });

        new cdk.CfnOutput(this, 'SSTCognitoIssuer', {
            value: oauth2ClientConfig.cognitoIssuer,
            description: 'Cognito Issuer URL for SST deployment',
            exportName: `dphoto-${environmentName}-sst-cognito-issuer`,
        });

        new cdk.CfnOutput(this, 'SSTCognitoClientId', {
            value: oauth2ClientConfig.userPoolClientId,
            description: 'Cognito Client ID for SST deployment',
            exportName: `dphoto-${environmentName}-sst-cognito-client-id`,
        });

        new cdk.CfnOutput(this, 'SSTCognitoClientSecret', {
            value: oauth2ClientConfig.userPoolClientSecret.unsafeUnwrap(),
            description: 'Cognito Client Secret for SST deployment',
            exportName: `dphoto-${environmentName}-sst-cognito-client-secret`,
        });
    }

}