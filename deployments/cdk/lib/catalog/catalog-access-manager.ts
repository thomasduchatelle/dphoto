import {Workload} from "../utils/workload";

export interface CatalogAccessManager {
    grantCatalogReadWriteAccess(grantee: Workload): void;

    grantCatalogReadAccess(authorizerLambda: Workload): void;
}