import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {StoragesConnectorConstruct} from '../constructs-storages/storages-connector-construct';

export interface UserEndpointsProps {
    environmentName: string;
    apiGateway: { httpApi: apigatewayv2.HttpApi };
    context: StoragesConnectorConstruct;
}

export class UserEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, {context, ...props}: UserEndpointsProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.apiGateway.httpApi,
        }

        const listUsers = createSingleRouteEndpoint(this, 'ListUsers', {
            ...endpointProps,
            functionName: 'list-users',
            path: '/api/v1/users',
            method: apigatewayv2.HttpMethod.GET,
        });
        context.grantReadToCatalogTable(listUsers.lambda);

        const listOwners = createSingleRouteEndpoint(this, 'ListOwners', {
            ...endpointProps,
            functionName: 'list-owners',
            path: '/api/v1/owners',
            method: apigatewayv2.HttpMethod.GET,
        });
        context.grantReadToStorageAndCache(listOwners.lambda);

    }
}
