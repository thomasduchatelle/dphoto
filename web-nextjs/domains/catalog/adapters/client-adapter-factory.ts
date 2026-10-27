import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {getAccessToken} from "@/libs/security/client";

export function newClientAdapterFactory(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(
        getAccessToken,
        async () => `/api/v1`,
    );
}
