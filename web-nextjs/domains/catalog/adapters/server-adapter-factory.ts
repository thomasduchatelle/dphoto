import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {getAccessTokenHolder} from "@/libs/security/session-service";

export function newServerSideRestCatalogAdapter(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(getAccessTokenHolder());
}
