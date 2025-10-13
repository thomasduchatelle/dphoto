import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import {Construct} from 'constructs';

export interface CognitoUserPoolConstructProps {
    environmentName: string;
    domainName: string;
    googleClientId: string;
    googleClientSecret: string;
}

export class CognitoUserPoolConstruct extends Construct {
    public readonly userPool: cognito.UserPool;
    public readonly userPoolClient: cognito.UserPoolClient;
    public readonly adminGroup: cognito.CfnUserPoolGroup;
    public readonly ownerGroup: cognito.CfnUserPoolGroup;
    public readonly visitorGroup: cognito.CfnUserPoolGroup;
    public readonly sessionTable: dynamodb.Table;

    constructor(scope: Construct, id: string, props: CognitoUserPoolConstructProps) {
        super(scope, id);

        // Create User Pool
        this.userPool = new cognito.UserPool(this, 'UserPool', {
            userPoolName: `dphoto-${props.environmentName}-users`,
            signInAliases: {
                email: true,
            },
            standardAttributes: {
                email: {
                    required: true,
                    mutable: false,
                },
                givenName: {
                    required: true,
                    mutable: true,
                },
                familyName: {
                    required: true,
                    mutable: true,
                },
            },
            autoVerify: {
                email: false, // Google already verifies email
            },
            selfSignUpEnabled: false, // Users must be pre-provisioned
            accountRecovery: cognito.AccountRecovery.NONE,
            removalPolicy: cdk.RemovalPolicy.RETAIN,
        });

        // Configure Google Identity Provider
        const googleProvider = new cognito.UserPoolIdentityProviderGoogle(this, 'GoogleProvider', {
            userPool: this.userPool,
            clientId: props.googleClientId,
            clientSecretValue: cdk.SecretValue.unsafePlainText(props.googleClientSecret),
            scopes: ['profile', 'email', 'openid'],
            attributeMapping: {
                email: cognito.ProviderAttribute.GOOGLE_EMAIL,
                givenName: cognito.ProviderAttribute.GOOGLE_GIVEN_NAME,
                familyName: cognito.ProviderAttribute.GOOGLE_FAMILY_NAME,
            },
        });

        // Create User Pool Client
        this.userPoolClient = new cognito.UserPoolClient(this, 'UserPoolClient', {
            userPool: this.userPool,
            userPoolClientName: `dphoto-${props.environmentName}-web-client`,
            authFlows: {
                userSrp: true,
            },
            oAuth: {
                flows: {
                    authorizationCodeGrant: true,
                },
                scopes: [
                    cognito.OAuthScope.OPENID,
                    cognito.OAuthScope.EMAIL,
                    cognito.OAuthScope.PROFILE,
                ],
                callbackUrls: [
                    `https://${props.domainName}/auth/callback`,
                ],
                logoutUrls: [
                    `https://${props.domainName}/auth/logout`,
                ],
            },
            generateSecret: true,
            accessTokenValidity: cdk.Duration.hours(1),
            refreshTokenValidity: cdk.Duration.days(30),
            idTokenValidity: cdk.Duration.hours(1),
            preventUserExistenceErrors: true,
            supportedIdentityProviders: [
                cognito.UserPoolClientIdentityProvider.GOOGLE,
            ],
        });

        // Ensure the client depends on the Google provider
        this.userPoolClient.node.addDependency(googleProvider);

        // Create User Pool Domain
        const userPoolDomain = this.userPool.addDomain('UserPoolDomain', {
            cognitoDomain: {
                domainPrefix: `dphoto-${props.environmentName}`,
            },
        });

        // Create User Groups
        this.adminGroup = new cognito.CfnUserPoolGroup(this, 'AdminGroup', {
            userPoolId: this.userPool.userPoolId,
            groupName: 'admins',
            description: 'Administrators with full system access',
        });

        this.ownerGroup = new cognito.CfnUserPoolGroup(this, 'OwnerGroup', {
            userPoolId: this.userPool.userPoolId,
            groupName: 'owners',
            description: 'Users with full control over their content',
        });

        this.visitorGroup = new cognito.CfnUserPoolGroup(this, 'VisitorGroup', {
            userPoolId: this.userPool.userPoolId,
            groupName: 'visitors',
            description: 'Users with read access to specific shared albums',
        });

        // Create DynamoDB table for OAuth session state
        this.sessionTable = new dynamodb.Table(this, 'SessionTable', {
            tableName: `dphoto-${props.environmentName}-auth-sessions`,
            partitionKey: {
                name: 'pk',
                type: dynamodb.AttributeType.STRING,
            },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            timeToLiveAttribute: 'ttl',
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        });

        // Outputs
        new cdk.CfnOutput(this, 'UserPoolId', {
            value: this.userPool.userPoolId,
            description: 'Cognito User Pool ID',
        });

        new cdk.CfnOutput(this, 'UserPoolClientId', {
            value: this.userPoolClient.userPoolClientId,
            description: 'Cognito User Pool Client ID',
        });

        new cdk.CfnOutput(this, 'UserPoolDomain', {
            value: userPoolDomain.domainName,
            description: 'Cognito User Pool Domain',
        });

        new cdk.CfnOutput(this, 'SessionTableName', {
            value: this.sessionTable.tableName,
            description: 'DynamoDB table for OAuth session state',
        });
    }
}
