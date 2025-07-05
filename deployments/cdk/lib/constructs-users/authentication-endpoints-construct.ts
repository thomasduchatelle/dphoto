import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {StoragesConnectorConstruct} from '../constructs-storages/storages-connector-construct';

export interface AuthenticationEndpointsProps {
    environmentName: string;
    apiGateway: { httpApi: apigatewayv2.HttpApi };
    context: StoragesConnectorConstruct;
    googleLoginClientId: string;
}

export class AuthenticationEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, {context, ...props}: AuthenticationEndpointsProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.apiGateway.httpApi,
        }

        const authToken = createSingleRouteEndpoint(this, 'OAuthToken', {
            ...endpointProps,
            functionName: 'oauth-token',
            path: '/oauth/token',
            method: apigatewayv2.HttpMethod.POST,
        });
        context.grantRWToCatalogTable(authToken.lambda)


        const logout = createSingleRouteEndpoint(this, 'OAuthLogout', {
            ...endpointProps,
            functionName: 'oauth-revoke',
            path: '/oauth/logout',
            method: apigatewayv2.HttpMethod.POST,
        });
        context.grantRWToCatalogTable(logout.lambda)

        createSingleRouteEndpoint(this, 'EnvConfig', {
            ...endpointProps,
            functionName: 'env-config',
            path: '/env-config.json',
            method: apigatewayv2.HttpMethod.GET,
            environment: {
                GOOGLE_LOGIN_CLIENT_ID: props.googleLoginClientId,
            }
        });
    }
}
