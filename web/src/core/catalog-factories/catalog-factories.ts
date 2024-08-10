
// CatalogFactory is stateless as well as any of the constructed instance.
import {DPhotoApplication} from "../application";
import {MediaViewLoader} from "../catalog";
import {CatalogAPIAdapter} from "../catalog-apis";

export class CatalogFactory {

    constructor(
        private readonly application: DPhotoApplication
    ) {
    }

    public mediaViewLoader(): MediaViewLoader {

    }

    public restAdapter(): CatalogAPIAdapter {

    }

}