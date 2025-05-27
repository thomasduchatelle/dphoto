import {DPhotoApplication} from "../../application";
import {CatalogViewerAction} from "../domain";

export interface CatalogFactoryArgs {
    app: DPhotoApplication
    dispatch: (action: CatalogViewerAction) => void
}