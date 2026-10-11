import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {runtimeApiPrefix} from "@/libs/requests/basepath";

export function newClientAdapterFactory(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(
        async () => undefined,
        async () => `${runtimeApiPrefix}/api/v1`,
        async () => runtimeApiPrefix,
    );
}
