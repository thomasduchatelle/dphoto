import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import {ManagedLoginVersion, UserPoolClient, UserPoolDomain} from 'aws-cdk-lib/aws-cognito';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';
import * as route53 from 'aws-cdk-lib/aws-route53';
import {Construct} from 'constructs';

export interface CognitoClientConstructProps {
    environmentName: string;
    userPool: cognito.IUserPool;
    rootDomain: string;
    domainName: string;
    cognitoDomainName: string;
    cognitoExtraRedirectURLs: string[]
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
            userPool: props.userPool,
            managedLoginVersion: ManagedLoginVersion.NEWER_MANAGED_LOGIN,
        })

        // Create DNS record for custom domain
        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.rootDomain
        });

        new route53.ARecord(this, 'CognitoDnsRecord', {
            zone: hostedZone,
            recordName: props.cognitoDomainName,
            target: route53.RecordTarget.fromAlias({
                bind: (): route53.AliasRecordTargetConfig => ({
                    dnsName: userPoolDomain.cloudFrontEndpoint,
                    hostedZoneId: hostedZone.hostedZoneId,
                })
            })
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
                    ...props.cognitoExtraRedirectURLs.map(url => `${url}/auth/callback`)
                ],
                logoutUrls: [
                    `https://${props.domainName}/`,
                    ...props.cognitoExtraRedirectURLs.map(url => `${url}/`)
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


        new cognito.CfnManagedLoginBranding(this, "ManagedLoginBranding", {
            userPoolId: props.userPool.userPoolId,
            clientId: this.userPoolClient.userPoolClientId,
            returnMergedResources: true,
            useCognitoProvidedValues: true,
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
