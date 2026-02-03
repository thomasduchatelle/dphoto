import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {ArchiveAccessManager} from "../archive/archive-access-manager";
import {CatalogAccessManager} from "../catalog/catalog-access-manager";

export interface AccessEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    archiveStore: ArchiveAccessManager;
    catalogStore: CatalogAccessManager;
    googleLoginClientId: string;
    jwtEncryptionKey: string;
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

        authToken.lambda.function.addEnvironment('DPHOTO_JWT_KEY_B64', props.jwtEncryptionKey)
        authToken.lambda.function.addEnvironment('DPHOTO_JWT_ISSUER', `https://${props.environmentName}.duchatelle/dphoto`)

        const logout = createSingleRouteEndpoint(this, 'OAuthLogout', {
            ...endpointProps,
            functionName: 'oauth-revoke',
            path: '/oauth/logout',
            method: apigatewayv2.HttpMethod.POST,
        });
        props.catalogStore.grantCatalogReadWriteAccess(logout.lambda);

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