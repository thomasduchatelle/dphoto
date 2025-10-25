import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {IHttpRouteAuthorizer} from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {createSingleRouteEndpoint} from '../utils/simple-go-endpoint';
import {CatalogStoreConstruct} from '../catalog/catalog-store-construct';
import {ArchiveStoreConstruct} from '../archive/archive-store-construct';

export interface AccessEndpointsConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    catalogStore: CatalogStoreConstruct;
    archiveStore: ArchiveStoreConstruct;
    authorizer?: IHttpRouteAuthorizer;
}

export class UserEndpointsConstruct extends Construct {
    constructor(scope: Construct, id: string, props: AccessEndpointsConstructProps) {
        super(scope, id);

        const endpointProps = {
            environmentName: props.environmentName,
            httpApi: props.httpApi,
            authorizer: props.authorizer,
        }

        const listUsers = createSingleRouteEndpoint(this, 'ListUsers', {
            ...endpointProps,
            functionName: 'list-users',
            path: '/api/v1/users',
            method: apigatewayv2.HttpMethod.GET,
        });
        props.catalogStore.grantReadAccess(listUsers.lambda);

        const listOwners = createSingleRouteEndpoint(this, 'ListOwners', {
            ...endpointProps,
            functionName: 'list-owners',
            path: '/api/v1/owners',
            method: apigatewayv2.HttpMethod.GET,
        });
        props.archiveStore.grantReadAccessToRawAndCacheMedias(listOwners.lambda);

    }
}
