import * as cdk from 'aws-cdk-lib';
import {Template} from 'aws-cdk-lib/assertions';
import {ApplicationStack} from './application-stack';
import {environments} from '../config/environments';
import {ArchiveStoreConstruct} from '../archive/archive-store-construct';
import {CatalogStoreConstruct} from '../catalog/catalog-store-construct';
import {ArchivistConstruct} from '../archive/archivist-construct';

// Mock the store constructs to provide test implementations
jest.mock('../archive/archive-store-construct');
jest.mock('../catalog/catalog-store-construct');
jest.mock('../archive/archivist-construct');

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

jest.mock('aws-cdk-lib/aws-s3-deployment', () => {
    const actual = jest.requireActual('aws-cdk-lib/aws-s3-deployment');
    return {
        ...actual,
        Source: {
            ...actual.Source,
            asset: jest.fn().mockImplementation(() => {
                return actual.Source.asset('bin/');
            }),
        }
    };
});

describe('DPhotoApplicationStack', () => {
    let app: cdk.App;
    let stack: ApplicationStack;
    let template: Template;
    let mockArchiveStore: jest.Mocked<ArchiveStoreConstruct>;
    let mockCatalogStore: jest.Mocked<CatalogStoreConstruct>;
    let mockArchivist: jest.Mocked<ArchivistConstruct>;

    beforeEach(() => {
        // Create mock store constructs
        mockArchiveStore = {
            grantReadAccessToRawAndCacheMedias: jest.fn(),
            grantWriteAccessToRawAndCachedMedias: jest.fn(),
        } as any;

        mockCatalogStore = {
            grantReadAccess: jest.fn(),
            grantReadWriteAccess: jest.fn(),
        } as any;

        mockArchivist = {
            grantAccessToAsyncArchivist: jest.fn(),
        } as any;

        app = new cdk.App();
        stack = new ApplicationStack(app, 'TestStack', {
            environmentName: 'test',
            config: environments.test,
            archiveStore: mockArchiveStore,
            catalogStore: mockCatalogStore,
            archivist: mockArchivist,
            env: {
                region: 'eu-west-1',
                account: '0123456789',
            },
        });
        template = Template.fromStack(stack);
    });

    test('lambda for the endpoint /oauth/token has all access to the catalog table (where refresh tokens are stored)', () => {
        const oauthTokenFunction = findLambdaByRoute(template, '/oauth/token', 'POST');

        expect(oauthTokenFunction).toBeDefined();
        expect(mockCatalogStore.grantReadWriteAccess).toHaveBeenCalled();
    });

    test('lambda for the endpoint /env-config.json has the environment variable GOOGLE_LOGIN_CLIENT_ID set', () => {
        const envConfigFunction = findLambdaByRoute(template, '/env-config.json', 'GET');

        expect(envConfigFunction).toBeDefined();

        assertLambdaEnvironmentVariables(envConfigFunction, {
            GOOGLE_LOGIN_CLIENT_ID: environments.test.googleLoginClientId,
        });
    });

    test('catalog endpoints are served by lambdas', () => {
        // Test key catalog endpoints
        const listAlbumsFunction = findLambdaByRoute(template, '/api/v1/albums', 'GET');
        expect(listAlbumsFunction).toBeDefined();

        const createAlbumsFunction = findLambdaByRoute(template, '/api/v1/albums', 'POST');
        expect(createAlbumsFunction).toBeDefined();

        const listMediasFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/albums/{folderName}/medias', 'GET');
        expect(listMediasFunction).toBeDefined();

        const deleteAlbumsFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/albums/{folderName}', 'DELETE');
        expect(deleteAlbumsFunction).toBeDefined();

        const shareAlbumFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/albums/{folderName}/shares/{email}', 'PUT');
        expect(shareAlbumFunction).toBeDefined();
    });

    test('archive endpoints are served by lambdas', () => {
        // Test archive endpoints
        const getMediaFunction = findLambdaByRoute(template, '/api/v1/owners/{owner}/medias/{mediaId}/{filename}', 'GET');
        expect(getMediaFunction).toBeDefined();

        // Verify it has higher memory allocation for media processing
        expect(getMediaFunction.Properties.MemorySize).toBe(1024);
        expect(getMediaFunction.Properties.Timeout).toBe(29);
    });

    test('user endpoints are served by lambdas', () => {
        // Test user endpoints
        const listUsersFunction = findLambdaByRoute(template, '/api/v1/users', 'GET');
        expect(listUsersFunction).toBeDefined();

        const listOwnersFunction = findLambdaByRoute(template, '/api/v1/owners', 'GET');
        expect(listOwnersFunction).toBeDefined();
    });

    test('all required API routes have Lambda Authoriser attached', () => {
        // Define the routes that MUST have authorizer
        const requiredAuthorizedRoutes = [
            { method: 'PUT', path: '/api/v1/owners/{owner}/albums/{folderName}/dates' },      // amend-album-dates
            { method: 'PUT', path: '/api/v1/owners/{owner}/albums/{folderName}/name' },       // amend-album-name
            { method: 'POST', path: '/api/v1/albums' },                                       // create-album
            { method: 'DELETE', path: '/api/v1/owners/{owner}/albums/{folderName}' },         // delete-album
            { method: 'GET', path: '/api/v1/owners/{owner}/medias/{mediaId}/{filename}' },    // get-media
            { method: 'GET', path: '/api/v1/albums' },                                        // list-albums
            { method: 'GET', path: '/api/v1/owners/{owner}/albums/{folderName}/medias' },     // list-medias
            { method: 'GET', path: '/api/v1/owners' },                                        // list-owners
            { method: 'GET', path: '/api/v1/users' },                                         // list-users
            { method: 'PUT', path: '/api/v1/owners/{owner}/albums/{folderName}/shares/{email}' },    // share-album
            { method: 'DELETE', path: '/api/v1/owners/{owner}/albums/{folderName}/shares/{email}' }, // share-album
        ];

        // Define routes that should NOT have authorizer (whitelist)
        const whitelistedRoutes = [
            { method: 'POST', path: '/oauth/token' },
            { method: 'POST', path: '/oauth/logout' },
            { method: 'GET', path: '/env-config.json' },
            { method: 'GET', path: '/version' },
        ];

        // Get all routes from the template
        const allRoutes = template.findResources('AWS::ApiGatewayV2::Route');

        // Check that required routes have authorizer
        requiredAuthorizedRoutes.forEach(({ method, path }) => {
            const routeKey = `${method} ${path}`;
            const matchingRoutes = Object.entries(allRoutes).filter(
                ([_, route]: [string, any]) => route.Properties.RouteKey === routeKey
            );

            expect(matchingRoutes.length).toBe(1);
            const [routeId, route] = matchingRoutes[0];
            
            expect(route.Properties.AuthorizerId).toBeDefined();
            expect(route.Properties.AuthorizerId.Ref).toBeDefined();
        });

        // Check that whitelisted routes do NOT have authorizer
        whitelistedRoutes.forEach(({ method, path }) => {
            const routeKey = `${method} ${path}`;
            const matchingRoutes = Object.entries(allRoutes).filter(
                ([_, route]: [string, any]) => route.Properties.RouteKey === routeKey
            );

            if (matchingRoutes.length > 0) {
                const [routeId, route] = matchingRoutes[0];
                expect(route.Properties.AuthorizerId).toBeUndefined();
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
