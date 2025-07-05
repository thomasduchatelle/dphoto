import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {StoragesConnectorConstruct} from '../constructs-storages/storages-connector-construct';
import {Duration} from 'aws-cdk-lib';

export interface ArchiveEndpointsProps {
    environmentName: string;
    apiGateway: { httpApi: apigatewayv2.HttpApi };
    context: StoragesConnectorConstruct;
}

export class ArchiveEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, {context, ...props}: ArchiveEndpointsProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.apiGateway.httpApi,
        }

        const getMedia = createSingleRouteEndpoint(this, 'GetMedia', {
            ...endpointProps,
            functionName: 'get-media',
            path: '/api/v1/owners/{owner}/medias/{mediaId}/{filename}',
            method: apigatewayv2.HttpMethod.GET,
            memorySize: 1024,
            timeout: Duration.seconds(29), // maximum allowed by API gateway
        });
        context.grantReadToCatalogTable(getMedia.lambda);
        context.grantReadToStorageAndCache(getMedia.lambda);
        context.grantTriggerOptimisationToArchive(getMedia.lambda);

    }
}
