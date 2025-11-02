import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {IHttpRouteAuthorizer} from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {Duration} from 'aws-cdk-lib';
import {ArchiveAccessManager} from "./archive-access-manager";
import {CatalogAccessManager} from "../catalog/catalog-access-manager";
import {ArchivistAccessManager} from "./archivist-access-manager";

export interface ArchiveEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    archiveStore: ArchiveAccessManager;
    catalogStore: CatalogAccessManager;
    archivist: ArchivistAccessManager;
    authorizer?: IHttpRouteAuthorizer;
    queryParamAuthorizer?: IHttpRouteAuthorizer;
}

export class ArchiveEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, props: ArchiveEndpointsConstructProps) {
        super(scope, id);

        const getMedia = createSingleRouteEndpoint(this, 'GetMedia', {
            environmentName: props.environmentName,
            functionName: 'get-media',
            httpApi: props.httpApi,
            path: '/api/v1/owners/{owner}/medias/{mediaId}/{filename}',
            method: apigatewayv2.HttpMethod.GET,
            memorySize: 1024,
            timeout: Duration.seconds(29), // maximum allowed by API gateway
            authorizer: props.queryParamAuthorizer,
        });

        props.catalogStore.grantCatalogReadAccess(getMedia.lambda);
        props.archiveStore.grantReadAccessToRawAndCacheMedias(getMedia.lambda);
        props.archivist.grantAccessToAsyncArchivist(getMedia.lambda);
    }
}
