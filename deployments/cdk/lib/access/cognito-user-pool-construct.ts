import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as ssm from 'aws-cdk-lib/aws-ssm';
import {Construct} from 'constructs';

export interface CognitoUserPoolConstructProps {
    environmentName: string;
    googleClientId: string;
}

export class CognitoUserPoolConstruct extends Construct {
    public readonly userPool: cognito.UserPool;
    public readonly adminsGroup: cognito.CfnUserPoolGroup;
    public readonly ownersGroup: cognito.CfnUserPoolGroup;
    public readonly visitorsGroup: cognito.CfnUserPoolGroup;

    constructor(scope: Construct, id: string, props: CognitoUserPoolConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // Create User Pool
        this.userPool = new cognito.UserPool(this, 'UserPool', {
            userPoolName: `${prefix}-users`,
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
                    mutable: false,
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
            removalPolicy: cdk.RemovalPolicy.RETAIN,
            advancedSecurityMode: cognito.AdvancedSecurityMode.ENFORCED,
        });

        cdk.Tags.of(this.userPool).add('Name', `${prefix}-user-pool`);

        // Retrieve Google Client Secret from SSM Parameter Store
        const googleClientSecret = ssm.StringParameter.valueForSecureStringParameter(
            this,
            `dphoto/cdk-input/googleClientSecret/${props.environmentName}`,
            1
        );

        // Configure Google Identity Provider
        new cognito.UserPoolIdentityProviderGoogle(this, 'GoogleProvider', {
            userPool: this.userPool,
            clientId: props.googleClientId,
            clientSecretValue: cdk.SecretValue.unsafePlainText(googleClientSecret),
            scopes: ['profile', 'email', 'openid'],
            attributeMapping: {
                email: cognito.ProviderAttribute.GOOGLE_EMAIL,
                givenName: cognito.ProviderAttribute.GOOGLE_GIVEN_NAME,
                familyName: cognito.ProviderAttribute.GOOGLE_FAMILY_NAME,
                profilePicture: cognito.ProviderAttribute.GOOGLE_PICTURE,
            },
        });

        // Create User Groups
        this.adminsGroup = new cognito.CfnUserPoolGroup(this, 'AdminsGroup', {
            userPoolId: this.userPool.userPoolId,
            groupName: 'admins',
            description: 'Administrators with full system access',
            precedence: 1,
        });

        this.ownersGroup = new cognito.CfnUserPoolGroup(this, 'OwnersGroup', {
            userPoolId: this.userPool.userPoolId,
            groupName: 'owners',
            description: 'Content owners with full access to their media',
            precedence: 2,
        });

        this.visitorsGroup = new cognito.CfnUserPoolGroup(this, 'VisitorsGroup', {
            userPoolId: this.userPool.userPoolId,
            groupName: 'visitors',
            description: 'Visitors with limited access to shared albums',
            precedence: 3,
        });

        // Outputs
        new cdk.CfnOutput(this, 'UserPoolId', {
            value: this.userPool.userPoolId,
            description: 'Cognito User Pool ID',
        });
    }
}
