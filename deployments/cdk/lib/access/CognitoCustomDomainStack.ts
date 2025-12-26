import {Construct} from "constructs";
import {IUserPool, ManagedLoginVersion, UserPoolDomain} from "aws-cdk-lib/aws-cognito";
import {ICertificate} from "aws-cdk-lib/aws-certificatemanager";
import {ARecord, HostedZone, RecordTarget} from "aws-cdk-lib/aws-route53";
import {UserPoolDomainTarget} from "aws-cdk-lib/aws-route53-targets";
import * as cdk from "aws-cdk-lib";
import {Stack} from "aws-cdk-lib";

/**
 * Creates a custom domain for a Cognito User Pool and sets up the necessary DNS records.
 *
 * This NEEDS to be done as part of the APPLICATION stack because AWS REQUIRES a ARecord on the parent domain, which is only created with the API Gateway.
 */
export class CognitoCustomDomainStack extends Stack {
    constructor(scope: Construct, id: string, {userPool, cognitoDomainName, rootDomain, cognitoCertificate, ...props}: {
        userPool: IUserPool,
        rootDomain: string,
        cognitoDomainName: string,
        cognitoCertificate: ICertificate,
    } & cdk.StackProps) {
        super(scope, id, props);

        const domain = new UserPoolDomain(this, 'CognitoDomain', {
            userPool: userPool,
            customDomain: {
                domainName: cognitoDomainName,
                certificate: cognitoCertificate,
            },
            managedLoginVersion: ManagedLoginVersion.NEWER_MANAGED_LOGIN,
        });

        // Create DNS record for custom domain
        const hostedZone = HostedZone.fromLookup(this, 'HostedZone', {
            domainName: rootDomain
        });

        new ARecord(this, 'CognitoDnsRecord', {
            zone: hostedZone,
            recordName: cognitoDomainName,
            target: RecordTarget.fromAlias(
                new UserPoolDomainTarget(domain)
            )
        });
    }
}