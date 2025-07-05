import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {ApiGatewayConstruct} from '../utils/api-gateway-construct';

export interface MetadataEndpointsProps {
    environmentName: string;
    apiGateway: ApiGatewayConstruct;
}

export class MetadataEndpointsConstruct extends Construct {

    constructor(scope: Construct, id: string, props: MetadataEndpointsProps) {
        super(scope, id);

        createSingleRouteEndpoint(this, 'NotFound', {
            environmentName: props.environmentName,
            functionName: 'not-found',
            httpApi: props.apiGateway.httpApi,
            path: '/api/{path+}',
            method: apigatewayv2.HttpMethod.ANY,
        });

        createSingleRouteEndpoint(this, 'Version', {
            environmentName: props.environmentName,
            functionName: 'version',
            httpApi: props.apiGateway.httpApi,
            path: '/api/v1/version',
            method: apigatewayv2.HttpMethod.GET,
        });
    }
}
