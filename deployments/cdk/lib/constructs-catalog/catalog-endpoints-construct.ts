import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {HttpApi} from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint, SimpleGoEndpoint} from '../utils/simple-go-endpoint';
import {StoragesConnectorConstruct} from '../constructs-storages/storages-connector-construct';

export interface CatalogEndpointsProps {
    environmentName: string;
    apiGateway: { httpApi: apigatewayv2.HttpApi };
    context: StoragesConnectorConstruct;
}

export class CatalogEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, {context, ...props}: CatalogEndpointsProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.apiGateway.httpApi,
        }

        this.readOnlyCatalogEndpoints(endpointProps, context);

        this.amendTimelineEndpoints(endpointProps, context);

        this.accessControlEndpoints(endpointProps, context);
    }

    private readOnlyCatalogEndpoints(endpointProps: { environmentName: string; httpApi: HttpApi }, context: StoragesConnectorConstruct) {
        const listAlbums = createSingleRouteEndpoint(this, 'ListAlbums', {
            ...endpointProps,
            functionName: 'list-albums',
            path: '/api/v1/albums',
            method: apigatewayv2.HttpMethod.GET,
        });
        context.grantReadToCatalogTable(listAlbums.lambda);

        const listMedias = createSingleRouteEndpoint(this, 'ListMedias', {
            ...endpointProps,
            functionName: 'list-medias',
            path: '/api/v1/owners/{owner}/albums/{folderName}/medias',
            method: apigatewayv2.HttpMethod.GET,
        });
        context.grantReadToCatalogTable(listMedias.lambda);
    }

    private amendTimelineEndpoints(endpointProps: { environmentName: string; httpApi: HttpApi }, context: StoragesConnectorConstruct) {
        const createAlbums = createSingleRouteEndpoint(this, 'CreateAlbums', {
            ...endpointProps,
            functionName: 'create-album',
            path: '/api/v1/albums',
            method: apigatewayv2.HttpMethod.POST,
        });
        context.grantRWToCatalogTable(createAlbums.lambda);
        context.grantTriggerOptimisationToArchive(createAlbums.lambda);

        const deleteAlbums = createSingleRouteEndpoint(this, 'DeleteAlbums', {
            ...endpointProps,
            functionName: 'delete-album',
            path: '/api/v1/owners/{owner}/albums/{folderName}',
            method: apigatewayv2.HttpMethod.DELETE,
        });
        context.grantRWToCatalogTable(deleteAlbums.lambda);
        context.grantTriggerOptimisationToArchive(deleteAlbums.lambda);

        const amendAlbumDates = createSingleRouteEndpoint(this, 'AmendAlbumDates', {
            ...endpointProps,
            functionName: 'amend-album-dates',
            path: '/api/v1/owners/{owner}/albums/{folderName}/dates',
            method: apigatewayv2.HttpMethod.PUT,
        });
        context.grantRWToCatalogTable(amendAlbumDates.lambda);
        context.grantTriggerOptimisationToArchive(amendAlbumDates.lambda);

        const amendAlbumName = createSingleRouteEndpoint(this, 'AmendAlbumName', {
            ...endpointProps,
            functionName: 'amend-album-name',
            path: '/api/v1/owners/{owner}/albums/{folderName}/name',
            method: apigatewayv2.HttpMethod.PUT,
        });
        context.grantRWToCatalogTable(amendAlbumName.lambda);
        context.grantTriggerOptimisationToArchive(amendAlbumName.lambda);
    }

    private accessControlEndpoints(endpointProps: { environmentName: string; httpApi: HttpApi }, context: StoragesConnectorConstruct) {
        const shareAlbum = new SimpleGoEndpoint(this, 'ShareAlbum', {
            ...endpointProps,
            functionName: 'share-album',
            routes: [
                {
                    path: '/api/v1/owners/{owner}/albums/{folderName}/shares/{email}',
                    method: apigatewayv2.HttpMethod.PUT,
                },
                {
                    path: '/api/v1/owners/{owner}/albums/{folderName}/shares/{email}',
                    method: apigatewayv2.HttpMethod.DELETE,
                }
            ]
        });
        context.grantRWToCatalogTable(shareAlbum.lambda);
    }
}
