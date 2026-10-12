import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {authenticatedFetch, runtimeApiPrefix} from "@/libs/requests";

export function newClientAdapterFactory(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(
        async () => undefined,
        async () => `${runtimeApiPrefix}/api/v1`,
        async () => runtimeApiPrefix,
        authenticatedFetch,
    );
}
