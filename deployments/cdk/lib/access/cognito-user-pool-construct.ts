import * as cdk from 'aws-cdk-lib';
import {aws_ssm} from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
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

        cdk.Tags.of(this.userPool).add('Name', `${prefix}-user-pool`);

        // SSM SecretString cannot be used, it triggers the error: "SSM Secure reference is not supported in: [AWS::Cognito::UserPoolIdentityProvider/Properties/ProviderDetails/client_secret]"
        const googleClientSecret = aws_ssm.StringParameter.valueForStringParameter(
            this,
            `/dphoto/cdk-input/googleClientSecret/${props.environmentName}`,
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
}
