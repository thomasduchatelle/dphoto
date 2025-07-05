import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {SimpleGoEndpoint} from './simple-go-endpoint';
import {ApiGatewayConstruct} from './api-gateway';

export interface MetadataEndpointsProps {
    environmentName: string;
    apiGateway: ApiGatewayConstruct;
}

export class MetadataEndpoints extends Construct {

    constructor(scope: Construct, id: string, props: MetadataEndpointsProps) {
        super(scope, id);

        new SimpleGoEndpoint(this, 'NotFound', {
            environmentName: props.environmentName,
            functionName: 'not-found',
            httpApi: props.apiGateway.httpApi,
            path: '/api/{path+}',
            method: apigatewayv2.HttpMethod.ANY
        });

        new SimpleGoEndpoint(this, 'Version', {
            environmentName: props.environmentName,
            functionName: 'version',
            httpApi: props.apiGateway.httpApi,
            path: '/api/v1/version',
            method: apigatewayv2.HttpMethod.GET
        });
    }
}
