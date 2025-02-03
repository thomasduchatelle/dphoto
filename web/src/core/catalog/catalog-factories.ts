import {DPhotoApplication} from "../application";
import {CatalogViewerLoader} from "./index";
import {CatalogAPIAdapter} from "./adapters/api";
import {MediaPerDayLoader} from "./domain";

export class CatalogFactory {

    constructor(
        private readonly application: DPhotoApplication
    ) {
    }

    public mediaViewLoader(): CatalogViewerLoader {
        const restAdapter = this.restAdapter();
        return new CatalogViewerLoader(restAdapter, new MediaPerDayLoader(restAdapter));
    }

    public restAdapter(): CatalogAPIAdapter {
        return new CatalogAPIAdapter(this.application.axiosInstance, this.application);
    }
}