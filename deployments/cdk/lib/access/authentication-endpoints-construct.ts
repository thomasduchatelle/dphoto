import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {CatalogAccessManager} from "../catalog/catalog-access-manager";

export interface AccessEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    catalogStore: CatalogAccessManager;
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
        props.catalogStore.grantCatalogReadWriteAccess(authToken.lambda);

        const logout = createSingleRouteEndpoint(this, 'OAuthLogout', {
            ...endpointProps,
            functionName: 'oauth-revoke',
            path: '/oauth/logout',
            method: apigatewayv2.HttpMethod.POST,
        });
        props.catalogStore.grantCatalogReadWriteAccess(logout.lambda);
    }
}
