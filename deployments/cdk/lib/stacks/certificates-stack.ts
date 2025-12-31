import * as cdk from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {LetsEncryptCertificateConstruct} from '../utils/letsencrypt-certificate-construct';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';

export interface CognitoCertificateStackProps extends cdk.StackProps {
    environmentName: string;
    certificateEmail: string
    domainNames: string[];
}

/** Create certificates in us-east-1: required by Cognito and CloudFront. */
export class CertificatesStack extends cdk.Stack {
    public readonly certificates: Record<string, ICertificate>;

    constructor(scope: Construct, id: string, props: CognitoCertificateStackProps) {
        super(scope, id, {
            ...props,
            crossRegionReferences: true,
        });

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "CertificatesStack");

        let certificates: Record<string, ICertificate> = {};
        for (const domainName of props.domainNames) {
            const normalisedName = domainName.split('.').map(name => name[0].toUpperCase() + name.slice(1)).slice(0, 2).join("");
            const certificateConstruct = new LetsEncryptCertificateConstruct(this, normalisedName, {
                environmentName: `${props.environmentName}-${normalisedName.toLowerCase()}-${this.region}`,
                domainName: domainName,
                certificateEmail: props.certificateEmail,
            });
            certificates[domainName] = certificateConstruct.certificate;
        }

        this.certificates = certificates;
    }
}
