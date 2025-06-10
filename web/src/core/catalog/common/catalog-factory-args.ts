import {DPhotoApplication} from "../../application";
import {CatalogViewerAction} from "../actions";

export interface CatalogFactoryArgs {
    app: DPhotoApplication
    dispatch: (action: CatalogViewerAction) => void
}