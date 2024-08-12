import {DPhotoApplication} from "../application";
import {MediaViewLoader} from "../catalog";
import {CatalogAPIAdapter} from "../catalog-apis";

// CatalogFactory is stateless as well as any of the constructed instance.
export class CatalogFactory {

    constructor(
        private readonly application: DPhotoApplication
    ) {
    }

    public mediaViewLoader(): MediaViewLoader {
        const restAdapter = this.restAdapter();
        return new MediaViewLoader(restAdapter, restAdapter);
    }

    public restAdapter(): CatalogAPIAdapter {
        return new CatalogAPIAdapter(this.application.axiosInstance, this.application);
    }

}