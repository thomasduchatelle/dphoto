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

    constructor(scope: Construct, id: string, props: CognitoCertificateStackProps) {
        super(scope, id, {
            ...props,
            crossRegionReferences: true,
        });

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "CognitoCertificateStack");

        const letsEncryptCertificate = new LetsEncryptCertificateConstruct(this, 'CognitoLetsEncryptCertificate', {
            environmentName: `${props.environmentName}-us-east-1`,
            domainName: props.config.cognitoDomainName,
            certificateEmail: props.config.certificateEmail,
            ssmParameterSuffix: 'cognitoDomainCertificationArn',
        });

        this.cognitoCertificate = letsEncryptCertificate.certificate;
    }
}
