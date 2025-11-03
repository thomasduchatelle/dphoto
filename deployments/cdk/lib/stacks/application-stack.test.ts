import * as cdk from 'aws-cdk-lib';
import {SecretValue} from 'aws-cdk-lib';
import {Template} from 'aws-cdk-lib/assertions';
import {ApplicationStack} from './application-stack';
import {environments} from '../config/environments';
import {computeLetsEncryptHash} from '../utils/letsencrypt-certificate-construct';
import {FakeArchiveAccessManager, FakeArchivistAccessManager, FakeCatalogAccessManager} from '../test/fakes/fake-access-managers';

function functionName(oauthTokenFunction: any) {
    const lambdaName = oauthTokenFunction.Properties.FunctionName || oauthTokenFunction.Properties.Description || 'unknown';
    return lambdaName;
}

describe('DPhotoApplicationStack', () => {
    let app: cdk.App;
    let stack: ApplicationStack;
    let template: Template;
    let fakeArchiveAccessManager: FakeArchiveAccessManager;
    let fakeArchivistAccessManager: FakeArchivistAccessManager;
    let fakeCatalogAccessManager: FakeCatalogAccessManager;

    beforeEach(async () => {
        await computeLetsEncryptHash();
        const mockCognitoCertificate = cdk.aws_certificatemanager.Certificate.fromCertificateArn(
            new cdk.Stack(),
            'MockCognitoCert',
            'arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012'
        );
        app = new cdk.App();
        fakeArchiveAccessManager = new FakeArchiveAccessManager();
        fakeArchivistAccessManager = new FakeArchivistAccessManager();
        fakeCatalogAccessManager = new FakeCatalogAccessManager();
        stack = new ApplicationStack(app, 'TestStack', {
            archiveAccessManager: fakeArchiveAccessManager,
            archivistAccessManager: fakeArchivistAccessManager,
            catalogAccessManager: fakeCatalogAccessManager,
            oauth2ClientConfig: {
                cognitoIssuer: "https://issuer-junit-tests-01.example.com",
                userPoolClientId: "0987654321",
                userPoolClientSecret: new SecretValue("super-secret-value"),
            },
            environmentName: 'test',
            config: environments.test,
            env: {
                region: 'eu-west-1',
                account: '0123456789',
            }
        });
        template = Template.fromStack(stack);
    });

    test('lambda for the endpoint /oauth/token has all access to the catalog table (where refresh tokens are stored)', () => {
        const oauthTokenFunction = findLambdaByRoute(template, '/oauth/token', 'POST');

        expect(oauthTokenFunction).toBeDefined();

        expect(fakeCatalogAccessManager.hasOnlyBeenGrantedCatalogReadWriteTo(functionName(oauthTokenFunction))).toBe('');
    });

    test('catalog endpoints are served by lambdas', () => {
        // Test key catalog endpoints
        const listAlbumsFunction = findLambdaByRoute(template, '/api/v1/albums', 'GET');
        expect(listAlbumsFunction).toBeDefined();
        expect(fakeCatalogAccessManager.hasOnlyBeenGrantedCatalogReadWriteTo(functionName(listAlbumsFunction))).toBe('');


        const createAlbumsFunction = findLambdaByRoute(template, '/api/v1/albums', 'POST');
        expect(createAlbumsFunction).toBeDefined();

        const listMediasFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/albums/{folderName}/medias', 'GET');
        expect(listMediasFunction).toBeDefined();

        const deleteAlbumsFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/albums/{folderName}', 'DELETE');
        expect(deleteAlbumsFunction).toBeDefined();

        const shareAlbumFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/albums/{folderName}/shares/{email}', 'PUT');
        expect(shareAlbumFunction).toBeDefined();
    });

    test('archive get-media endpoint is served by a lambda with read+write access', () => {
        // Test archive endpoints
        const getMediaFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/medias/{mediaId}/{filename}', 'GET');
        expect(getMediaFunction).toBeDefined();

        // Verify it has higher memory allocation for media processing
        expect(getMediaFunction.Properties.MemorySize).toBe(1024);
        expect(getMediaFunction.Properties.Timeout).toBe(29);

        expect(fakeCatalogAccessManager.hasBeenGrantedForCatalogRead(functionName(getMediaFunction))).toBe('');
        expect(fakeCatalogAccessManager.hasOnlyBeenGrantedCatalogReadWriteTo([functionName(getMediaFunction)])).toBe('');
    });

    test('user endpoints are served by lambdas', () => {
        // Test user endpoints
        const listUsersFunction = findLambdaByRoute(template, '/api/v1/users', 'GET');
        expect(listUsersFunction).toBeDefined();

        const listOwnersFunction = findLambdaByRoute(template, '/api/v1/owners', 'GET');
        expect(listOwnersFunction).toBeDefined();
    });

    test('all API routes have Lambda Authoriser attached unless whitelisted', () => {
        // Define routes that should NOT have authorizer (whitelist)
        const whitelistedRoutes = [
            {method: 'POST', path: '/oauth/token'},
            {method: 'POST', path: '/oauth/logout'},
            {method: 'GET', path: '/api/v1/version'},
            {method: 'ANY', path: '/api/{path+}'},
            {method: 'ANY', path: '/{proxy+}'},
            {method: '$default', path: ''},
        ];

        // Get all routes from the template
        const allRoutes = template.findResources('AWS::ApiGatewayV2::Route');

        Object.entries(allRoutes).forEach(([routeId, route]: [string, any]) => {
            const routeKey = route.Properties.RouteKey;
            const [method, ...pathParts] = routeKey.split(' ');
            const path = pathParts.join(' ');
            const isWhitelisted = whitelistedRoutes.some(r => r.method === method && r.path === path);
            try {
                if (isWhitelisted) {
                    expect(route.Properties.AuthorizerId).toBeUndefined();
                } else {
                    // console.log(`routeId: ${method} ${path} [${routeId}]`)
                    expect(route.Properties.AuthorizerId).toBeDefined();
                    expect(route.Properties.AuthorizerId.Ref).toBeDefined();
                }
            } catch (e) {
                throw new Error(`Route ${method} ${path} [${routeId}] failed authorizer check: ${e}`);
            }
        });
    });
});

function assertLambdaEnvironmentVariables(lambdaFunction: any, expectedVariables: Record<string, string>): void {
    const environment = lambdaFunction.Properties.Environment?.Variables;
    expect(environment).toMatchObject(expectedVariables);
}

function getIntegrationId(routeResource: { [p: string]: any }, method: string, path: string) {
    const target = routeResource.Properties?.Target;

    if (!target || !target['Fn::Join']) {
        throw new Error(`Target not found for route: ${method} ${path}`);
    }

    const joinArray = target['Fn::Join'][1];
    const integrationRef = joinArray.find((item: any) => item && item.Ref);
    const integrationId = integrationRef?.Ref;

    if (!integrationId) {
        throw new Error(`Integration ID not found in target for route: ${method} ${path}`);
    }
    return integrationId;
}

function findLambdaByRoute(template: Template, path: string, method: string = 'POST'): any {
    const routes = template.findResources('AWS::ApiGatewayV2::Route', {
        Properties: {
            RouteKey: `${method} ${path}`,
        },
    })

    if (Object.keys(routes).length === 0) {
        throw new Error(`No routes found for ${method} ${path} ; defined routes:\n${JSON.stringify(template.findResources('AWS::ApiGatewayV2::Route'), null, 2)}`);
    }

    const routeResource = routes[Object.keys(routes)[0]];
    const integrationId = getIntegrationId(routeResource, method, path);

    const integrations = template.findResources('AWS::ApiGatewayV2::Integration');
    const integration = integrations[integrationId];

    if (!integration) {
        throw new Error(`Integration resource not found: ${integrationId}`);
    }

    const lambdaArn = integration.Properties?.IntegrationUri;
    if (!lambdaArn || !lambdaArn['Fn::GetAtt']) {
        throw new Error(`Lambda ARN not found in integration: ${integrationId}`);
    }

    const lambdaLogicalId = lambdaArn['Fn::GetAtt'][0];

    const lambdaFunctions = template.findResources('AWS::Lambda::Function');
    const lambdaFunction = lambdaFunctions[lambdaLogicalId];

    if (!lambdaFunction) {
        throw new Error(`Lambda function not found: ${lambdaLogicalId}`);
    }

    return lambdaFunction;
}
