import {DPhotoApplication} from "../../application";
import {Action} from "@/libs/daction";
import {CatalogViewerState} from "../language";

export interface CatalogFactoryArgs {
    app: DPhotoApplication;
    dispatch: (action: Action<CatalogViewerState, any>) => void;
}
