import {CfnOutput, Stack, StackProps} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {Nextjs} from 'cdk-nextjs-standalone';
import {CognitoStackExports} from './cognito-stack';

export interface AppRouterStackProps extends StackProps {
    cognitoConfig: CognitoStackExports;
}

export class AppRouterStack extends Stack {
    constructor(scope: Construct, id: string, props: AppRouterStackProps) {
        super(scope, id, props);

        const nextjs = new Nextjs(this, 'nextjs', {
            nextjsPath: '../../web-nextjs',
            skipBuild: true,
            streaming: true,
        });

        new CfnOutput(this, "CloudFrontDistributionDomain", {
            value: nextjs.distribution.distributionDomain,
        });

        new CfnOutput(this, "CloudFrontDistributionId", {
            value: nextjs.distribution.distributionId,
            description: 'CloudFront Distribution ID for SST deployment',
            exportName: `${this.stackName}-DistributionId`,
        });

        new CfnOutput(this, "CognitoIssuer", {
            value: props.cognitoConfig.cognitoIssuer,
            description: 'Cognito Issuer URL',
            exportName: `${this.stackName}-CognitoIssuer`,
        });

        new CfnOutput(this, "CognitoClientId", {
            value: props.cognitoConfig.userPoolClientId,
            description: 'Cognito User Pool Client ID',
            exportName: `${this.stackName}-CognitoClientId`,
        });

        new CfnOutput(this, "CognitoClientSecret", {
            value: props.cognitoConfig.userPoolClientSecret.unsafeUnwrap(),
            description: 'Cognito User Pool Client Secret',
            exportName: `${this.stackName}-CognitoClientSecret`,
        });
    }
}