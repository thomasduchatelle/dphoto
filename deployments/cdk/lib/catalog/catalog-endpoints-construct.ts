import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {HttpApi} from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint, SimpleGoEndpoint} from '../utils/simple-go-endpoint';
import {CatalogStoreConstruct} from './catalog-store-construct';
import {ArchivistConstruct} from '../archive/archivist-construct';
import {ArchiveStoreConstruct} from "../archive/archive-store-construct";

export interface CatalogEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    catalogStore: CatalogStoreConstruct;
    archiveStore: ArchiveStoreConstruct;
    archiveMessaging: ArchivistConstruct;
}

export class CatalogEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, props: CatalogEndpointsConstructProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.httpApi,
        }

        this.readOnlyCatalogEndpoints(endpointProps, props.catalogStore);
        this.amendTimelineEndpoints(endpointProps, props.catalogStore, props.archiveMessaging, props.archiveStore);
        this.accessControlEndpoints(endpointProps, props.catalogStore);
    }

    private readOnlyCatalogEndpoints(endpointProps: { environmentName: string; httpApi: HttpApi }, catalogStore: CatalogStoreConstruct) {
        const listAlbums = createSingleRouteEndpoint(this, 'ListAlbums', {
            ...endpointProps,
            functionName: 'list-albums',
            path: '/api/v1/albums',
            method: apigatewayv2.HttpMethod.GET,
        });
        catalogStore.grantReadAccess(listAlbums.lambda);

        const listMedias = createSingleRouteEndpoint(this, 'ListMedias', {
            ...endpointProps,
            functionName: 'list-medias',
            path: '/api/v1/owners/{owner}/albums/{folderName}/medias',
            method: apigatewayv2.HttpMethod.GET,
        });
        catalogStore.grantReadAccess(listMedias.lambda);
    }

    private amendTimelineEndpoints(endpointProps: {
        environmentName: string;
        httpApi: HttpApi
    }, catalogStore: CatalogStoreConstruct, archivist: ArchivistConstruct, archiveStore: ArchiveStoreConstruct) {
        const createAlbums = createSingleRouteEndpoint(this, 'CreateAlbums', {
            ...endpointProps,
            functionName: 'create-album',
            path: '/api/v1/albums',
            method: apigatewayv2.HttpMethod.POST,
        });
        catalogStore.grantReadWriteAccess(createAlbums.lambda);
        archiveStore.grantReadAccessToRawAndCacheMedias(createAlbums.lambda);
        archivist.grantAccessToAsyncArchivist(createAlbums.lambda);

        const deleteAlbums = createSingleRouteEndpoint(this, 'DeleteAlbums', {
            ...endpointProps,
            functionName: 'delete-album',
            path: '/api/v1/owners/{owner}/albums/{folderName}',
            method: apigatewayv2.HttpMethod.DELETE,
        });
        catalogStore.grantReadWriteAccess(deleteAlbums.lambda);
        archiveStore.grantReadAccessToRawAndCacheMedias(deleteAlbums.lambda);
        archivist.grantAccessToAsyncArchivist(deleteAlbums.lambda);

        const amendAlbumDates = createSingleRouteEndpoint(this, 'AmendAlbumDates', {
            ...endpointProps,
            functionName: 'amend-album-dates',
            path: '/api/v1/owners/{owner}/albums/{folderName}/dates',
            method: apigatewayv2.HttpMethod.PUT,
        });
        catalogStore.grantReadWriteAccess(amendAlbumDates.lambda);
        archiveStore.grantReadAccessToRawAndCacheMedias(amendAlbumDates.lambda);
        archivist.grantAccessToAsyncArchivist(amendAlbumDates.lambda);

        const amendAlbumName = createSingleRouteEndpoint(this, 'AmendAlbumName', {
            ...endpointProps,
            functionName: 'amend-album-name',
            path: '/api/v1/owners/{owner}/albums/{folderName}/name',
            method: apigatewayv2.HttpMethod.PUT,
        });
        catalogStore.grantReadWriteAccess(amendAlbumName.lambda);
        archiveStore.grantReadAccessToRawAndCacheMedias(amendAlbumName.lambda);
        archivist.grantAccessToAsyncArchivist(amendAlbumName.lambda);
    }

    private accessControlEndpoints(endpointProps: { environmentName: string; httpApi: HttpApi }, catalogStore: CatalogStoreConstruct) {
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
        catalogStore.grantReadWriteAccess(shareAlbum.lambda);
    }
}
