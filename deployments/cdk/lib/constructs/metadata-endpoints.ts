import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {SimpleGoEndpoint} from './simple-go-endpoint';
import {ApiGatewayConstruct} from './api-gateway';

export interface MetadataEndpointsProps {
    environmentName: string;
}

export class MetadataEndpoints extends Construct {
    private readonly notFoundEndpoint: SimpleGoEndpoint;
    private readonly versionEndpoint: SimpleGoEndpoint;

    constructor(scope: Construct, id: string, props: MetadataEndpointsProps) {
        super(scope, id);

        this.notFoundEndpoint = new SimpleGoEndpoint(this, 'NotFound', {
            environmentName: props.environmentName,
            functionName: 'not-found',
            path: '/api/{path+}',
            method: apigatewayv2.HttpMethod.ANY
        });

        this.versionEndpoint = new SimpleGoEndpoint(this, 'Version', {
            environmentName: props.environmentName,
            functionName: 'version',
            path: '/api/v1/version',
            method: apigatewayv2.HttpMethod.GET
        });
    }

    public addToApiGateway(apiGateway: ApiGatewayConstruct): void {
        this.notFoundEndpoint.addToApiGateway(apiGateway.httpApi, '/api/{path+}', apigatewayv2.HttpMethod.ANY);
        this.versionEndpoint.addToApiGateway(apiGateway.httpApi, '/api/v1/version', apigatewayv2.HttpMethod.GET);
    }
}
