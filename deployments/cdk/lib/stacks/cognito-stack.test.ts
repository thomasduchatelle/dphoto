import * as cdk from 'aws-cdk-lib';
import {Template} from 'aws-cdk-lib/assertions';
import {CognitoStack} from './cognito-stack';
import {environments} from '../config/environments';
import {FakeCatalogAccessManager} from '../test/fakes/fake-access-managers';

jest.mock('aws-cdk-lib/aws-lambda', () => {
    const actual = jest.requireActual('aws-cdk-lib/aws-lambda');
    return {
        ...actual,
        Code: {
            ...actual.Code,
            fromAsset: jest.fn().mockImplementation(() => {
                return actual.Code.fromAsset('bin/');
            }),
        }
    };
});

describe('CognitoStack', () => {
    let app: cdk.App;
    let stack: CognitoStack;
    let template: Template;
    let fakeCatalogAccessManager: FakeCatalogAccessManager;

    beforeEach(() => {
        app = new cdk.App();
        fakeCatalogAccessManager = new FakeCatalogAccessManager();

        stack = new CognitoStack(app, 'TestCognitoStack', {
            environmentName: 'test',
            config: environments.test,
            catalogStore: fakeCatalogAccessManager,
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

        const userPools = template.findResources('AWS::Cognito::UserPool');
        const userPool = Object.values(userPools)[0];
        expect(userPool.Properties.UserPoolAddOns).toBeUndefined();
    });

    test('creates required user groups: admin, owner, and visitor', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
            GroupName: 'admin',
        });

        template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
            GroupName: 'owner',
        });

        template.hasResourceProperties('AWS::Cognito::UserPoolGroup', {
            GroupName: 'visitor',
        });
    });

    test('attaches a V2 Pre Token Generation trigger to the user pool', () => {
        template.hasResourceProperties('AWS::Cognito::UserPool', {
            LambdaConfig: {
                PreTokenGenerationConfig: {
                    LambdaVersion: 'V2_0',
                },
            },
        });
    });

    test('grants the pre-token Lambda read access to the catalog table', () => {
        expect(fakeCatalogAccessManager.hasBeenGrantedForCatalogRead('dphoto-test-cognito-pre-token')).toBe('');
    });

    test('configures Google as identity provider', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolIdentityProvider', {
            ProviderType: 'Google',
        });
    });

    test('exports correct web environment variables', () => {
        const envVars = stack.getWebEnvironmentVariables();
        expect(envVars.cognitoIssuer).toMatch(/^https:\/\/cognito-idp\./);
        expect(envVars.userPoolClientId).toBeDefined();
        expect(envVars.userPoolClientSecret).toBeDefined();
    });

    test('User Pool Client callback and Logout URLs include main domain and extra URLs', () => {
        template.hasResourceProperties('AWS::Cognito::UserPoolClient', {
            CallbackURLs: [
                'https://dphoto.example.com/nextjs/auth/callback',
                'http://localhost:3210/nextjs/auth/callback'
            ],
            LogoutURLs: [
                'https://dphoto.example.com/nextjs/auth/logout',
                'http://localhost:3210/nextjs/auth/logout'
            ],
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
