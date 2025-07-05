import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {SimpleGoEndpoint} from './simple-go-endpoint';
import {ApiGatewayConstruct} from './api-gateway';

export interface MetadataEndpointsProps {
    environmentName: string;
    apiGateway: ApiGatewayConstruct;
}

export class MetadataEndpoints extends Construct {
    private readonly notFoundEndpoint: SimpleGoEndpoint;
    private readonly versionEndpoint: SimpleGoEndpoint;

    constructor(scope: Construct, id: string, props: MetadataEndpointsProps) {
        super(scope, id);

        this.notFoundEndpoint = new SimpleGoEndpoint(this, 'NotFound', {
            environmentName: props.environmentName,
            functionName: 'not-found',
            httpApi: props.apiGateway.httpApi,
            path: '/api/{path+}',
            method: apigatewayv2.HttpMethod.ANY
        });

        this.versionEndpoint = new SimpleGoEndpoint(this, 'Version', {
            environmentName: props.environmentName,
            functionName: 'version',
            httpApi: props.apiGateway.httpApi,
            path: '/api/v1/version',
            method: apigatewayv2.HttpMethod.GET
        });
    }
}
