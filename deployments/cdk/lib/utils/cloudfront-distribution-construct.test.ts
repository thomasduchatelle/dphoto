import {App, Stack} from 'aws-cdk-lib';
import {Template} from 'aws-cdk-lib/assertions';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {CloudFrontDistributionConstruct} from '../utils/cloudfront-distribution-construct';

describe('CloudFrontDistributionConstruct', () => {
    let app: App;
    let stack: Stack;
    let httpApi: apigatewayv2.HttpApi;

    beforeEach(() => {
        app = new App();
        stack = new Stack(app, 'TestStack', {
            env: {region: 'us-east-1', account: '123456789012'}
        });
        httpApi = new apigatewayv2.HttpApi(stack, 'TestApi');
    });

    test('creates CloudFront distribution with API origin', () => {
        new CloudFrontDistributionConstruct(stack, 'TestCF', {
            environmentName: 'test',
            domainName: 'test.example.com',
            httpApi: httpApi,
        });

        const template = Template.fromStack(stack);
        
        template.hasResourceProperties('AWS::CloudFront::Distribution', {
            DistributionConfig: {
                Enabled: true,
            }
        });
    });

    test('configures /api/* path with no-cache policy', () => {
        new CloudFrontDistributionConstruct(stack, 'TestCF', {
            environmentName: 'test',
            domainName: 'test.example.com',
            httpApi: httpApi,
        });

        const template = Template.fromStack(stack);
        
        template.hasResourceProperties('AWS::CloudFront::CachePolicy', {
            CachePolicyConfig: {
                Comment: 'Policy to never cache API responses',
                DefaultTTL: 0,
                MaxTTL: 0,
                MinTTL: 0,
            }
        });
    });

    test('exports distribution ID as CloudFormation output', () => {
        new CloudFrontDistributionConstruct(stack, 'TestCF', {
            environmentName: 'test',
            domainName: 'test.example.com',
            httpApi: httpApi,
        });

        const template = Template.fromStack(stack);
        
        template.hasOutput('DistributionId', {
            Description: 'CloudFront Distribution ID',
            Export: {
                Name: 'dphoto-test-distribution-id'
            }
        });
    });

    test('configures API origin to forward all headers, cookies, and query strings', () => {
        new CloudFrontDistributionConstruct(stack, 'TestCF', {
            environmentName: 'test',
            domainName: 'test.example.com',
            httpApi: httpApi,
        });

        const template = Template.fromStack(stack);
        
        template.hasResourceProperties('AWS::CloudFront::OriginRequestPolicy', {
            OriginRequestPolicyConfig: {
                Comment: 'Policy to forward all headers, cookies, and query strings to API',
                CookiesConfig: {
                    CookieBehavior: 'all'
                },
                HeadersConfig: {
                    HeaderBehavior: 'all'
                },
                QueryStringsConfig: {
                    QueryStringBehavior: 'all'
                }
            }
        });
    });
});
