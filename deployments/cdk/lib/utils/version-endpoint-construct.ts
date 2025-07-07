import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from './simple-go-endpoint';

export interface AccessEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
}

export class VersionEndpointConstruct extends Construct {
    constructor(scope: Construct, id: string, props: AccessEndpointsConstructProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.httpApi,
        }
        createSingleRouteEndpoint(this, 'Version', {
            ...endpointProps,
            functionName: 'version',
            path: '/api/v1/version',
            method: apigatewayv2.HttpMethod.GET,
        });
    }
}
