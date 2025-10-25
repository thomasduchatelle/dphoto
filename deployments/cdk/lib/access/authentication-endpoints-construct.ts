import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {CatalogStoreConstruct} from '../catalog/catalog-store-construct';
import {ArchiveStoreConstruct} from '../archive/archive-store-construct';

export interface AccessEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    catalogStore: CatalogStoreConstruct;
    archiveStore: ArchiveStoreConstruct;
}

export class AuthenticationEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, props: AccessEndpointsConstructProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.httpApi,
        }

        const authToken = createSingleRouteEndpoint(this, 'OAuthToken', {
            ...endpointProps,
            functionName: 'oauth-token',
            path: '/oauth/token',
            method: apigatewayv2.HttpMethod.POST,
        });
        props.catalogStore.grantReadWriteAccess(authToken.lambda);

        const logout = createSingleRouteEndpoint(this, 'OAuthLogout', {
            ...endpointProps,
            functionName: 'oauth-revoke',
            path: '/oauth/logout',
            method: apigatewayv2.HttpMethod.POST,
        });
        props.catalogStore.grantReadWriteAccess(logout.lambda);
    }
}
