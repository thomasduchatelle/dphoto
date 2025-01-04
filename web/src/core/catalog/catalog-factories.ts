import {DPhotoApplication} from "../application";
import {CatalogViewerLoader} from "./index";
import {CatalogAPIAdapter} from "./adapters/api";
import {MediaPerDayLoader, SelectAlbumHandler} from "./domain";

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

    // TODO removes selectAlbumHandler and manage it the same way as albumFilterHandler
    public selectAlbumHandler(): SelectAlbumHandler {
        return new SelectAlbumHandler(new MediaPerDayLoader(this.restAdapter()));
    }
}