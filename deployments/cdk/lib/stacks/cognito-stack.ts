import * as cdk from 'aws-cdk-lib';
import {aws_ssm} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {EnvironmentConfig} from '../config/environments';
import {ICertificate} from 'aws-cdk-lib/aws-certificatemanager';
import * as cognito from "aws-cdk-lib/aws-cognito";
import {ManagedLoginVersion, UserPoolClient} from "aws-cdk-lib/aws-cognito";
import {ARecord, HostedZone, RecordTarget} from "aws-cdk-lib/aws-route53";
import {UserPoolDomainTarget} from "aws-cdk-lib/aws-route53-targets";

export interface CognitoStackProps extends cdk.StackProps {
    environmentName: string;
    config: EnvironmentConfig;
    cognitoCertificate: ICertificate;
}

export interface CognitoStackExports {
    cognitoIssuer: string;
    userPoolClientId: string;
    userPoolClientSecret: cdk.SecretValue;
}

export class CognitoStack extends cdk.Stack {
    public readonly userPool: cognito.UserPool;
    private userPoolClient: UserPoolClient;

    constructor(scope: Construct, id: string, props: CognitoStackProps) {
        super(scope, id, {
            ...props,
            crossRegionReferences: true,
        });

        cdk.Tags.of(this).add('CreatedBy', 'cdk');
        cdk.Tags.of(this).add('Application', 'dphoto');
        cdk.Tags.of(this).add('Environment', props.environmentName);
        cdk.Tags.of(this).add('Stack', "CognitoStack");

        const prefix = `dphoto-${props.environmentName}`;

        this.userPool = this.createUserPool(prefix);
        this.addGroups()

        const googleProvider = this.addGoogleSocialIdentityProviders(props.environmentName, props.config.googleLoginClientId);

        this.addCustomDomain(props.config.rootDomain, props.config.cognitoDomainName, props.cognitoCertificate);

        this.userPoolClient = this.createDPhotoClient(prefix, props.config.domainName, props.config.cognitoExtraRedirectURLs);
        this.userPoolClient.node.addDependency(googleProvider)
    }

    /**
     * Returns environment variables needed by the web application
     */
    public getWebEnvironmentVariables(): CognitoStackExports {
        return {
            cognitoIssuer: `https://cognito-idp.${this.region}.amazonaws.com/${this.userPool.userPoolId}`,
            userPoolClientId: this.userPoolClient.userPoolClientId,
            userPoolClientSecret: this.userPoolClient.userPoolClientSecret!,
        };
    }

    private createUserPool(prefix: string) {
        const userPool = new cognito.UserPool(this, 'UserPool', {
            userPoolName: `${prefix}-users`,
            featurePlan: cognito.FeaturePlan.ESSENTIALS,
            selfSignUpEnabled: false,
            signInAliases: {
                email: true,
            },
            autoVerify: {
                email: true,
            },
            standardAttributes: {
                email: {
                    required: true,
                    mutable: true, // required to be mutable for social identity providers (Google)
                },
                givenName: {
                    required: false,
                    mutable: true,
                },
                familyName: {
                    required: false,
                    mutable: true,
                },
                profilePicture: {
                    required: false,
                    mutable: true,
                },
            },
            passwordPolicy: {
                minLength: 6,
                requireLowercase: false,
                requireUppercase: false,
                requireDigits: false,
                requireSymbols: false,
            },
            accountRecovery: cognito.AccountRecovery.EMAIL_ONLY,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        });

        cdk.Tags.of(userPool).add('Name', `${prefix}-user-pool`);

        new cdk.CfnOutput(this, 'UserPoolId', {
            value: userPool.userPoolId,
            description: 'Cognito User Pool ID',
            exportName: `cognito-user-pool-id`,
        });

        return userPool
    }

    private addGroups() {
        this.userPool.addGroup("Admin", {
            groupName: 'admin',
            description: 'Administrators with full system access',
            precedence: 1,
        });
        this.userPool.addGroup("Owner", {
            groupName: 'owner',
            description: 'Content owners with full access to their media',
            precedence: 2,
        });
        this.userPool.addGroup("Visitor", {
            groupName: 'visitor',
            description: 'Visitors with limited access to shared albums',
            precedence: 3,
        });
    }

    private addGoogleSocialIdentityProviders(environmentName: string, googleLoginClientId: string) {

        // SSM SecretString cannot be used, it triggers the error: "SSM Secure reference is not supported in: [AWS::Cognito::UserPoolIdentityProvider/Properties/ProviderDetails/client_secret]"
        const googleClientSecret = aws_ssm.StringParameter.valueForStringParameter(
            this,
            `/dphoto/cdk-input/googleClientSecret/${environmentName}`,
        );

        // Configure Google Identity Provider
        const googleProvider = new cognito.UserPoolIdentityProviderGoogle(this, 'GoogleProvider', {
            userPool: this.userPool,
            clientId: googleLoginClientId,
            clientSecretValue: cdk.SecretValue.unsafePlainText(googleClientSecret),
            scopes: ['profile', 'email', 'openid'],
            attributeMapping: {
                email: cognito.ProviderAttribute.GOOGLE_EMAIL,
                givenName: cognito.ProviderAttribute.GOOGLE_GIVEN_NAME,
                familyName: cognito.ProviderAttribute.GOOGLE_FAMILY_NAME,
                profilePicture: cognito.ProviderAttribute.GOOGLE_PICTURE,
            },
        });
        this.userPool.registerIdentityProvider(
            googleProvider
        );
        return googleProvider
    }

    private addCustomDomain(rootDomain: string, cognitoDomainName: string, cognitoCertificate: ICertificate) {
        const domain = this.userPool.addDomain("LoginDomain", {
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

    private createDPhotoClient(prefix: string, domainName: string, cognitoExtraRedirectURLs: string[]): UserPoolClient {
        const userPoolClient = this.userPool.addClient('WebUIClient', {
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
                    `https://${domainName}/auth/callback`,
                    ...cognitoExtraRedirectURLs.map(url => `${url}/auth/callback`)
                ],
                logoutUrls: [
                    `https://${domainName}/`,
                    ...cognitoExtraRedirectURLs.map(url => `${url}/`)
                ],
            },
            supportedIdentityProviders: [
                cognito.UserPoolClientIdentityProvider.GOOGLE,
            ],
            accessTokenValidity: cdk.Duration.hours(1),
            idTokenValidity: cdk.Duration.hours(1),
            refreshTokenValidity: cdk.Duration.days(30),
            preventUserExistenceErrors: true,

        });


        new cognito.CfnManagedLoginBranding(this, 'DefaultBranding', {
            userPoolId: this.userPool.userPoolId,
            clientId: userPoolClient.userPoolClientId,
            useCognitoProvidedValues: true,
        });

        return userPoolClient;
    }
}
