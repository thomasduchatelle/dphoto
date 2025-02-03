import {DPhotoApplication} from "../application";
import {CatalogAPIAdapter} from "./adapters/api";

export class CatalogFactory {

    constructor(
        private readonly application: DPhotoApplication
    ) {
    }

    public restAdapter(): CatalogAPIAdapter {
        return new CatalogAPIAdapter(this.application.axiosInstance, this.application);
    }
}