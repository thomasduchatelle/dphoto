import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {SimpleGoEndpoint} from './simple-go-endpoint';

export interface AuthenticationEndpointsProps {
    environmentName: string;
    apiGateway: { httpApi: apigatewayv2.HttpApi };
    googleLoginClientId: string;
}

export class AuthenticationEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, props: AuthenticationEndpointsProps) {
        super(scope, id);

        new SimpleGoEndpoint(this, 'OAuthToken', {
            environmentName: props.environmentName,
            functionName: 'oauth-token',
            httpApi: props.apiGateway.httpApi,
            path: '/oauth/token',
            method: apigatewayv2.HttpMethod.POST
        });

        new SimpleGoEndpoint(this, 'OAuthLogout', {
            environmentName: props.environmentName,
            functionName: 'oauth-revoke',
            httpApi: props.apiGateway.httpApi,
            path: '/oauth/logout',
            method: apigatewayv2.HttpMethod.POST
        });

        new SimpleGoEndpoint(this, 'EnvConfig', {
            environmentName: props.environmentName,
            functionName: 'env-config',
            httpApi: props.apiGateway.httpApi,
            path: '/env-config.json',
            method: apigatewayv2.HttpMethod.GET,
            environment: {
                GOOGLE_LOGIN_CLIENT_ID: props.googleLoginClientId,
            }
        });
    }
}
