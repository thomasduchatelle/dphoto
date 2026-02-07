import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";

export function newClientAdapterFactory(): FetchCatalogAdapter {
    // TODO Get the token from memory, or make a call to the server to get it / refresh it.
    // TODO use default /api/v1 prefix as URL
    // TODO then return the instance of FetchCatalogAdapter.
    throw new Error("newClientAdapterFactory() hasn't been implemented yet.")
}
