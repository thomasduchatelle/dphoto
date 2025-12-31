import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {LetsEncryptCertificateConstruct} from '../utils/letsencrypt-certificate-construct';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';

export interface CloudFrontCertificateStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
}

/** Create certificate in us-east-1 which is required for CloudFront custom domain. */
export class CloudFrontCertificateStack extends cdk.Stack {
    public readonly cloudFrontCertificate: ICertificate;

    constructor(scope: Construct, id: string, props: CloudFrontCertificateStackProps) {
        super(scope, id, {
            ...props,
            crossRegionReferences: true,
        });

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "CloudFrontCertificateStack");

        const letsEncryptCertificate = new LetsEncryptCertificateConstruct(this, 'CloudFrontLetsEncryptCertificate', {
            environmentName: `${props.environmentName}-us-east-1`,
            domainName: props.config.nextjsDomainName,
            certificateEmail: props.config.certificateEmail,
            ssmParameterSuffix: 'cloudFrontDomainCertificationArn',
        });

        this.cloudFrontCertificate = letsEncryptCertificate.certificate;
    }
}
