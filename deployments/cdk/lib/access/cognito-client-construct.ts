import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import {UserPoolClient, UserPoolDomain} from 'aws-cdk-lib/aws-cognito';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53_targets from 'aws-cdk-lib/aws-route53-targets';
import {Construct} from 'constructs';

export interface CognitoClientConstructProps {
    environmentName: string;
    userPool: cognito.IUserPool;
    rootDomain: string;
    domainName: string;
    cognitoDomainName: string;
    cognitoExtraRedirectDomains: string[]
    certificate: ICertificate;
    cognitoCertificate: ICertificate;
}

export class CognitoClientConstruct extends Construct {
    public readonly userPoolClient: cognito.UserPoolClient;
    public readonly cognitoDomainName: string;

    constructor(scope: Construct, id: string, props: CognitoClientConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;
        this.cognitoDomainName = props.cognitoDomainName;

        // Create User Pool Domain with custom domain
        const userPoolDomain = new UserPoolDomain(this, "UserPoolDomain", {
            customDomain: {
                domainName: props.cognitoDomainName,
                certificate: props.cognitoCertificate,
            },
            userPool: props.userPool
        })

        // Create DNS record for custom domain
        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.rootDomain
        });

        new route53.ARecord(this, 'CognitoDnsRecord', {
            zone: hostedZone,
            recordName: props.cognitoDomainName,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.UserPoolDomainTarget(userPoolDomain)
            )
        });

        // Create User Pool Client with secret for SSR
        this.userPoolClient = new UserPoolClient(this, 'UserPoolClient', {
            userPool: props.userPool,
            userPoolClientName: `${prefix}-web-client`,
            generateSecret: true,
            authFlows: {
                userPassword: false,
                userSrp: false,
                custom: false,
            },
            oAuth: {
                flows: {
                    authorizationCodeGrant: true,
                },
                scopes: [
                    cognito.OAuthScope.EMAIL,
                    cognito.OAuthScope.OPENID,
                    cognito.OAuthScope.PROFILE,
                ],
                callbackUrls: [
                    `https://${props.domainName}/auth/callback`,
                    ...props.cognitoExtraRedirectDomains.map(domain => `https://${domain}/auth/callback`)
                ],
                logoutUrls: [
                    `https://${props.domainName}/`,
                    ...props.cognitoExtraRedirectDomains.map(domain => `https://${domain}/`)
                ],
            },
            supportedIdentityProviders: [
                cognito.UserPoolClientIdentityProvider.GOOGLE,
            ],
            accessTokenValidity: cdk.Duration.hours(1),
            idTokenValidity: cdk.Duration.hours(1),
            refreshTokenValidity: cdk.Duration.days(30),
            preventUserExistenceErrors: true
        });

        // Outputs
        new cdk.CfnOutput(this, 'UserPoolClientId', {
            value: this.userPoolClient.userPoolClientId,
            description: 'Cognito User Pool Client ID',
        });

        new cdk.CfnOutput(this, 'UserPoolDomainOutput', {
            value: userPoolDomain.domainName,
            description: 'Cognito User Pool Domain',
        });
    }
}
