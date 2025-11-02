import * as cdk from 'aws-cdk-lib';
import {Match, Template} from 'aws-cdk-lib/assertions';
import {CognitoStack} from './cognito-stack';
import {environments} from '../config/environments';

describe('CognitoStack', () => {
    let app: cdk.App;
    let stack: CognitoStack;
    let template: Template;
    let mockCognitoCertificate: cdk.aws_certificatemanager.ICertificate;

    beforeEach(() => {
        app = new cdk.App();
        
        // Create a temporary stack for the certificate import
        const tempStack = new cdk.Stack(app, 'TempStack');
        mockCognitoCertificate = cdk.aws_certificatemanager.Certificate.fromCertificateArn(
            tempStack,
            'MockCognitoCert',
            'arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012'
        );

        stack = new CognitoStack(app, 'TestCognitoStack', {
            environmentName: 'test',
            config: environments.test,
            cognitoCertificate: mockCognitoCertificate,
            env: {
                account: '123456789012',
                region: 'eu-west-1'
            }
        });
        template = Template.fromStack(stack);
    });

    test('creates User Pool with ESSENTIALS feature plan (no advanced security)', () => {
        template.hasResourceProperties('AWS::Cognito::UserPool', {
            UserPoolName: 'dphoto-test-users',
            UserPoolTier: 'ESSENTIALS'
        });
    });

    test('User Pool does not have advanced security mode', () => {
        const userPools = template.findResources('AWS::Cognito::UserPool');
        const userPool = Object.values(userPools)[0];
        
        // ESSENTIALS plan should not have UserPoolAddOns
        expect(userPool.Properties.UserPoolAddOns).toBeUndefined();
    });

    test('creates required user groups', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
            GroupName: 'admin',
            Description: 'Administrators with full system access',
            Precedence: 1
        });

        template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
            GroupName: 'owner',
            Description: 'Content owners with full access to their media',
            Precedence: 2
        });

        template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
            GroupName: 'visitor',
            Description: 'Visitors with limited access to shared albums',
            Precedence: 3
        });
    });

    test('configures Google as identity provider', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolIdentityProvider', {
            ProviderName: 'Google',
            ProviderType: 'Google',
            AttributeMapping: {
                email: 'email',
                given_name: 'given_name',
                family_name: 'family_name',
                picture: 'picture'
            }
        });
    });

    test('creates User Pool Client with OAuth2 configuration', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolClient', {
            UserPoolId: {
                Ref: Match.anyValue()
            },
            GenerateSecret: true,
            AllowedOAuthFlows: ['code'],
            AllowedOAuthScopes: ['email', 'openid', 'profile'],
            SupportedIdentityProviders: ['Google']
        });
    });

    test('creates custom domain for Cognito hosted UI', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolDomain', {
            Domain: 'login.dphoto.example.com',
            CustomDomainConfig: {
                CertificateArn: 'arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012'
            }
        });
    });

    test('creates DNS A record for custom domain', () => {
        template.hasResourceProperties('AWS::Route53::RecordSet', {
            Name: 'login.dphoto.example.com.exmaple.com.',  // Note: test env has typo in rootDomain
            Type: 'A'
        });
    });

    test('configures managed login branding', () => {
        template.hasResourceProperties('AWS::Cognito::ManagedLoginBranding', {
            ReturnMergedResources: true,
            UseCognitoProvidedValues: true
        });
    });

    test('exports correct web environment variables', () => {
        const envVars = stack.getWebEnvironmentVariables();
        expect(envVars.cognitoIssuer).toMatch(/^https:\/\/cognito-idp\./);
        expect(envVars.userPoolClientId).toBeDefined();
        expect(envVars.userPoolClientSecret).toBeDefined();
    });

    test('creates CloudFormation outputs for Cognito resources', () => {
        const outputs = Object.keys(template.findOutputs('*'));
        
        expect(outputs).toContain('UserPoolId');
        expect(outputs).toContain('UserPoolClientId');
        expect(outputs).toContain('CognitoDomain');
        expect(outputs).toContain('CognitoIssuer');
    });

    test('User Pool Client callback URLs include main domain and extra URLs', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolClient', {
            CallbackURLs: [
                'https://dphoto.example.com/auth/callback',
                'http://localhost:3210/auth/callback'
            ]
        });
    });

    test('User Pool Client logout URLs include main domain and extra URLs', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolClient', {
            LogoutURLs: [
                'https://dphoto.example.com/',
                'http://localhost:3210/'
            ]
        });
    });

    test('all resources have correct tags', () => {
        const expectedTagsArray = [
            {Key: 'Application', Value: 'dphoto'},
            {Key: 'CreatedBy', Value: 'cdk'},
            {Key: 'Environment', Value: 'test'},
            {Key: 'Stack', Value: 'CognitoStack'}
        ];

        // Test Cognito User Pool tags (object format)
        const userPools = template.findResources('AWS::Cognito::UserPool');
        Object.values(userPools).forEach((userPool: any) => {
            expectedTagsArray.forEach(tag => {
                expect(userPool.Properties.UserPoolTags).toHaveProperty(tag.Key, tag.Value);
            });
        });
    });
});
