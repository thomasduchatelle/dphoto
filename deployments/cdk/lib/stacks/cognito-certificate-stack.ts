import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {LetsEncryptCertificateConstruct} from '../utils/letsencrypt-certificate-construct';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';

export interface CognitoCertificateStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

/** Create certificate in us-east-1 which is required for Cognito custom domain. */
export class CognitoCertificateStack extends cdk.Stack {
    public readonly cognitoCertificate: ICertificate;

    constructor(scope: Construct, id: string, {environmentName, config, ...props}: CognitoCertificateStackProps) {
        super(scope, id, props);

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', environmentName);
        cdk.Tags.of(this).add('Stack', "CognitoCertificateStack");

        const letsEncryptCertificate = new LetsEncryptCertificateConstruct(this, 'CognitoLetsEncryptCertificate', {
            environmentName: `${environmentName}-us-east-1`,
            domainName: config.cognitoDomainName,
            certificateEmail: config.certificateEmail,
            ssmParameterSuffix: 'cognitoDomainCertificationArn',
        });

        this.cognitoCertificate = letsEncryptCertificate.certificate;
    }
}
