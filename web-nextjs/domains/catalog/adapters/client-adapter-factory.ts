import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";

class ClientAccessTokenHolder {
    async getAccessToken(): Promise<string | undefined> {
        return undefined;
    }
}

export function newClientAdapterFactory(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(new ClientAccessTokenHolder());
}
