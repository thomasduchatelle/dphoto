import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {getAccessTokenHolder} from "@/libs/security";

export class CatalogFactory {
    public restAdapter(): FetchCatalogAdapter {
        return new FetchCatalogAdapter(getAccessTokenHolder());
    }
}
