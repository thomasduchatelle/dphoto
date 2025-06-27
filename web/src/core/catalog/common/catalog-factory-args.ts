import {DPhotoApplication} from "../../application";
import {Action} from "src/libs/daction";
import {CatalogViewerState} from "../language";

export interface CatalogFactoryArgs {
    app: DPhotoApplication;
    dispatch: (action: Action<CatalogViewerState, any>) => void;
}
