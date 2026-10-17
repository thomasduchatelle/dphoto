import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {authenticatedFetch, getAccessToken} from "@/libs/security/client";
import {runtimeApiPrefix} from "@/libs/requests";

export function newClientAdapterFactory(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(
        getAccessToken,
        async () => `${runtimeApiPrefix}/api/v1`,
        async () => runtimeApiPrefix,
        authenticatedFetch,
    );
}
